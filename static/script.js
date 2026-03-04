async function reactToPost(postId, endpoint) {
    const url = `/api/posts/${postId}/${endpoint}`;
    try {
        const response = await fetch(url, {
            method: 'POST'
        });

        if (!response.ok) {
            // Handle server errors (e.g., status 400, 500)
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
    } catch (error) {
        console.error('Error sending vote:', error);
    }
}

async function reactToComment(commentId, endpoint) {
    const url = `/api/comments/${commentId}/${endpoint}`;
    try {
        const response = await fetch(url, {
            method: 'POST'
        });

        if (!response.ok) {
            // Handle server errors (e.g., status 400, 500)
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
    } catch (error) {
        console.error('Error sending vote:', error);
    }
}

// Function to handle button clicks
function handlePostReactionsClick(event) {
    const button = event.target;
    const postContainer = button.closest('.post');
    const postId = postContainer.getAttribute('data-post-id');
    
    if (!postId) return;

    let endpoint;
    if (button.classList.contains('like-btn')) {
        endpoint = 'like';
    } else if (button.classList.contains('dislike-btn')) {
        endpoint = 'dislike';
    } else {
        return;
    }

    reactToPost(postId, endpoint)
        .then(data => {
            // Update the UI with the new count and selection
        })
        .catch(err => {});
}

// Function to handle button clicks
function handleCommentReactionsClick(event) {
    const button = event.target;
    const commentContainer = button.closest('.comment');
    const commentId = commentContainer.getAttribute('data-comment-id');

    if (!commentId) return;

    let endpoint;
    if (button.classList.contains('comment-like-btn')) {
        endpoint = 'like';
    } else if (button.classList.contains('comment-dislike-btn')) {
        endpoint = 'dislike';
    } else {
        return;
    }

    reactToComment(commentId, endpoint)
        .then(data => {
            // Update the UI with the new count and selection
        })
        .catch(err => {});
}

// Add event listeners to buttons
document.addEventListener('DOMContentLoaded', () => {
    // post reaction buttons
    document.querySelectorAll('.like-btn').forEach(button => {
        button.addEventListener('click', handlePostReactionsClick);
    });
    document.querySelectorAll('.dislike-btn').forEach(button => {
        button.addEventListener('click', handlePostReactionsClick);
    });

    // comment reaction buttons
    document.querySelectorAll('.comment-like-btn').forEach(button => {
        button.addEventListener('click', handleCommentReactionsClick);
    });
    document.querySelectorAll('.comment-dislike-btn').forEach(button => {
        button.addEventListener('click', handleCommentReactionsClick);
    });
});
