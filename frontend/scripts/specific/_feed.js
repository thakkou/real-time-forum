import { CommentResolver,CreatComment } from "../../api/comments.js";
import { getPosts } from "../../api/posts.js";
function setupHomePage() {
    //get all posts 10 by 10

    document.querySelectorAll('.filter-btn').forEach((button) => {
    button.addEventListener('click', function() {
        this.classList.toggle('active');
        const hiddenInputs = {
        'my-creat-postes': document.getElementById('input-my-creat-postes'),
        'my-liked-post': document.getElementById('input-my-liked-post'),
        };
        Object.keys(hiddenInputs).forEach((name) => {
        if (hiddenInputs[name]) hiddenInputs[name].value = '';
        });
        const activeButtons = Array.from(document.querySelectorAll('.filter-btn.active'));
        activeButtons.forEach((btn) => {
        if (hiddenInputs[btn.name]) {
            hiddenInputs[btn.name].value = 'true';
        }
        });
    });
    });

}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', setupHomePage);
} else {
    setupHomePage();
}

//this will be clean up  ASAP

//     function reactToPost(postId, endpoint) {
//     const url = `/api/posts/${postId}/${endpoint}`;
//     return fetch(url, { method: 'POST' })
//         .then(response => {
//         if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);
//         window.location.reload();
//         })
//         .catch(error => console.error('Error sending reaction:', error));
//     }

//     function handlePostReactionsClick(event) {
//     const button = event.target;
//     const postContainer = button.closest('.post');
//     const postId = postContainer?.getAttribute('data-post-id');
//     if (!postId) return;
    
//     let endpoint;
//     if (button.classList.contains('like-btn')) endpoint = 'like';
//     else if (button.classList.contains('dislike-btn')) endpoint = 'dislike';
//     else return;
    
//     reactToPost(postId, endpoint);
//     }

//     function reactToComment(commentId, endpoint) {
//     const url = `/api/comments/${commentId}/${endpoint}`;
//     return fetch(url, { method: 'POST' })
//         .then(response => {
//         if (!response.ok) throw new Error(`HTTP error! Status: ${response.status}`);
//         window.location.reload();
//         })
//         .catch(error => console.error('Error sending reaction:', error));
//     }

//     function handleCommentReactionsClick(event) {
//     const button = event.target;
//     const commentContainer = button.closest('.comment');
//     const commentId = commentContainer?.getAttribute('data-comment-id');
//     if (!commentId) return;
    
//     let endpoint;
//     if (button.classList.contains('comment-like-btn')) endpoint = 'like';
//     else if (button.classList.contains('comment-dislike-btn')) endpoint = 'dislike';
//     else return;
    
//     reactToComment(commentId, endpoint);
//     }

//     document.querySelectorAll('.like-btn, .dislike-btn').forEach(button => {
//     button.removeEventListener('click', handlePostReactionsClick);
//     button.addEventListener('click', handlePostReactionsClick);
//     });
    
//     document.querySelectorAll('.comment-like-btn, .comment-dislike-btn').forEach(button => {
//     button.removeEventListener('click', handleCommentReactionsClick);
//     button.addEventListener('click', handleCommentReactionsClick);
//     });
// }