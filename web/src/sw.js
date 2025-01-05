import { clientsClaim } from "workbox-core";
import { precacheAndRoute } from "workbox-precaching";

precacheAndRoute(self.__WB_MANIFEST);
self.skipWaiting();
clientsClaim();

self.addEventListener("fetch", (event) => {
	if (
		event.request.method === "POST" &&
		event.request.url.endsWith("/save-bookmark") // Remove trailing slash
	) {
		event.respondWith(
			(async () => {
				const formData = await event.request.formData();
				const link = formData.get("url") || "";
				const title = formData.get("title") || "";
				const text = formData.get("text") || "";

				if (!link) {
					return new Response("URL is required.", {
						status: 400,
						statusText: "Bad Request",
					});
				}

				try {
					const resp = await saveBookmark(link);
					return Response.redirect("/bookmarks", 303);
				} catch (error) {
					return new Response(error.message, {
						status: 500,
						statusText: "Internal Server Error",
					});
				}
			})(),
		);
	}
});

async function saveBookmark(link) {
	// Save the bookmark to the server and return the response URL.
	return fetch("/api/v1/bookmarks", {
		method: "POST",
		body: JSON.stringify({ url: link }),
		credentials: "include",
		headers: {
			"Content-Type": "application/json",
		},
	})
		.then((response) => {
			if (!response.ok) {
				throw new Error(`Failed to save the bookmark: ${response.statusText}`);
			}
			return response.json();
		})
		.then((data) => {
			return data.data;
		});
}
