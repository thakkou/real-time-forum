const serverURI = env.serverUri;

export const getPosts = async ({
  categories = [],
  isLiked = false,
  isCreatedByMe = false,
} = {}) => {
  const params = new URLSearchParams();
  

  categories.forEach((c) => params.append("categories", c));

  if (isLiked) {
    params.append("my-liked-posts", "true");
  }

  if (isCreatedByMe) {
    params.append("my-creat-posts", "true");
  }

  const response = await fetch(
    `${serverURI}/api/posts?${params.toString()}`,
    {
      method: "GET",
      credentials: "include",
    }
  );

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || "Failed to fetch posts");
  }

  return data;
};

export const getPostByID = async ({ id }) => {
  const response = await fetch(
    `${serverURI}/api/posts/getPost/${id}`,
    {
      method: "GET",
      credentials: "include",
    }
  );

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || "Failed to fetch post");
  }

  return data;
};

export const CreatePost = async ({ data }) => {
  const response = await fetch(
    `${serverURI}/api/posts/create`,
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

  const result = await response.json();

  if (!response.ok) {
    throw new Error(result.message || "Failed to create post");
  }

  return result;
};

export const PostResolver = async ({ id, type }) => {
  const method = type === "delete" ? "DELETE" : "POST";

  const response = await fetch(
    `${serverURI}/api/posts/${id}/${type}`,
    {
      method,
      credentials: "include",
    }
  );

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || `${type} failed`);
  }

  return data;
};