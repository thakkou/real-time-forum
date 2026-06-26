export const Conversation = (conversation) => (`
    <div class="user-item" data-user-id="${userId}" data-username="${username}">
        <div class="user-avatar">
            <div class="avatar-circle">NV</div>
            <span class="online-dot ${isOnline ? 'online' : 'offline'}"></span>
        </div>
        <div class="user-info">
            <div class="user-name-row">
                <span class="user-name">${username}</span>
                <span class="user-time">${timeAgo}</span>
            </div>
            <div class="user-preview">${last_message}</div>
        </div>
    </div>
`);

// NV: abbreviation of the name !
