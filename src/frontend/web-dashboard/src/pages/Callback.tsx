import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { initConfigProcess } from "../app";

/**
 * Callback page that handles OAuth redirect after login
 * Processes the authorization code and redirects user back to their original destination
 */
export function CallbackPage() {
  const navigate = useNavigate();

  const handleCallback = async () => {
    try {
      const auth = await initConfigProcess;
      const user = await auth.signinRedirectCallback();

      console.log("[Callback] Login successful, redirecting...");

      // Redirect to the page user was trying to access
      const state = user?.state as { url?: string } | undefined;
      const redirectUrl = state?.url || "/dashboard";
      const url = new URL(redirectUrl, window.location.origin);

      // Avoid infinite redirect loop
      if (url.pathname !== "/callback") {
        navigate(url.pathname + url.search, { replace: true });
      } else {
        navigate("/dashboard", { replace: true });
      }
    } catch (error) {
      console.error("[Callback] Error processing callback:", error);

      // Handle specific error cases
      if (
        error instanceof Error &&
        error.message === "No matching state found in storage"
      ) {
        // State mismatch, redirect to home
        navigate("/dashboard", { replace: true });
      } else {
        // Other errors, show error message or redirect to home
        navigate("/dashboard", { replace: true });
      }
    }
  };

  useEffect(() => {
    handleCallback();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
        height: "100vh",
        backgroundColor: "#f5f5f5",
      }}
    >
      <h2 style={{ color: "#333", marginBottom: "16px" }}>
        Processing login...
      </h2>
      <p style={{ color: "#666" }}>
        Please wait while we complete your authentication.
      </p>
    </div>
  );
}
