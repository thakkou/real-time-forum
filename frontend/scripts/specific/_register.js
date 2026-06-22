function setupRegisterPage() {
    const form = document.getElementById("register-form");
    if (!form) return;
    
    form.addEventListener("submit", async (e) => {
    e.preventDefault();
    
    const btn = document.getElementById("register-btn");
    const errorBox = document.getElementById("form-error");
    const formData = e.target;
    
    const nickname = formData.nickname.value;
    const first_name = formData.first_name.value;
    const last_name = formData.last_name.value;
    const age = parseInt(formData.age.value, 10);
    const gender = formData.querySelector('input[name="gender"]:checked');
    const email = formData.email.value;
    const password = formData.password.value;
    const confirm_password = formData.confirm_password.value;
    
    errorBox.style.display = "none";
    errorBox.textContent = "";
    
    if (!nickname || nickname.trim().length < 2) {
        errorBox.textContent = "Nickname must be at least 2 characters";
        errorBox.style.display = "block";
        return;
    }
    if (!first_name || first_name.trim().length < 1) {
        errorBox.textContent = "First name is required";
        errorBox.style.display = "block";
        return;
    }
    if (!last_name || last_name.trim().length < 1) {
        errorBox.textContent = "Last name is required";
        errorBox.style.display = "block";
        return;
    }
    if (!age || isNaN(age) || age < 13 || age > 120) {
        errorBox.textContent = "Age must be a number between 13 and 120";
        errorBox.style.display = "block";
        return;
    }
    if (!gender) {
        errorBox.textContent = "Please select a gender";
        errorBox.style.display = "block";
        return;
    }
    if (!email || !email.includes("@") || !email.includes(".")) {
        errorBox.textContent = "Please enter a valid email address";
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
    
    console.log(JSON.stringify({ 
        nickname, first_name, last_name, age, 
        gender: gender.value, email, password, confirm_password 
    }));
    
    btn.disabled = true;
    btn.textContent = "Registering...";
    });
}

if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', setupRegisterPage);
} else {
    setupRegisterPage();
}