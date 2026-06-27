import { ws } from './websocket.js';
import { showToast } from './toast.js';
import { saveOnlineUsers ,getOnlineUsers} from '../scripts/main.js';

export const routes = { // turn it to map !
    '/': {
        method: 'GET',
        name:"home",
        page: () => import('../pages/feed.js'),
        auth: true,
    },

    '/feed': {
        method: 'GET',
        name:"feeds",
        page: () => import('../pages/feed.js'), // duplicated
        auth: true,
    },

    '/post/:id': { // with id !!!
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
    const matched = matchRoute(path);

    // if (!matched) {
    //     history.pushState({}, '', '/');
    //     return null;
    // }

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
        const matched = matchRoute(location.pathname);

        if (!matched) {
            document.body.innerHTML = '<h1>404</h1>';
            return;
        }

        const loader = matched.route.page;

        if (!loader) {
            document.body.innerHTML = '<h1>404</h1>';
            return;
        }

        const page = await loader();

        document.querySelector('#app').innerHTML =
            await page.render({
                ...data,
                params: matched.params
            });
    },

    async init() {

        window.addEventListener('popstate', async () => {
            const nickname = await guard(location.pathname);
            this.render({ nickname });
        });
    

        const nickname = await guard(location.pathname);
        if (nickname) {
            ws.connect();

            const onlineUsers = getOnlineUsers();

            ws.on("init", (data) => {
                console.log("init users:", data);
                data.forEach(id => onlineUsers.add(id));
                saveOnlineUsers(onlineUsers);
                this.renderOnlineUsers?.([...onlineUsers]);
            });

            ws.on("client_connect", (userId) => {
                console.log("user connected:", userId);
                onlineUsers.add(userId);
                saveOnlineUsers(onlineUsers);
                this.renderOnlineUsers?.([...onlineUsers]);
            });

            ws.on("client_disconnect", (userId) => {
                console.log("user disconnected:", userId);
                onlineUsers.delete(userId);
                saveOnlineUsers(onlineUsers);
                this.renderOnlineUsers?.([...onlineUsers]);
            });

            ws.on("new_post",(data)=>{
                console.log(data)
            });

            ws.on("new_message", (data) => {
                console.log("new message:", data);
                showToast(data.text, "success");
            });

            ws.on("typing", (data) => {
                console.log("someone is typing:", data.userId);
            });
        }

        // the page is fully rendered first, then the specific scripts are loded after !
        await this.render({ nickname });
        // Load the page-specific script
        // await loadPageScript(window.location.pathname.slice(1)); // feed default

        const scriptName = location.pathname.split('/')[1] || 'feed';
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
  feed: () => import('../scripts/_feed.js'),
  login: () => import('../scripts/_login.js'),
  register: () => import('../scripts/_register.js'),
  chat: () => import('../scripts/_chat.js'), // used temporarely
  post: () => import('../scripts/_post.js'),
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

async function loadPageScript(pageName) {
//   if (window.currentPageScript && typeof window.currentPageScript.cleanup === 'function') {
//     window.currentPageScript.cleanup();
//   } // ?!
  
  if (pageScripts[pageName]) {
    const script = await pageScripts[pageName]();
    await script.setup();
    window.currentPage = pageName; // need to be done before !!!
  }
}