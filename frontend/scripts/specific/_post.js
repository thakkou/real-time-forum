
import { getPostByID } from "../../api/posts.js";
import { Post } from "../../components/Post.js";

const setupPostPage = async () => {
  try {
    const id = getPostIdFromURL();

    if (!id) {
      console.error("No post id found in URL");
      return;
    }

    console.log("Loading post:", id);

    const data = await getPostByID({ id });

    const container = document.getElementById("post-detaille");

    container.innerHTML = Post(data);
  } catch (err) {
    console.error("Failed to load post:", err.message);
  }
};

function getPostIdFromURL() {
  const parts = window.location.pathname.split('/');
  return parts[2];
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', setupPostPage);
} else {
    setupPostPage();
}





