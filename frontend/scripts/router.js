export const routes = { // turn it to map !
    '/': {
        method: 'GET',
        page: () => import('../pages/feed.js'),
        auth: true,
    },

    '/feed': {
        method: 'GET',
        page: () => import('../pages/feed.js'), // duplicated
        auth: true,
    },

    '/post': { // with id !!!
        method: 'GET',
        page: () => import('../pages/post.js'),
        auth: true,
    },

    '/login': {
        method: 'GET',
        page: () => import('../pages/login.js'),
        auth: false,
    },

    '/register': {
        method: 'GET',
        page: () => import('../pages/register.js'),
        auth: false,
    },

    // 'error': () => import('../pages/error.js'),

    '/chat': {
        method: 'GET',
        page: () => import('../pages/chat.js'),
        auth: true,
    }
};

// code that need to implement 'navigate' method: (form actions)
// . logout form in Header
// . login form in LoginForm
// . post creation in PostCreationForm
// . comment creation in Post
// . register form in RegisterForm

async function guard(path) {
    const requiresAuth = routes[path].auth;
    const me = await isAuthenticated();
    if (requiresAuth && !me.authenticated) {
        path = '/login';
    } else if (!requiresAuth && me.authenticated) {
        path = '/';
    }
    history.pushState({}, '', path);
    return me.nickname;
}

import { isAuthenticated } from '../services/auth.js';
export const router = {
    async navigate(path) {
        // check if auth (do also for init())
        const nickname = await guard(path);

        await this.render({ nickname: nickname });

        // Load the page-specific script
        const scriptName = path.slice(1) || 'feed';
        await loadPageScript(scriptName);
    },

    async render(data = {}) {
        const requiresAuth = routes[location.pathname]
        const loader = routes[location.pathname].page; // page was not found

        if (!loader) {
            document.body.innerHTML = '<h1>404</h1>';
            return;
        }

        const page = await loader();
        document.querySelector('#app').innerHTML =
            await page.render(data);
    },

    async init() {
        window.addEventListener(
            'popstate',
            async () => {
                const nickname = await guard(location.pathname)
                this.render({ nickname: nickname })
            }
        );

        const nickname = await guard(location.pathname)
        await this.render({ nickname: nickname });
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