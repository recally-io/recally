import { createLazyFileRoute, Navigate } from "@tanstack/react-router";

import CTA from "@/components/landing/cta";
import Features from "@/components/landing/features";
import Footer from "@/components/landing/footer";
import Header from "@/components/landing/header";
import Hero from "@/components/landing/hero";
import Pricing from "@/components/landing/pricing";
import Testimonials from "@/components/landing/testimonials";
import { useUser } from "@/lib/apis/auth";
import { ROUTES } from "@/lib/router";

export const Route = createLazyFileRoute("/")({
	component: Index,
});

function Index() {
	const { user } = useUser();

	if (user) {
		return (
			<Navigate
				to={ROUTES.BOOKMARKS}
				search={{
					page: 1,
					filters: [],
					query: "",
				}}
			/>
		);
	}

	return (
		<div className="min-h-screen bg-background font-sans">
			<Header />
			<main>
				<Hero />
				<Features />
				<Testimonials />
				<Pricing />
				<CTA />
			</main>
			<Footer />
		</div>
	);
}
