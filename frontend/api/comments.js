const serverURI = env.serverUri;
import { handleAuthError } from "./helperApi.js";

export const CreatComment = async ({ data }) => {
  const response = await fetch(
    `${serverURI}/comments/create`,
    {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        postId: data.postId,
        text: data.text,
      }),
    }
  );
    if (handleAuthError(response)) return;


  const result = await response.json();

  if (!response.ok) {
    throw new Error(result.message || "Failed to create comment");
  }

  return result;
};

export const CommentResolver = async ({ id, type }) => {
  const method = type === "delete" ? "DELETE" : "POST";

  const response = await fetch(
    `${serverURI}/comments/${id}/${type}`,
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