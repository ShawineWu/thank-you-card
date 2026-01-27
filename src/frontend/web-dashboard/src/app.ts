import { setupConfig } from "./config/setupConfig";
import { AuthService } from "./services/AuthService";

/**
 * Initialize configuration and authentication
 * This promise should be awaited before rendering the app
 */
export const initConfigProcess = setupConfig().then(() => {
  return new AuthService({
    authServer: AUTH_SERVER,
    clientID: CLIENT_ID,
    loginCallbackUrl: LOGIN_CALLBACK_URL,
    logoutCallbackUrl: LOGOUT_CALLBACK_URL,
  });
});
