import { User } from "oidc-client-ts";

/**
 * Checks if the user's access token is valid and not expired
 * Considers token expired 5 minutes before actual expiry for safety margin
 */
export function isTokenValid(user: User | null): boolean {
  if (!user || !user.access_token) {
    return false;
  }

  // Check if token is expired
  const now = Math.floor(Date.now() / 1000);
  const expiresAt = user.expires_at || 0;

  // Consider token expired 5 minutes (300 seconds) before actual expiry
  return expiresAt > now + 300;
}
