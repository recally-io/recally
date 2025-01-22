import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
	title: "Recally ",
	description: "Documents for Recally",
	base: "/docs/",
	srcDir: "src",
	lastUpdated: true,
	themeConfig: {
		// https://vitepress.dev/reference/default-theme-config
		logo: "/logo.svg",
		search: {
			provider: "local",
		},
		nav: [
			{
				text: "Recally Cloud",
				link: "https://recally.io",
				target: "_self",
				rel: "sponsored",
			},
      {
				text: "API Reference",
				link: "https://recally.io/swagger/index.html",
				target: "_self",
				rel: "sponsored",
			},
			{
				text: "Policies",
				items: [
					{ text: "Privacy", link: "/privacy-policy" },
					{ text: "Terms", link: "/terms-of-service" },
				],
			},
		],

		sidebar: [
      {
        text: "Introduction",
        link: "/introduction",
      },
			{
				text: "Tutorials",
				items: [],
			},
			{
				text: "How-to guides",
				items: [],
			},
			{
				text: "Reference",
				items: [],
			},
		],

		socialLinks: [
			{ icon: "github", link: "https://github.com/recally-io/recally" },
      { icon: 'twitter', link: "https://twitter.com/recally_io" },
		],
		footer: {
			message:
				'Released under the <a href="https://github.com/recally-io/recally/blob/main/LICENSE">Recally License</a>.',
			copyright:
				'Copyright © 2025-present <a href="https://github.com/recally-io">Recally</a>',
		},
	},
});
