export const Toast = (message, type = 'success') => (`
    <div class="toast ${type}">
        ${message}
    </div>
`);