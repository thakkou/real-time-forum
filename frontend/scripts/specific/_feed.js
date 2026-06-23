import { CommentResolver,CreatComment } from "../../api/comments.js";
import {logout} from '../../api/auth.js'
import { getPosts ,PostResolver} from "../../api/posts.js";
import { Post } from "../../components/Post.js";
let page = 1;
let loading = false;
let hasMore = true;

function throttle(fn, delay = 200) {
  let last = 0;
  return (...args) => {
    const now = Date.now();
    if (now - last < delay) return;
    last = now;
    fn(...args);
  };
}

let offset = 0;

async function fetchPosts() {
  if (loading || !hasMore) return;

  loading = true;

  const params = new URLSearchParams(window.location.search);

  const categories = params.getAll("categories");
  const isLiked = params.has("my-liked-posts");
  const isCreatedByMe = params.has("my-creat-posts");

  try {
    const res = await getPosts({
      offset,
      limit: 5,
      categories,
      isLiked,
      isCreatedByMe,
    });

    const posts = res.data;
    console.log("the posts",posts)

    if (!posts || posts.length === 0) {
      hasMore = false;
      loading = false;
      return;
    }

    renderPosts(posts);

    offset += posts.length;
  } catch (err) {
    console.error("Failed to load posts:", err);
  }

  loading = false;
}


function renderPosts(posts) {
  const container = document.querySelector(".posts");

  console.log("rendering posts:", posts);

  const empty = container.querySelector(".no-post");
  if (empty) empty.remove();

  container.insertAdjacentHTML(
    "beforeend",
    posts.map(Post).join("")
  );
}


async function handleAction(postId, type) {
  try {
    const res = await PostResolver(postId, type);
    const json = await res.json();

    console.log(json.message);

    // refresh or update UI later if needed
  } catch (err) {
    console.error(err);
  }
}

async function handleComment(postId, text) {
  try {
    const res = await CreatComment(postId, text);
    const json = await res.json();

    console.log(json.message);
  } catch (err) {
    console.error(err);
  }
}

function setupHomePage() {
  // initial load
  fetchPosts();
    // give DOM time to settle before scroll logic
  setTimeout(() => {
    fetchPosts(); // safety initial fetch (important fix)
  }, 50);

  // infinite scroll (throttled)
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

  //event of logout
  document.addEventListener("click", (e) => {
  const logoutBtn = e.target.closest("#logout-btn");

  if (logoutBtn) {
    handleLogout();
  }
});

  // event delegation (like / dislike / comments)
  document.addEventListener("click", async (e) => {
    const likeBtn = e.target.closest(".like-btn");
    const dislikeBtn = e.target.closest(".dislike-btn");
    const commentBtn = e.target.closest(".comment-btn");
    const sendCommentBtn = e.target.closest(".send-comment");

    if (likeBtn) {
      const id = likeBtn.dataset.id;
      await handleAction(id, "like");
    }

    if (dislikeBtn) {
      const id = dislikeBtn.dataset.id;
      await handleAction(id, "dislike");
    }

    if (commentBtn) {
      const id = commentBtn.dataset.id;
      const box = document.getElementById(`comments-${id}`);
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
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", setupHomePage);
} else {
  setupHomePage();
}




//logout event
async function handleLogout() {
  try {
    await logout(); // call API

    // clear client state if needed
    localStorage.clear();

    // redirect to login page
    window.location.href = "/login";
  } catch (err) {
    console.error("Logout failed:", err);
  }
}