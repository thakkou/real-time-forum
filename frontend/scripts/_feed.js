import { getPosts, PostResolver, CreatePost } from "../api/posts.js";
import { Post } from "../components/Post.js";
import { showToast } from "../services/toast.js";
import { router } from "../services/router.js";
/* ======================
   STATE
====================== */
const state = {
  posts: [],
lastId: null,
  limite: 10,
  loading: false,
};
function resetState(){
 state.posts = [];
  state.lastId = null;
  state.loading = false;
}

/* ======================
   INIT
====================== */
/* ======================
   EVENT HANDLER REFERENCES
====================== */
let cachedScrollHandler = null;
let cachedClickHandler = null;
let cachedToggleHandler = null;
let cachedCreateSubmitHandler = null;
let cachedFilterSubmitHandler = null;

export function cleanup() {
  console.log("Cleaning up old feed event listeners...");

  if (cachedScrollHandler) {
    console.log(cachedScrollHandler)
    window.removeEventListener("scroll", cachedScrollHandler);
  }
  if (cachedClickHandler) {
    document.removeEventListener("click", cachedClickHandler);
  }

  const details = document.getElementById("create-post-details");
  if (details && cachedToggleHandler) {
    details.removeEventListener("toggle", cachedToggleHandler);
  }

  const createPostForm = document.getElementById("create-post-form");
  if (createPostForm && cachedCreateSubmitHandler) {
    createPostForm.removeEventListener("submit", cachedCreateSubmitHandler);
  }

  const filterForm = document.getElementById("filter-form");
  if (filterForm && cachedFilterSubmitHandler) {
    filterForm.removeEventListener("submit", cachedFilterSubmitHandler);
  }
}


export function setup() {
resetState()
cleanup(); // Wipe out any lingering event listeners before binding fresh ones

  fetchPosts();
  setupEvents();
}

function resetFeed() {
  state.lastId = null;
  state.posts = [];
  document.querySelector(".posts").innerHTML = "";
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
  const isLiked = params.get("my-liked-post") === "true";
  const isCreatedByMe = params.get("my-creat-postes") === "true";

  try {
    console.log("getPosts",state.lastId)
    const res = await getPosts({
      lastId: state.lastId, // Pass lastId instead of offset
      limit: state.limite,
      categories,
      isLiked,
      isCreatedByMe,
    });

    const posts = res.data;

    if (posts?.length) {
      state.posts.push(...posts);
      console.log(posts)
      // The oldest post in this freshly fetched batch becomes our new anchor point
      state.lastId = posts[posts.length - 1].Id; 
      
      renderPosts(posts);
    }
  } catch (err) {
    console.error("Failed to load posts:", err);
  } finally {
    state.loading = false;
  }
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
    showToast(err.message || "Failed to create post", "error");
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

}




/* ======================
   EVENTS
====================== */

/* ======================
   EVENTS SETUP
====================== */
function setupEvents() {
  const details = document.getElementById("create-post-details");
  const createPostForm = document.getElementById("create-post-form");
  const filterForm = document.getElementById("filter-form");

  // 1. Details toggle listener
  cachedToggleHandler = () => {
    console.log(details.open ? "Form opened" : "Form closed");
  };
  details?.addEventListener("toggle", cachedToggleHandler);

  // 2. Window Scroll infinite loading listener
  cachedScrollHandler = throttle(() => {
    const scrollTop = window.scrollY;
    const windowHeight = window.innerHeight;
    const docHeight = document.documentElement.scrollHeight;

    if (scrollTop + windowHeight >= docHeight - 200) {
      fetchPosts();
    }
  }, 200);
  window.addEventListener("scroll", cachedScrollHandler);

  // 3. Post Form creation submit listener
  cachedCreateSubmitHandler = (e) => {
    e.preventDefault();
    handleCreatePost(e.target);
  };
  createPostForm?.addEventListener("submit", cachedCreateSubmitHandler);

  // 4. Global Click delegation handler
  cachedClickHandler = async (e) => {
    const post = e.target.closest(".post");
    const likeBtn = e.target.closest(".like-btn");
    const dislikeBtn = e.target.closest(".dislike-btn");
    const deleteBtn = e.target.closest(".delete-btn");
    const createdBtn = e.target.closest("[name='my-creat-postes']");
    const likedBtn = e.target.closest("[name='my-liked-post']");

    // Filter Buttons triggers
    if (createdBtn) return toggleFilter("my-creat-postes", createdBtn);
    if (likedBtn) return toggleFilter("my-liked-post", likedBtn);

    // Like Action
    if (likeBtn) {
      const res = await handleAction(likeBtn.dataset.id, "like");
      if (res?.message === "liked") {
        updatePostUI(likeBtn.dataset.id, "like", res.data);
      }
      return;
    }

    // Dislike Action
    if (dislikeBtn) {
      const res = await handleAction(dislikeBtn.dataset.id, "dislike");
      if (res?.message === "disliked") {
        updatePostUI(dislikeBtn.dataset.id, "dislike", res.data);
      }
      return;
    }

    // Delete Action
    if (deleteBtn) {
      const res = await handleAction(deleteBtn.dataset.id, "delete");
      if (res?.message === "deleted") {
        updatePostUI(deleteBtn.dataset.id, "delete");
        showToast("deleted post", "success");
      }
      return;
    }

    // Row Click Redirection (Ignore if clicking general buttons/actions)
    if (post) {
      if (e.target.closest("button") || likeBtn || dislikeBtn || deleteBtn) {
        return;
      }
    cleanup()
      
      router.navigate(`/post/${post.dataset.postId}`);
    }
  };
  document.addEventListener("click", cachedClickHandler);

  // 5. Filter Form submission listener
  cachedFilterSubmitHandler = (e) => {
    e.preventDefault();

    const params = new URLSearchParams(new FormData(e.target));
    const url = new URL(window.location);

    window.history.pushState({}, "", `${url.pathname}?${params.toString()}`);

    resetFeed();
    fetchPosts();
  };
  filterForm?.addEventListener("submit", cachedFilterSubmitHandler);
}

/* ======================
   UPDATE UI
====================== */

export function updatePostUI(postId, action, data) {
  console.log("start update the post UI", postId, action, data);

  const post = document.querySelector(`.post[data-post-id="${postId}"]`);
  if (!post) return;

  if (action === "delete") {
    post.remove();
    return;
  }

  const likeBtn = post.querySelector(".like-btn");
  const dislikeBtn = post.querySelector(".dislike-btn");

  const likeCount = post.querySelector(".like-count");
  const dislikeCount = post.querySelector(".dislike-count");

  // Update counts
  if (data) {
    if (likeCount) likeCount.innerText = data.likes;
    if (dislikeCount) dislikeCount.innerText = data.dislikes;

    // Reset both buttons
    likeBtn?.classList.remove("active");
    dislikeBtn?.classList.remove("active");

    // Activate the correct one
    switch (data.isLike) {
      case 1:
        likeBtn?.classList.add("active");
        break;

      case -1:
        dislikeBtn?.classList.add("active");
        break;

      case 0:
      default:
        // Neither button is active
        break;
    }
  }
}