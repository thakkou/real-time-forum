// const routes = {
//     '/': () => import('../pages/login.js'),
//     '/dashboard': () => import('../pages/dashboard.js'),
//     '/chat': () => import('../pages/chat.js')
// };

export const routes = { // turn it to map !
    '/': {
        method: 'GET',
        page: () => import('../pages/login.js')
    },

    '/dashboard': {
        method: 'GET',
        page: () => import('../pages/dashboard.js'),
        auth: true
    },

    '/chat': {
        method: 'GET',
        page: () => import('../pages/chat.js'),
        auth: true
    }
};

export const router = {
    async navigate(path) {
        history.pushState({}, '', path);
        await this.render();
    },

    async render() {
        const loader = routes[location.pathname];

        if (!loader) {
            document.body.innerHTML = '<h1>404</h1>';
            return;
        }

        const page = await loader();

        document.querySelector('#app').innerHTML =
            await page.render();
    },

    async init() {
        window.addEventListener(
            'popstate',
            () => this.render()
        );

        await this.render();
    }
};
