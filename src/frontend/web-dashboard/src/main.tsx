import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.tsx";
import { initConfigProcess } from "./app";

// Initialize configuration and authentication before rendering
initConfigProcess
  .then(() => {
    console.log("[App] Configuration and authentication initialized");

    createRoot(document.getElementById("root")!).render(
      <StrictMode>
        <App />
      </StrictMode>,
    );
  })
  .catch((error) => {
    console.error("[App] Failed to initialize:", error);

    // Show error message to user
    const root = document.getElementById("root");
    if (root) {
      root.innerHTML = `
        <div style="display: flex; flex-direction: column; justify-content: center; align-items: center; height: 100vh; font-family: sans-serif;">
          <h2 style="color: #dc3545;">Failed to Initialize Application</h2>
          <p style="color: #666;">Please refresh the page or contact support if the problem persists.</p>
          <pre style="background: #f5f5f5; padding: 16px; border-radius: 4px; margin-top: 16px; max-width: 600px; overflow: auto;">
            ${error.message || "Unknown error"}
          </pre>
        </div>
      `;
    }
  });
