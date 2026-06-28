import { Header } from '../components/Header.js';

/**
 * Main Layout Renderer
 * Simply sets up the shell structural wrappers.
 * All data loading, success states, and 404 errors are managed by setup().
 */
export async function render(data = {}) {
  return `
    ${Header(data.nickname)}
    <div class="post-detail-wrapper">
      <div class="loading-state" style="font-family: var(--mono); font-size: 0.75rem; color: var(--text-muted); padding: 3rem 2rem; text-align: center;">
        ▶ FETCHING_POST_INDEX...
      </div>
    </div>
  `;
}