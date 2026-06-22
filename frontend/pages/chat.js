// import { socket } from '../core/websocket.js';

import { Header } from '../components/Header.js';

export async function render() {
  const header = Header();
  return `
    ${header}
 
    <!-- MESSAGES LAYOUT -->
    <div class="messages-container">
 
      <!-- LEFT PANEL: User list -->
      <aside class="users-panel" id="usersPanel">
        <div class="users-panel-header">
          <h2>Messages</h2>
          <span class="online-count" id="onlineCount">● 0 online</span>
        </div>
 
        <div class="users-search">
          <input
            type="text"
            id="userSearch"
            placeholder="Search users..."
            autocomplete="off"
          />
        </div>
 
        <div class="users-list" id="usersList">
          <!-- SECTION: Online users with recent messages (sorted by last msg) -->
          <div class="section-divider">Online</div>
 
          <!-- Example online user with a message -->
          <div class="user-item unread active" data-user-id="2" data-username="k0r3y">
            <div class="user-avatar">
              <div class="avatar-circle">KR</div>
              <span class="online-dot online"></span>
            </div>
            <div class="user-info">
              <div class="user-name-row">
                <span class="user-name">k0r3y</span>
                <span class="user-time">2m ago</span>
              </div>
              <div class="user-preview">yo did you see the new post?</div>
            </div>
          </div>
 
          <div class="user-item" data-user-id="5" data-username="l3af">
            <div class="user-avatar">
              <div class="avatar-circle">LF</div>
              <span class="online-dot online"></span>
            </div>
            <div class="user-info">
              <div class="user-name-row">
                <span class="user-name">l3af</span>
                <span class="user-time">15m ago</span>
              </div>
              <div class="user-preview">thanks!</div>
            </div>
          </div>
 
          <!-- SECTION: Offline users -->
          <div class="section-divider">Offline</div>
 
          <div class="user-item" data-user-id="3" data-username="n0vax">
            <div class="user-avatar">
              <div class="avatar-circle">NV</div>
              <span class="online-dot offline"></span>
            </div>
            <div class="user-info">
              <div class="user-name-row">
                <span class="user-name">n0vax</span>
                <span class="user-time">2h ago</span>
              </div>
              <div class="user-preview">see you tomorrow</div>
            </div>
          </div>
 
          <div class="user-item" data-user-id="4" data-username="axl">
            <div class="user-avatar">
              <div class="avatar-circle">AX</div>
              <span class="online-dot offline"></span>
            </div>
            <div class="user-info">
              <div class="user-name-row">
                <span class="user-name">axl</span>
                <span class="user-time"></span>
              </div>
              <div class="user-preview" style="color: var(--text-dim); font-style: italic;">No messages yet</div>
            </div>
          </div>
 
          <div class="user-item" data-user-id="6" data-username="byte">
            <div class="user-avatar">
              <div class="avatar-circle">BY</div>
              <span class="online-dot offline"></span>
            </div>
            <div class="user-info">
              <div class="user-name-row">
                <span class="user-name">byte</span>
                <span class="user-time"></span>
              </div>
              <div class="user-preview" style="color: var(--text-dim); font-style: italic;">No messages yet</div>
            </div>
          </div>
        </div>
      </aside>
 
      <!-- RIGHT PANEL: Chat area -->
      <main class="chat-panel" id="chatPanel">
 
        <!-- Mobile toggle button (shown inside chat area on mobile) -->
        <div style="display:flex; padding: 0.5rem 1rem; border-bottom: 1px solid var(--border); background: var(--surface);" id="mobileBar">
          <button class="mobile-users-toggle" id="mobileUsersToggle">
            ☰ Users
          </button>
        </div>
 
        <!-- Chat with k0r3y (default open for demo) -->
        <div id="chatView">
          <div class="chat-header">
            <button class="back-btn" id="backBtn">← Back</button>
            <div class="user-avatar">
              <div class="avatar-circle" id="chatAvatarInitials">KR</div>
              <span class="online-dot online" id="chatStatusDot"></span>
            </div>
            <div class="chat-header-info">
              <div class="chat-header-name" id="chatHeaderName">k0r3y</div>
              <div class="chat-header-status online" id="chatHeaderStatus">● Online</div>
            </div>
          </div>
 
          <div class="chat-messages" id="chatMessages">
            <!-- Load-more sentinel at top -->
            <div class="load-more-indicator" id="loadMoreIndicator">Loading older messages...</div>
 
            <!-- Day separator -->
            <div class="day-separator"><span>Yesterday</span></div>
 
            <!-- Message group: from them -->
            <div class="message-group theirs">
              <div class="message-sender">k0r3y</div>
              <div class="message-row">
                <div class="message-bubble">hey, saw your post about the Go project</div>
                <div class="message-meta">
                  <span>Jun 20 · 18:41</span>
                </div>
              </div>
              <div class="message-row">
                <div class="message-bubble">did you end up using goroutines for the handlers?</div>
                <div class="message-meta">
                  <span>Jun 20 · 18:42</span>
                </div>
              </div>
            </div>
 
            <!-- Message group: from me -->
            <div class="message-group mine">
              <div class="message-sender">you</div>
              <div class="message-row">
                <div class="message-bubble">yeah, each request spins up its own goroutine via the stdlib mux</div>
                <div class="message-meta">Jun 20 · 18:45</div>
              </div>
              <div class="message-row">
                <div class="message-bubble">it's actually pretty clean, no external deps</div>
                <div class="message-meta">Jun 20 · 18:45</div>
              </div>
            </div>
 
            <div class="day-separator"><span>Today</span></div>
 
            <div class="message-group theirs">
              <div class="message-sender">k0r3y</div>
              <div class="message-row">
                <div class="message-bubble">yo did you see the new post?</div>
                <div class="message-meta">Jun 21 · 09:12</div>
              </div>
            </div>
          </div>
 
          <!-- Input -->
          <div class="chat-input-area">
            <div class="chat-input-form">
              <div class="chat-input-wrapper">
                <textarea
                  class="chat-input"
                  id="messageInput"
                  placeholder="Message k0r3y..."
                  rows="1"
                  maxlength="1000"
                ></textarea>
              </div>
              <button class="send-btn" id="sendBtn">Send ↵</button>
            </div>
            <div class="ws-status">
              <span class="ws-dot connected" id="wsDot"></span>
              <span id="wsStatusText">Connected</span>
            </div>
          </div>
        </div>
 
        <!-- Empty state (shown when no user is selected) — hidden by default in demo -->
        <div class="chat-empty" id="chatEmpty" style="display: none;">
          <p>Pick a conversation from the left<br>or start a new one.</p>
        </div>
 
      </main>
    </div>
  `;
}

// document.addEventListener(
//     'ws:new_message',
//     e => {
//         console.log(e.detail);
//     }
// );