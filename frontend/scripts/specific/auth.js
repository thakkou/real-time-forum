const login = async (credentials) => {
    const {
        nickname,
        first_name,
        last_name,
        age,
        gender,
        email,
        password
    } = credentials;

    // Regex validation rules
    const rules = {
        nickname: /^[a-zA-Z0-9_]{3,20}$/,
        first_name: /^[A-Za-z]{2,30}$/,
        last_name: /^[A-Za-z]{2,30}$/,
        gender: /^(male|female)$/i,
        email: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
        password: /^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[@$!%*?&]).{8,}$/
    };

    if (!rules.nickname.test(nickname)) {
        throw new Error("Invalid nickname");
    }

    if (!rules.first_name.test(first_name)) {
        throw new Error("Invalid first name");
    }

    if (!rules.last_name.test(last_name)) {
        throw new Error("Invalid last name");
    }

   const parsedAge = parseInt(age);

if (isNaN(parsedAge) || parsedAge < 18 || parsedAge > 99) {
    throw new Error("Age must be between 18 and 99");
}

    if (!rules.gender.test(gender)) {
        throw new Error("Invalid gender");
    }

    if (!rules.email.test(email)) {
        throw new Error("Invalid email");
    }

    if (!rules.password.test(password)) {
        throw new Error(
            "Password must contain uppercase, lowercase, number, special character and be at least 8 chars"
        );
    }

   try{
    const res= await fetch(`${API_uri}/register`,{
        method: "POST",
            },
        credentials)
        console.log(res)
    } catch(err) {
      throw new Error(
            err
        );
   }

    return "Login validation passed";
};