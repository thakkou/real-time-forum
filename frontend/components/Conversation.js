import { formatTime } from '../scripts/helpers.js';

export const Conversation = (item) => {
    const u = item.profile;
    const c = item.conversation;

    const div = document.createElement("div");

    div.className = "user-item";
    if (c.unreadCount > 0) div.classList.add("unread");

    div.dataset.userId = u.id;
    div.dataset.conversationId = c.conversationId;

    div.innerHTML = `
        <div class="user-avatar">
            <div class="avatar-circle">${u.nickname.slice(0,2)}</div>
            <span class="online-dot ${c.lastSeen ? "offline" : "online"}"></span>
        </div>

        <div class="user-info">
            <div class="user-name-row">
                <span class="user-name">${u.nickname}</span>
                <span class="user-time">${formatTime(c.date)}</span>
            </div>

            <div class="user-preview">
                ${c.lastMessage || "No messages yet"}
            </div>
        </div>
    `;
    return div;
}

// <div class="user-item ${c.unreadCount > 0 ? 'unread' : ''}" data-user-id="${u.id}" data-conversation-id="${c.conversationId}" data-username="${username}">
//     <div class="user-avatar">
//         <div class="avatar-circle">${u.nickname.slice(0,2)}</div>
//         <span class="online-dot ${c.lastSeen ? "offline" : "online"}"></span>
//     </div>
//     <div class="user-info">
//         <div class="user-name-row">
//             <span class="user-name">${u.nickname}</span>
//             <span class="user-time">${formatTime(c.date)}</span>
//         </div>

//         <div class="user-preview">
//             ${c.lastMessage || "No messages yet."}
//         </div>
//     </div>
// </div>

// user avatar:
// 1. NV: abbreviation of the name !
// 2. isonline ?!

// format time for timeago

