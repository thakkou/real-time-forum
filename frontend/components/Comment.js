import { sanitize } from "../scripts/helpers.js";
export const Comment = (comment) => (`
    <div class="comment" data-comment-id="${comment.Id}">
        <div class="comment-header">
            <div>
                <div class="comment-username">
                    ${comment.Nickname}
                </div>

                <span class="comment-meta">
                    ${comment.TimeAgo}
                </span>
            </div>

            ${ window.profile?.id == comment.UserId ? `
                <button
                    class="comment-delete-btn btn small danger"
                    type="button"
                    data-id="${comment.Id}"
                    >
                    <i class="fa-solid fa-trash" style="color: rgb(255, 255, 255);"></i>
                </button>
            ` : "" }
        </div>

        <div class="comment-text">
            ${sanitize(comment.Text)}
        </div>

       <div class="comment-actions">
    <button
        class="comment-like-btn ${
        comment.IsLiked === 1 ? "active" : ""
        }"
        data-id="${comment.Id}"
    >
        👍 <span class="comment-like-count">${comment.LikeCount}</span>
    </button>

    <button
        class="comment-dislike-btn ${
        comment.IsLiked === -1 ? "active" : ""
        }"
        data-id="${comment.Id}"
    >
        👎 <span class="comment-dislike-count">${comment.DislikeCount}</span>
    </button>
</div>
    </div>
`);