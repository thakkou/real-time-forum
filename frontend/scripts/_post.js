
import { Header } from '../components/Header.js';
import { Post } from "../components/Post.js";
import { Comment } from "../components/Comment.js";
import { PostNotFound } from "../components/PostNotFound.js";
import { getPostByID, PostResolver } from "../api/posts.js";
import { CommentResolver, CreatComment } from "../api/comments.js";
import { updatePostUI } from './_feed.js';
import { showToast } from '../services/toast.js';
/* ================================================================
   INITIALIZATION & RENDER LIEFOCYCLE
   ================================================================ */




   /* ================================================================
   UI STATE
================================================ */
const ui = {
  wrapper: null,
  post: null,

  likeBtn: null,
  dislikeBtn: null,
  deleteBtn: null,

  likeCount: null,
  dislikeCount: null,
  commentCount: null,

  commentsList: null,
  commentForm: null,
  loadMoreBtn: null,
};

const state = {
  postId: null,
  commentLimit: 10,
  commentLastId: null,
  hasMoreComments: false,
  isLoadingComments: false,
};

function cacheUI() {
  ui.post = document.querySelector(".post");

  if (!ui.post) return;

  ui.likeBtn = ui.post.querySelector(".like-btn");
  ui.dislikeBtn = ui.post.querySelector(".dislike-btn");
  ui.deleteBtn = ui.post.querySelector(".delete-btn");

  ui.likeCount = ui.post.querySelector(".like-count");
  ui.dislikeCount = ui.post.querySelector(".dislike-count");
  ui.commentCount = ui.post.querySelector(".comment-count");

  ui.commentsList = ui.post.querySelector(".comments-list");
  ui.commentForm = ui.post.querySelector("#comment-form");
}

export async function setup() {
  console.log("setup post id")
  try {
    setupEventListeners();
    await setupPostPage(); 
  } catch (err) {
    console.error("Failed to load page:", err.message);
  }
}

export async function render(data = {}) {
  // Returns immediate base structural wrapper layout
  return `
    ${Header(data.nickname)}
    <main class="content">
      <div class="post-detail-wrapper">
        <div class="loading-state" style="font-family: var(--mono); font-size: 0.75rem; color: var(--text-muted); padding: 2rem;">
          ▶ LOADING_POST_DATA...
        </div>
      </div>
    </main>
  `;
}

/* ================================================================
   CORE UI RENDERER
   ================================================================ */
async function setupPostPage() {
  const postId = getPostIdFromURL();
  const wrapper = document.querySelector(".post-detail-wrapper");
  if (!wrapper) return;

  if (!postId) {
    wrapper.innerHTML = PostNotFound();
    return;
  }

  state.postId = postId;
  state.commentLastId = null;
  state.hasMoreComments = true;

  try {
    await loadPostPage(postId, false);
  } catch (err) {
    console.error("Failed to synchronize layout view:", err);
    wrapper.innerHTML = PostNotFound();
  }
}


function updateCommentCount(change) {
  const countEl = document.querySelector(".comment-count");
  if (!countEl) return;

  const current = parseInt(countEl.textContent, 10) || 0;
  countEl.textContent = Math.max(0, current + change);
}

async function loadPostPage(postId, append = false) {
  const res = await getPostByID({
    id: postId,
    commentLimit: state.commentLimit,
    commentLastId: append ? state.commentLastId : null,
  });

  if (!res?.data) {
    throw new Error("Failed to load post data");
  }

  if (!append) {
    const wrapper = document.querySelector(".post-detail-wrapper");
    wrapper.innerHTML = Post(res.data, { withComments: true });
    ui.wrapper = wrapper;
    cacheUI();
    attachLoadMoreHandler();
  }

  const comments = res.data.Comments || [];
  renderComments(comments, append);

  if (comments.length > 0) {
    const lastComment = comments[comments.length - 1];
    state.commentLastId = lastComment.Id;
    state.hasMoreComments = comments.length === state.commentLimit;
  } else {
    state.hasMoreComments = false;
  }

  updateLoadMoreButton();
}

function renderComments(comments, append = false) {
  if (!ui.commentsList) return;

  const html = comments.map(Comment).join("");

  if (!append) {
    ui.commentsList.innerHTML = html || `
      <div class="comment">
        <span class="comment-meta">No comments yet.</span>
      </div>
    `;
    return;
  }

  const fallbackNode = ui.commentsList.querySelector(".comment-meta");
  if (fallbackNode && fallbackNode.textContent.includes("No comments yet")) {
    ui.commentsList.innerHTML = html;
  } else {
    ui.commentsList.insertAdjacentHTML("beforeend", html);
  }
}

function attachLoadMoreHandler() {
  ui.loadMoreBtn = ui.post.querySelector("#loadMoreCommentsBtn");
  if (!ui.loadMoreBtn) return;

  ui.loadMoreBtn.addEventListener("click", loadMoreComments);
  updateLoadMoreButton();
}

function updateLoadMoreButton() {
  if (!ui.loadMoreBtn) return;
  ui.loadMoreBtn.style.display = state.hasMoreComments ? "block" : "none";
}

async function loadMoreComments() {
  if (!state.hasMoreComments || state.isLoadingComments) return;

  state.isLoadingComments = true;
  if (ui.loadMoreBtn) {
    ui.loadMoreBtn.disabled = true;
  }

  try {
    await loadPostPage(state.postId, true);
  } catch (err) {
    console.error("Failed to load more comments:", err);
  } finally {
    state.isLoadingComments = false;
    if (ui.loadMoreBtn) {
      ui.loadMoreBtn.disabled = false;
    }
  }
}

/* ================================================================
   EVENTS HANDLER (Delegation Mode)
   ================================================================ */
function setupEventListeners() {
  
  document.addEventListener("click", async (e) => {
    const likeBtn = e.target.closest(".like-btn");
    const dislikeBtn = e.target.closest(".dislike-btn");
    const commentLikeBtn = e.target.closest(".comment-like-btn");
    const commentDislikeBtn = e.target.closest(".comment-dislike-btn");
    const commentDeleteBtn = e.target.closest(".comment-delete-btn"); 
    const deleteBtn = e.target.closest(".delete-btn");


    // Post Like/Dislike
    if (likeBtn || dislikeBtn) {
      e.stopPropagation();
      const id = (likeBtn || dislikeBtn).dataset.id;
      const type = likeBtn ? "like" : "dislike";
      try {
        const data =  await PostResolver({ id, type });
        if(data){
           updatePostUI(id,type,data.data)
        }
      } catch (err) {
        console.error(err);
        showToast(err.message || "Action failed", "error");
      }
      return;
    }

    // Comment Like/Dislike
    if (commentLikeBtn || commentDislikeBtn) {

      const id = (commentLikeBtn || commentDislikeBtn).dataset.id;
      const type = commentLikeBtn ? "like" : "dislike";
      try {
        const data = await CommentResolver({ id, type });

updateCommentUI(id,type,data.data) 
     } catch (err) {
        console.error(err);
        showToast(err.message || "Action failed", "error");
      }
      return;
    }

    // INSTANT: Comment Delete Action
    if (commentDeleteBtn) {
      e.stopPropagation();
      const id = commentDeleteBtn.dataset.id;
      if (!confirm("Delete this comment?")) return;

      const commentTarget = commentDeleteBtn.closest(".comment"); 
      if (commentTarget) {
        commentTarget.remove(); 
        updateCommentCount(-1);      }

      try {
        await CommentResolver({ id, type: "delete" });
      } catch (err) {
        console.error("Server failed to delete comment:", err);
        showToast(err.message || "Could not remove comment from server. Reloading feed...", "error");
        await setupPostPage(); 
      }
      return;
    }

    // INSTANT: Post Deletion
    if (deleteBtn) {
      e.stopPropagation();
      const postId = getPostIdFromURL();
      if (!confirm("Delete this post?")) return;

      const postTarget = deleteBtn.closest(".post"); 
      if (postTarget) {
        postTarget.remove();
      }

      try {
        await PostResolver({ id: postId, type: "delete" });
        navigate("/"); 
      } catch (err) {
        console.error("Server failed to delete post:", err);
        showToast(err.message || "Failed to delete post from database. Reloading...", "error");
        await setupPostPage(); 
      }
      return;
    }
  });

  // INSTANT: Form Submission handler (Comment Creation)
  document.addEventListener("submit", async (e) => {
    const form = e.target.closest("#comment-form");
    if (!form) return;
    
    e.preventDefault();
    const input = form.querySelector('input[name="comment"]');
    const text = input.value.trim();

    if (!text) return;
    input.value = "";

    try {
      const res = await CreatComment({
        data: {
          postId: parseInt(form.dataset.postId, 10) || form.dataset.postId,
          text,
        },
      });

      if (res && res.data) {
        appendCommentToUI(res.data);
      } else {
        await setupPostPage(); 
      }

    } catch (err) {
      console.error("Comment creation failed:", err);
      showToast(err.message || "Failed to post comment. Please try again.", "error");
      await setupPostPage(); 
    }
  });
}

/* ================================================================
   DOM MUTATION COMPONENT
   ================================================================ */
function appendCommentToUI(comment) {
  const commentsListContainer = document.querySelector(".comments-list");
  if (!commentsListContainer) return;

  const fallbackNode = commentsListContainer.querySelector(".comment-meta");
  if (fallbackNode && fallbackNode.textContent.includes("No comments yet")) {
    commentsListContainer.innerHTML = "";
  }

  const commentFormat = {
    Id: comment.id,
    UserId: comment.userId,
    Nickname: comment.nickname,
    Text: comment.text,
    TimeAgo: "now",
    LikeCount: 0,
    DislikeCount: 0,
    IsLiked: 0
  };

  const commentHTML = Comment(commentFormat);
  const insertPosition = commentsListContainer.firstChild ? "afterbegin" : "beforeend";
  commentsListContainer.insertAdjacentHTML(insertPosition, commentHTML);
  updateCommentCount(1);
}

/* ================================================================
   UTILS
   ================================================================ */
export function getPostIdFromURL() {
  const parts = window.location.pathname.split("/");
  return parts[2];
}
const updateCommentUI = (id, type, data) => {
  const comment = document.querySelector(`.comment[data-id="${id}"]`) 
    || document.querySelector(`.comment[data-comment-id="${id}"]`);

    console.log("update coment ui",data)
  if (!comment) return;

  // Delete case
  if (type === "delete") {
    comment.remove();
    return;
  }
  console.log(comment)

  const likeBtn = comment.querySelector(".comment-like-btn");
  const dislikeBtn = comment.querySelector(".comment-dislike-btn");

  const likeCount = comment.querySelector(".comment-like-count");
  const dislikeCount = comment.querySelector(".comment-dislike-count");
console.log(likeCount)
console.log(dislikeCount)
  // Update counters
  if (likeCount) {
    console.log("likes",data.likes)
    likeCount.textContent = data.likes;
  }

  if (dislikeCount) {
        console.log("dislikes",data.dislikes)

    dislikeCount.textContent = data.dislikes;
  }

  // Remove old active state
  likeBtn?.classList.remove("active");
  dislikeBtn?.classList.remove("active");

  // Set current reaction
  if (data.userReaction === "like") {
    likeBtn?.classList.add("active");
  }

  if (data.userReaction === "dislike") {
    dislikeBtn?.classList.add("active");
  }
};