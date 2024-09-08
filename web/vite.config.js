import react from "@vitejs/plugin-react";
import { resolve } from "path";
import { defineConfig } from "vite";
import { VitePWA } from "vite-plugin-pwa";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    VitePWA({
      workbox: {
        maximumFileSizeToCacheInBytes: 10000000,
      },
      registerType: "autoUpdate",
      devOptions: {
        enabled: true,
        type: "module",
      },
      manifest: {
        name: "Vibrain",
        short_name: "Vibrain",
        description: "Vibrain",
        theme_color: "#000000",
        icons: [
          {
            src: "pwa-192x192.png",
            sizes: "192x192",
            type: "image/png",
          },
          {
            src: "pwa-512x512.png",
            sizes: "512x512",
            type: "image/png",
          },
        ],
      },
    }),
  ],
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
    target: "esnext",
    rollupOptions: {
      input: {
        main: resolve(__dirname, "index.html"),
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
