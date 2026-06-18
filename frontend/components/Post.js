export const Post = () => (`
    <article class="post" data-post-id="{{.Id}}">
        <div class="post-header">
        <div class="delete-block">
            <h3>{{.Title}}</h3>

            {{if and $.IsLoggedIn (eq $.User.Id .UserId)}}
            <button class="delete-btn btn small danger" type="submit">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512">
                <path d="M136.7 5.9L128 32 32 32C14.3 32 0 46.3 0 64S14.3 96 32 96l384 0c17.7 0 32-14.3 32-32s-14.3-32-32-32l-96 0-8.7-26.1C306.9-7.2 294.7-16 280.9-16L167.1-16c-13.8 0-26 8.8-30.4 21.9zM416 144L32 144 53.1 467.1C54.7 492.4 75.7 512 101 512L347 512c25.3 0 46.3-19.6 47.9-44.9L416 144z"/>
                </svg>
            </button>
            {{end}}
        </div>

        <!-- Categories moved here, inside post-header -->
        <div class="post-categories">
            {{range .Categories}}
            <span class="category">{{.}}</span>
            {{end}}
        </div>

        <span class="post-meta">Posted by {{.Username}} · {{.TimeAgo}}</span>
        </div>

        {{if .Image}}
        <div class="post-image-container"> <!-- style="margin: 1rem 0;"> -->
            <div class="image-preview active">
                <img
                src="{{.Image}}"
                alt="Post image"
                onerror="this.onerror = null; this.src = '/assets/default.jpg';"
                >
            </div>
        </div>
        {{end}}

        <div class="post-body">
        <pre>{{.Text}}</pre>
        </div>

        <div class="post-actions">
        <button class="like-btn{{if eq .IsLiked 1}} active{{end}}">
            👍 {{.LikeCount}}
        </button>
        <button class="dislike-btn{{if eq .IsLiked -1}} active{{end}}">
            👎 {{.DislikeCount}}
        </button>
        </div>

        <!-- COMMENTS -->
        <div class="comments">
        <h4>Comments</h4>

        {{range .Comments}}
        <div class="comment" data-comment-id="{{.Id}}">
            <div class="delete-block">
            <div>
                <span class="comment-meta">Posted by {{.Username}} · {{.TimeAgo}}</span>
                <p class="comment-text">{{.Text}}</p>
            </div>
            {{if and $.IsLoggedIn (eq $.User.Id .UserId)}}
                <button class="comment-delete-btn btn small danger" type="submit">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512">
                    <path d="M136.7 5.9L128 32 32 32C14.3 32 0 46.3 0 64S14.3 96 32 96l384 0c17.7 0 32-14.3 32-32s-14.3-32-32-32l-96 0-8.7-26.1C306.9-7.2 294.7-16 280.9-16L167.1-16c-13.8 0-26 8.8-30.4 21.9zM416 144L32 144 53.1 467.1C54.7 492.4 75.7 512 101 512L347 512c25.3 0 46.3-19.6 47.9-44.9L416 144z"/>
                </svg>
                </button></div>
            {{end}}
            </div>
            <div class="comment-actions">
            <button class="comment-like-btn small{{if eq .IsLiked 1}} active{{end}}">
                👍 {{.LikeCount}}
            </button>
            <button class="comment-dislike-btn small{{if eq .IsLiked -1}} active{{end}}">
                👎 {{.DislikeCount}}
            </button>
            </div>
        </div>
        {{end}}

        <!-- Add Comment (REGISTERED USERS ONLY) -->
        <div class="add-comment registered-only">
            <form action="/api/comments/create" method="POST">
            <input type="hidden" name="postId" value="{{.Id}}" />
            <input name="text" type="text" placeholder="Write a comment..." maxlength="1000" required/>
            <button type="submit" class="btn small primary">Comment</button>
            </form>
        </div>
        </div>
    </article>    
`);