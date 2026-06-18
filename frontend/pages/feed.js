export const Home = () => (`
  ${header}

  <div class="container">
    <!-- FILTER SECTION (desktop sidebar) -->
    <aside class="sidebar">
      <h3>Filter Posts</h3>


      <form method="GET" action="/" id="filter-form">
        <!-- My Posts / Liked Posts Buttons -->
        {{if .IsLoggedIn}}
        <div class="filter-group registered-only">
          <button type="button" name="my-creat-postes" class="filter-btn">
            My Created Posts
          </button>

          <button type="button" name="my-liked-post" class="filter-btn">
            Posts I Liked
          </button>
        </div>
        {{end}}

        <!-- Categories -->
        <div class="filter-group">
          <h4>By Category</h4>
          <input
            type="hidden"
            name="my-creat-postes"
            id="input-my-creat-postes"
          />
          <input
            type="hidden"
            name="my-liked-post"
            id="input-my-liked-post"
          />

          <label><input type="checkbox" name="categories" value="General" /> General</label>
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

          <button type="submit" class="filter-btn">Apply Filters</button>
          <a href="/" class="clear-btn">Clear Filters</a>
        </div>
      </form>
    </aside>

    <!-- MAIN CONTENT -->
    <main class="content">

      ${postCreationForm}

      <!-- MOBILE FILTER (visible only on small screens, below create-post) -->
      <aside class="sidebar-mobile">
        <details>
          <summary>
            <h3>Filter Posts</h3>
          </summary>
          <form method="GET" action="/">
            <div class="filter-group registered-only">
              <button type="button" name="my-creat-postes" class="filter-btn">My Created Posts</button>
              <button type="button" name="my-liked-post" class="filter-btn">Posts I Liked</button>
            </div>
            <div class="filter-group">
              <h4>By Category</h4>
              <input type="hidden" name="my-creat-postes" id="input-my-creat-postes-mobile" />
              <input type="hidden" name="my-liked-post" id="input-my-liked-post-mobile" />
              <div class="mobile-categories-grid">
                <label><input type="checkbox" name="categories" value="General" /> General</label>
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
              <button type="submit" class="filter-btn">Apply Filters</button>
              <a href="/" class="clear-btn">Clear Filters</a>
            </div>
          </form>
        </details>
      </aside>

      <!-- POSTS LIST -->
      <section class="posts">
        {{if gt (len .Posts) 0}}

        {{range .Posts}}
        ${post}
        {{end}}
      {{else}}
        <h1 class="no-post">No posts right now</h1>
      {{end}}
      </section>
    </main>
  </div>
`);

//   <head>
//     <meta charset="UTF-8" />
//     <meta name="viewport" content="width=device-width, initial-scale=1.0" />
//     <link rel="icon" type="image/x-icon" href="/assets/favicon.ico">
//     <title>01 Forum</title>
//     <link rel="stylesheet" href="/assets/style.css" />
//     <script src="/assets/script.js" defer></script>
//     <!-- Modern best practice: 
//      1. Script downloads in parallel
//      2. Executes after HTML is fully parsed
//      3. Does NOT block rendering -->
//   </head>