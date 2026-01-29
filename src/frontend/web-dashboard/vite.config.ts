import path from "path";
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: true,
    port: 3000,
    proxy: {
      "/cards": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => "/api/v1" + path,
      },
      "/values": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => "/api/v1" + path,
      },
      "/statistics": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => "/api/v1" + path,
      },
      "/analytics": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => "/api/v1" + path,
      },
      "/employees": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => "/api/v1" + path,
      },
    },
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
