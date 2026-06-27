import { getPostByID, PostResolver } from "../api/posts.js";
import { CommentResolver, CreatComment } from "../api/comments.js";

export async function setup() {
  try {
    setupEventListeners();
  } catch (err) {
    console.error("Failed to load post:", err.message);
  }
};

function setupEventListeners() {
  setupPostActions();
  setupCommentActions();
  setupCreateComment();
  setupDeletePost();
}

function setupPostActions() {
  document.querySelector(".like-btn")?.addEventListener("click", async (e) => {
    e.stopPropagation();

    const id = e.currentTarget.dataset.id;

    try {
      await PostResolver({
        id,
        type: "like",
      });

      setupPostPage();
    } catch (err) {
      console.error(err);
    }
  });

  document.querySelector(".dislike-btn")?.addEventListener("click", async (e) => {
    e.stopPropagation();

    const id = e.currentTarget.dataset.id;

    try {
      await PostResolver({
        id,
        type: "dislike",
      });

      setupPostPage();
    } catch (err) {
      console.error(err);
    }
  });
}

function setupCommentActions() {
  document.querySelectorAll(".comment-like-btn").forEach((btn) => {
    btn.addEventListener("click", async (e) => {
      const id = e.currentTarget.dataset.id;

      try {
        await CommentResolver({
          id,
          type: "like",
        });

        setupPostPage();
      } catch (err) {
        console.error(err);
      }
    });
  });

  document.querySelectorAll(".comment-dislike-btn").forEach((btn) => {
    btn.addEventListener("click", async (e) => {
      const id = e.currentTarget.dataset.id;

      try {
        await CommentResolver({
          id,
          type: "dislike",
        });

        setupPostPage();
      } catch (err) {
        console.error(err);
      }
    });
  });
}

function setupCreateComment() {
  const form = document.getElementById("comment-form");

  if (!form) return;

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const input = form.querySelector('input[name="comment"]');

    const text = input.value.trim();

    if (!text) return;

    try {
      await CreatComment({
        data: {
          postId: form.dataset.postId,
          text,
        },
      });

      input.value = "";

      setupPostPage();
    } catch (err) {
      console.error(err);
    }
  });
}

function setupDeletePost() {
  const btn = document.querySelector(".delete-btn");

  if (!btn) return;

  btn.addEventListener("click", async (e) => {
    e.stopPropagation();

    const postId = getPostIdFromURL();

    if (!confirm("Delete this post?")) return;

    try {
      await PostResolver({
        id: postId,
        type: "delete",
      });

      navigate("/");
    } catch (err) {
      console.error(err);
    }
  });
}

export function getPostIdFromURL() {
  const parts = window.location.pathname.split("/");
  return parts[2];
}