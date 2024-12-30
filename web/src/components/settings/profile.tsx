import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useUser } from "@/lib/apis/auth";

export function ProfileSettings() {
	const { user } = useUser();

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

			<div className="flex justify-end">
				{/* TODO: actual update */}
				<Button>Save Changes</Button>
			</div>
		</div>
	);
}
