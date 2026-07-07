import { router } from '../services/router.js';

// const app = document.getElementById("app");

window.navigate = router.navigate.bind(router);
window.error = router.error.bind(router);

router.init();
