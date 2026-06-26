export const Message = (message) => (`
    <div class="message-row">
        <div class="message-bubble">
            ${text}
        </div>
        <div class="message-meta">
            <span>${time}</span>
        </div>
    </div>
`);

// time format: Jun 20 · 18:41