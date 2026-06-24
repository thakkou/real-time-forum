// import { socket } from '../core/websocket.js';

import { Header } from '../components/Header.js';

export async function render(data = {}) {
  const header = Header(data.nickname);
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
 
      <!-- RIGHT PANEL: Chat area --><main class="chat-panel" id="chatPanel">
       
        <div id="mobileBar">
          <button class="mobile-users-toggle" id="mobileUsersToggle">
            ☰ Users
          </button>
        </div>

        <div class="chat-empty" id="chatEmpty">
          <p>Pick a conversation from the left<br>or start a new one.</p>
        </div>

        <div id="chatView" style="display:none; flex-direction: column; flex: 1; overflow: hidden;">
          
          <div class="chat-header">
            <button class="back-btn" id="backBtn">← Back</button>

            <div class="user-avatar">
              <div class="avatar-circle" id="chatAvatarInitials"></div>
              <span class="online-dot offline" id="chatStatusDot"></span>
            </div>

            <div class="chat-header-info">
              <div class="chat-header-name" id="chatHeaderName">Select a chat</div>
              <div class="chat-header-status offline" id="chatHeaderStatus">
                ● Offline
              </div>
            </div>
          </div>

          <div class="chat-messages" id="chatMessages"></div>

          <div class="chat-input-area">
            <div class="chat-input-form">
              <div class="chat-input-wrapper">
                <textarea
                  class="chat-input"
                  id="messageInput"
                  placeholder="Message..."
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
      </main>
    </div>
  `;
}
