import react from "@vitejs/plugin-react";
import path from "node:path";
import { defineConfig } from "vite";

export default defineConfig({
	plugins: [react()],
	build: {
		rollupOptions: {
			input: {
				main: path.resolve(__dirname, "index.html"),
				bookmarks: path.resolve(__dirname, "bookmarks.html"),
				auth: path.resolve(__dirname, "auth.html"),
			},
		},
	},
	resolve: {
		alias: {
			"@": path.resolve(__dirname, "./src"),
		},
	},
	server: {
		proxy: {
			"/api": {
				target: "http://localhost:1323",
				changeOrigin: true,
			},
		},
	},
});
