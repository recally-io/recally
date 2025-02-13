import { TanStackRouterVite } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";
import path from "node:path";
import { defineConfig } from "vite";
import { VitePWA } from "vite-plugin-pwa";

export default defineConfig({
	plugins: [
		react(),
		TanStackRouterVite(),
		VitePWA({
			registerType: "autoUpdate",
			workbox: {
				cleanupOutdatedCaches: true,
			},
			includeAssets: [
				"favicon.ico",
				"apple-touch-icon.png",
				"maskable-icon-512x512.png",
			],
			strategies: "injectManifest",
			injectManifest: {
				maximumFileSizeToCacheInBytes: 5000000
			},
			srcDir: "src",
			filename: "sw.js",
			manifest: {
				name: "Recally",
				short_name: "Recally",
				display: "fullscreen",
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
					{
						src: "pwa-512x512.png",
						sizes: "512x512",
						type: "image/png",
						purpose: "any",
					},
					{
						src: "maskable-icon-512x512.png",
						sizes: "512x512",
						type: "image/png",
						purpose: "maskable",
					},
				],
				share_target: {
					action: "/save-bookmark",
					method: "POST",
					enctype: "multipart/form-data",
					params: {
						title: "title",
						text: "text",
						url: "url",
					},
				},
			},
		}),
	],
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
			"/docs/": {
				target: "http://localhost:1323",
				changeOrigin: true,
			},
			"/swagger/": {
				target: "http://localhost:1323",
				changeOrigin: true,
			},
		},
	},
});
