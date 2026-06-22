import { router } from './router.js';

const app = document.getElementById("app");

// Store loaded scripts to avoid duplicates
const loadedScripts = new Map();


// ========================
// GLOBAL FUNCTIONS
// ========================
window.navigate = router.navigate.bind(router); // navigate
window.reactToPost = async function(postId, endpoint) {
  const url = `/api/posts/${postId}/${endpoint}`;
  try {
    const response = await fetch(url, { method: 'POST' });
    if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);
    window.location.reload();
  } catch (error) {
    console.error('Error sending reaction:', error);
  }
};
window.reactToComment = async function(commentId, endpoint) {
  const url = `/api/comments/${commentId}/${endpoint}`;
  try {
    const response = await fetch(url, { method: 'POST' });
    if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);
    window.location.reload();
  } catch (error) {
    console.error('Error sending reaction:', error);
  }
};

// ========================
// INITIALIZE APP
// ========================
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

showToast('Message received!');
showToast('User connected', 'success');
showToast('Connection lost', 'error');
showToast('New notification', 'warning', 5000);

// socket.onmessage = (event) => {
//     showToast('New message received');
// };