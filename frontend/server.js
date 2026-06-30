import http from 'http';
import fs from 'fs';
import path from 'path';

const PORT = 3000;

// Map file extensions to MIME types
const MIME_TYPES = {
    '.html': 'text/html',
    '.css': 'text/css',
    '.js': 'application/javascript',
    '.ico': 'image/x-icon',
};

const server = http.createServer((req, res) => { // async
    try {
        // const path = req.url.split('?')[0];

        // public and scripts
        // const SCRIPTS = path.resolve('./scripts');
        const decodedPath = decodeURIComponent(req.url).split('?')[0];

        const ext = path.extname(decodedPath);
        const mimeType = MIME_TYPES[ext] || 'application/octet-stream';

        // =========================
        // 1. STATIC FILES
        // =========================

        const filePath = path.join(process.cwd(), decodedPath);

        const isFileRequest = // what this variable mean ? what else can be ?!
            decodedPath.includes('.') || decodedPath === '/favicon.ico';

        const isAllowed = path.resolve(decodedPath).startsWith('/styles/') ||
            path.resolve(decodedPath).startsWith('/api/') ||
            path.resolve(decodedPath).startsWith('/services/') ||
            path.resolve(decodedPath).startsWith('/scripts/') ||
            path.resolve(decodedPath).startsWith('/components/') ||
            path.resolve(decodedPath).startsWith('/pages/');

        // let route = routes[decodedPath];
        // if (!route || route.method !== req.method) route = null;

        // routes.find(
        //     r => r.method === req.method &&
        //         r.path === req.url.split('?')[0]
        // );

        if (isFileRequest && isAllowed) {
            fs.readFile(filePath, (err, data) => {
                if (err) {
                    res.writeHead(404, { 'Content-Type': 'text/plain' });
                    res.end('File not found');
                    return;
                }

                const ext = path.extname(filePath);

                res.writeHead(200, {
                    'Content-Type': MIME_TYPES[ext] || 'application/octet-stream'
                });

                res.end(data);
            });

            return;
        }

        // =========================
        // 2. SPA FALLBACK (IMPORTANT)
        // =========================
        fs.readFile('./index.html', (err, data) => {
            if (err) {
                res.writeHead(500, { 'Content-Type': 'text/plain' });
                res.end('500 - Server Error');
                return;
            }

            res.writeHead(200, { 'Content-Type': 'text/html' });
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

        // res.writeHead(500, { 'Content-Type': 'text/plain' });
        // res.end(`500 — ${err.message}`);
    }
});

server.listen(PORT, () => {
    console.log(`Http server running at http://localhost:${PORT}`);
});