package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/adamroach/webrd/pkg/auth"
	"github.com/adamroach/webrd/pkg/capture"
	"github.com/adamroach/webrd/pkg/config"
	"github.com/adamroach/webrd/pkg/hid"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
	"github.com/google/uuid"
)

type Server struct {
	MakeVideoCapturer func() (capture.VideoCapturer, error)
	MakeAudioCapturer func() (capture.AudioCapturer, error)
	MakeKeyboard      func() (hid.Keyboard, error)
	MakeMouse         func() (hid.Mouse, error)
	Authenticator     auth.Authenticator
	mu                sync.RWMutex // mutex to protect access to sessions
	sessions          map[uuid.UUID]*Session
	serverError       chan (error)
	config            *config.Config
}

func (s *Server) Run(config *config.Config) error {
	s.config = config
	s.sessions = make(map[uuid.UUID]*Session)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(httprate.LimitByIP(10, 1*time.Second)) // Prevent password brute-force attacks

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(s, w, r)
	})

	r.Post("/v1/login", s.Login)

	// All other paths serve from the filesystem -- TODO convert to go:embed
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pkg/server/html/index.html")
	})
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		root := http.Dir("./pkg/server/html")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})

	s.serverError = make(chan error)
	for _, bindAddress := range config.BindAddresses {
		go s.listenAndServe(bindAddress, r)
	}
	return <-s.serverError
}

func (s *Server) listenAndServe(address string, r *chi.Mux) {
	log.Printf("Server listening on %s\n", address)
	if s.config.Tls.Enabled {
		err := CheckCert(s.config.Tls.CertFile, s.config.Tls.KeyFile)
		if err != nil {
			log.Printf("TLS certs could not be validated or created: %v\n", err)
			s.serverError <- err
			return
		}
		log.Printf("TLS enabled, using cert %s and key %s\n", s.config.Tls.CertFile, s.config.Tls.KeyFile)
		s.serverError <- http.ListenAndServeTLS(address, s.config.Tls.CertFile, s.config.Tls.KeyFile, r)
	} else {
		log.Printf("DANGER: TLS DISABLED -- THIS ALLOWS ANYONE ON YOUR LOCAL NETWORK TO SPY ON YOUR KEYSTROKES\n")
		s.serverError <- http.ListenAndServe(address, r)
	}
}

func (s *Server) NewSession(messageChannel MessageChannel) (*Session, error) {
	var videoCapturer capture.VideoCapturer
	var audioCapturer capture.AudioCapturer
	var keyboard hid.Keyboard
	var mouse hid.Mouse
	var err error

	err = s.waitForUserAuth(messageChannel)
	if err != nil {
		return nil, err
	}

	if s.MakeVideoCapturer != nil {
		videoCapturer, err = s.MakeVideoCapturer()
		if err != nil {
			return nil, fmt.Errorf("could not create video capturer: %v", err)
		}
	}

	if s.MakeAudioCapturer != nil {
		audioCapturer, err = s.MakeAudioCapturer()
		if err != nil {
			return nil, fmt.Errorf("could not create audio capturer: %v", err)
		}
	}

	if s.MakeKeyboard != nil {
		keyboard, err = s.MakeKeyboard()
		if err != nil {
			return nil, fmt.Errorf("could not create keyboard: %v", err)
		}
	}

	if s.MakeMouse != nil {
		mouse, err = s.MakeMouse()
		if err != nil {
			return nil, fmt.Errorf("could not create mouse: %v", err)
		}
	}

	videoEncoder, err := NewVideoEncoder(videoCapturer, s.config.Video.Bitrate, s.config.Video.Framerate)
	if err != nil {
		return nil, fmt.Errorf("could not create video encoder: %v", err)
	}
	videoSender := NewVideoSender(videoEncoder)
	webRTCConnection, err := NewWebRTCConnection(
		WithVideoSender(videoSender),
		WithICEServers(s.config.IceServers),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create WebRTC connection: %v", err)
	}

	session := &Session{
		ID:               uuid.New(),
		Server:           s,
		WebRTCConnection: webRTCConnection,
		MessageChannel:   messageChannel,
		VideoCapturer:    videoCapturer,
		AudioCapturer:    audioCapturer,
		Keyboard:         keyboard,
		Mouse:            mouse,
	}

	err = session.Start()
	if err != nil {
		return nil, fmt.Errorf("could not start session: %v", err)
	}

	s.mu.Lock()
	if s.sessions == nil {
		s.sessions = make(map[uuid.UUID]*Session)
	}
	s.sessions[session.ID] = session
	s.mu.Unlock()

	return session, nil
}

func (s *Server) waitForUserAuth(messageChannel MessageChannel) error {
	for {
		message, err := messageChannel.Receive()
		if err != nil {
			if err == io.EOF {
				err = errors.New("connection closed before authentication")
				log.Printf("%v\n", err)
				return err
			}
			log.Printf("could not receive message: %v\n", err)
			return err
		}
		m, ok := message.(*AuthMessage)
		if !ok {
			messageChannel.Send(&AuthFailureMessage{
				Type:  TypeAuthFailure,
				Error: "session is not authenticated yet",
			})
			continue
		}
		username, err := s.Authenticator.ValidateToken(m.Token)
		if err != nil {
			log.Printf("could not validate token: %v\n", err)
			err = messageChannel.Send(&AuthFailureMessage{
				Type:  TypeAuthFailure,
				Error: err.Error(),
			})
			if err != nil {
				log.Printf("could not send auth failure message: %v\n", err)
			}
			continue
		}
		log.Printf("user %s authenticated\n", username)
		return nil
	}
}

func (s *Server) GetSession(id uuid.UUID) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}

func (s *Server) EndSession(id uuid.UUID) error {
	s.mu.Lock()
	session, ok := s.sessions[id]
	s.mu.Unlock()
	if !ok {
		return fmt.Errorf("session not found")
	}

	err := session.Close()
	if err != nil {
		return fmt.Errorf("could not close session: %v", err)
	}

	return nil
}

func (s *Server) removeSession(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, session.ID)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	if s.Authenticator == nil {
		http.Error(w, "Authenticator not set", http.StatusInternalServerError)
		return
	}

	var loginBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	username := loginBody.Username
	password := loginBody.Password

	token, err := s.Authenticator.Authenticate(username, password)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"token": token,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("User %s logged in successfully", username)
}
