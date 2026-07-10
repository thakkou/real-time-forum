export function formatTime(dateStr) {
  if (!dateStr) return "";
  return new Date(dateStr).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });
}
export function sanitize(str) {
  if (str == null) {
    return null;
  }

  return str
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}
export function formatDate(input) {
    const date = input instanceof Date ? input : new Date(input);

    // Handle invalid dates
    if (isNaN(date.getTime())) {
        return "";
    }

    const now = new Date();

    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const d = new Date(date.getFullYear(), date.getMonth(), date.getDate());

    const diff = Math.floor((today - d) / 86400000);

    if (diff === 0) return "Today";
    if (diff === 1) return "Yesterday";

    return date.toLocaleDateString([], {
        month: "short",
        day: "numeric",
    });
}