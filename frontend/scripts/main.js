
import { router } from '../services/router.js';

const app = document.getElementById("app");

window.navigate = router.navigate.bind(router); // navigate

router.init();
