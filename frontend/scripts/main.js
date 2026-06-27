
import { router } from '../services/router.js';

const app = document.getElementById("app");
window.env = {
    serverUri: "http://localhost:8080/api",
    wsUri:"ws://localhost:8080/ws"
};

// Store loaded scripts to avoid duplicates
const loadedScripts = new Map();




window.navigate = router.navigate.bind(router); // navigate

router.init();




