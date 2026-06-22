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
        await loadPageScript(path.slice(1));
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
        await loadPageScript('feed');
        // loaded first time, must be :
        // 1. chnaged depending on app state (first page)
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