import { Button } from "@/components/ui/button";
import { ROUTES } from "@/lib/router";
import { Link } from "@tanstack/react-router";

export default function CTA() {
	return (
		<section className="py-24 px-6 bg-primary text-primary-foreground">
			<div className="container text-center mx-auto">
				<h2 className="text-3xl font-bold mb-6">
					Ready to transform your reading experience?
				</h2>
				<p className="text-xl mb-8 max-w-2xl mx-auto">
					Join thousands of readers who have already discovered the power of
					Recally. Start your free trial today!
				</p>
				<Button size="lg" variant="secondary">
					<Link href={ROUTES.AUTH_REGISTER}> Start Your Free Trial </Link>
				</Button>
			</div>
		</section>
	);
}
