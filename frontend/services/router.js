import { isAuthenticated } from '../services/auth.js';
import { ws } from './websocket.js';
import { showToast } from './toast.js';
import { reRender, reRenderMessages, updateTheConv } from '../scripts/_chat.js';
import { Header } from '../components/Header.js';
import { handleIncomingTypingEvent } from '../scripts/_chat.js';
import { logout } from '../api/auth.js';
export const onlineUsers = new Set()
let me = null;

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

async function guard(path, pushToHistory = true) {
    console.log("qguards",path)
    const matched = matchRoute(path);

    if (!matched) {
        await router.error(404, "page not found");
        return null;
    }

    const requiresAuth = matched.route.auth;
    me = await isAuthenticated();
    
    let targetPath = path;

    // Check auth overrides
    if (requiresAuth && !me.authenticated) {
        targetPath = "/login";
    } else if (!requiresAuth && me.authenticated) {
        targetPath = "/";
    }

    // If the path was hijacked by auth guard rules
    if (targetPath !== path) {
        // Force replace the history state so the address bar corrects itself 
        // even if we arrived here via a back-button action
        history.replaceState({}, "", targetPath);
    } else if (pushToHistory) {
        // Regular forward navigation
        history.pushState({}, "", targetPath);
    }
console.log(me)
    return me.nickname;
}
async function handleLogout() {
  try {
    await logout();
    localStorage.clear();
    window.location.href = "/login";
  } catch (err) {
    console.error("Logout failed:", err);
  }
}


function setupHeader(nickname) {
    const header = document.getElementById("header");

    if (!header) return;

    if (nickname) {
        header.innerHTML = Header(nickname);

        const logoutBtn = header.querySelector("#logout-btn");

        if (logoutBtn) {
            logoutBtn.addEventListener("click", handleLogout);
        }

    } else {
        header.innerHTML = "";
    }
}
function extractPath(path) {
    const parts = path.split("/").filter(Boolean);
    // no path => default page
    if (parts.length === 0) {
        return "feed";
    }

    // only /post/:id is allowed
    if (parts.length === 2  &&  parts[0] === "post") {
        return "post";
    }
    if(parts.length>=2){
        return null
    }

    return parts[0]
}

export const router = {

  async error(status, message) {
    const errorPage = await import('../pages/error.js');

    document.querySelector('#app').innerHTML =
        await errorPage.render({
            status,
            message,
        });
},
  async navigate(path) {

    const nickname = await guard(path);

    if (nickname === null) {
        return;
    }

    setupHeader(nickname);

    await this.render({ nickname });

    const scriptName = extractPath(path);
console.log("start navigate to ",scriptName,"from path")
    await loadPageScript(scriptName);
},

    async render(data = {}) {
        const matched = matchRoute(location.pathname);

        if (!matched) {
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
    // 1. Tell guard NOT to call pushState again
    const nickname = await guard(location.pathname, false); 
    
    if (nickname === null) return;

    // 2. Refresh the layout
    setupHeader(nickname);
    await this.render({ nickname });

    // 3. CRITICAL: Reload the scripts so event listeners bind!
    const scriptName = extractPath(location.pathname);
    await loadPageScript(scriptName);
});
    

        const nickname = await guard(location.pathname);
        setupHeader(nickname);

        if (nickname) {
            console.log("user connect and set the header ")
            ws.connect();


            ws.on("init", (data) => {
                console.log("init users:", data);
                data.forEach(id => onlineUsers.add(id));
           
            });

            ws.on("client_connect", (userId) => {
                console.log("user connected:", userId);
                onlineUsers.add(userId);
                reRender("connect",userId)
            });

            ws.on("client_disconnect", (userId) => {
                console.log("user disconnected:", userId);
                onlineUsers.delete(userId);
    reRender("disconnect",userId)
            });


            ws.on("new_post",(data)=>{
            });

            ws.on("new_message", (data) => {
                const isMe = data.isMine
                const isNew=data.isNewConversation
               

                    updateTheConv(data,isNew)
                

                 if(!isMe){
                    console.log("append me ")

              showToast(data.text, "success");

             reRenderMessages(data,false)

                 }else{
                    console.log("append him ")

                    reRenderMessages(data,true)
                 }


 });

          ws.on("typing:start", (data) => {
  console.log("someone is start typing:", data);

  handleIncomingTypingEvent({
    ...data,
    is_typing: true
  });
});

ws.on("typing:stop", (data) => {
  console.log("someone is stop typing:", data);

  handleIncomingTypingEvent({
    ...data,
    is_typing: false
  });
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