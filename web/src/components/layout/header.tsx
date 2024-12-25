import { Button } from "@/components/ui/button";
import { useUser } from "@/lib/apis/auth";
import { ROUTES } from "@/lib/router";
import { Link, useNavigate } from "react-router-dom";
import { UserNav } from "./user-nav";

export default function Header() {
	const { user } = useUser();
	const navigate = useNavigate();

	return (
		<div className="container mx-auto flex px-2 h-16 items-center justify-between">
			<div className="flex items-center space-x-4">
				<Link to={ROUTES.HOME} className="flex items-center space-x-2">
					<span className="font-bold text-xl">Vibrain</span>
				</Link>
				{/* <nav className="hidden md:flex items-center space-x-4">
					<Link
						to="/features"
						className="text-sm font-medium hover:text-primary"
					>
						Features
					</Link>
					<Link
						to="/pricing"
						className="text-sm font-medium hover:text-primary"
					>
						Pricing
					</Link>
					<Link to="/blog" className="text-sm font-medium hover:text-primary">
						Blog
					</Link>
				</nav> */}
			</div>
			<div className="flex items-center">
				{user ? (
					<UserNav />
				) : (
					<Button
						variant="default"
						size="sm"
						onClick={() => navigate(ROUTES.LOGIN)}
					>
						Sign In
					</Button>
				)}
			</div>
		</div>
	);
}
