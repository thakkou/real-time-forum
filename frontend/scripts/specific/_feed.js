import { CommentResolver, CreatComment } from "../../api/comments.js";
import { logout } from "../../api/auth.js";
import { getPosts, PostResolver,CreatePost } from "../../api/posts.js";
import { Post } from "../../components/Post.js";

/* ======================
   STATE
====================== */
const state = {
  offset: 0,
  loading: false,
  hasMore: true,
};

/* ======================
   UTILS
====================== */
function throttle(fn, delay = 200) {
  let last = 0;
  return (...args) => {
    const now = Date.now();
    if (now - last < delay) return;
    last = now;
    fn(...args);
  };
}

function resetFeed() {
  state.offset = 0;
  state.hasMore = true;
  document.querySelector(".posts").innerHTML = "";
}

/* ======================
   API ACTIONS
====================== */
async function fetchPosts() {
  if (state.loading || !state.hasMore) return;

  state.loading = true;

  const params = new URLSearchParams(window.location.search);

  const categories = params.getAll("categories");
  const isLiked = params.get("my-liked-post") === "true";
  const isCreatedByMe = params.get("my-creat-postes") === "true";

  try {
    const res = await getPosts({
      offset: state.offset,
      limit: 30,
      categories,
      isLiked,
      isCreatedByMe,
    });

    const posts = res.data;
    console.log("posts",posts)
    if (!posts?.length) {
      state.hasMore = false;
      return;
    }

    renderPosts(posts);
    state.offset += posts.length;
  } catch (err) {
    console.error("Failed to load posts:", err);
  } finally {
    state.loading = false;
  }
}

async function handleAction(postId, type) {
  try {
    const res = await PostResolver({ id: postId, type });
    console.log(res)
    return res; 
  } catch (err) {
    console.error(err);
  }
}

async function handleComment(postId, text) {
  try {
    const res = await CreatComment(postId, text);
    await res.json();
  } catch (err) {
    console.error(err);
  }
}


async function handleCreatePost(form) {
  const formData = new FormData(form);

  const data = {
    title: formData.get("title"),
    text: formData.get("text"),
    categories: formData.getAll("categories"), // if checkbox/multi-select
  };

  try {
    const result = await CreatePost({ data });

    console.log("Post created:", result);

    // reset feed and reload
    resetFeed();
    fetchPosts();

    form.reset();
  } catch (err) {
    console.error("Create post failed:", err.message);
  }
}

/* ======================
   UI RENDER
====================== */
function renderPosts(posts) {
  const container = document.querySelector(".posts");

  const empty = container.querySelector(".no-post");
  if (empty) empty.remove();

  container.insertAdjacentHTML(
    "beforeend",
    posts.map(Post).join("")
  );
}

/* ======================
   FILTERS
====================== */
function toggleFilter(name, button) {
  const input = document.querySelector(`input[name='${name}']`);
  const isActive = button.classList.contains("active");

  button.classList.toggle("active");
  input.value = isActive ? "" : "true";

  resetFeed();
  fetchPosts();
}

/* ======================
   LOGOUT
====================== */
async function handleLogout() {
  try {
    await logout();
    localStorage.clear();
    window.location.href = "/login";
  } catch (err) {
    console.error("Logout failed:", err);
  }
}

/* ======================
   EVENTS
====================== */
function setupEvents() {
  /* Scroll */
  window.addEventListener(
    "scroll",
    throttle(() => {
      const scrollTop = window.scrollY;
      const windowHeight = window.innerHeight;
      const docHeight = document.documentElement.scrollHeight;

      if (scrollTop + windowHeight >= docHeight - 200) {
        fetchPosts();
      }
    }, 200)
  );

  /*creat a post */
  const createPostForm = document.getElementById("create-post-form");

if (createPostForm) {
  createPostForm.addEventListener("submit", (e) => {
    e.preventDefault();
    handleCreatePost(e.target);
  });
}

/* click the buttons */
document.addEventListener("click", (e) => {
  const post = e.target.closest(".post");

  if (!post) return;

  if (
    e.target.closest("button") ||
    e.target.closest(".like-btn") ||
    e.target.closest(".dislike-btn") ||
    e.target.closest(".comment-btn") ||
    e.target.closest(".send-comment")
  ) {
    return;
  }

  navigate(`/post/${post.dataset.postId}`);
});

  /* Click delegation */
  document.addEventListener("click", async (e) => {
    const likeBtn = e.target.closest(".like-btn");
    const dislikeBtn = e.target.closest(".dislike-btn");
    const deleteBtn = e.target.closest(".delete-btn");
    const commentBtn = e.target.closest(".comment-btn");
    const sendCommentBtn = e.target.closest(".send-comment");

    const createdBtn = e.target.closest("[name='my-creat-postes']");
    const likedBtn = e.target.closest("[name='my-liked-post']");
    const logoutBtn = e.target.closest("#logout-btn");

    if (createdBtn) return toggleFilter("my-creat-postes", createdBtn);
    if (likedBtn) return toggleFilter("my-liked-post", likedBtn);
    if (logoutBtn) return handleLogout();

if (likeBtn) {
  const res = await handleAction(likeBtn.dataset.id, "like");

  if (res?.message === "liked") {
    updatePostUI(likeBtn.dataset.id, "like", res.data);
  }
}

if (dislikeBtn) {
  const res = await handleAction(dislikeBtn.dataset.id, "dislike");

  if (res?.message === "disliked") {
    updatePostUI(dislikeBtn.dataset.id, "dislike", res.data);
  }
}

if (deleteBtn) {
  const res = await handleAction(deleteBtn.dataset.id, "delete");

  if (res?.message === "deleted") {
    updatePostUI(deleteBtn.dataset.id, "delete");
  }
}
    if (commentBtn) {
      const box = document.getElementById(`comments-${commentBtn.dataset.id}`);
      box.style.display = box.style.display === "none" ? "block" : "none";
    }

    if (sendCommentBtn) {
      const id = sendCommentBtn.dataset.id;
      const input = document.querySelector(`#comments-${id} .comment-input`);

      if (input.value.trim()) {
        await handleComment(id, input.value);
        input.value = "";
      }
    }
  });

  /* Filter form */
  document
    .getElementById("filter-form")
    .addEventListener("submit", (e) => {
      e.preventDefault();

      const params = new URLSearchParams(new FormData(e.target));
      const url = new URL(window.location);

      window.history.pushState({}, "", `${url.pathname}?${params.toString()}`);

      resetFeed();
      fetchPosts();
    });
}

/* ======================
   INIT
====================== */
function setupHomePage() {
  fetchPosts();
  setTimeout(fetchPosts, 50); // safety double fetch
  setupEvents();
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", setupHomePage);
} else {
  setupHomePage();
}


function updatePostUI(postId, action, data) {
  const post = document.querySelector(`.post[data-post-id="${postId}"]`);
  if (!post) return;

  const likeBtn = post.querySelector(".like-btn");
  const dislikeBtn = post.querySelector(".dislike-btn");

  const likeCount = post.querySelector(".like-count");
  const dislikeCount = post.querySelector(".dislike-count");

  if (action === "like") {
    likeBtn?.classList.add("active");
    dislikeBtn?.classList.remove("active");

    if (data) {
      if (likeCount) likeCount.innerText = data.likes;
      if (dislikeCount) dislikeCount.innerText = data.dislikes;
    }
  }

  if (action === "dislike") {
    dislikeBtn?.classList.add("active");
    likeBtn?.classList.remove("active");

    if (data) {
      if (likeCount) likeCount.innerText = data.likes;
      if (dislikeCount) dislikeCount.innerText = data.dislikes;
    }
  }

  if (action === "delete") {
    post.remove();
  }
}