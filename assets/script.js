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
        console.error('Error sending reaction:', error);
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
        console.error('Error sending reaction:', error);
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
            window.location.reload();
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
            window.location.reload();
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

// ================== Filters event listener ======================

document.querySelectorAll('.filter-btn').forEach((button) => {
    button.addEventListener('click', function () {
        this.classList.toggle('active');

        const hiddenInputs = {
            'my-creat-postes': document.getElementById('input-my-creat-postes'),
            'my-liked-post': document.getElementById('input-my-liked-post'),
        };

        Object.keys(hiddenInputs).forEach((name) => {
            if (hiddenInputs[name]) hiddenInputs[name].value = '';
        });

        const activeButtons = Array.from(
            document.querySelectorAll('.filter-btn.active'),
        );
        activeButtons.forEach((btn) => {
            if (hiddenInputs[btn.name]) {
            hiddenInputs[btn.name].value = 'true';
            }
        });
    });
});

// Checks file size & and previews image
function previewImage(event) {
  const file = event.target.files[0];
  const preview = document.getElementById("imagePreview");
  const img = document.getElementById("previewImg");

  if (!file) return;

  // Size check (20MB)
  if (file.size > 20 * 1024 * 1024) {
    event.target.value = "";
    document.getElementById("image-error").style.display = "block";
    document.getElementById("image-error").textContent = "Max file size: 20Mb";
    return;
  } else {
    document.getElementById("image-error").style.display = "none";
  }

  // check Mime type (images only)
  // All standard image MIME types are under the image/ umbrella.
  // Caveats: some browsers (or drag & drop cases) may give empty file.type
  if (!file.type || !file.type.startsWith("image/")) {
    event.target.value = "";
    document.getElementById("image-error").style.display = "block";
    document.getElementById("image-error").textContent = "Image file types only";
    return;
  } else {
    document.getElementById("image-error").style.display = "none";
  }
  
  // Preview
  const reader = new FileReader();
  reader.onload = function(e) {
    img.src = e.target.result;
    preview.classList.add("active");
  };

  reader.readAsDataURL(file);
}