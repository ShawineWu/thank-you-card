import { initConfigProcess } from "../app";

/**
 * Logout button component
 * Initiates the logout flow when clicked
 */
export function LogoutButton() {
  async function handleLogout() {
    try {
      const auth = await initConfigProcess;
      await auth.logout();
    } catch (error) {
      console.error("[Logout] Error during logout:", error);
    }
  }

  return (
    <button
      onClick={handleLogout}
      style={{
        padding: "8px 16px",
        backgroundColor: "#dc3545",
        color: "white",
        border: "none",
        borderRadius: "4px",
        cursor: "pointer",
        fontSize: "14px",
      }}
    >
      Logout
    </button>
  );
}
