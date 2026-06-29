// import { logout } from "../api/auth.js";
import { getPosts, PostResolver, CreatePost } from "../api/posts.js";
import { Post } from "../components/Post.js";
import { showToast } from "../services/toast.js";

/* ======================
   STATE
====================== */
const state = {
  posts: [],
  offset: 0,
  limite: 15,
  loading: false,
};

/* ======================
   INIT
====================== */
export function setup() {
  fetchPosts();
  setupEvents();
}

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

/* ======================
   API ACTIONS
====================== */

async function fetchPosts() {
  if (state.loading) return;

  state.loading = true;

  const params = new URLSearchParams(window.location.search);
  const categories = params.getAll("categories");
  const isLiked = params.get("my-liked-posts") === "true";
  const isCreatedByMe = params.get("my-creat-posts") === "true";

  try {
    const res = await getPosts({
      offset: state.offset,
      limit: state.limite,
      categories,
      isLiked,
      isCreatedByMe,
    });

    const posts = res.data;

    if (posts?.length) {
      state.posts.push(...posts);
      state.offset += posts.length;
    }
    renderPosts(state.posts);
  } catch (err) {
    console.error("Failed to load posts:", err);
  } finally {
    state.loading = false;
  }
}

function resetFeed() {
  state.offset = 0;
  state.posts = [];
  document.querySelector(".posts").innerHTML = "";
}

async function handleAction(postId, type) {
  try {
    return await PostResolver({ id: postId, type });
  } catch (err) {
    console.error(err);
  }
}

function prependPostToUI(post) {
  const container = document.querySelector(".posts");
  if (!container) return;

  const empty = container.querySelector(".no-post");
  if (empty) empty.remove();

  container.insertAdjacentHTML("afterbegin", Post(post));
}

/* ======================
   CREATE POST
====================== */

async function handleCreatePost(form) {
  const formData = new FormData(form);

  const data = {
    title: formData.get("title"),
    text: formData.get("text"),
    categories: formData.getAll("categories"),
  };

  try {
    const result = await CreatePost({ data });

    if (result && result.data) {
      const newPost = result.data;

      state.posts.unshift(newPost);
      state.offset += 1;

      const details = document.getElementById("create-post-details");
      if (details) details.open = false;

      showToast("Post added", "success");

      prependPostToUI(newPost);
    } else {
      resetFeed();
      fetchPosts();
    }

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

  container.insertAdjacentHTML("beforeend", posts.map(Post).join(""));
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
   EVENTS
====================== */

function setupEvents() {
  const details = document.getElementById("create-post-details");

  details?.addEventListener("toggle", () => {
    console.log(details.open ? "Form opened" : "Form closed");
  });

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

  const createPostForm = document.getElementById("create-post-form");

  createPostForm?.addEventListener("submit", (e) => {
    e.preventDefault();
    handleCreatePost(e.target);
  });

  document.addEventListener("click", (e) => {
    const post = e.target.closest(".post");

    if (!post) return;

    if (
      e.target.closest("button") ||
      e.target.closest(".like-btn") ||
      e.target.closest(".dislike-btn")
    ) {
      return;
    }

    navigate(`/post/${post.dataset.postId}`);
  });

  document.addEventListener("click", async (e) => {
    const likeBtn = e.target.closest(".like-btn");
    const dislikeBtn = e.target.closest(".dislike-btn");
    const deleteBtn = e.target.closest(".delete-btn");

    const createdBtn = e.target.closest("[name='my-creat-posts']");
    const likedBtn = e.target.closest("[name='my-liked-posts']");
    // const logoutBtn = e.target.closest("#logout-btn");

    if (createdBtn) return toggleFilter("my-creat-posts", createdBtn);
    if (likedBtn) return toggleFilter("my-liked-posts", likedBtn);
    // if (logoutBtn) return handleLogout();

    if (likeBtn) {
      const res = await handleAction(likeBtn.dataset.id, "like");
      if (res?.message === "liked") {
        updatePostUI(likeBtn.dataset.id, "like", res.data);
        showToast("liked post", "success");
      }
    }

    if (dislikeBtn) {
      const res = await handleAction(dislikeBtn.dataset.id, "dislike");
      if (res?.message === "disliked") {
        updatePostUI(dislikeBtn.dataset.id, "dislike", res.data);
        showToast("disliked post", "success");
      }
    }

    if (deleteBtn) {
      const res = await handleAction(deleteBtn.dataset.id, "delete");
      if (res?.message === "deleted") {
        updatePostUI(deleteBtn.dataset.id, "delete");
        showToast("deleted post", "success");
      }
    }
  });

  document
    .getElementById("filter-form")
    ?.addEventListener("submit", (e) => {
      e.preventDefault();

      const params = new URLSearchParams(new FormData(e.target));
      const url = new URL(window.location);

      window.history.pushState({}, "", `${url.pathname}?${params.toString()}`);

      resetFeed();
      fetchPosts();
    });
}

/* ======================
   UPDATE UI
====================== */

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