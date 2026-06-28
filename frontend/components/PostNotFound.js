export const PostNotFound = () => `
  <div class="card error-card">
    <div style="font-family: var(--mono); font-size: 0.65rem; color: var(--accent); margin-bottom: 1rem; letter-spacing: 0.15em;">
      ⚡ ERROR_404 // RESOURCE_NOT_FOUND
    </div>
    
    <h1>404</h1>
    <p>The post you are searching for has been removed, deleted, or never existed in the database index.</p>
    
    <div style="margin-top: 2rem; border-top: 1px solid var(--border); padding-top: 1rem;">
      <a onclick="navigate('/')" style="cursor: pointer; display: inline-flex; align-items: center; gap: 0.5rem;">
        <span>← RETURN TO FEED</span>
      </a>
    </div>
  </div>
`;