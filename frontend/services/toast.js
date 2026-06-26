import { Toast } from '../components/Toast.js';

export function showToast(message, type = 'success', duration = 5000) {
    // available types: success, error, warning...
    const container = document.getElementById('toast-container');

    const toast = Toast(message, type);

    // var doc = new DOMParser().parseFromString(toast, "text/xml");
    container.appendChild(toast);

    requestAnimationFrame(() => {
        toast.classList.add('show');
    });

    setTimeout(() => {
        toast.classList.remove('show');

        setTimeout(() => {
            toast.remove();
        }, 300);
    }, duration);
}