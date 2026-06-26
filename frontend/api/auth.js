const serverURI = env.serverUri;

export const login = async (credentials) => {
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

    if (!response.ok) {
      throw new Error(data.message || "Login failed");
    }

    console.log(data);
    return data;
  } catch (error) {
    console.error("Login error:", error);
    throw error;
  }
};

export const register = async ({ data }) => {
  try {
    const response = await fetch(`${serverURI}/register`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });

    const result = await response.json();

    if (!response.ok) {
      throw new Error(result.message || "Registration failed");
    }

    return result;
  } catch (error) {
    console.error("Registration error:", error);
    throw error;
  }
};

export const logout = async () => {
  try {
    const response = await fetch(`${serverURI}/logout`, {
      method: "POST",
      credentials: "include",
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || "Logout failed");
    }

    return data;
  } catch (error) {
    console.error("Logout error:", error);
    throw error;
  }
};

// const login = async (credentials) => {
//     const {
//         nickname,
//         first_name,
//         last_name,
//         age,
//         gender,
//         email,
//         password
//     } = credentials;

//     // Regex validation rules
//     const rules = {
//         nickname: /^[a-zA-Z0-9_]{3,20}$/,
//         first_name: /^[A-Za-z]{2,30}$/,
//         last_name: /^[A-Za-z]{2,30}$/,
//         gender: /^(male|female)$/i,
//         email: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
//         password: /^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[@$!%*?&]).{8,}$/
//     };

//     if (!rules.nickname.test(nickname)) {
//         throw new Error("Invalid nickname");
//     }

//     if (!rules.first_name.test(first_name)) {
//         throw new Error("Invalid first name");
//     }

//     if (!rules.last_name.test(last_name)) {
//         throw new Error("Invalid last name");
//     }

//    const parsedAge = parseInt(age);

// if (isNaN(parsedAge) || parsedAge < 18 || parsedAge > 99) {
//     throw new Error("Age must be between 18 and 99");
// }

//     if (!rules.gender.test(gender)) {
//         throw new Error("Invalid gender");
//     }

//     if (!rules.email.test(email)) {
//         throw new Error("Invalid email");
//     }

//     if (!rules.password.test(password)) {
//         throw new Error(
//             "Password must contain uppercase, lowercase, number, special character and be at least 8 chars"
//         );
//     }

//    try{
//     const res= await fetch(`${API_uri}/register`,{
//         method: "POST",
//             },
//         credentials)
//         console.log(res)
//     } catch(err) {
//       throw new Error(
//             err
//         );
//    }

//     return "Login validation passed";
// };