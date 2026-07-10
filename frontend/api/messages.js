
const serverURI = env.serverUri;
import { handleAuthError } from "./helperApi.js";

export const createMessage = async ({
  receiverId,
  text,
  conversationId = null,
}) => {
  const response = await fetch(
    `${serverURI}/messages`,
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

  // Check for unauthorized access before parsing the response body
  if (handleAuthError(response)) return;

  const result = await response.json();

  if (!response.ok) {
    throw new Error(result.message || "Failed to send message");
  }

  return result;
};