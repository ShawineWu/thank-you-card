import { useState } from "react";
import { LogoutButton } from "../components/LogoutButton";
import reactLogo from "../assets/react.svg";
import viteLogo from "/vite.svg";

/**
 * Home page - main application page after authentication
 */
export function HomePage() {
  const [count, setCount] = useState(0);

  return (
    <div style={{ position: "relative" }}>
      {/* Logout button in top right */}
      <div style={{ position: "absolute", top: "20px", right: "20px" }}>
        <LogoutButton />
      </div>

      {/* Original Vite template content */}
      <div>
        <a href="https://vite.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Thank You Card - Web Dashboard</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>Authentication is working! You are logged in.</p>
        <p style={{ fontSize: "14px", color: "#666" }}>
          Edit <code>src/pages/Home.tsx</code> to build your dashboard
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </div>
  );
}
