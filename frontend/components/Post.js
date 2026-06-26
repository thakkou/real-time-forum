import { Comment } from "./Comment.js";

export const Post = (post, options = { withComments: false }) => `
	${ options.withComments ? `<a onclick="navigate('/')" style="font-size: 12px;">← Back to Home</a>` : '' }
	<article class="post ${ options.withComments ? 'detailed-post' : '' }" data-post-id="${post.Id}">
		<div class="post-header">
			<div class="delete-block">
				<h3>${post.Title}</h3>

				${window.profile?.id == post.UserId ? `
					<button class="delete-btn btn small danger" type="button">
					<i class="fa-solid fa-trash" style="color: rgb(255, 255, 255);"></i>
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
				<div class="image-preview active">
					<img src="${post.Image}"
						alt="Post image"
						onerror="this.onerror = null; this.src = '/assets/default.jpg';"
					/>
				</div>
			</div>
		` : ""}

		<div class="post-body">
			<pre>${post.Text}</pre>
		</div>

		<div class="post-actions">
			<button class="like-btn ${post.IsLiked === 1 ? "active" : ""}" data-id="${post.Id}">
			👍 <span class="like-count">${post.LikeCount}</span>
			</button>

			<button class="dislike-btn ${post.IsLiked === -1 ? "active" : ""}" data-id="${post.Id}">
			👎 <span class="dislike-count">${post.DislikeCount}</span>
			</button>
		</div>

		${ options.withComments ? `
			<section class="comments">
				<h4>Comments (${post.Comments?.length || 0})</h4>

				<div class="comments-list">
				${ post.Comments?.length ?
					post.Comments.map(Comment).join("") : `
						<div class="comment">
							<span class="comment-meta">
								No comments yet.
							</span>
						</div>
					` }
				</div>

				<div class="add-comment">
					<form id="comment-form" data-post-id="${post.Id}">
						<input
						type="text"
						name="comment"
						placeholder="Write a comment..."
						maxlength="500"
						required
						/>

						<button
						type="submit"
						class="btn primary small"
						>
						Comment
						</button>
					</form>
				</div>
			</section>
		` : '' }
	</article>
`;