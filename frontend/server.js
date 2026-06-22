import http from 'http';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

import { routes } from './scripts/router.js'; // router
// import { socket } from './core/websocket.js';

const PORT = 3000;

// Map file extensions to MIME types
const MIME_TYPES = {
    '.html': 'text/html',
    '.css': 'text/css',
    '.js': 'application/javascript',
    '.ico': 'image/x-icon',
};


const server = http.createServer(async (req, res) => {
    
    try {
        // const path = req.url.split('?')[0];

        // public and scripts
        // const SCRIPTS = path.resolve('./scripts');
        const decodedPath = decodeURIComponent(req.url).split('?')[0];
        const ext = path.extname(decodedPath);
        const mimeType = MIME_TYPES[ext] || 'application/octet-stream';

        if (path.resolve(decodedPath) === '/favicon.ico' ||
            path.resolve(decodedPath).startsWith('/styles/') ||
              path.resolve(decodedPath).startsWith('/api/') ||
            path.resolve(decodedPath).startsWith('/scripts/') ||
            path.resolve(decodedPath).startsWith('/components/') ||
            path.resolve(decodedPath).startsWith('/pages/')) {
            fs.readFile(`.${decodedPath}`, (err, data) => {
                if (err) throw err;

                res.writeHead(200, { 'Content-Type': mimeType });
                res.end(data);
            });
            return;
        }


        let route = routes[decodedPath];
        if (!route || route.method !== req.method) route = null;
        
        // routes.find(
        //     r => r.method === req.method &&
        //         r.path === req.url.split('?')[0]
        // );

        if (!route) { // need to throw error instead !
            res.writeHead(404, {
                'Content-Type': 'application/json'
            });

            res.end(JSON.stringify({
                error: 'Not Found'
            }));

            return;
        }

        // await route.handler(req, res);
        // always serve index.html
        // const ext = path.extname(filePath);
        // const mimeType = MIME_TYPES[ext] || 'application/octet-stream';

        // To run the server with node from anywhere. (needs to handle all files)
        // const __filename = fileURLToPath(import.meta.url);
        // const __dirname = path.dirname(__filename);
        // const indexPath = path.join(__dirname, 'index.html');
        fs.readFile('./index.html', (err, data) => { // indexPath instead
            if (err) throw err;

            res.writeHead(200, { 'Content-Type': 'text/html' }); // mimeType });
            res.end(data);
        });
    } catch (err) {
        console.error(err);

        if (err.code === 'ENOENT') {
            res.writeHead(404, { 'Content-Type': 'text/plain' });
            res.end(`404 — File not found: ${urlPath}`);
        } else {
            res.writeHead(500, { 'Content-Type': 'text/plain' });
            res.end(`500 — Internal server error: ${err.message}`);
        }
    }

    // ==========================================

    // await router.init();

    // if (localStorage.getItem('token')) {
    //     socket.connect();
    // }

    // ===========================================

    // await handle(req, res);
});

server.listen(PORT, () => {
    console.log(`Http server running at http://localhost:${PORT}`);
});