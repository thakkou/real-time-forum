import { login } from "../api/auth.js";
import { showToast } from "../services/toast.js";
import { router } from "../services/router.js";
import { ws } from "../services/websocket.js";

// login not working and refreshes at the first time when redirected from register, but works after !!!
export function setup() {
	const form = document.querySelector("form");
	const errorBox = document.getElementById("login-error");
	const btn = document.getElementById("login-btn");
	if (!form) return;
	// Fill credentials for test users
	const testButtons = document.querySelectorAll(".test-user-btn");
testButtons.forEach((btn) => {
	btn.addEventListener("click", () => {
		form.identifier.value = btn.dataset.user;
		form.password.value = "password123";
		form.identifier.focus();
	});
});

	form.addEventListener("submit", async (e) => {
		e.preventDefault();

		const identifier = form.identifier.value;
		const password = form.password.value;

		btn.disabled = true;
        btn.textContent = "Logging in...";

		try {
			const resp = await login({ identifier, password });

			console.log("login successful");
			showToast('Login successful!', 'success');
			// Replace the login entry in history so back button doesn't return to login
			history.replaceState({}, "", "/");
			ws.connect()
			await router.navigate("/");
			
			// Store user globally
			window.user = resp.data; // ?!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		} catch (err) {
			errorBox.style.display = "block";
			errorBox.textContent = err.message || "Login failed"; // err.message for debugging
			showToast(err.message || "Login failed", "error");
		} finally {
            btn.disabled = false;
            btn.textContent = "Login";
        }
	});
};