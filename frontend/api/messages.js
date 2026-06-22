const serverURI = env.serverUri;

export const createMessage = async ({
  receiverId,
  text,
  conversationId = null,
}) => {
  const response = await fetch(
    `${serverURI}/api/messages`,
    {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        receiver_id: receiverId,
        text,
        conversation_id: conversationId,
      }),
    }
  );

  const result = await response.json();

  if (!response.ok) {
    throw new Error(result.message || "Failed to send message");
  }

  return result;
};