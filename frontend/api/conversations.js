const serverURI = window.env.serverUri;
import { handleAuthError } from "./helperApi.js";
export const getConversations = async ({
  offset = 0,
  limit = 30,
} = {}) => {
  const response = await fetch(
    `${serverURI}/conversations?offset=${offset}&limit=${limit}`,
    {
      method: "GET",
      credentials: "include",
    }
  );
  if (handleAuthError(response)) return;

  const result = await response.json();

  if (!response.ok) {
    throw new Error(result.message || "Failed to fetch conversations");
  }

  return result;
};

export const getConversationById = async (
  conversationId,
  pagination={
    offset :0,
    limit :10,
  } 
) => {
  
  const response = await fetch(
    `${serverURI}/conversation/${conversationId}?offset=${pagination.offset}&limit=${pagination.limit}`,
    {
      method: "GET",
      credentials: "include",
    }
  );

  const result = await response.json();

  if (!response.ok) {
    throw new Error(result.message || "Failed to fetch conversation");
  }

  return result;
};