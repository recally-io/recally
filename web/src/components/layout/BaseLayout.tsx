import Footer from "./footer";
import Header from "./header";

interface BaseLayoutProps {
	children: React.ReactNode;
}

export function BaseLayout({ children }: BaseLayoutProps) {
	return (
		<div className="flex flex-col min-h-screen bg-background">
			<header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
				<Header />
			</header>

			<main className="flex-grow container py-8">{children}</main>

			<footer className="sticky bottom-0 z-50 w-full border-t w-full bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
				<Footer />
			</footer>
		</div>
	);
}
