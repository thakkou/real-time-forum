
import {login} from "../../api/auth.js"

const setupLoginPage = () => {
  console.log("kisk");

  const form = document.querySelector("form");
  const errorBox = document.getElementById("login-error");

  if (!form) return;

 form.addEventListener("submit", async (e) => {
  e.preventDefault();

  const identifier = form.identifier.value;
  const password = form.password.value;

  try {
    const res = await login({ identifier, password });

    console.log("login success", res);

    // Store user globally
    window.user = res.data;

    // Redirect to home
    window.location.href = "/";

  } catch (err) {
    errorBox.style.display = "block";
    errorBox.textContent = err.message || "Login failed";
  }
});
};

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', setupLoginPage);
} else {
    setupLoginPage();
}




