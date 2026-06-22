function setupLoginPage() {
    console.log("Login page loaded");
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', setupLoginPage);
} else {
    setupLoginPage();
}