import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ROUTES } from "@/lib/router";
import { Link } from "@tanstack/react-router";
import { Check } from "lucide-react";

const plans = [
	{
		name: "Beta Access",
		price: "Free",
		description: "Get early access to all features while we're in beta",
		features: [
			"Unlimited saves",
			"AI-powered summaries",
			"Full-text search",
			"Self-hosting option",
			"Mobile app access",
			"Priority support",
			"Help shape the future of Recally",
		],
	},
];

export default function Pricing() {
	return (
		<section id="pricing" className="py-24 px-6 bg-muted">
			<div className="container mx-auto">
				<h2 className="text-3xl font-bold text-center mb-4">
					Try Recally Beta
				</h2>
				<p className="text-center text-muted-foreground mb-12">
					Free while in beta. We'll announce pricing with plenty of notice.
				</p>
				<div className="max-w-md mx-auto">
					<Card className="border-primary">
						<CardHeader>
							<CardTitle>{plans[0].name}</CardTitle>
							<CardDescription>{plans[0].description}</CardDescription>
						</CardHeader>
						<CardContent>
							<p className="text-4xl font-bold mb-4">
								{plans[0].price}
								<span className="text-sm font-normal text-muted-foreground">
									/month
								</span>
							</p>
							<ul className="space-y-2">
								{plans[0].features.map((feature) => (
									<li key={feature} className="flex items-center">
										<Check className="w-4 h-4 mr-2 text-green-500" />
										{feature}
									</li>
								))}
							</ul>
						</CardContent>
						<CardFooter>
							<Button className="w-full" variant="default">
								<Link href={ROUTES.AUTH_REGISTER}> Start for free</Link>
							</Button>
						</CardFooter>
					</Card>
				</div>
			</div>
		</section>
	);
}
