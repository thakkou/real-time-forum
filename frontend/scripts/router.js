export const templates = {
  feed: () => import('../pages/feed.js'),
  login: () => import('../pages/login.js'),
  register: () => import('../pages/register.js'),
  error: () => import('../pages/error.js'),
  chat: () => import('../pages/chat.js'),
};

// ========================
// PAGE-SPECIFIC SCRIPTS
// ========================
const pageScripts = {
  feed: () => {
    function setupHomePage() {
      document.querySelectorAll('.filter-btn').forEach((button) => {
        button.addEventListener('click', function() {
          this.classList.toggle('active');
          const hiddenInputs = {
            'my-creat-postes': document.getElementById('input-my-creat-postes'),
            'my-liked-post': document.getElementById('input-my-liked-post'),
          };
          Object.keys(hiddenInputs).forEach((name) => {
            if (hiddenInputs[name]) hiddenInputs[name].value = '';
          });
          const activeButtons = Array.from(document.querySelectorAll('.filter-btn.active'));
          activeButtons.forEach((btn) => {
            if (hiddenInputs[btn.name]) {
              hiddenInputs[btn.name].value = 'true';
            }
          });
        });
      });

      function reactToPost(postId, endpoint) {
        const url = `/api/posts/${postId}/${endpoint}`;
        return fetch(url, { method: 'POST' })
          .then(response => {
            if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);
            window.location.reload();
          })
          .catch(error => console.error('Error sending reaction:', error));
      }

      function handlePostReactionsClick(event) {
        const button = event.target;
        const postContainer = button.closest('.post');
        const postId = postContainer?.getAttribute('data-post-id');
        if (!postId) return;
        
        let endpoint;
        if (button.classList.contains('like-btn')) endpoint = 'like';
        else if (button.classList.contains('dislike-btn')) endpoint = 'dislike';
        else return;
        
        reactToPost(postId, endpoint);
      }

      function reactToComment(commentId, endpoint) {
        const url = `/api/comments/${commentId}/${endpoint}`;
        return fetch(url, { method: 'POST' })
          .then(response => {
            if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);
            window.location.reload();
          })
          .catch(error => console.error('Error sending reaction:', error));
      }

      function handleCommentReactionsClick(event) {
        const button = event.target;
        const commentContainer = button.closest('.comment');
        const commentId = commentContainer?.getAttribute('data-comment-id');
        if (!commentId) return;
        
        let endpoint;
        if (button.classList.contains('comment-like-btn')) endpoint = 'like';
        else if (button.classList.contains('comment-dislike-btn')) endpoint = 'dislike';
        else return;
        
        reactToComment(commentId, endpoint);
      }

      document.querySelectorAll('.like-btn, .dislike-btn').forEach(button => {
        button.removeEventListener('click', handlePostReactionsClick);
        button.addEventListener('click', handlePostReactionsClick);
      });
      
      document.querySelectorAll('.comment-like-btn, .comment-dislike-btn').forEach(button => {
        button.removeEventListener('click', handleCommentReactionsClick);
        button.addEventListener('click', handleCommentReactionsClick);
      });
    }

    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', setupHomePage);
    } else {
      setupHomePage();
    }
  },

  login: () => {
    function setupLoginPage() {
      console.log("Login page loaded");
    }
    
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', setupLoginPage);
    } else {
      setupLoginPage();
    }
  },

  register: () => {
    function setupRegisterPage() {
      const form = document.getElementById("register-form");
      if (!form) return;
      
      form.addEventListener("submit", async (e) => {
        e.preventDefault();
        
        const btn = document.getElementById("register-btn");
        const errorBox = document.getElementById("form-error");
        const formData = e.target;
        
        const nickname = formData.nickname.value;
        const first_name = formData.first_name.value;
        const last_name = formData.last_name.value;
        const age = parseInt(formData.age.value, 10);
        const gender = formData.querySelector('input[name="gender"]:checked');
        const email = formData.email.value;
        const password = formData.password.value;
        const confirm_password = formData.confirm_password.value;
        
        errorBox.style.display = "none";
        errorBox.textContent = "";
        
        if (!nickname || nickname.trim().length < 2) {
          errorBox.textContent = "Nickname must be at least 2 characters";
          errorBox.style.display = "block";
          return;
        }
        if (!first_name || first_name.trim().length < 1) {
          errorBox.textContent = "First name is required";
          errorBox.style.display = "block";
          return;
        }
        if (!last_name || last_name.trim().length < 1) {
          errorBox.textContent = "Last name is required";
          errorBox.style.display = "block";
          return;
        }
        if (!age || isNaN(age) || age < 13 || age > 120) {
          errorBox.textContent = "Age must be a number between 13 and 120";
          errorBox.style.display = "block";
          return;
        }
        if (!gender) {
          errorBox.textContent = "Please select a gender";
          errorBox.style.display = "block";
          return;
        }
        if (!email || !email.includes("@") || !email.includes(".")) {
          errorBox.textContent = "Please enter a valid email address";
          errorBox.style.display = "block";
          return;
        }
        if (!password || password.length < 6) {
          errorBox.textContent = "Password must be at least 6 characters";
          errorBox.style.display = "block";
          return;
        }
        if (password !== confirm_password) {
          errorBox.textContent = "Passwords do not match";
          errorBox.style.display = "block";
          return;
        }
        
        console.log(JSON.stringify({ 
          nickname, first_name, last_name, age, 
          gender: gender.value, email, password, confirm_password 
        }));
        
        btn.disabled = true;
        btn.textContent = "Registering...";
      });
    }
    
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', setupRegisterPage);
    } else {
      setupRegisterPage();
    }
  }
};

function loadPageScript(pageName) {
  if (window.currentPageScript && typeof window.currentPageScript.cleanup === 'function') {
    window.currentPageScript.cleanup();
  }
  
  if (pageScripts[pageName]) {
    pageScripts[pageName]();
    window.currentPage = pageName;
  }
}


// ========================
// ROUTER WITH HASH SUPPORT
// ========================
export function getRouteFromHash() {
  // Get the hash without the # symbol
  const hash = window.location.hash.slice(1);
  // Check if the hash corresponds to a valid route
  if (hash && templates[hash]) {
    return hash;
  }
  // Default to home
  return 'feed';
}

export function updateHash(route) {
  window.location.hash = route; // why hash ?!
}

export async function navigate(page, options = { updateHash: true }) {
  // Update URL hash if needed
  if (options.updateHash) {
    updateHash(page);
  }
  
  // Render the template
  const mod = templates[page]();
  app.innerHTML = await (await mod).render();
  
  // Load the page-specific script
  loadPageScript(page);
}

// Handle browser back/forward buttons
export function handleRouteChange() {
  const route = getRouteFromHash();
  navigate(route, { updateHash: false });
}