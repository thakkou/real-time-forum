import { formatTime, formatDate, sanitize } from '../scripts/helpers.js';

export const Conversation = (item) => {
    console.log(item)
    const u = item.profile;
    const c = item.conversation;

    const div = document.createElement("div");

    div.className = "user-item";
    div.dataset.userId = u.id;
    div.dataset.conversationId = c.conversationId;

    div.innerHTML = `
        <div class="user-avatar">
            <div class="avatar-circle">${sanitize(u.nickname.slice(0,2))}</div>
            <span class="online-dot ${c.lastSeen ? "offline" : "online"}"></span>
        </div>

        <div class="user-info">
            <div class="user-name-row">
                <span class="user-name">${sanitize(u.nickname)}</span>

                <span class="user-time">
                    <span>${formatTime(c.date)}</span>
                    <span class="user-date">${c.conversationId ? formatDate(c.date): "*"}</span>
                </span>
            </div>

            <div class="user-preview">
                ${sanitize(c.lastMessage) || "No messages yet"}
            </div>
        </div>
    `;

    return div;
};