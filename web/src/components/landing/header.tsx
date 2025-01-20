import { Button } from "@/components/ui/button";
import { ROUTES } from "@/lib/router";
import { SiGithub } from "@icons-pack/react-simple-icons";
import { Link } from "@tanstack/react-router";

export default function Header() {
	return (
		<header className="py-4 px-6 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 sticky top-0 z-40 w-full border-b border-border/40">
			<div className="container flex items-center justify-between">
				<Link href="/" className="flex items-center space-x-2">
					{/* <img src="/logo.svg" alt="Recally" className="w-8 h-8" /> */}
					<span className="text-2xl font-bold">Recally</span>
				</Link>
				<nav className="hidden md:flex items-center space-x-6">
					<Link
						href="#features"
						className="text-sm font-medium hover:underline underline-offset-4"
					>
						Features
					</Link>
					<Link
						href="#pricing"
						className="text-sm font-medium hover:underline underline-offset-4"
					>
						Pricing
					</Link>
					<a
						href="/docs/"
						target="_blank"
						className="text-sm font-medium hover:underline underline-offset-4"
						rel="noreferrer"
					>
						Docs
					</a>
					<a
						href="https://github.com/recally-io/recally"
						target="_blank"
						rel="noopener noreferrer"
						className="flex gap-2 items-center text-sm font-medium hover:underline underline-offset-4"
					>
						<SiGithub className="w-4 h-4" />
					</a>
				</nav>
				<div className="flex items-center space-x-4">
					<Button variant="ghost">
						<Link href={ROUTES.AUTH_LOGIN}> Log in </Link>
					</Button>
					<Button>
						<Link href={ROUTES.AUTH_REGISTER}> Sign up </Link>
					</Button>
				</div>
			</div>
		</header>
	);
}
