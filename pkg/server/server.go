package server

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/adamroach/webrd/pkg/capture"
	"github.com/adamroach/webrd/pkg/hid"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

type Server struct {
	MakeVideoCapturer func() (capture.VideoCapturer, error)
	MakeAudioCapturer func() (capture.AudioCapturer, error)
	MakeKeyboard      func() (hid.Keyboard, error)
	MakeMouse         func() (hid.Mouse, error)
	mu                sync.RWMutex // mutex to protect access to sessions
	sessions          map[uuid.UUID]*Session
}

func (s *Server) Run() error {
	const listen = ":8080" // TODO make this configurable
	s.sessions = make(map[uuid.UUID]*Session)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(s, w, r)
	})

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

	fmt.Printf("Server listening on %s\n", listen)
	return http.ListenAndServe(listen, r)
}

func (s *Server) NewSession(messageChannel MessageChannel) (*Session, error) {
	var videoCapturer capture.VideoCapturer
	var audioCapturer capture.AudioCapturer
	var keyboard hid.Keyboard
	var mouse hid.Mouse
	var err error

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

	videoEncoder, err := NewVideoEncoder(videoCapturer, 8_000_000, 30) // TODO make this configurable
	if err != nil {
		return nil, fmt.Errorf("could not create video encoder: %v", err)
	}
	videoSender := NewVideoSender(videoEncoder)
	webRTCConnection, err := NewWebRTCConnection(
		WithVideoSender(videoSender),
		WithTURNServers([]string{"stun:stun.l.google.com:19302", "stun:stun1.l.google.com:19302"}), // TODO make this configurable
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
