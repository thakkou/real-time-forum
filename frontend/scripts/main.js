import { router } from './router.js';

const app = document.getElementById("app");

// Store loaded scripts to avoid duplicates
const loadedScripts = new Map();


// ========================
// GLOBAL FUNCTIONS
// ========================
//later 
window.env = {
    serverUri: "http://localhost:8080/api",
    wsUri:"ws://localhost:8080/ws"
};

window.navigate = router.navigate.bind(router); // navigate

router.init();




