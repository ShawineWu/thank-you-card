type Env = "local" | "dev" | "test" | "prod";

/**
 * Loads environment-specific configuration and injects it into the window object
 * Configuration is determined by the current hostname
 */
export async function setupConfig(): Promise<void> {
  // Map hostname to environment
  const hostnameToEnv: Record<string, Env> = {
    localhost: "local",
    "thank-you-card-dev.castlery.com": "dev",
    "thank-you-card-test.castlery.com": "test",
    "thank-you-card.castlery.com": "prod",
  };

  const hostname = window.location.hostname;
  const env = hostnameToEnv[hostname] || "local";

  // Dynamic import of environment-specific config
  const configs = {
    local: () => import("../config.local.json"),
    dev: () => import("../config.dev.json"),
    test: () => import("../config.test.json"),
    prod: () => import("../config.prod.json"),
  };

  const { default: config } = await configs[env]();

  // Inject configuration into window object as read-only properties
  Object.keys(config).forEach((key) => {
    Object.defineProperty(window, key, {
      value: config[key as keyof typeof config],
      configurable: false,
      enumerable: false,
      writable: false,
    });
  });

  console.log(`[Config] Loaded ${env} environment configuration`);
}
