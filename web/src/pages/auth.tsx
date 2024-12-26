import { Button } from "@/components/ui/button";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { useAuth } from "@/lib/apis/auth";
import { Github, Mail } from "lucide-react";
import { parseAsString, useQueryState } from "nuqs";
import type React from "react";
import { useState } from "react";
import { createRoot } from "react-dom/client";

import { ROUTES } from "@/lib/router";
import App from "./app-basic";

interface AuthFormData {
	email: string;
	password: string;
	confirmPassword?: string;
	name?: string;
}

const AuthMode = {
	Login: "login",
	Register: "register",
};

export default function AuthPage() {
	// "login" or "register"
	const [mode, _] = useQueryState(
		"mode",
		parseAsString.withDefault(AuthMode.Login),
	);
	const isLogin = mode === AuthMode.Login;

	const [formData, setFormData] = useState<AuthFormData>({
		email: "",
		password: "",
		...(isLogin ? {} : { confirmPassword: "", name: "" }),
	});

	const { login, register, oauthLogin } = useAuth();

	const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		const { name, value } = e.target;
		setFormData((prevState) => ({
			...prevState,
			[name]: value,
		}));
	};

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		try {
			if (isLogin) {
				await login({
					email: formData.email,
					password: formData.password,
				});
			} else {
				await register({
					email: formData.email,
					password: formData.password,
					username: formData.name || "",
				});
			}
			window.location.href = "/";
		} catch (error) {
			console.error(`${mode} failed:`, error);
		}
	};

	const handleOAuth = (provider: string) => {
		oauthLogin(provider);
	};

	return (
		<div className="flex items-center justify-center min-h-screen p-4">
			<Card className="w-full max-w-md">
				<CardHeader>
					<CardTitle>{isLogin ? "Login" : "Sign Up"}</CardTitle>
					<CardDescription>
						{isLogin ? "Welcome back" : "Create your account"}
					</CardDescription>
				</CardHeader>
				<CardContent>
					<div className="grid grid-cols-2 gap-4">
						<Button
							variant="outline"
							onClick={() => handleOAuth("Google")}
							className="w-full"
						>
							<svg
								className="mr-2 h-4 w-4"
								aria-hidden="true"
								focusable="false"
								data-prefix="fab"
								data-icon="google"
								role="img"
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 0 488 512"
							>
								<path
									fill="currentColor"
									d="M488 261.8C488 403.3 391.1 504 248 504 110.8 504 0 393.2 0 256S110.8 8 248 8c66.8 0 123 24.5 166.3 64.9l-67.5 64.9C258.5 52.6 94.3 116.6 94.3 256c0 86.5 69.1 156.6 153.7 156.6 98.2 0 135-70.4 140.8-106.9H248v-85.3h236.1c2.3 12.7 3.9 24.9 3.9 41.4z"
								></path>
							</svg>
							Google
						</Button>
						<Button
							variant="outline"
							onClick={() => handleOAuth("GitHub")}
							className="w-full"
						>
							<Github className="mr-2 h-4 w-4" />
							GitHub
						</Button>
					</div>
					<Separator className="my-4" />
					<form onSubmit={handleSubmit}>
						<div className="space-y-4">
							{!isLogin && (
								<div className="space-y-2">
									<Label htmlFor="name">Name</Label>
									<Input
										id="name"
										name="name"
										type="text"
										placeholder="John Doe"
										required
										value={formData.name || ""}
										onChange={handleInputChange}
									/>
								</div>
							)}
							<div className="space-y-2">
								<Label htmlFor="email">Email</Label>
								<Input
									id="email"
									name="email"
									type="email"
									placeholder="your@email.com"
									required
									value={formData.email}
									onChange={handleInputChange}
								/>
							</div>
							<div className="space-y-2">
								<Label htmlFor="password">Password</Label>
								<Input
									id="password"
									name="password"
									type="password"
									required
									value={formData.password}
									onChange={handleInputChange}
								/>
							</div>
							{!isLogin && (
								<div className="space-y-2">
									<Label htmlFor="confirmPassword">Confirm Password</Label>
									<Input
										id="confirmPassword"
										name="confirmPassword"
										type="password"
										required
										value={formData.confirmPassword || ""}
										onChange={handleInputChange}
									/>
								</div>
							)}
						</div>
						<Button type="submit" className="w-full mt-4">
							<Mail className="mr-2 h-4 w-4" />
							{isLogin ? "Log in" : "Sign up"} with Email
						</Button>
					</form>
					<div className="mt-4 text-center text-sm">
						<span className="text-muted-foreground">
							{isLogin
								? "Don't have an account? "
								: "Already have an account? "}
						</span>
						<a
							href={isLogin ? ROUTES.SIGNUP : ROUTES.SIGNUP}
							className="text-primary hover:underline"
						>
							{isLogin ? "Sign up" : "Log in"}
						</a>
					</div>
				</CardContent>
			</Card>
		</div>
	);
}

createRoot(document.getElementById("root")!).render(
	<App>
		<AuthPage />
	</App>,
);
