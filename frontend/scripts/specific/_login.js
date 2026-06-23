
// import {login} from "../../api/auth.js"
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

const serverURI = env.serverUri;

 const login = async (credentials) => {
  try {
    console.log("login the api", credentials);
    console.log(`${serverURI}/login`);

    const response = await fetch(`${serverURI}/login`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        identifier: credentials.identifier,
        password: credentials.password,
      }),
    });
    const data = await response.json();
    console.log("the response",data,data.status_code)

    if (data.status_code !=200) {
      throw new Error(data.message || "Login failed");
    }

    return data;
  } catch (error) {
    console.error("Login error:", error);
    throw error;
  }
};



