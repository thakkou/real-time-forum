import { Header } from '../components/Header.js';

export async function render(data = {}) {


  const header = Header(data.nickname);
  return `
    ${header}
    <div id="post-detaille">
    </div>
    
    `;
}