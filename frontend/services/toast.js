import { Toast } from '../components/Toast.js';


export function showToast(
    message,
    type = "info",
    title = "",
    duration = 5000,
    position = "bottom" // "top" or "bottom"
) {
    const container = document.getElementById("toast-container");

    // Set container position
    container.classList.toggle("top", position === "top");
    container.classList.toggle("bottom", position !== "top");

    // Max 2 toasts
    while (container.children.length >= 2) {
        container.firstElementChild.remove();
    }

    const toast = Toast(message, type, title);
    container.appendChild(toast);

    requestAnimationFrame(() => {
        toast.classList.add("show");
    });

    setTimeout(() => {
        toast.classList.remove("show");

        setTimeout(() => {
            toast.remove();
        }, 300);
    }, duration);
}