import { Button } from "@/components/ui/button";
import { ROUTES } from "@/lib/router";
import { Link } from "@tanstack/react-router";
import { BookmarkPlus } from "lucide-react";

export const AuthBanner = () => {
	return (
		<div className="sticky top-0 z-10 border-b border-border bg-gradient-to-r from-primary/5 via-background to-primary/5">
			<div className="container mx-auto max-w-4xl flex items-center justify-between py-4 px-4 animate-fade-in">
				<div className="flex items-center gap-3">
					<div className="hidden sm:flex">
						<BookmarkPlus className="h-5 w-5 text-primary animate-bounce" />
					</div>
					<div className="space-y-0.5">
						<p className="text-sm font-medium">Want to save this for later?</p>
						<p className="text-xs text-muted-foreground">
							Join Recally to create your personal reading collection
						</p>
					</div>
				</div>
				<div className="space-x-2 shrink-0">
					<Button variant="outline" size="sm" asChild>
						<Link to={ROUTES.AUTH_LOGIN}>Log in</Link>
					</Button>
					<Button size="sm" className="bg-primary" asChild>
						<Link to={ROUTES.AUTH_REGISTER}>Sign up free</Link>
					</Button>
				</div>
			</div>
		</div>
	);
};
