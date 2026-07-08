import { isAuthenticated } from '../services/auth.js';
import { ws } from './websocket.js';

export const routes = {
    '/': {
        method: 'GET',
        name:"home",
        page: () => import('../pages/feed.js'),
        auth: true,
    },

    '/feed': {
        method: 'GET',
        name:"feeds",
        page: () => import('../pages/feed.js'),
        auth: true,
    },

    '/post/:id': {
        method: 'GET',
        name:"post detaills",
        page: () => import('../pages/post.js'),
        auth: true,
    },

    '/login': {
        method: 'GET',
        name:"login",
        page: () => import('../pages/login.js'),
        auth: false,
    },

    '/register': {
        method: 'GET',
        name:"register",
        page: () => import('../pages/register.js'),
        auth: false,
    },

    '/chat': {
        method: 'GET',
        page: () => import('../pages/chat.js'),
        auth: true,
    }
};

// ========================
// PAGE-SPECIFIC SCRIPTS
// ========================
const pageScripts = {
  feed: () => import('/scripts/_feed.js'),
  login: () => import('/scripts/_login.js'),
  register: () => import('/scripts/_register.js'),
  chat: () => import('/scripts/_chat.js'),
  post: () => import('/scripts/_post.js'),
};

function matchRoute(path) {
    for (const route in routes) {
        const paramNames = [];

        const regexPath = route.replace(/:([^/]+)/g, (_, key) => {
            paramNames.push(key);
            return '([^/]+)';
        });

        const regex = new RegExp(`^${regexPath}$`);
        const match = path.match(regex);

        if (match) {
            const params = {};
            paramNames.forEach((name, i) => {
                params[name] = match[i + 1];
            });

            return { route: routes[route], params };
        }
    }

    return null;
}

async function guard(path) {
    console.log(path)
    const matched = matchRoute(path);

    if (!matched) {
        return null;
    }

    const requiresAuth = matched.route.auth;
    const me = await isAuthenticated();
    if (requiresAuth && !me.authenticated) {
        path = '/login';
    } else if (!requiresAuth && me.authenticated) {
        path = '/';
    }
    history.pushState({}, '', path);
    return me.nickname;
}


export const router = {
    async error(status, message) {
        const errorPage = await import('../pages/error.js');
        document.querySelector('#app').innerHTML =
            await errorPage.render({
                status: status,
                message: message,
            });
        return;
    },

    async navigate(path) {
        const nickname = await guard(path);
        await this.render({ nickname });
    },

    async render(data = {}) {
        const matched = matchRoute(location.pathname);
        if (!matched) {
            this.error(404, "Page Not Found");
            return;
        }

        const loader = matched.route.page;
        const page = await loader();
        document.querySelector('#app').innerHTML =
            await page.render({
                ...data,
                params: matched.params
            }); // should await error 404 for post not found !

        const scriptName = location.pathname.split('/')[1] || 'feed';
        await this.loadPageScript(scriptName);
    },

    async init() {
        window.addEventListener('popstate', async () => {
            this.navigate(location.pathname);
        });
        this.navigate(location.pathname);

        if (location.pathname !== '/login' && location.pathname !== '/register')
            ws.connect();
    },

    async loadPageScript(pageName) {
        if (pageName !== 'login' && pageName !== 'register') {
            const globalScript = await import('../scripts/_global.js');
            await globalScript.setup();
        }

        if (pageScripts[pageName]) {
            const script = await pageScripts[pageName]();
            await script.setup();
            window.currentPage = pageName; // need to be done before !!!
        }
    }
};