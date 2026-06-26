import { login } from "../../api/auth.js";
import { showToast } from "../../services/toast.js";
import { router } from "../router.js";

// login not working and refreshes at the first time when redirected from register, but works after !!!
const setupLoginPage = () => {
	const form = document.querySelector("form");
	const errorBox = document.getElementById("login-error");
	const btn = document.getElementById("login-btn");
	if (!form) return;

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
			router.navigate('/');

			// Store user globally
			window.user = resp.data; // ?!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		} catch (err) {
			errorBox.style.display = "block";
			errorBox.textContent = err.message || "Login failed"; // err.message for debugging
		} finally {
            btn.disabled = false;
            btn.textContent = "Login";
        }
	});
};

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', setupLoginPage);
} else {
    setupLoginPage();
}