import { Button } from "@/components/ui/button";
import { ROUTES } from "@/lib/router";
import { SiTelegram } from "@icons-pack/react-simple-icons";
import { Link } from "@tanstack/react-router";

export default function Hero() {
	return (
		<section className="py-24 px-6 text-center">
			<h1 className="text-4xl font-extrabold tracking-tight sm:text-5xl md:text-6xl">
				Save what matters. Recall what counts.
			</h1>
			<p className="mt-6 text-lg text-muted-foreground max-w-3xl mx-auto">
				Save interesting articles, newsletters, and documents. Read and annotate
				them on any device. Recall valuable insights with powerful search and
				organization.
			</p>
			<div className="mt-10 flex flex-col items-center gap-6">
				<div className="flex justify-center gap-4">
					<Button size="lg">
						<Link href={ROUTES.AUTH_REGISTER}>Start for free</Link>
					</Button>
					<Button size="lg" variant="outline" asChild>
						<a
							href="https://t.me/RecallyReaderBot"
							target="_blank"
							className="flex items-center gap-2"
							rel="noreferrer"
						>
							<SiTelegram className="w-5 h-5" />
							Try with Telegram
						</a>
					</Button>
				</div>
				<div className="text-sm text-muted-foreground">
					<p>
						Quick start: Send any link to our{" "}
						<Link
							href="https://t.me/RecallyReaderBot"
							target="_blank"
							className="text-primary hover:underline underline-offset-4"
						>
							Telegram bot
						</Link>{" "}
						and get instant AI summaries
					</p>
				</div>
			</div>
			<div className="mt-16">
				{/* <img
					src="/logo.svg?height=400&width=800"
					alt="Readwise Reader interface"
					className="rounded-lg shadow-xl"
				/> */}
			</div>
		</section>
	);
}
