import { templates, getRouteFromHash, updateHash, navigate, handleRouteChange } from './router.js';

const app = document.getElementById("app");

// Store loaded scripts to avoid duplicates
const loadedScripts = new Map();




// ========================
// GLOBAL FUNCTIONS
// ========================
window.navigate = navigate;
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
// Listen for hash changes (back/forward buttons)
window.addEventListener('hashchange', handleRouteChange);

// Initialize the app based on current hash
const initialRoute = getRouteFromHash();
navigate(initialRoute, { updateHash: false });
