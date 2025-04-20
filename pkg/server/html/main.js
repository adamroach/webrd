
window.addEventListener('load', async () => {
    const videoElement = document.getElementById('video');
    const peerConnection = new RTCPeerConnection({
        iceServers: [
            {
              urls: 'stun:stun.l.google.com:19302' // TODO receive from server
            }
          ]
    });
    const websocket = new WebSocket('/ws');

    videoElement.addEventListener('pointermove', (event) => {
        const rect = videoElement.getBoundingClientRect();
        const x = event.offsetX;
        const y = event.offsetY;

        websocket.send(JSON.stringify({
            type: 'mouse_move',
            x: Math.round(x),
            y: Math.round(y)
        }));
    });

    peerConnection.ontrack = (event) => {
        videoElement.srcObject = event.streams[0];
        videoElement.muted = true;
        videoElement.autoplay = true;
        videoElement.play();
    };

    sendKeyEvent = function(event) {
        console.log("Key event", event);
        websocket.send(JSON.stringify({
            type: 'keyboard',
            event: {
                key: event.key,
                code: event.code,
                location: event.location,
                keyDown: event.type === 'keydown',
            }
        }));
        event.preventDefault();
    };
    document.body.addEventListener('keydown', sendKeyEvent);
    document.body.addEventListener('keyup', sendKeyEvent);

    sendMouseButtonEvent = function(event) {
        console.log("Mouse button event", event);
        websocket.send(JSON.stringify({
            type: 'mouse_button',
            button: event.button,
            x: event.clientX,
            y: event.clientY,
            down: event.type === 'mousedown',
        }));
        event.preventDefault();
    };
    videoElement.addEventListener('mousedown', sendMouseButtonEvent);
    videoElement.addEventListener('mouseup', sendMouseButtonEvent);

    videoElement.addEventListener('wheel', (event) => {
        websocket.send(JSON.stringify({
            type: 'mouse_wheel',
            deltaX: event.deltaX,
            deltaY: event.deltaY,
            deltaZ: event.deltaZ,
        }));
        event.preventDefault();
    });

    websocket.onmessage = async (event) => {
        console.log("Received message", event.data)
        const message = JSON.parse(event.data);

        if (message.type === 'offer') {
            await peerConnection.setRemoteDescription(new RTCSessionDescription({
                type: 'offer',
                sdp: message.sdp,
            }));
            const answer = await peerConnection.createAnswer();
            await peerConnection.setLocalDescription(answer);
            console.log("Sending answer", answer);

            websocket.send(JSON.stringify({
                type: 'answer',
                sdp: answer.sdp
            }));
        }
    };

    peerConnection.onicecandidate = (event) => {
        if (event.candidate) {
            console.log("Sending ICE Candidate:" , event)
            websocket.send(JSON.stringify({
                type: 'candidate',
                candidate: event.candidate
            }));
        }
    };

    peerConnection.oniceconnectionstatechange = e => {
        console.log("ICE Connection State:",  peerConnection.iceConnectionState)
    }
});