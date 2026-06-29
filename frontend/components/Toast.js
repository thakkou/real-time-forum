export const Toast = (message, type = 'info', title = '') => {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;

    if (title) {
        // Create a structured container to house both title and message text
        toast.innerHTML = `
            <div class="toast-title">${title}</div>
            <div class="toast-message">${message}</div>
        `;
    } else {
        // Fallback layout if no title is provided
        toast.innerHTML = `<div class="toast-message">${message}</div>`;
    }

    return toast;
};