import { Header } from '../components/Header.js';
import { Post } from '../components/Post.js';

export async function render() {
  const header = Header();
  const post = Post();
  return `
    ${header}
    ${post}
    `;
}