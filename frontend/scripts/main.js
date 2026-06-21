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
