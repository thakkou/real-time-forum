import { Header } from '../components/Header.js';
import { Post } from '../components/Post.js';

export async function render(data = {}) {
  const header = Header(data.nickname);
  const post = Post();
  return `
    ${header}
    ${post}
    `;
}