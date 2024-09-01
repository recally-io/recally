import { resolve } from "path";
import { defineConfig } from "vite";

import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  // esbuild: {
  //   supported: {
  //     "top-level-await": true,
  //   },
  // },
  // optimizeDeps: {
  //   esbuildOptions: {
  //     target: "esnext",
  //   },
  // },
  build: {
    // target: "esnext",
    rollupOptions: {
      input: {
        main: resolve(__dirname, "index.html"),
        assistants: resolve(__dirname, "assistants.html"),
        threads: resolve(__dirname, "threads.html"),
        auth: resolve(__dirname, "auth.html"),
      },
    },
  },
  server: {
    // open: "/assistants.html",
    // https://vitejs.dev/config/server-options#server-proxy
    proxy: {
      "/api": {
        target: "http://localhost:1323",
        changeOrigin: true,
        secure: false,
      },
    },
  },
});
