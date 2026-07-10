const serverURI = env.serverUri;
import { handleAuthError } from "./helperApi.js";

export const getPosts = async ({
  lastId = null,
  limit = 30,
  categories = [],
  isLiked = false,
  isCreatedByMe = false,
} = {}) => {
  const params = new URLSearchParams();

  if (lastId !== null) {
    params.append("lastId", lastId);
  }
  params.append("limit", limit);

  categories.forEach((c) => params.append("categories", c));

  if (isLiked) {
    params.append("my-liked-posts", "true");
  }

  if (isCreatedByMe) {
    params.append("my-creat-posts", "true");
  }

  const response = await fetch(
    `${serverURI}/posts?${params.toString()}`,
    {
      method: "GET",
      credentials: "include",
    }
  );

  // Check for unauthorized access before parsing JSON
  if (handleAuthError(response)) return;

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || "Failed to fetch posts");
  }

  return data;
};

export const getPostByID = async ({ id, commentLimit = 10, commentLastId = null } = {}) => {
  const params = new URLSearchParams();

  if (commentLimit !== null) {
    params.append("commentLimit", commentLimit);
  }
  if (commentLastId !== null) {
    params.append("commentLastId", commentLastId);
  }

  const url = `${serverURI}/posts/${id}${params.toString() ? `?${params.toString()}` : ""}`;
  const response = await fetch(url, {
    method: "GET",
    credentials: "include",
  });

  if (handleAuthError(response)) return;

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || "Failed to fetch post");
  }

  return data;
};

export const CreatePost = async ({ data }) => {
  const response = await fetch(
    `${serverURI}/posts/create`,
    {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        title: data.title,
        text: data.text,
        categories: data.categories,
      }),
    }
  );

  if (handleAuthError(response)) return;

  const result = await response.json();

  if (!response.ok) {
    throw new Error(result.message || "Failed to create post");
  }

  return result;
};

export const PostResolver = async ({ id, type }) => {
  console.log("call api to", type, "for", id);
  const method = type === "delete" ? "DELETE" : "POST";

  const response = await fetch(
    `${serverURI}/posts/${id}/${type}`,
    {
      method,
      credentials: "include",
    }
  );

  if (handleAuthError(response)) return;

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || `${type} failed`);
  }

  return data;
};