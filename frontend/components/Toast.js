export const Toast = (message, type = 'success') => {
    const toast = document.createElement('div');
    toast.className = "toast ${type}";
    toast.innerText = message;
    return toast;
}

// (`
//     <div class="toast ${type}">
//         ${message}
//     </div>
// `);