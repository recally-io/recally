import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { BookOpen, Bot, Search, Shield } from "lucide-react";

const features = [
	{
		title: "Save Anything, Find Everything",
		description:
			"One-click saving of articles, PDFs, videos, and podcasts with smart search capability.",
		icon: Search,
	},
	{
		title: "AI That Makes Sense",
		description:
			"Smart summaries, intelligent tagging, and chat with your documents powered by AI.",
		icon: Bot,
	},
	{
		title: "Your Content, Your Control",
		description:
			"Self-host option available, no tracking, and flexible organization system.",
		icon: Shield,
	},
	{
		title: "Beautiful Reading Experience",
		description:
			"Saved content looks great and is easily accessible across all your devices.",
		icon: BookOpen,
	},
];

export default function Features() {
	return (
		<section id="features" className="py-24 px-6 bg-muted">
			<div className="container mx-auto">
				<h2 className="text-3xl font-bold text-center mb-12">Key Features</h2>
				<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
					{features.map((feature) => (
						<Card key={feature.title.toLowerCase().replace(/\s+/g, "-")}>
							<CardHeader>
								<feature.icon className="w-10 h-10 mb-4 text-primary" />
								<CardTitle>{feature.title}</CardTitle>
							</CardHeader>
							<CardContent>
								<CardDescription>{feature.description}</CardDescription>
							</CardContent>
						</Card>
					))}
				</div>
			</div>
		</section>
	);
}
