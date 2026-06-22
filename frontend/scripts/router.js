export const routes = { // turn it to map !
    '/': {
        method: 'GET',
        page: () => import('../pages/feed.js')
    },

    '/feed': {
        method: 'GET',
        page: () => import('../pages/feed.js'), // duplicated
        // auth: true
    },

    '/post': { // with id !!!
        method: 'GET',
        page: () => import('../pages/post.js'),
    },

    '/login': {
        method: 'GET',
        page: () => import('../pages/login.js')
    },

    '/register': {
        method: 'GET',
        page: () => import('../pages/register.js')
    },

    // 'error': () => import('../pages/error.js'),

    '/chat': {
        method: 'GET',
        page: () => import('../pages/chat.js'),
        // auth: true
    }
};

// code that need to implement 'navigate' method: (form actions)
// . logout form in Header
// . login form in LoginForm
// . post creation in PostCreationForm
// . comment creation in Post
// . register form in RegisterForm

export const router = {
    async navigate(path) {
        history.pushState({}, '', path);
        await this.render();

        // Load the page-specific script
        const scriptName = path.slice(1) || 'feed';
        await loadPageScript(scriptName);
    },

    async render() {
        const loader = routes[location.pathname].page; // page was not found

        if (!loader) {
            document.body.innerHTML = '<h1>404</h1>';
            return;
        }

        const page = await loader();
        // console.log(page)

        document.querySelector('#app').innerHTML =
            await page.render();
    },

    async init() {
        window.addEventListener(
            'popstate',
            () => this.render()
        );

        await this.render();
        // Load the page-specific script
        // await loadPageScript(window.location.pathname.slice(1)); // feed default
        const scriptName = window.location.pathname.slice(1) || 'feed';
        await loadPageScript(scriptName);
        // loaded first time, must be :
        // 1. chnaged depending on app state (first page) x
        // 2. not loaded if already exists (same in navigate) -> is default behavior maybe !?
    }
};

// ========================
// PAGE-SPECIFIC SCRIPTS
// ========================
const pageScripts = {
  feed: () => import('./specific/_feed.js'),
  login: () => import('./specific/_login.js'),
  register: () => import('./specific/_register.js'),
  chat: () => import('./specific/_chat.js'),
  // single post
};

async function loadPageScript(pageName) {
  if (window.currentPageScript && typeof window.currentPageScript.cleanup === 'function') {
    window.currentPageScript.cleanup();
  } // ?!
  
  if (pageScripts[pageName]) {
    const script = await pageScripts[pageName]();
    window.currentPage = pageName;
  }
}