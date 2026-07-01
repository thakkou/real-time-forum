import { formatTime } from './helpers.js';
import { getConversations, getConversationById } from "../api/conversations.js";
import { createMessage } from "../api/messages.js";
import { onlineUsers } from "../services/router.js";
import { Conversation } from "../components/Conversation.js";
import { ws } from '../services/websocket.js';

/* =========================
   STATE MANAGEMENT
========================= */

const state = {
  currentConversationId: null,
  currentReceiverId: null,
  isTyping: false,         // Tracks if the counter-party is typing
  showTyping: false,        // FORCE OVERRIDE: Set to true to see the animation all the time!
  
  // 🟢 NEW LOCAL TYPING TRACKERS
  isSelfTyping: false,     // Tracks if YOU are currently typing
  selfTypingTimeout: null, // References the debouncer timer instance
  
  // 🔄 Pagination States
  isLoadingOlder: false,
  hasMoreMessages: true,
  limit: 10,
  offset: 0,
};

/* =========================
   DOM ELEMENTS CACHE
========================= */
const dom = {
  usersList: null,
  chatMessages: null,
  chatHeaderName: null,
  chatHeaderStatus: null,
  messageInput: null,
  sendBtn: null,
  chatEmpty: null,
  chatView: null,
};

const HTML_CHARS = { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#039;" };

/* =========================
   INITIALIZATION
========================= */
export async function setup() {
  cacheDom();

  if (!dom.usersList) {
    console.error("Chat DOM not found — page not rendered yet");
    return;
  }

  setupMobileUI();
  setupSendMessage();
  setupScrollPagination();

  try {
    const resp = await getConversations();
    const data = resp.data || [];
    renderUsers(data);

    if (data.length > 0) {
      await openConversation(data[0]);
    } else {
      showEmptyState();
    }
  } catch (err) {
    console.error("Failed to initialize conversation list:", err);
    showEmptyState();
  }
}

/* =========================
   LOCAL USER TYPING EMITTER
========================= */
function handleLocalTypingActivity() {
  if (!state.currentReceiverId || !state.currentConversationId) return;

  if (!state.isSelfTyping) {
    state.isSelfTyping = true;
console.log('start typing')
console.log(window.profile)
    ws.send({
      event_type: "typing:start",
      data: {
        userId:window.profile.id,
        conversationId: state.currentConversationId,
        receiverId: state.currentReceiverId,
      }
    });
  }

  clearTimeout(state.selfTypingTimeout);

  state.selfTypingTimeout = setTimeout(() => {
    stopLocalTypingNotification();
  }, 1500);
}

function stopLocalTypingNotification() {
  if (!state.isSelfTyping) return;
  state.isSelfTyping = false;

  clearTimeout(state.selfTypingTimeout);
  console.log('stop typing')

  ws.send({
    event_type: "typing:stop",
    data: {
              userId:window.profile.id,

      conversationId: state.currentConversationId,
      receiverId: state.currentReceiverId,
    }
  })
}

function cacheDom() {
  dom.usersList = document.getElementById("usersList");
  dom.chatMessages = document.getElementById("chatMessages");
  dom.chatHeaderName = document.getElementById("chatHeaderName");
  dom.chatHeaderStatus = document.getElementById("chatHeaderStatus");
  dom.messageInput = document.getElementById("messageInput");
  dom.sendBtn = document.getElementById("sendBtn");
  dom.chatEmpty = document.getElementById("chatEmpty");
  dom.chatView = document.getElementById("chatView");
}

/* =========================
   MOBILE RESPONSIVENESS
========================= */
function setupMobileUI() {
  const usersPanel = document.getElementById("usersPanel");
  const mobileUsersToggle = document.getElementById("mobileUsersToggle");
  const backBtn = document.getElementById("backBtn");

  mobileUsersToggle?.addEventListener("click", () => usersPanel?.classList.add("mobile-open"));
  backBtn?.addEventListener("click", () => usersPanel?.classList.remove("mobile-open"));
}

/* =========================
   USER ROSTER RENDERING
========================= */
function renderUsers(items) {
  const online = [];
  const offline = [];

  dom.usersList.innerHTML = ""; 

  items.forEach((item) => {
    const userId = String(item.profile.id);
    onlineUsers.has(userId) ? online.push(item) : offline.push(item);
  });

  dom.usersList.appendChild(createSectionTitle("Online"));
  online.forEach(renderUserItem);

  dom.usersList.appendChild(createSectionTitle("Offline"));
  offline.forEach(renderUserItem);

  updateOnlineCountText();
}

function createSectionTitle(text) {
  const div = document.createElement("div");
  div.className = "section-divider";
  div.textContent = text;
  return div;
}

function renderUserItem(item) {
  const el = Conversation(item);
  const userId = String(item.profile.id);
  const isOnline = onlineUsers.has(userId);

  el.classList.toggle("is-online-user", isOnline);
  el.classList.toggle("is-offline-user", !isOnline);

  el.addEventListener("click", () => {
    openConversation(item);
    document.getElementById("usersPanel")?.classList.remove("mobile-open");
  });

  dom.usersList.appendChild(el);
}

export const reRender = (type, userId) => {
  console.log("Start UI-only rerender:", type, userId);
  
  const targetId = String(userId);
  const isOnline = type === "connect" || type === "register";

  if (isOnline) {
    onlineUsers.add(targetId);
  } else if (type === "disconnect") {
    onlineUsers.delete(targetId);
  }

  if (!dom.usersList) return;

  let targetUserItem = null;
  const conversationItems = dom.usersList.querySelectorAll(".conversation-item, [data-user-id]"); 
  
  for (let item of conversationItems) {
    if (item.getAttribute("data-user-id") === targetId || item.dataset.userId === targetId) {
      targetUserItem = item;
      break;
    }
  }

  if (targetUserItem) {
    const dividers = Array.from(dom.usersList.querySelectorAll(".section-divider"));
    const onlineHeader = dividers.find(d => d.textContent === "Online");
    const offlineHeader = dividers.find(d => d.textContent === "Offline");

    if (isOnline && onlineHeader) {
      onlineHeader.insertAdjacentElement("afterend", targetUserItem);
      targetUserItem.classList.add("is-online-user");
      targetUserItem.classList.remove("is-offline-user");
    } else if (!isOnline && offlineHeader) {
      offlineHeader.insertAdjacentElement("afterend", targetUserItem);
      targetUserItem.classList.add("is-offline-user");
      targetUserItem.classList.remove("is-online-user");
    }
  }

  updateOnlineCountText();

  if (state.currentReceiverId && String(state.currentReceiverId) === targetId) {
    const chatStatusDot = document.getElementById("chatStatusDot");

    if (dom.chatHeaderStatus) {
      dom.chatHeaderStatus.textContent = isOnline ? "● Online" : "● Offline";
      dom.chatHeaderStatus.className = `chat-header-status ${isOnline ? "online" : "offline"}`;
    }
    
    if (chatStatusDot) {
      chatStatusDot.className = `online-dot ${isOnline ? "online" : "offline"}`;
    }
  }
};

function updateOnlineCountText() {
  const countEl = document.getElementById("onlineCount");
  if (countEl) {
    countEl.textContent = `● ${onlineUsers.size} online`;
  }
}

/* =========================
   CORE CHAT VIEW LOGIC
========================= */
async function openConversation(item) {
  const { profile: user, conversation: chat } = item;
//clear the msg
  if (dom.messageInput) {
    dom.messageInput.value = "";
  }

  state.currentReceiverId = user.id;
  state.currentConversationId = chat?.conversationId ?? null;
  dom.chatHeaderName.textContent = user.nickname;
  
  // Reset pagination flags for the newly opened chat
  state.offset = 0;
  state.hasMoreMessages = true;
  state.isLoadingOlder = false;

  const chatAvatarInitials = document.getElementById("chatAvatarInitials");
  if (chatAvatarInitials) {
    chatAvatarInitials.textContent = user.nickname.slice(0, 2);
  }

  const chatStatusDot = document.getElementById("chatStatusDot");
  if (chat?.lastSeen) {
    dom.chatHeaderStatus.textContent = "● Offline";
    dom.chatHeaderStatus.className = "chat-header-status offline";
    if (chatStatusDot) chatStatusDot.className = "online-dot offline";
  } else {
    dom.chatHeaderStatus.textContent = "● Online";
    dom.chatHeaderStatus.className = "chat-header-status online";
    if (chatStatusDot) chatStatusDot.className = "online-dot online";
  }

  dom.chatView.style.display = "flex";
  dom.chatEmpty.style.display = "none";
  dom.chatMessages.innerHTML = "";

  if (!state.currentConversationId) {
    renderEmptyConversation(user);
    // Even if empty, trigger our rendering method to evaluate showTyping configuration
    evaluateTypingIndicatorState();
    return;
  }

  await loadInitialMessages(user.id);
}

async function loadInitialMessages(receiverId) {
  try {
    const res = await getConversationById(state.currentConversationId, {
      limit: state.limit,
      offset: state.offset
    });
    
    const messages = res.data?.messages || [];
    renderMessages(messages, receiverId);
    
    state.offset += state.limit;
  } catch (err) {
    console.error("Error loading chat history:", err);
  }
}

function renderMessages(messages, receiverId) {
  dom.chatMessages.innerHTML = "";
  [...messages].reverse().forEach((m) => {
    appendMessage(m, m.sender_id !== receiverId);
  });
  
  // Apply our override/typing evaluation after historical list renders
  evaluateTypingIndicatorState();
}

function appendMessage(m, mine = false, prepend = false) {
  const group = document.createElement("div");
  group.className = `message-group ${mine ? "mine" : "theirs"}`;
  group.innerHTML = `
    <div class="message-sender">${mine ? "you" : "them"}</div>
    <div class="message-row">
      <div class="message-bubble">${escapeHTML(m.text)}</div>
      <div class="message-meta">${formatTime(m.created_at)}</div>
    </div>`;

  if (prepend) {
    dom.chatMessages.insertBefore(group, dom.chatMessages.firstChild);
  } else {
    const typingIndicator = document.getElementById("typingIndicator");
    if (typingIndicator) {
      dom.chatMessages.insertBefore(group, typingIndicator);
    } else {
      dom.chatMessages.appendChild(group);
    }
    scrollChatToBottom();
  }
}

function scrollChatToBottom() {
  if (dom.chatMessages) {
    dom.chatMessages.scrollTop = dom.chatMessages.scrollHeight;
  }
}

/* =========================
   SCROLL UP PAGINATION
========================= */
function setupScrollPagination() {
  let throttleTimer = false;

  if (!dom.chatMessages) return;

  dom.chatMessages.addEventListener("scroll", () => {
    if (throttleTimer) return;
    
    throttleTimer = true;
    setTimeout(() => { throttleTimer = false; }, 250);

    const pos = dom.chatMessages.scrollTop;
    const maxScrollUp = dom.chatMessages.scrollHeight - dom.chatMessages.clientHeight;

    const isNearTop = pos <= 15 && pos >= 0;
    const isNearTopReversed = Math.abs(pos) >= (maxScrollUp - 15) && pos < 0;

    if (isNearTop || isNearTopReversed) {
      console.log(`🎯 Top reached! Fetching -> Limit: ${state.limit}, Offset: ${state.offset}`);
      loadMoreMessages();
    }
  });
}

async function loadMoreMessages() {
  if (!state.currentConversationId || state.isLoadingOlder || !state.hasMoreMessages) return;

  state.isLoadingOlder = true;
  console.log(`📡 Fetching older data context -> Limit: ${state.limit}, Offset: ${state.offset}`);

  try {
    const previousScrollHeight = dom.chatMessages.scrollHeight;

    const res = await getConversationById(state.currentConversationId,{offset:state.offset + state.limit,limit:state.limit});
    state.offset = state.offset + state.limit;
    const olderMessages = res.data?.messages || [];

    if (olderMessages.length === 0) {
      state.hasMoreMessages = false;
      console.log("🏁 No more historical messages left on the server.");
    } else {
      olderMessages.forEach((m) => {
        appendMessage(m, m.sender_id !== state.currentReceiverId, true);
      });

      state.offset += state.limit;
      dom.chatMessages.scrollTop = dom.chatMessages.scrollHeight - previousScrollHeight;
    }
  } catch (err) {
    console.error("Failed loading historical message pagination blocks:", err);
  } finally {
    state.isLoadingOlder = false;
  }
}

/* =========================
   MESSAGE DELIVERY & SEND
========================= */
/* =========================
   MESSAGE DELIVERY & SEND
========================= */
function setupSendMessage() {
  dom.sendBtn.addEventListener("click", sendMessage);
  
  dom.messageInput.addEventListener("keydown", (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  });

  // 🟢 Listens to input value shifts to capture backspaces, edits, and characters
  dom.messageInput.addEventListener("input", () => {
    handleLocalTypingActivity();
  });
}

async function sendMessage() {
  const text = dom.messageInput.value.trim();
  if (!text || !state.currentReceiverId) return;

  // Stop your typing state immediately
  stopLocalTypingNotification();

  dom.messageInput.value = "";
  // 🟢 Removed the tempMessage and appendMessage from here completely!

  try {
    const res = await createMessage({ 
      receiverId: state.currentReceiverId, 
      text, 
      conversationId: state.currentConversationId 
    });
    if (!state.currentConversationId && res?.data?.conversationId) {
      state.currentConversationId = res.data.conversationId;
    }
  } catch (err) { 
    console.error("send failed", err); 
  }
}




export const reRenderMessages = (data) => {
  // 1. Add console logs to inspect exactly what keys and types are coming from the server
  console.log("=== Realtime Message Received ===");
  console.log("Full WS payload data:", data);
  console.log("sender_id from server:", data.sender_id, "Type:", typeof data.sender_id);
  console.log("Your profile ID:", window.profile?.id, "Type:", typeof window.profile?.id);

  const { conversation_id, text, created_at } = data;
  if (!dom.chatMessages) return;

  // 2. Safely extract sender ID handling both snake_case or camelCase properties just in case
  const incomingSenderId = data.sender_id 
  // 3. Force both to strings and compare cleanly
  const isMine = String(incomingSenderId) === String(window.profile?.id);
  console.log("Calculated isMine evaluation result:",incomingSenderId, isMine,window.profile.id);
  console.log("=================================");

  // 4. Only append if this belongs to your currently open chat screen
  if (String(state.currentConversationId) === String(conversation_id)) {
    
    if (String(incomingSenderId) === String(state.currentReceiverId)) {
      setPartnerTyping(false);
    }

    const incomingMsg = {
      sender_id: incomingSenderId,
      text: text,
      created_at: created_at || new Date().toISOString()
    };
    
    appendMessage(incomingMsg, isMine);
  }
};

/* =========================
   TYPING INDICATOR CONTROL
========================= */
function evaluateTypingIndicatorState() {
  const currentNickname = dom.chatHeaderName?.textContent || "them";
  
  // If showTyping is manually turned on, render it immediately
  if (state.showTyping || state.isTyping) {
    console.log(`💬 [Typing State]: Active (${state.showTyping ? 'Forced Overwrite' : 'Network Event'})`);
    renderTypingIndicator(currentNickname);
  } else {
    console.log(`🚫 [Typing State]: Stopped / Hidden.`);
    removeTypingIndicator();
  }
}

export function setPartnerTyping(isTyping, nickname = "them") {
  state.isTyping = isTyping;
  evaluateTypingIndicatorState();
}

// Global hook so you can quickly switch state variables right out of your web console!
export function toggleOverrideTyping(forceVisible) {
  state.showTyping = forceVisible;
  evaluateTypingIndicatorState();
}

function renderTypingIndicator(nickname) {
  if (!dom.chatMessages) return;
  if (document.getElementById("typingIndicator")) return;

  const group = document.createElement("div");
  group.className = "message-group theirs";
  group.id = "typingIndicator";
  group.innerHTML = `
    <div class="message-sender">${escapeHTML(nickname)}</div>
    <div class="message-row">
      <div class="typing-indicator-container">
        <div class="typing-dots">
          <span></span><span></span><span></span>
        </div>
      </div>
    </div>`;

  dom.chatMessages.appendChild(group);
  scrollChatToBottom();
}

function removeTypingIndicator() {
  const indicator = document.getElementById("typingIndicator");
  if (indicator) {
    indicator.remove();
  }
}

export const handleIncomingTypingEvent = (data) => {
  // 🟢 Extract the correct camelCase fields coming from your server logs
  const { conversationId, userId, is_typing } = data;
  
  console.log("the data coming", data);
  
  // 🟢 Compare against the correct variable keys
  if (
    String(state.currentConversationId) === String(conversationId) &&
    String(state.currentReceiverId) === String(userId)
  ) {
    const activeName = dom.chatHeaderName ? dom.chatHeaderName.textContent : "them";
    setPartnerTyping(is_typing, activeName);
  }
};

/* =========================
   HELPERS & RENDER LAYOUTS
========================= */
function escapeHTML(str) {
  return str.replace(/[&<>"']/g, (m) => HTML_CHARS[m] || m);
}

function renderEmptyConversation(user) {
  dom.chatMessages.innerHTML = `
    <div class="chat-empty-state" style="margin: auto; text-align: center;">
      <div class="day-separator"><span>New chat with ${user.nickname}</span></div>
      <p style="font-family: var(--mono); font-size: 0.72rem; color: var(--text-dim);">No messages yet — say hi 👋</p>
    </div>`;
  scrollChatToBottom();
}

function showEmptyState() {
  dom.chatView.style.display = "none";
  dom.chatEmpty.style.display = "flex";
}