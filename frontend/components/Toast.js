import {sanitize} from "../scripts/helpers.js"
export const Toast = (message, type = "info", title = "") => {
    const MAX_LENGTH = 25;
message=sanitize(message)
    const truncatedMessage =
        message.length > MAX_LENGTH
            ? message.slice(0, MAX_LENGTH) + "..."
            : message;

    const toast = document.createElement("div");
    toast.className = `toast ${type}`;

    if (title) {
        toast.innerHTML = `
            <div class="toast-header">
                <div class="toast-title">${title}</div>
            </div>
            <div class="toast-message">${truncatedMessage}</div>
        `;
    } else {
        toast.innerHTML = `
            <div class="toast-message">${truncatedMessage}</div>
        `;
    }

    return toast;
};