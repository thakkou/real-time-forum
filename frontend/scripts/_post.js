/* ================================================================
   IMPORTS
   ================================================================ */
import { Header } from '../components/Header.js';
import { Post } from "../components/Post.js";
import { Comment } from "../components/Comment.js";
import { PostNotFound } from "../components/PostNotFound.js";
import { getPostByID, PostResolver } from "../api/posts.js";
import { CommentResolver, CreatComment } from "../api/comments.js";

/* ================================================================
   INITIALIZATION & RENDER LIEFOCYCLE
   ================================================================ */
export async function setup() {
  console.log(1)
  try {  
      await setupPostPage(); 

    setupEventListeners();
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

  try {
    const res = await getPostByID({ id: postId });
    
    if (res && res.data) {
      wrapper.innerHTML = Post(res.data, { withComments: true });
    } else {
      wrapper.innerHTML = PostNotFound();
    }
  } catch (err) {
    console.error("Failed to synchronize layout view:", err);
    wrapper.innerHTML = PostNotFound();
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
        await PostResolver({ id, type });
        await setupPostPage();
      } catch (err) {
        console.error(err);
      }
      return;
    }

    // Comment Like/Dislike
    if (commentLikeBtn || commentDislikeBtn) {
      const id = (commentLikeBtn || commentDislikeBtn).dataset.id;
      const type = commentLikeBtn ? "like" : "dislike";
      try {
        await CommentResolver({ id, type });
        await setupPostPage();
      } catch (err) {
        console.error(err);
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
      }

      try {
        await CommentResolver({ id, type: "delete" });
      } catch (err) {
        console.error("Server failed to delete comment:", err);
        alert("Could not remove comment from server. Reloading feed...");
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
        alert("Failed to delete post from database. Reloading...");
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
      alert("Failed to post comment. Please try again.");
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
  commentsListContainer.insertAdjacentHTML("beforeend", commentHTML);
}

/* ================================================================
   UTILS
   ================================================================ */
export function getPostIdFromURL() {
  const parts = window.location.pathname.split("/");
  return parts[2];
}