export const PostCreationForm = () => (`
    <section class="create-post">
        <details id="create-post-details">
            <summary>
                <h2>Create New Post</h2>
            </summary>
            <form id="create-post-form">
                <!-- Title -->
                <input type="text" name="title" placeholder="Post title" maxlength="255" required />

                <!-- Text -->
                <textarea
                    name="text"
                    placeholder="Write your post..."
                    maxlength="1000"
                    required
                ></textarea>

                <!-- Categories -->
                <div class="categories">
                    <label><input type="checkbox" name="categories" value="General" checked /> General</label>
                    <label><input type="checkbox" name="categories" value="Lifestyle" /> Lifestyle</label>
                    <label><input type="checkbox" name="categories" value="Health & Fitness" /> Health & Fitness</label>
                    <label><input type="checkbox" name="categories" value="Travel" /> Travel</label>
                    <label><input type="checkbox" name="categories" value="Food & Cooking" /> Food & Cooking</label>
                    <label><input type="checkbox" name="categories" value="Education" /> Education</label>
                    <label><input type="checkbox" name="categories" value="Business" /> Business</label>
                    <label><input type="checkbox" name="categories" value="Finance" /> Finance</label>
                    <label><input type="checkbox" name="categories" value="Entertainment" /> Entertainment</label>
                    <label><input type="checkbox" name="categories" value="Sports" /> Sports</label>
                    <label><input type="checkbox" name="categories" value="Personal Dev" /> Personal Dev</label>
                    <label><input type="checkbox" name="categories" value="Culture" /> Culture</label>
                    <label><input type="checkbox" name="categories" value="News" /> News</label>
                </div>

                <button type="submit" class="btn primary">Publish</button>
            </form>
        </details>
    </section>
`);

// <!-- Image upload -->
// <div class="image-upload">
//     <label class="upload-label">
//     <input 
//         type="file" 
//         name="image" 
//         accept="image/*"
//         onchange="previewImage(event)"
//     />
//     <!-- accept="image/*" only filter the file picker UI, but doesnt guarantee it is an image -->
//     <span class="upload-span" >+ Add Image (max 20MB)</span>
//     </label>

//     <div id="image-error" class="form-error" style="display: none; width: 100%; text-align: center;"></div>

//     <div class="post-image-container">
//     <div class="image-preview" id="imagePreview">
//         <span class="placeholder">No image selected</span>
//         <img id="previewImg" />
//     </div>
//     </div>
// </div>