class Client {
    constructor() {
        this.auth = new Auth();
        this.websocket = new WebSocket("/ws");
        this.videoElement = document.getElementById("video");
        this.peerConnection = null;
        this.authed = false;
    }

    async start() {
        try {
            await this.startWebsocket();
        } catch (e) {
            alert(`Could not set up websocket: ${e}`);
        }
        await this.login();
    }

    async login(message) {
        while (!this.authed) {
            try {
                let e = await this.authUser(message);
                this.websocket.send(
                    JSON.stringify({
                        type: "auth",
                        token: e.token,
                    }),
                );
                this.authed = true;
            } catch (e) {
                console.log("could not authenticate user:", e);
                message = `<font color="red">${e.reason}</font>`;
            }
        }
    }

    async authUser(message) {
        return new Promise((accept, reject) => {
            this.auth.on("login", accept);
            this.auth.on("failure", reject);
            this.auth.login(message);
        });
    }

    async startWebsocket() {
        return new Promise((accept, reject) => {
            this.websocket.addEventListener("open", accept);
            this.websocket.addEventListener("error", reject);
            this.websocket.addEventListener(
                "message",
                this.handleMessage.bind(this),
            );
        });
    }

    async setupPeerConnection(offer) {
        this.peerConnection = new RTCPeerConnection({
            iceServers: offer.iceServers,
        });

        this.peerConnection.ontrack = (event) => {
            this.videoElement.srcObject = event.streams[0];
            this.videoElement.muted = true;
            this.videoElement.autoplay = true;
            this.videoElement.play();
        };

        this.peerConnection.onicecandidate = (event) => {
            if (event.candidate) {
                console.log("Sending ICE Candidate:", event.candidate);
                this.websocket.send(
                    JSON.stringify({
                        type: "candidate",
                        candidate: event.candidate,
                    }),
                );
            }
        };

        this.peerConnection.oniceconnectionstatechange = (e) => {
            console.log(
                "ICE Connection State:",
                this.peerConnection.iceConnectionState,
            );
        };

        await this.peerConnection.setRemoteDescription(
            new RTCSessionDescription({
                type: "offer",
                sdp: offer.sdp,
            }),
        );
        const answer = await this.peerConnection.createAnswer();
        await this.peerConnection.setLocalDescription(answer);

        return {
            type: "answer",
            sdp: answer.sdp,
        };
    }

    captureInput() {
        this.videoElement.addEventListener("pointermove", (event) => {
            const rect = this.videoElement.getBoundingClientRect();
            const x = event.offsetX;
            const y = event.offsetY;

            this.websocket.send(
                JSON.stringify({
                    type: "mouse_move",
                    x: Math.round(x),
                    y: Math.round(y),
                }),
            );
        });

        const sendKeyEvent = (event) => {
            console.log("Key event", event);
            this.websocket.send(
                JSON.stringify({
                    type: "keyboard",
                    event: {
                        key: event.key,
                        code: event.code,
                        location: event.location,
                        keyDown: event.type === "keydown",
                    },
                }),
            );
            event.preventDefault();
        };
        document.body.addEventListener("keydown", sendKeyEvent);
        document.body.addEventListener("keyup", sendKeyEvent);

        const sendMouseButtonEvent = (event) => {
            console.log("Mouse button event", event);
            this.websocket.send(
                JSON.stringify({
                    type: "mouse_button",
                    button: event.button,
                    x: event.clientX,
                    y: event.clientY,
                    down: event.type === "mousedown",
                }),
            );
            event.preventDefault();
        };
        this.videoElement.addEventListener("mousedown", sendMouseButtonEvent);
        this.videoElement.addEventListener("mouseup", sendMouseButtonEvent);

        this.videoElement.addEventListener("wheel", (event) => {
            this.websocket.send(
                JSON.stringify({
                    type: "mouse_wheel",
                    deltaX: event.deltaX,
                    deltaY: event.deltaY,
                    deltaZ: event.deltaZ,
                }),
            );
            event.preventDefault();
        });
    }

    async handleMessage(event) {
        console.log("Received message", event.data);
        const message = JSON.parse(event.data);

        switch (message.type) {
            case "offer":
                const answer = await this.setupPeerConnection(message);
                console.log("Sending answer", answer);
                this.websocket.send(JSON.stringify(answer));
                if (answer.type === "answer") {
                    this.captureInput(); // maybe wait until after connection succeeds?
                }
                break;
            case "auth_failure":
                auth.reset();
                await this.login(`<font color="red">${message.error}</font>`);
                break;
            default:
                console.log("Received unexpected message type:", message);
        }
    }
}
