export const Toast = (message, type = "info", title = "") => {
    const toast = document.createElement("div");
    toast.className = `toast ${type}`;

    if (title) {
        toast.innerHTML = `
            <div class="toast-header">
                <div class="toast-title">${title}</div>
            </div>
            <div class="toast-message">${message}</div>
        `;
    } else {
        toast.innerHTML = `
            <div class="toast-message">${message}</div>
        `;
    }

    return toast;
};