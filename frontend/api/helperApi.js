import { router } from "../services/router.js";
import { showToast } from "../services/toast.js";

// Helper function to handle unauthorized errors centrally
export function handleAuthError(response) {
  if (response.status === 401) {
    showToast("Please login", "error"); // Changed to "error" style since it's a failure
    router.navigate("/login");
    return true;
  }
  return false;
}