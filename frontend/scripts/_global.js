import { logout } from "../api/auth.js";

/* ======================
   LOGOUT
====================== */

async function handleLogout() {
  try {
    await logout();
    localStorage.clear();
    window.location.href = "/login";
  } catch (err) {
    console.error("Logout failed:", err);
  }
}

export function setup() {
    document.getElementById('logout-btn').addEventListener("click", handleLogout);
}