import { formatTime } from './helpers.js';

import { getConversations, getConversationById } from "../api/conversations.js";
import { createMessage } from "../api/messages.js";
import { getOnlineUsers } from "./main.js";
import { Conversation } from "../components/Conversation.js";


let currentConversationId = null;
let currentReceiverId = null;

const usersList = document.getElementById("usersList");
const chatMessages = document.getElementById("chatMessages");
const chatHeaderName = document.getElementById("chatHeaderName");
const chatHeaderStatus = document.getElementById("chatHeaderStatus");
const messageInput = document.getElementById("messageInput");
const sendBtn = document.getElementById("sendBtn");
const chatEmpty = document.getElementById("chatEmpty");
const chatView = document.getElementById("chatView");

export async function setup() {
  try {
		const resp = await getConversations();
		const data = resp.data || [];

		renderUsers(data);

		if (data.length > 0) {
			const first = data[0];
			if (first.conversation.conversationId) {
				openConversation(first);
			}
		} else {
			showEmptyState();
		}

		setupSendMessage();
	} catch (err) {
		console.error("chat init error:", err);
	}
};

////////////////////////////////////

function renderUsers(items) {
	const onlineUsers = getOnlineUsers();
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
