import http from 'http';
import fs from 'fs';
import path from 'path';

const PORT = 3000;

const MIME_TYPES = {
    '.html': 'text/html',
    '.css': 'text/css',
    '.js': 'application/javascript',
    '.ico': 'image/x-icon',
};

const server = http.createServer((req, res) => {
    try {
        const decodedPath = decodeURIComponent(req.url).split('?')[0];

        const ext = path.extname(decodedPath);
        const mimeType = MIME_TYPES[ext] || 'application/octet-stream';

        // =========================
        // 1. STATIC FILES
        // =========================

const filePath = path.join(process.cwd(), decodedPath);

const isFileRequest =
    decodedPath.includes('.') || decodedPath === '/favicon.ico';

if (isFileRequest) {
    fs.readFile(filePath, (err, data) => {
        if (err) {
            res.writeHead(404, { 'Content-Type': 'text/plain' });
            res.end('File not found');
            return;
        }

        const ext = path.extname(filePath);

        const mimeTypes = {
            '.html': 'text/html',
            '.css': 'text/css',
            '.js': 'application/javascript',
            '.ico': 'image/x-icon',
        };

        res.writeHead(200, {
            'Content-Type': mimeTypes[ext] || 'application/octet-stream'
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

        res.writeHead(500, { 'Content-Type': 'text/plain' });
        res.end(`500 — ${err.message}`);
    }
});

server.listen(PORT, () => {
    console.log(`Http server running at http://localhost:${PORT}`);
});