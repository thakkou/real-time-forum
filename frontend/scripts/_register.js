import { register } from "../api/auth.js";
import { showToast } from "../services/toast.js";
import { router } from "../services/router.js";

export function setup() {
    const form = document.querySelector("form");
    if (!form) return;

    const btn = document.getElementById("register-btn");
    const errorBox = document.getElementById("register-error");

    form.addEventListener("submit", async (e) => {
        e.preventDefault();

        const formData = e.target;

        const nickname = formData.nickname.value.trim();
        const first_name = formData.firstname.value.trim();
        const last_name = formData.lastname.value.trim();
        const age = parseInt(formData.age.value, 10);
        const gender = formData.querySelector('input[name="gender"]:checked');
        const email = formData.email.value.trim();
        const password = formData.password.value;
        const confirm_password = formData.confirm_password.value;

        errorBox.style.display = "none";
        errorBox.textContent = "";

        // Validation
        if (!nickname || nickname.length < 2) {
            errorBox.textContent = "Nickname must be at least 2 characters";
            errorBox.style.display = "block";
            return;
        }

        if (!first_name) {
            errorBox.textContent = "First name is required";
            errorBox.style.display = "block";
            return;
        }

        if (!last_name) {
            errorBox.textContent = "Last name is required";
            errorBox.style.display = "block";
            return;
        }

        if (!age || isNaN(age) || age < 13 || age > 120) {
            errorBox.textContent = "Age must be between 13 and 120";
            errorBox.style.display = "block";
            return;
        }

        if (!gender) {
            errorBox.textContent = "Please select a gender";
            errorBox.style.display = "block";
            return;
        }

        if (!email || !email.includes("@")) {
            errorBox.textContent = "Please enter a valid email";
            errorBox.style.display = "block";
            return;
        }

        if (!password || password.length < 6) {
            errorBox.textContent = "Password must be at least 6 characters";
            errorBox.style.display = "block";
            return;
        }

        if (password !== confirm_password) {
            errorBox.textContent = "Passwords do not match";
            errorBox.style.display = "block";
            return;
        }

        btn.disabled = true;
        btn.textContent = "Registering...";

        try {
            const resp = await register({
                nickname,
                first_name,
                last_name,
                age,
                gender: gender.value,
                email,
                password,
                confirm_password,
            });

            console.log("registration successful");
            showToast('Registration successful!', 'success');
            router.navigate('/login');
        } catch (err) {
            errorBox.style.display = "block";
            errorBox.textContent = err.message || "Registration failed"; // err.message for debugging (but script is loaded !)
        } finally {
            btn.disabled = false;
            btn.textContent = "Register";
        }
    });
}