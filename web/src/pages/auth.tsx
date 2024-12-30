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
import { useAuth, useUser } from "@/lib/apis/auth";
import { ROUTES } from "@/lib/router";
import { SiGithub } from "@icons-pack/react-simple-icons";
import { Mail } from "lucide-react";
import { parseAsString, useQueryState } from "nuqs";
import type React from "react";
import { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
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

const OAuthProviders = [
	{
		name: "GitHub",
		icon: SiGithub,
	},
	// {
	// 	name: "Google",
	// 	icon: SiGoogle,
	// },
	// {
	// 	name: "Telegram",
	// 	icon: SiTelegram,
	// },
];

export default function AuthPage() {
	// "login" or "register"
	const [mode, _] = useQueryState(
		"mode",
		parseAsString.withDefault(AuthMode.Login),
	);
	const isLogin = mode === AuthMode.Login;

	// oauth callback state and code
	const [code] = useQueryState("code");
	const [state] = useQueryState("state");

	// email and password login form data
	const [formData, setFormData] = useState<AuthFormData>({
		email: "",
		password: "",
		...(isLogin ? {} : { confirmPassword: "", name: "" }),
	});

	const { login, register, oauthLogin, oauthCallback } = useAuth();

	const { user } = useUser();

	useEffect(() => {
		const handleCallback = async () => {
			if (state && code) {
				console.log("OAuth Callback:", state, code);
				await handleOAuthCallback();
			}
		};
		handleCallback();
	}, [state, code]);

	if (user) {
		window.location.href = ROUTES.HOME;
		return null;
	}

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

	const handleOAuthLogin = async (provider: string) => {
		const resp = await oauthLogin(provider);
		console.log("OAuth Login response:", resp);
		window.location.href = resp.url;
	};

	const handleOAuthCallback = async () => {
		if (state && code) {
			const provider = state.split(":")[1];
			console.log("OAuth Callback:", provider, code);
			const data = await oauthCallback(provider.toLowerCase(), code);
			console.log("OAuth Callback data:", data);
			// TODO: redirect to previous page
			window.location.href = "/";
		}
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
						{OAuthProviders.map((provider) => {
							return (
								<Button
									key={provider.name}
									variant="outline"
									onClick={async () => await handleOAuthLogin(provider.name)}
									className="w-full"
								>
									<provider.icon className="mr-2 h-4 w-4" />
									{provider.name}
								</Button>
							);
						})}
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
