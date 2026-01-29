import { UserManager, WebStorageStateStore, Log, User } from "oidc-client-ts";

export interface AuthServiceConfig {
  authServer: string;
  clientID: string;
  loginCallbackUrl: string;
  logoutCallbackUrl: string;
}

/**
 * Authentication service wrapper around oidc-client-ts
 * Handles OAuth2/OIDC authentication flow with Hydra
 */
export class AuthService {
  private userManager: UserManager;
  private isLogining: boolean = false;
  private renewTokenPromise: Promise<User | null> | null = null;

  constructor(config: AuthServiceConfig) {
    const settings = {
      authority: config.authServer,
      client_id: config.clientID,
      response_type: "code",
      scope: "offline_access offline openid",
      redirect_uri: config.loginCallbackUrl,
      post_logout_redirect_uri: config.logoutCallbackUrl,
      userStore: new WebStorageStateStore({ store: localStorage }),
      automaticSilentRenew: true,
    };

    this.userManager = new UserManager(settings);
    this.userManager.clearStaleState();

    // Configure logging
    Log.setLevel(Log.INFO);
    Log.setLogger(console);
  }

  /**
   * Initiates the login flow by redirecting to the authorization server
   */
  async login(): Promise<void> {
    if (this.isLogining) {
      return;
    }

    this.isLogining = true;
    const url = new URL(window.location.href);
    const args = { state: { url: url.href } };

    // Avoid redirect loop
    if (window.location.pathname !== "/callback") {
      return this.userManager.signinRedirect(args);
    }
  }

  /**
   * Handles the callback from the authorization server
   * Processes the authorization code and retrieves tokens
   */
  async signinRedirectCallback(): Promise<User> {
    const user = await this.userManager.signinRedirectCallback();
    this.isLogining = false;
    return user;
  }

  /**
   * Renews the access token using the refresh token
   * Prevents multiple concurrent renewal requests
   */
  async renewToken(): Promise<User | null> {
    if (this.renewTokenPromise) {
      return this.renewTokenPromise;
    }

    const silentRenew = this.userManager.signinSilent();
    const timeout = new Promise<null>((_, reject) =>
      setTimeout(() => reject(new Error("Token renewal timed out")), 10000),
    );

    this.renewTokenPromise = Promise.race([silentRenew, timeout])
      .then((user) => {
        if (!user) return null;
        return user as User;
      })
      .finally(() => {
        this.renewTokenPromise = null;
      });

    return this.renewTokenPromise;
  }

  /**
   * Gets the current user from storage
   */
  async getUser(): Promise<User | null> {
    return this.userManager.getUser();
  }

  /**
   * Initiates the logout flow
   */
  async logout(): Promise<void> {
    return this.userManager.signoutRedirect();
  }

  /**
   * Clears the current user from storage
   */
  async clearUser(): Promise<void> {
    return this.userManager.removeUser();
  }
}
