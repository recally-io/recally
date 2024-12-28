import react from "@vitejs/plugin-react";
import path from "node:path";
import { defineConfig } from "vite";
import { VitePWA } from "vite-plugin-pwa";

export default defineConfig({
	plugins: [
		react(),
		VitePWA({
			registerType: "autoUpdate",
			workbox: {
        cleanupOutdatedCaches: false
      },
			includeAssets: ["favicon.ico", "apple-touch-icon.png", "maskable-icon-512x512.png"],
			manifest: {
				name: "Vibrain",
				short_name: "Vibrain",
				description: "Save what matters, recall when it counts.",
				theme_color: "#ffffff",
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
