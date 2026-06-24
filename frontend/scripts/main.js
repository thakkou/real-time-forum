import { router } from './router.js';

const app = document.getElementById("app");

// Store loaded scripts to avoid duplicates
const loadedScripts = new Map();


// ========================
// GLOBAL FUNCTIONS
// ========================
//later 
window.env = {
    serverUri: "http://localhost:8080/api",
    wsUri:"ws://localhost:8080/ws"
};

window.navigate = router.navigate.bind(router); // navigate

router.init();


// =============================================================

function showToast(message, type = 'success', duration = 3000) {
    const container = document.getElementById('toast-container');

    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;

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

// showToast('Message received!');
// showToast('User connected', 'success');
// showToast('Connection lost', 'error');
// showToast('New notification', 'warning', 5000);

// socket.onmessage = (event) => {
//     showToast('New message received');
// };