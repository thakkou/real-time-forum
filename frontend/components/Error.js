export const Error = (code, message) => (`
    <div class="card error-card">
        <h1>${code}</h1>
        <p>${message}</p>
        <a href="/">← Back to Home</a>
    </div>
`);