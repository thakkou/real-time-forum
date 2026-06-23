export const Post = (post) => `
<article class="post" data-post-id="${post.Id}" onclick="navigate('/post/${post.Id}')">
  <div class="post-header">
    <div class="delete-block">
      <h3>${post.Title}</h3>

      ${post.IsOwner ? `
        <button class="delete-btn btn small danger" type="button">
          🗑
        </button>
      ` : ""}
    </div>

    <div class="post-categories">
      ${post.Categories.map(c => `<span class="category">${c}</span>`).join("")}
    </div>

    <span class="post-meta">
      Posted by ${post.Nickname} · ${post.TimeAgo}
    </span>
  </div>

  ${post.Image ? `
    <div class="post-image-container">
      <img src="${post.Image}" alt="Post image" />
    </div>
  ` : ""}

  <div class="post-body">
    <pre>${post.Text}</pre>
  </div>

  <div class="post-actions">
    <button class="like-btn ${post.IsLiked === 1 ? "active" : ""}">
      👍 ${post.LikeCount}
    </button>

    <button class="dislike-btn ${post.IsLiked === -1 ? "active" : ""}">
      👎 ${post.DislikeCount}
    </button>
  </div>
</article>
`;