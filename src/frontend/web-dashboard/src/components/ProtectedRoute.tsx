import { useEffect, useState, type ReactNode } from "react";
import { useLocation } from "react-router-dom";
import { initConfigProcess } from "../app";
import { isTokenValid } from "../utils/token";
import { useAuthStore } from "../store/authStore";

interface ProtectedRouteProps {
  children: ReactNode;
}

/**
 * Route guard component that ensures user is authenticated
 * Automatically refreshes tokens if expired
 * Redirects to login if authentication fails
 */
export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { user, setUser, setLoading, isLoading } = useAuthStore();
  const [isInit, setIsInit] = useState(false);
  const location = useLocation();

  useEffect(() => {
    const checkAuth = async () => {
      console.log("[Auth] Starting checkAuth...");
      try {
        const auth = await initConfigProcess;
        console.log("[Auth] Config initialized, getting user...");
        const currentUser = await auth.getUser();
        console.log("[Auth] User retrieved:", currentUser ? "Found" : "Null");

        if (!currentUser) {
          console.log("[Auth] No user found, redirecting to login...");
          await auth.login();
          return;
        }

        if (!isTokenValid(currentUser)) {
          // Token expired, try to refresh
          const now = Math.floor(Date.now() / 1000);
          const expiresAt = currentUser?.expires_at || 0;
          console.log(
            `[Auth] Token expired or invalid. ExpiresAt: ${expiresAt}, Now: ${now}, Diff: ${expiresAt - now}`,
          );
          console.log("[Auth] Attempting refresh...");
          const refreshedUser = await auth.renewToken();
          console.log(
            "[Auth] Token renewed:",
            refreshedUser ? "Success" : "Failed",
          );
          setUser(refreshedUser);
        } else {
          console.log("[Auth] Token is valid, setting user...");
          setUser(currentUser);
        }
      } catch (error) {
        // Token invalid or refresh failed, redirect to login
        console.log(
          "[Auth] Authentication failed, redirecting to login",
          error,
        );
        const auth = await initConfigProcess;
        await auth.clearUser();
        await auth.login();
      } finally {
        console.log("[Auth] checkAuth finally block - setting loading false");
        setLoading(false);
        setIsInit(true);
      }
    };

    checkAuth();
  }, [location.pathname, setUser, setLoading]);

  if (isLoading || !isInit) {
    return (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
          background: "var(--bg-color)",
          color: "var(--text-primary)",
        }}
      >
        <div className="loader">Loading...</div>
      </div>
    );
  }

  return user ? <>{children}</> : null;
}
