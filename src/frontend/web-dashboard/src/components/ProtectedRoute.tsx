import { useEffect, useState, type ReactNode } from "react";
import { useLocation } from "react-router-dom";
import { initConfigProcess } from "../app";
import { isTokenValid } from "../utils/token";

interface ProtectedRouteProps {
  children: ReactNode;
}

/**
 * Route guard component that ensures user is authenticated
 * Automatically refreshes tokens if expired
 * Redirects to login if authentication fails
 */
export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  const location = useLocation();

  const checkAuth = async () => {
    try {
      const auth = await initConfigProcess;
      const user = await auth.getUser();

      if (!isTokenValid(user)) {
        // Token expired, try to refresh
        console.log("[Auth] Token expired, refreshing...");
        await auth.renewToken();
        setIsAuthenticated(true);
      } else {
        setIsAuthenticated(true);
      }
    } catch {
      // Token invalid or refresh failed, redirect to login
      console.log("[Auth] Authentication failed, redirecting to login");
      const auth = await initConfigProcess;
      await auth.login();
    }
  };

  useEffect(() => {
    checkAuth();
  }, [location.pathname]);

  if (isAuthenticated === null) {
    return (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
          fontSize: "18px",
          color: "#666",
        }}
      >
        Loading...
      </div>
    );
  }

  return isAuthenticated ? <>{children}</> : null;
}
