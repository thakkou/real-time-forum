// ── Helpers ────────────────────────────────────────────────────
function initials(name) {
    return name.slice(0, 2).toUpperCase();
}

function formatTime(date) {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false });
}

function formatDate(date) {
    const now = new Date();
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const d = new Date(date.getFullYear(), date.getMonth(), date.getDate());
    const diff = (today - d) / 86400000;
    if (diff === 0) return 'Today';
    if (diff === 1) return 'Yesterday';
    return date.toLocaleDateString([], { month: 'short', day: 'numeric' });
}

// ── State ──────────────────────────────────────────────────────
let currentUserId = null;
let currentUsername = '';
let isLoadingHistory = false;
let allHistoryLoaded = false;
let messageHistory = {}; // userId -> [{text, mine, ts}]
let ws = null;

// Seed demo data
messageHistory['2'] = [
    { text: 'hey, saw your post about the Go project', mine: false, ts: new Date(Date.now() - 86400000 - 3600000), sender: 'k0r3y' },
    { text: 'did you end up using goroutines for the handlers?', mine: false, ts: new Date(Date.now() - 86400000 - 3540000), sender: 'k0r3y' },
    { text: 'yeah, each request spins up its own goroutine via the stdlib mux', mine: true, ts: new Date(Date.now() - 86400000 - 3300000), sender: 'you' },
    { text: "it's actually pretty clean, no external deps", mine: true, ts: new Date(Date.now() - 86400000 - 3299000), sender: 'you' },
    { text: 'yo did you see the new post?', mine: false, ts: new Date(Date.now() - 7080000), sender: 'k0r3y' },
];

// ── DOM refs ───────────────────────────────────────────────────
const chatMessages  = document.getElementById('chatMessages');
const messageInput  = document.getElementById('messageInput');
const sendBtn       = document.getElementById('sendBtn');
const loadMoreEl    = document.getElementById('loadMoreIndicator');
const wsDot         = document.getElementById('wsDot');
const wsStatusText  = document.getElementById('wsStatusText');
const usersList     = document.getElementById('usersList');
const onlineCount   = document.getElementById('onlineCount');
const userSearch    = document.getElementById('userSearch');
const usersPanel    = document.getElementById('usersPanel');
const mobileToggle  = document.getElementById('mobileUsersToggle');
const backBtn       = document.getElementById('backBtn');
const mobileBar     = document.getElementById('mobileBar');

// ── User list click ────────────────────────────────────────────
usersList.addEventListener('click', e => {
    const item = e.target.closest('.user-item');
    if (!item) return;

    document.querySelectorAll('.user-item').forEach(el => el.classList.remove('active'));
    item.classList.add('active');
    item.classList.remove('unread');

    const userId = item.dataset.userId;
    const username = item.dataset.username;
    const isOnline = item.querySelector('.online-dot').classList.contains('online');

    openChat(userId, username, isOnline);

    // On mobile, close panel
    usersPanel.classList.remove('mobile-open');
});

function openChat(userId, username, isOnline) {
    currentUserId = userId;
    currentUsername = username;
    allHistoryLoaded = false;

    document.getElementById('chatHeaderName').textContent = username;
    document.getElementById('chatAvatarInitials').textContent = initials(username);

    const statusEl = document.getElementById('chatHeaderStatus');
    const dotEl = document.getElementById('chatStatusDot');
    if (isOnline) {
        statusEl.textContent = '● Online';
        statusEl.className = 'chat-header-status online';
        dotEl.className = 'online-dot online';
    } else {
        statusEl.textContent = '○ Offline';
        statusEl.className = 'chat-header-status offline';
        dotEl.className = 'online-dot offline';
    }

    messageInput.placeholder = `Message ${username}...`;

    renderMessages(userId);
    chatMessages.scrollTop = chatMessages.scrollHeight;
}

function renderMessages(userId) {
    // Clear existing messages (keep load indicator)
    while (chatMessages.children.length > 1) {
        chatMessages.removeChild(chatMessages.lastChild);
    }

    const msgs = messageHistory[userId] || [];
    let lastDay = null;

    msgs.forEach(msg => {
        const msgDate = new Date(msg.ts);
        const dayLabel = formatDate(msgDate);

        if (dayLabel !== lastDay) {
        const sep = document.createElement('div');
        sep.className = 'day-separator';
        sep.innerHTML = `<span>${dayLabel}</span>`;
        chatMessages.appendChild(sep);
        lastDay = dayLabel;
        }

        appendMessageEl(msg.text, msg.mine, msg.sender, msg.ts);
    });
}

function appendMessageEl(text, mine, sender, ts) {
    const date = new Date(ts);
    const group = document.createElement('div');
    group.className = `message-group ${mine ? 'mine' : 'theirs'}`;

    const senderEl = document.createElement('div');
    senderEl.className = 'message-sender';
    senderEl.textContent = mine ? 'you' : sender;

    const row = document.createElement('div');
    row.className = 'message-row';

    const bubble = document.createElement('div');
    bubble.className = 'message-bubble';
    bubble.textContent = text;

    const meta = document.createElement('div');
    meta.className = 'message-meta';
    meta.textContent = `${formatDate(date)} · ${formatTime(date)}`;

    row.appendChild(bubble);
    row.appendChild(meta);
    group.appendChild(senderEl);
    group.appendChild(row);
    chatMessages.appendChild(group);
}

// ── Send message ───────────────────────────────────────────────
function sendMessage() {
    const text = messageInput.value.trim();
    if (!text || !currentUserId) return;

    const ts = new Date();
    const msg = { text, mine: true, sender: 'you', ts };

    if (!messageHistory[currentUserId]) messageHistory[currentUserId] = [];
    messageHistory[currentUserId].push(msg);

    const dayLabel = formatDate(ts);
    const prevMsgs = messageHistory[currentUserId];
    const prevDay = prevMsgs.length > 1 ? formatDate(new Date(prevMsgs[prevMsgs.length - 2].ts)) : null;

    if (dayLabel !== prevDay) {
        const sep = document.createElement('div');
        sep.className = 'day-separator';
        sep.innerHTML = `<span>${dayLabel}</span>`;
        chatMessages.appendChild(sep);
    }

    appendMessageEl(text, true, 'you', ts);
    chatMessages.scrollTop = chatMessages.scrollHeight;

    // Update user list preview
    const item = document.querySelector(`.user-item[data-user-id="${currentUserId}"]`);
    if (item) {
        item.querySelector('.user-preview').textContent = text;
        item.querySelector('.user-time').textContent = 'now';
    }

    // Send over WebSocket
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: 'message', toUserId: currentUserId, text }));
    }

    messageInput.value = '';
    messageInput.style.height = 'auto';
}

sendBtn.addEventListener('click', sendMessage);

messageInput.addEventListener('keydown', e => {
    if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();
        sendMessage();
    }
});

// Auto-resize textarea
messageInput.addEventListener('input', () => {
    messageInput.style.height = 'auto';
    messageInput.style.height = Math.min(messageInput.scrollHeight, 120) + 'px';
});

// ── Infinite scroll with throttle ──────────────────────────────
let lastScrollTime = 0;
const SCROLL_THROTTLE = 300; // ms

chatMessages.addEventListener('scroll', () => {
    const now = Date.now();
    if (now - lastScrollTime < SCROLL_THROTTLE) return;
    lastScrollTime = now;

    if (chatMessages.scrollTop < 60 && !isLoadingHistory && !allHistoryLoaded) {
        loadOlderMessages();
    }
});

function loadOlderMessages() {
    isLoadingHistory = true;
    loadMoreEl.classList.add('visible');
    loadMoreEl.textContent = 'Loading older messages...';

    const prevScrollHeight = chatMessages.scrollHeight;

    // Simulate API call — replace with real fetch
    setTimeout(() => {
        // Demo: generate 10 fake older messages
        const msgs = messageHistory[currentUserId] || [];
        const oldestTs = msgs.length ? new Date(msgs[0].ts) : new Date();
        const fakeOlder = [];

        for (let i = 10; i >= 1; i--) {
        const ts = new Date(oldestTs.getTime() - i * 600000);
        fakeOlder.push({
            text: `older message #${i} in history`,
            mine: i % 3 === 0,
            sender: i % 3 === 0 ? 'you' : currentUsername,
            ts,
        });
        }

        // Prepend to history
        messageHistory[currentUserId] = [...fakeOlder, ...msgs];

        // Re-render and restore scroll position
        renderMessages(currentUserId);
        chatMessages.scrollTop = chatMessages.scrollHeight - prevScrollHeight;

        loadMoreEl.classList.remove('visible');
        isLoadingHistory = false;

        // Demo: mark all history loaded after 2 loads
        if (fakeOlder[0].text.includes('#10')) allHistoryLoaded = false;
    }, 800);
}

// ── User search / filter ───────────────────────────────────────
let searchDebounceTimer = null;

userSearch.addEventListener('input', () => {
    clearTimeout(searchDebounceTimer);
    searchDebounceTimer = setTimeout(() => {
        const q = userSearch.value.trim().toLowerCase();
        document.querySelectorAll('.user-item').forEach(item => {
        const name = item.dataset.username.toLowerCase();
        item.style.display = name.includes(q) || !q ? '' : 'none';
        });
    }, 200);
});

// ── Mobile panel toggle ────────────────────────────────────────
mobileToggle.addEventListener('click', () => {
    usersPanel.classList.toggle('mobile-open');
});

backBtn.addEventListener('click', () => {
    usersPanel.classList.add('mobile-open');
});

// ── WebSocket ──────────────────────────────────────────────────
function connectWebSocket() {
    // Replace with your real WS endpoint
    // ws = new WebSocket('wss://yourdomain.com/ws');

    // Demo simulation — remove this block in production and uncomment above
    ws = {
        readyState: 1, // WebSocket.OPEN
        send: (data) => {
        console.log('[WS SEND]', JSON.parse(data));
        // Simulate echo reply after 1s from the other user (demo only)
        const parsed = JSON.parse(data);
        if (parsed.type === 'message') {
            setTimeout(() => simulateIncoming(parsed.toUserId, `got it: "${parsed.text.slice(0, 20)}..."`), 1200);
        }
        },
    };

    setWsStatus('connected');

    /* Real WebSocket event handlers (uncomment when using real WS):

    ws.onopen = () => setWsStatus('connected');

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);

        if (data.type === 'message') {
        const fromUserId = String(data.fromUserId);
        const msg = { text: data.text, mine: false, sender: data.senderName, ts: new Date(data.timestamp) };

        if (!messageHistory[fromUserId]) messageHistory[fromUserId] = [];
        messageHistory[fromUserId].push(msg);

        // If this is the open chat, render it
        if (fromUserId === currentUserId) {
            appendMessageEl(msg.text, false, msg.sender, msg.ts);
            chatMessages.scrollTop = chatMessages.scrollHeight;
        } else {
            // Mark unread
            const item = document.querySelector(`.user-item[data-user-id="${fromUserId}"]`);
            if (item) item.classList.add('unread');
        }

        // Update user preview
        const item = document.querySelector(`.user-item[data-user-id="${fromUserId}"]`);
        if (item) {
            item.querySelector('.user-preview').textContent = data.text;
            item.querySelector('.user-time').textContent = 'now';
        }
        }

        if (data.type === 'user_status') {
        updateUserStatus(data.userId, data.online);
        }

        if (data.type === 'online_users') {
        updateOnlineCount(data.count);
        }
    };

    ws.onclose = () => {
        setWsStatus('disconnected');
        setTimeout(connectWebSocket, 3000); // reconnect
    };

    ws.onerror = () => setWsStatus('error');
    */
}

function simulateIncoming(fromUserId, text) {
    const userId = String(fromUserId);
    const msg = { text, mine: false, sender: currentUsername, ts: new Date() };
    if (!messageHistory[userId]) messageHistory[userId] = [];
    messageHistory[userId].push(msg);

    if (userId === currentUserId) {
        appendMessageEl(msg.text, false, msg.sender, msg.ts);
        chatMessages.scrollTop = chatMessages.scrollHeight;
    } else {
        const item = document.querySelector(`.user-item[data-user-id="${userId}"]`);
        if (item) item.classList.add('unread');
    }

    const item = document.querySelector(`.user-item[data-user-id="${userId}"]`);
    if (item) {
        item.querySelector('.user-preview').textContent = text;
        item.querySelector('.user-time').textContent = 'now';
    }
}

function setWsStatus(state) {
    const states = {
        connected:    { label: 'Connected', cls: 'connected' },
        connecting:   { label: 'Connecting...', cls: 'connecting' },
        disconnected: { label: 'Disconnected — retrying', cls: '' },
        error:        { label: 'Connection error', cls: '' },
    };
    const s = states[state] || states.disconnected;
    wsDot.className = `ws-dot ${s.cls}`;
    wsStatusText.textContent = s.label;
}

function updateUserStatus(userId, online) {
    const item = document.querySelector(`.user-item[data-user-id="${userId}"]`);
    if (!item) return;
    const dot = item.querySelector('.online-dot');
    dot.className = `online-dot ${online ? 'online' : 'offline'}`;

    const onlineItems = document.querySelectorAll('.user-item .online-dot.online').length;
    onlineCount.textContent = `● ${onlineItems} online`;

    if (currentUserId === String(userId)) {
        document.getElementById('chatStatusDot').className = `online-dot ${online ? 'online' : 'offline'}`;
        const s = document.getElementById('chatHeaderStatus');
        s.textContent = online ? '● Online' : '○ Offline';
        s.className = `chat-header-status ${online ? 'online' : 'offline'}`;
    }
}

function updateOnlineCount(count) {
    onlineCount.textContent = `● ${count} online`;
}

// ── Init ───────────────────────────────────────────────────────
connectWebSocket();