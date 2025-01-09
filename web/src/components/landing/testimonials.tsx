import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Card, CardContent, CardFooter } from "@/components/ui/card";

const testimonials = [
	{
		quote:
			"Finally found a tool that doesn't make saving articles feel like a chore! The AI summaries are spot-on and help me remember why I saved stuff in the first place.",
		author: "Alex Chen",
		role: "Software Engineer",
		avatar: "AC",
	},
	{
		quote:
			"Love how I can self-host it! Been using it for my research papers and the search actually understands what I'm looking for. Game changer for my PhD work.",
		author: "Sarah Mitchell",
		role: "PhD Student",
		avatar: "SM",
	},
	{
		quote:
			"I was skeptical about another 'AI-powered' tool, but Recally's features actually make sense. The tagging suggestions are surprisingly accurate!",
		author: "Marcus Torres",
		role: "Content Writer",
		avatar: "MT",
	},
];

export default function Testimonials() {
	return (
		<section className="py-24 px-6">
			<div className="container mx-auto">
				<h2 className="text-3xl font-bold text-center mb-12">
					What Our Users Say
				</h2>
				<div className="grid grid-cols-1 md:grid-cols-3 gap-6">
					{testimonials.map((testimonial) => (
						<Card key={testimonial.author} className="text-center">
							<CardContent className="pt-6">
								<p className="text-muted-foreground italic mb-4">
									"{testimonial.quote}"
								</p>
							</CardContent>
							<CardFooter className="flex flex-col items-center">
								<Avatar className="w-12 h-12 mb-2">
									<AvatarFallback>{testimonial.avatar}</AvatarFallback>
								</Avatar>
								<div>
									<p className="font-semibold">{testimonial.author}</p>
									<p className="text-sm text-muted-foreground">
										{testimonial.role}
									</p>
								</div>
							</CardFooter>
						</Card>
					))}
				</div>
			</div>
		</section>
	);
}
