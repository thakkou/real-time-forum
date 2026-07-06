// Runs in sepaRATE EXECUTION CONTEXT

const ports = []; // new Set(); -> must change methods and use [...x] to use it ! 
let socket = null;
let onlineUsersCache = []; // <-- new: worker-level source of truth

function broadcast(message) {
    for (let i = ports.length - 1; i >= 0; i--) {
        try {
            ports[i].postMessage(message);
        } catch (e) {
            console.warn('Dead port, removing', e);
            ports.splice(i, 1);
        }
    }
}

function connect(wsUri) {
    if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
        return;
    }
    socket = new WebSocket(wsUri);

    socket.onopen = () => {
        console.log('Worker WS connected');
        broadcast({ type: '__open' });
    };
    socket.onclose = (event) => {
        console.log('Worker WS disconnected — code:', event.code, 'reason:', event.reason);
        broadcast({ type: event.reason === 'logged out' ? '__logged_out' : '__close' });
        socket = null;
    };
    socket.onerror = (err) => {
        console.error(err);
    };
    socket.onmessage = (event) => {
        let msg;
        try {
            msg = JSON.parse(event.data);
        } catch (e) {
            console.error('Bad WS payload:', event.data);
            return;
        }

        // keep the cache in sync with server events (added part to handle new tabs)
        if (msg.event_type === 'init') {
            onlineUsersCache = msg.data;
        } else if (msg.event_type === 'client_connect') {
            if (!onlineUsersCache.includes(msg.data)) {
                onlineUsersCache.push(msg.data);
            }
        } else if (msg.event_type === 'client_disconnect') {
            onlineUsersCache = onlineUsersCache.filter(id => id !== msg.data);
        }
        // end

        broadcast({ type: '__message', payload: msg });
    };
}

onconnect = (event) => {
    const port = event.ports[0];
    port.id = crypto.randomUUID();
    ports.push(port);
    port.start();

    port.onmessage = (e) => {
        const msg = e.data;
        switch (msg.type) {
            case 'connect':
                if (socket && socket.readyState === WebSocket.OPEN) {
                    // socket already exists — this port missed the real "init"
                    // event from the server, so replay our cached state to it (second added part)
                    port.postMessage({ type: '__open' });
                    port.postMessage({
                        type: '__message',
                        payload: { event_type: 'init', data: onlineUsersCache }
                    });
                    // end
                } else {
                    connect(msg.wsUri);
                }
                break;
            case 'send':
                console.log('[worker] port received message:', e.data);
                if (socket && socket.readyState === WebSocket.OPEN) {
                    socket.send(JSON.stringify(msg.payload));
                    // echo to this user's OTHER tabs immediately (no server round-trip)
                    if (msg.payload.event_type === 'new_message') {
                        console.log(12)
                        ports.forEach(p => {
                            if (p === port) return; // skip the sender's own tab
                            try {
                                p.postMessage({
                                    type: '__message',
                                    payload: {
                                        event_type: 'message_sent',
                                        data: msg.payload.data
                                    }
                                });
                            } catch (err) {
                                console.warn('Dead port during echo, will be pruned on next broadcast', err);
                            }
                        });
                    }
                }
                break;
            case 'logout':
                console.log('[worker] logout requested, closing socket');
                if (socket) {
                    socket.close(1000, 'user logout');
                    socket = null;
                }
                onlineUsersCache = [];
                broadcast({ type: '__logged_out' }); // <-- distinct from __close
                break;
            case 'disconnect':
                const idx = ports.indexOf(port);
                if (idx !== -1) ports.splice(idx, 1);
                if (ports.length === 0 && socket) {
                    socket.close();
                    socket = null;
                }
                break;
        }
    };
};