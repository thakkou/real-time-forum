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