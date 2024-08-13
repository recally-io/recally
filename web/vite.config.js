import { resolve } from "path";
import { defineConfig } from "vite";

import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, "index.html"),
        assistants: resolve(__dirname, "assistants.html"),
        threads: resolve(__dirname, "threads.html"),
      },
    },
  },
});
