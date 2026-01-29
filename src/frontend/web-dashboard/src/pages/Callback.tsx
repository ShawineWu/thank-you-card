import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { initConfigProcess } from "../app";

/**
 * Callback page that handles OAuth redirect after login
 * Processes the authorization code and redirects user back to their original destination
 */
export function CallbackPage() {
  const navigate = useNavigate();
  const processedRef = useRef(false);
  const [error, setError] = useState<string | null>(null);

  const handleCallback = async () => {
    // Prevent double execution in React Strict Mode
    if (processedRef.current) return;
    processedRef.current = true;

    try {
      console.log("[Callback] Processing callback...");
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
    } catch (err) {
      console.error("[Callback] Error processing callback:", err);

      let errorMessage = "Authentication failed. Please try again.";

      if (err instanceof Error) {
        if (err.message === "No matching state found in storage") {
          errorMessage =
            "Login session expired or invalid state. Please try logging in again.";
        } else {
          errorMessage = err.message;
        }
      }

      setError(errorMessage);
    }
  };

  useEffect(() => {
    handleCallback();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (error) {
    return (
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
          backgroundColor: "#f5f5f5",
          padding: "20px",
          textAlign: "center",
        }}
      >
        <h2 style={{ color: "#d32f2f", marginBottom: "16px" }}>Login Failed</h2>
        <p style={{ color: "#666", marginBottom: "24px", maxWidth: "400px" }}>
          {error}
        </p>
        <button
          onClick={() => navigate("/dashboard", { replace: true })}
          style={{
            padding: "10px 20px",
            backgroundColor: "#1976d2",
            color: "white",
            border: "none",
            borderRadius: "4px",
            cursor: "pointer",
            fontSize: "16px",
          }}
        >
          Return to Dashboard
        </button>
      </div>
    );
  }

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
