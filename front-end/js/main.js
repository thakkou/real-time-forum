// ========================
// APP STATE & ROUTER
// ========================
const app = document.getElementById("app");

// Store loaded scripts to avoid duplicates
const loadedScripts = new Map();

// Page templates (HTML only, no scripts)
const templates = {
  home: () => `
<header class="navbar">
  <div class="logo">01Forum</div>
  <div class="auth-buttons">
    <button onclick="navigate('login')" class="btn login">Login</button>
    <button onclick="navigate('register')" class="btn register">Register</button>
  </div>
</header>

<div class="container">
  <aside class="sidebar">
    <h3>Filter Posts</h3>
    <form method="GET" action="/">
      <div class="filter-group">
        <h4>By Category</h4>
        <label><input type="checkbox" name="categories" value="General" /> General</label>
        <label><input type="checkbox" name="categories" value="Lifestyle" /> Lifestyle</label>
        <label><input type="checkbox" name="categories" value="Health & Fitness" /> Health & Fitness</label>
        <label><input type="checkbox" name="categories" value="Travel" /> Travel</label>
        <button type="submit" class="filter-btn">Apply Filters</button>
        <button onclick="navigate('home')" class="clear-btn">Clear Filters</button>
      </div>
    </form>
  </aside>

  <main class="content">
    <section class="create-post">
      <h2>Create New Post</h2>
      <form action="/api/posts/create" method="POST" enctype="multipart/form-data">
        <input type="text" name="title" placeholder="Post title" required />
        <textarea name="text" placeholder="Write your post..." required></textarea>
        <div class="categories">
          <label><input type="checkbox" name="categories" value="General" checked /> General</label>
          <label><input type="checkbox" name="categories" value="Lifestyle" /> Lifestyle</label>
          <label><input type="checkbox" name="categories" value="Health & Fitness" /> Health & Fitness</label>
          <label><input type="checkbox" name="categories" value="Travel" /> Travel</label>
        </div>
        <button type="submit" class="btn primary">Publish</button>
      </form>
    </section>

    <section class="posts">
      <article class="post">
        <div class="post-header">
          <h3>Example Post Title</h3>
          <span class="post-meta">Posted by John · 2 hours ago</span>
        </div>
        <div class="post-body">
          <p>This is an example post content...</p>
        </div>
        <div class="post-actions">
          <button class="like-btn">👍 10</button>
          <button class="dislike-btn">👎 2</button>
        </div>
        <div class="comments">
          <h4>Comments</h4>
          <div class="comment">
            <span class="comment-meta">Posted by Alice · 1 hour ago</span>
            <p class="comment-text">Nice post!</p>
          </div>
          <div class="add-comment">
            <form action="/api/comments/create" method="POST">
              <input name="text" type="text" placeholder="Write a comment..." required />
              <button type="submit" class="btn small primary">Comment</button>
            </form>
          </div>
        </div>
      </article>
    </section>
  </main>
</div>
`,

  login: () => `
<div class="card">
  <h1>Login</h1>
  <p class="subtitle">Welcome back</p>
  <form action="/login" method="POST">
    <div class="field">
      <label>Email or Username</label>
      <input type="text" name="email" placeholder="Email or Username" required autocomplete="off" />
    </div>
    <div class="field">
      <label>Password</label>
      <input type="password" name="password" placeholder="Password" required autocomplete="current-password" />
    </div>
    <button type="submit">Login</button>
  </form>
  <div class="link-row">No account? <button onclick="navigate('register')">Register</button></div>
</div>
`,

  register: () => `
<div class="card">
  <h1>Register</h1>
  <p class="subtitle">Create your forum account</p>
  <div id="form-error" style="display:none;" class="form-error"></div>

  <form id="register-form">
    <div class="field">
      <label>Nickname *</label>
      <input type="text" name="nickname" required minlength="2" maxlength="30" autocomplete="off" />
    </div>
    
    <div class="field-row">
      <div class="field">
        <label>First Name *</label>
        <input type="text" name="first_name" required minlength="1" maxlength="50" autocomplete="given-name" />
      </div>
      <div class="field">
        <label>Last Name *</label>
        <input type="text" name="last_name" required minlength="1" maxlength="50" autocomplete="family-name" />
      </div>
    </div>
    
    <div class="field-row">
      <div class="field">
        <label>Age *</label>
        <input type="number" name="age" required min="13" max="120" step="1" />
      </div>
      <div class="field">
        <label>Gender *</label>
        <div class="gender-group">
          <label class="gender-option">
            <input type="radio" name="gender" value="Male" /> Male
          </label>
          <label class="gender-option">
            <input type="radio" name="gender" value="Female" /> Female
          </label>
        </div>
      </div>
    </div>
    
    <div class="field">
      <label>Email *</label>
      <input type="email" name="email" required autocomplete="email" maxlength="100" />
    </div>
    
    <div class="field">
      <label>Password *</label>
      <input type="password" name="password" required minlength="6" maxlength="128" autocomplete="new-password" />
    </div>
    
    <div class="field">
      <label>Confirm Password *</label>
      <input type="password" name="confirm_password" required minlength="6" maxlength="128" autocomplete="new-password" />
    </div>
    
    <button type="submit" id="register-btn">Register</button>
  </form>

  <div class="link-row">
    Already have an account? <button onclick="navigate('login')">Login</button>
  </div>
</div>

<style>
  .field-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }
  
  .gender-group {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    background: var(--surface-2, #222222);
    border: 1px solid var(--border, #2e2e2e);
    border-radius: var(--radius, 2px);
    padding: 0.6rem 0.9rem;
  }
  
  .gender-option {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    font-family: var(--mono, 'Space Mono', monospace);
    font-size: 0.7rem;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--text-muted, #7a7a7a);
    cursor: pointer;
    transition: color 0.15s ease;
    margin: 0;
  }
  
  .gender-option:hover {
    color: var(--text, #f0ede6);
  }
  
  .gender-option input[type="radio"] {
    accent-color: var(--accent, #e8ff47);
    width: 13px;
    height: 13px;
    margin: 0;
    cursor: pointer;
  }
  
  .gender-option:has(input:checked) {
    color: var(--accent, #e8ff47);
  }
  
  input[type="number"] {
    width: 100%;
    background: var(--surface-2, #222222);
    border: 1px solid var(--border, #2e2e2e);
    border-radius: var(--radius, 2px);
    padding: 0.65rem 0.9rem;
    color: var(--text, #f0ede6);
    font-family: var(--sans, 'IBM Plex Sans', sans-serif);
    font-size: 0.9rem;
    transition: border-color 0.15s;
    outline: none;
  }
  
  input[type="number"]:focus {
    border-color: var(--accent, #e8ff47);
  }
  
  input[type="number"]::-webkit-inner-spin-button,
  input[type="number"]::-webkit-outer-spin-button {
    opacity: 0.5;
  }
  
  @media (max-width: 600px) {
    .field-row {
      grid-template-columns: 1fr;
      gap: 0.75rem;
    }
    .gender-group {
      flex-direction: column;
      gap: 0.5rem;
    }
  }
</style>
`
};

// ========================
// PAGE-SPECIFIC SCRIPTS
// ========================
const pageScripts = {
  home: () => {
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

// ========================
// ROUTER WITH HASH SUPPORT
// ========================
function getRouteFromHash() {
  // Get the hash without the # symbol
  const hash = window.location.hash.slice(1);
  // Check if the hash corresponds to a valid route
  if (hash && templates[hash]) {
    return hash;
  }
  // Default to home
  return 'home';
}

function updateHash(route) {
  window.location.hash = route;
}

function loadPageScript(pageName) {
  if (window.currentPageScript && typeof window.currentPageScript.cleanup === 'function') {
    window.currentPageScript.cleanup();
  }
  
  if (pageScripts[pageName]) {
    pageScripts[pageName]();
    window.currentPage = pageName;
  }
}

function navigate(page, options = { updateHash: true }) {
  // Update URL hash if needed
  if (options.updateHash) {
    updateHash(page);
  }
  
  // Render the template
  app.innerHTML = templates[page]();
  
  // Load the page-specific script
  loadPageScript(page);
}

// Handle browser back/forward buttons
function handleRouteChange() {
  const route = getRouteFromHash();
  navigate(route, { updateHash: false });
}

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