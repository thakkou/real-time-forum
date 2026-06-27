import { Header } from '../components/Header.js';
import { getPostByID } from "../../api/posts.js";
import { getPostIdFromURL } from "../scripts/_post.js";
import { Post } from "../../components/Post.js";

export async function render(data = {}) {
  const id = getPostIdFromURL();
  if (!id) {
    console.error("No post id found in URL");
    return;
  }

  const response = await getPostByID({ id });
  return `
    ${Header(data.nickname)}
    <div class="post-detail-wrapper">
      ${Post(response.data, { withComments: true })}
    </div>
  `;
}