import { Link } from "@tanstack/react-router";

export default function Footer() {
	return (
		<footer className="py-12 px-6 bg-background">
			<div className="container mx-auto">
				<div className="grid grid-cols-2 md:grid-cols-4 gap-8">
					<div>
						<h3 className="font-semibold mb-4">Product</h3>
						<ul className="space-y-2">
							<li>
								<Link
									href="#"
									className="text-sm text-muted-foreground hover:text-foreground"
								>
									Features
								</Link>
							</li>
							<li>
								<Link
									href="#"
									className="text-sm text-muted-foreground hover:text-foreground"
								>
									Pricing
								</Link>
							</li>
							<li>
								<Link
									href="#"
									className="text-sm text-muted-foreground hover:text-foreground"
								>
									Integrations
								</Link>
							</li>
						</ul>
					</div>
					<div>
						<h3 className="font-semibold mb-4">Company</h3>
						<ul className="space-y-2">
							<li>
								<Link
									href="#"
									className="text-sm text-muted-foreground hover:text-foreground"
								>
									About
								</Link>
							</li>
							<li>
								<Link
									href="#"
									className="text-sm text-muted-foreground hover:text-foreground"
								>
									Blog
								</Link>
							</li>
							<li>
								<Link
									href="#"
									className="text-sm text-muted-foreground hover:text-foreground"
								>
									Careers
								</Link>
							</li>
						</ul>
					</div>
					<div>
						<h3 className="font-semibold mb-4">Resources</h3>
						<ul className="space-y-2">
							<li>
								<Link
									href="#"
									className="text-sm text-muted-foreground hover:text-foreground"
								>
									Documentation
								</Link>
							</li>
							<li>
								<Link
									href="#"
									className="text-sm text-muted-foreground hover:text-foreground"
								>
									Help Center
								</Link>
							</li>
							<li>
								<Link
									href="#"
									className="text-sm text-muted-foreground hover:text-foreground"
								>
									Community
								</Link>
							</li>
						</ul>
					</div>
					<div>
						<h3 className="font-semibold mb-4">Legal</h3>
						<ul className="space-y-2">
							<li>
								<Link
									href="#"
									className="text-sm text-muted-foreground hover:text-foreground"
								>
									Privacy Policy
								</Link>
							</li>
							<li>
								<Link
									href="#"
									className="text-sm text-muted-foreground hover:text-foreground"
								>
									Terms of Service
								</Link>
							</li>
							<li>
								<Link
									href="#"
									className="text-sm text-muted-foreground hover:text-foreground"
								>
									Cookie Policy
								</Link>
							</li>
						</ul>
					</div>
				</div>
				<div className="mt-12 pt-8 border-t border-border text-center text-sm text-muted-foreground">
					© {new Date().getFullYear()} Recally. All rights reserved.
				</div>
			</div>
		</footer>
	);
}
