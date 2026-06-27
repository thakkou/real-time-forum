import { formatTime } from './helpers.js';

import { getConversations, getConversationById } from "../api/conversations.js";
import { createMessage } from "../api/messages.js";
import { onlineUsers } from "../services/router.js";
import { Conversation } from "../components/Conversation.js";


let currentConversationId = null;
let currentReceiverId = null;

let usersList;
let chatMessages;
let chatHeaderName;
let chatHeaderStatus;
let messageInput;
let sendBtn;
let chatEmpty;
let chatView;

export async function setup() {
  usersList = document.getElementById("usersList");
  chatMessages = document.getElementById("chatMessages");
  chatHeaderName = document.getElementById("chatHeaderName");
  chatHeaderStatus = document.getElementById("chatHeaderStatus");
  messageInput = document.getElementById("messageInput");
  sendBtn = document.getElementById("sendBtn");
  chatEmpty = document.getElementById("chatEmpty");
  chatView = document.getElementById("chatView");

  if (!usersList) {
    console.error("Chat DOM not found — page not rendered yet");
    return;
  }

  const resp = await getConversations();
  const data = resp.data || [];

  renderUsers(data);

  if (data.length > 0) {
    openConversation(data[0]);
  } else {
    showEmptyState();
  }

  setupSendMessage();
}

export const reRender = (type, userId) => {
  console.log("Start UI-only rerender:", type, userId);
  
  const targetId = String(userId);
  const isOnline = type === "connect" || type === "register";

  // 1. Sync the router's online tracking set
  if (isOnline) {
    onlineUsers.add(targetId);
  } else if (type === "disconnect") {
    onlineUsers.delete(targetId);
  }

  // 2. Find the user item in the DOM list
  // We look for all elements inside usersList to find the one matching our user ID
  const usersList = document.getElementById("usersList");
  if (!usersList) return;

  // Assumes your Conversation component attaches a data attribute or you can find it. 
  // If your Conversation component doesn't have an ID, we'll look through them.
  let targetUserItem = null;
  const conversationItems = usersList.querySelectorAll(".conversation-item, [data-user-id]"); 
  
  // Find the existing DOM element for this user
  for (let item of conversationItems) {
    // Adjust this line if your Conversation component uses a different way to store the ID
    if (item.getAttribute("data-user-id") === targetId || item.dataset.userId === targetId) {
      targetUserItem = item;
      break;
    }
  }

  // 3. Move the element to the correct section if found
  if (targetUserItem) {
    // Find our Section Dividers inside the list
    const dividers = Array.from(usersList.querySelectorAll(".section-divider"));
    const onlineHeader = dividers.find(d => d.textContent === "Online");
    const offlineHeader = dividers.find(d => d.textContent === "Offline");

    if (isOnline && onlineHeader) {
      // Insert right after the Online header
      onlineHeader.insertAdjacentElement("afterend", targetUserItem);
    } else if (!isOnline && offlineHeader) {
      // Insert right after the Offline header
      offlineHeader.insertAdjacentElement("afterend", targetUserItem);
    }
  }

  // 4. Update Header status immediately if this is the currently open conversation
  if (currentReceiverId && String(currentReceiverId) === targetId) {
    const chatHeaderStatus = document.getElementById("chatHeaderStatus");
    const chatStatusDot = document.getElementById("chatStatusDot");

    if (chatHeaderStatus) {
      chatHeaderStatus.textContent = isOnline ? "● Online" : "● Offline";
      chatHeaderStatus.className = `chat-header-status ${isOnline ? "online" : "offline"}`;
    }
    
    if (chatStatusDot) {
      chatStatusDot.className = `online-dot ${isOnline ? "online" : "offline"}`;
    }
  }
};
////////////////////////////////////

function renderUsers(items) {

	const online = [], offline = [];
	
	usersList.innerHTML = "";
	items.forEach((item) => {
		const userId = String(item.profile.id); // force string
		onlineUsers.has(userId) ? online.push(item) : offline.push(item);
	});

	usersList.appendChild(sectionTitle("Online"));
	online.forEach(renderUserItem);

	usersList.appendChild(sectionTitle("Offline"));
	offline.forEach(renderUserItem);
}

function sectionTitle(text) {
	const div = document.createElement("div");
	div.className = "section-divider";
	div.textContent = text;
	return div;
}

function renderUserItem(item) {
	const conv = Conversation(item);
	conv.addEventListener("click", () => openConversation(item));
	usersList.appendChild(conv);
}

///////////////////////////////////////////

function renderEmptyConversation(user) {
  chatMessages.innerHTML = `
    <div class="chat-empty-state" style="margin: auto; text-align: center;">
      <div class="day-separator"><span>New chat with ${user.nickname}</span></div>
      <p style="font-family: var(--mono); font-size: 0.72rem; color: var(--text-dim);">
        No messages yet — say hi 👋
      </p>
    </div>
  `;
  chatMessages.scrollTop = chatMessages.scrollHeight;
}

async function openConversation(item) {
	const u = item.profile;
	const c = item.conversation;

	currentReceiverId = u.id;
	currentConversationId = c?.conversationId ?? null;

	// UI header updates
	chatHeaderName.textContent = u.nickname;
	
	// Update Header Avatar Initials
	const chatAvatarInitials = document.getElementById("chatAvatarInitials");
	if (chatAvatarInitials) {
		chatAvatarInitials.textContent = u.nickname.slice(0, 2);
	}

	// Update Status and Dots matching CSS definitions
	const chatStatusDot = document.getElementById("chatStatusDot");
	if (c?.lastSeen) {
		chatHeaderStatus.textContent = "● Offline";
		chatHeaderStatus.className = "chat-header-status offline";
		if (chatStatusDot) chatStatusDot.className = "online-dot offline";
	} else {
		chatHeaderStatus.textContent = "● Online";
		chatHeaderStatus.className = "chat-header-status online";
		if (chatStatusDot) chatStatusDot.className = "online-dot online";
	}

	chatView.style.display = "flex";
	chatEmpty.style.display = "none";

	chatMessages.innerHTML = "";

	// NO conversation yet → show empty chat state
	if (!currentConversationId) {
		renderEmptyConversation(u);
		return;
	}

	const res = await getConversationById(currentConversationId);
	const messages = res.data?.messages || [];

	renderMessages(messages, u.id);
}

function setupSendMessage() {
  sendBtn.addEventListener("click", sendMessage);

  messageInput.addEventListener("keydown", (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  });
}

function renderMessages(messages, receiverId) {
  chatMessages.innerHTML = "";

  messages.reverse().forEach((m) => {
    const isMine = m.sender_id !== receiverId;

    const group = document.createElement("div");
    group.className = `message-group ${isMine ? "mine" : "theirs"}`;

    group.innerHTML = `
      <div class="message-sender">${isMine ? "you" : "them"}</div>
      <div class="message-row">
        <div class="message-bubble">${escapeHTML(m.text)}</div>
        <div class="message-meta">${formatTime(m.created_at)}</div>
      </div>
    `;

    chatMessages.appendChild(group);
  });

  chatMessages.scrollTop = chatMessages.scrollHeight;
}

function appendMessage(m, mine = false) {
  const group = document.createElement("div");
  group.className = `message-group ${mine ? "mine" : "theirs"}`;

  group.innerHTML = `
    <div class="message-sender">${mine ? "you" : "them"}</div>
    <div class="message-row">
      <div class="message-bubble">${escapeHTML(m.text)}</div>
      <div class="message-meta">${formatTime(m.created_at)}</div>
    </div>
  `;

  chatMessages.appendChild(group);
  chatMessages.scrollTop = chatMessages.scrollHeight;
}

async function sendMessage() {
  const text = messageInput.value.trim();
  if (!text || !currentReceiverId) return;

  messageInput.value = "";

  // ✅ optimistic message (correct shape)
  const tempMessage = {
    sender_id: "me",
    text,
    created_at: new Date().toISOString(),
  };

  appendMessage(tempMessage, true);

  try {
    const res = await createMessage({
      receiverId: currentReceiverId,
      text,
      conversationId: currentConversationId, // may be null
    });

    // ✅ IMPORTANT: backend should return conversationId (if new)
    if (!currentConversationId && res?.data?.conversationId) {
      currentConversationId = res.data.conversationId;
    }

  } catch (err) {
    console.error("send failed", err);
  }
}

function escapeHTML(str) {
  return str.replace(/[&<>"']/g, (m) => ({
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    '"': "&quot;",
    "'": "&#039;",
  }[m]));
}

function showEmptyState() {
  chatView.style.display = "none";
  chatEmpty.style.display = "flex";
}
