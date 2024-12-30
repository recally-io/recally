import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useToast } from "@/hooks/use-toast";
import { useUser } from "@/lib/apis/auth";
import Cookies from "js-cookie";
import { useState } from "react";

export function ProfileSettings() {
	const { user } = useUser();
	const [showLinkingInfo, setShowLinkingInfo] = useState(false);
	const token = Cookies.get("token");
	const { toast } = useToast();
	const [isCopied, setIsCopied] = useState(false);

	const handleTelegramLink = () => {
		setShowLinkingInfo(true);
		// window.open(BOT_URL, '_blank');
	};

	const handleCopyToken = async () => {
		try {
			await navigator.clipboard.writeText(`/linkaccount ${token}`);
			setIsCopied(true);
			toast({
				description: "Token copied to clipboard!",
				duration: 2000,
			});
			setTimeout(() => setIsCopied(false), 2000);
		} catch (err) {
			toast({
				variant: "destructive",
				description: "Failed to copy token.",
			});
		}
	};

	return (
		<div className="flex-1 space-y-8 max-w-2xl">
			<Card>
				<CardContent className="pt-6">
					<div className="space-y-6">
						<div className="flex items-center gap-6">
							<Avatar className="h-20 w-20">
								<AvatarImage src={user?.avatar} />
								<AvatarFallback>V</AvatarFallback>
							</Avatar>
							<Button variant="outline">Change avatar</Button>
						</div>

						<div className="grid gap-4">
							<div className="grid gap-2">
								<Label htmlFor="name">Username</Label>
								<Input id="name" defaultValue={user?.username} />
							</div>

							<div className="grid gap-2">
								<Label htmlFor="email">Email</Label>
								<Input id="email" type="email" defaultValue={user?.email} />
							</div>
						</div>
					</div>
				</CardContent>
			</Card>

			<Card>
				<CardContent className="pt-6">
					<h2 className="text-lg font-medium mb-4">Change Password</h2>
					<div className="space-y-4">
						<div className="grid gap-2">
							<Label htmlFor="current-password">Current Password</Label>
							<Input id="current-password" type="password" />
						</div>

						<div className="grid gap-2">
							<Label htmlFor="new-password">New Password</Label>
							<Input id="new-password" type="password" />
						</div>

						<div className="grid gap-2">
							<Label htmlFor="confirm-password">Confirm New Password</Label>
							<Input id="confirm-password" type="password" />
						</div>
						{/* TODO: actual update */}
						<Button className="mt-4">Update Password</Button>
					</div>
				</CardContent>
			</Card>

			<Card>
				<CardContent className="pt-6">
					<h2 className="text-lg font-medium mb-4">Link Telegram Bot</h2>
					<div className="space-y-4">
						<Button onClick={handleTelegramLink}>Link Telegram</Button>
						{showLinkingInfo && (
							<div className="mt-4 p-4 bg-muted rounded-lg space-y-4">
								<p className="text-sm">Please follow these steps:</p>
								<ol className="list-decimal list-inside space-y-2">
									<li>Open the RecallyReader telegram bot</li>
									<li>Send this code to the bot:</li>
									<div className="relative">
										<code className="block mt-2 font-mono text-lg p-2 bg-background break-all">{`/linkaccount ${token}`}</code>
										<Button
											variant="outline"
											size="sm"
											className="absolute top-2 right-2"
											onClick={handleCopyToken}
										>
											{isCopied ? "Copied!" : "Copy"}
										</Button>
									</div>
								</ol>
							</div>
						)}
					</div>
				</CardContent>
			</Card>

			<div className="flex justify-end">
				{/* TODO: actual update */}
				<Button>Save Changes</Button>
			</div>
		</div>
	);
}
