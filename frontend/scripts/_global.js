import { logout } from "../api/auth.js";
// import { router } from "../services/router.js";
import { ws } from "../services/websocket.js";

/* ======================
   LOGOUT
====================== */

async function handleLogout() {
  try {
    ws.logout();
    await logout();
    localStorage.clear(); // why ?!
    // router.navigate('/login');
    // not used because we need to remove all previous scripts !
    // + some event listeners still work after that !
    window.location.href = "/login";
  } catch (err) {
    console.error("Logout failed:", err);
  }
}

export function setup() {
    document.getElementById('logout-btn').addEventListener("click", handleLogout);
}