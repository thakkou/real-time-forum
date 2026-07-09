
export async function render(data = {}) {
  return `
    <div class="post-detail-wrapper">
      <div class="loading-state" style="font-family: var(--mono); font-size: 0.75rem; color: var(--text-muted); padding: 3rem 2rem; text-align: center;">
        ▶ FETCHING_POST_INDEX...
      </div>
    </div>
  `;
}