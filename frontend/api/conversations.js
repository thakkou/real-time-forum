const serverURI = env.serverUri;

export const getConversations = async ({
  offset = 0,
  limit = 30,
} = {}) => {
  const response = await fetch(
    `${serverURI}/api/conversations?offset=${offset}&limit=${limit}`,
    {
      method: "GET",
      credentials: "include",
    }
  );

  const result = await response.json();

  if (!response.ok) {
    throw new Error(result.message || "Failed to fetch conversations");
  }

  return result;
};

export const getConversationById = async (
  conversationId,
  {
    offset = 0,
    limit = 10,
  } = {}
) => {
  const response = await fetch(
    `${serverURI}/api/conversation/${conversationId}?offset=${offset}&limit=${limit}`,
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