import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "@/components/ui/popover";
import { Switch } from "@/components/ui/switch";
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table";
import { useToast } from "@/hooks/use-toast";
import type { CreateApiKeyInput } from "@/lib/apis/auth";
import { useApiKeys, useApiKeysMutations } from "@/lib/apis/auth";
import { useState } from "react";

export function ApiKeysSettings() {
	const [name, setName] = useState("");
	const [prefix, setPrefix] = useState("");
	const [neverExpires, setNeverExpires] = useState(false);
	const [expiresAt, setExpiresAt] = useState<Date>(() => {
		const date = new Date();
		date.setDate(date.getDate() + 90);
		return date;
	});
	const [open, setOpen] = useState(false);
	const { toast } = useToast();
	const { createApiKey, deleteApiKey } = useApiKeysMutations();
	const { keys } = useApiKeys();

	const formatExpirationText = (date?: Date | null) => {
		if (!date) return "Set expiration";
		if (date.getFullYear() === 9999) return "Never expires";
		return date.toLocaleDateString(undefined, {
			year: "numeric",
			month: "long",
			day: "numeric",
		});
	};

	const handleCreate = async () => {
		try {
			const input: CreateApiKeyInput = {
				name,
				prefix,
				expires_at: neverExpires ? new Date("9999-12-31") : expiresAt,
			};
			const newKey = await createApiKey(input);
			toast({
				description: "API key created successfully!",
				duration: 2000,
			});
			setName("");
			setPrefix("");
			setExpiresAt(new Date(Date.now() + 90 * 24 * 60 * 60 * 1000));
			setNeverExpires(false);
			setOpen(false);
			return newKey;
		} catch (error) {
			toast({
				variant: "destructive",
				description: "Failed to create API key",
				duration: 3000,
			});
		}
	};

	const handleDelete = async (id: string) => {
		try {
			await deleteApiKey(id);
			toast({
				description: "API key deleted successfully!",
				duration: 2000,
			});
		} catch (error) {
			toast({
				variant: "destructive",
				description: "Failed to delete API key",
				duration: 3000,
			});
		}
	};

	const handleDateSelect = (date: Date | undefined) => {
		if (date) {
			setExpiresAt(date);
		}
	};

	return (
		<div className="flex-1 space-y-8 max-w-4xl">
			<Card>
				<CardHeader>
					<div className="flex justify-between items-center">
						<CardTitle>API Keys</CardTitle>
						<Dialog open={open} onOpenChange={setOpen}>
							<DialogTrigger asChild>
								<Button>Create New Key</Button>
							</DialogTrigger>
							<DialogContent>
								<DialogHeader>
									<DialogTitle>Create New API Key</DialogTitle>
								</DialogHeader>
								<div className="space-y-4 py-4">
									<div className="grid gap-2">
										<Label htmlFor="name">Key Name</Label>
										<Input
											id="name"
											placeholder="Enter key name"
											value={name}
											onChange={(e) => setName(e.target.value)}
										/>
									</div>
									<div className="grid gap-2">
										<Label htmlFor="prefix">Key Prefix</Label>
										<Input
											id="prefix"
											placeholder="Enter key prefix (optional)"
											value={prefix}
											onChange={(e) => setPrefix(e.target.value)}
										/>
									</div>
									<div className="grid gap-2">
										<div className="flex items-center justify-between">
											<Label htmlFor="never-expires">Never Expires</Label>
											<Switch
												id="never-expires"
												checked={neverExpires}
												onCheckedChange={setNeverExpires}
											/>
										</div>
									</div>
									{!neverExpires && (
										<div className="grid gap-2">
											<Label htmlFor="expires">Expiration Date</Label>
											<Popover>
												<PopoverTrigger asChild>
													<Button
														variant="outline"
														className={
															!expiresAt ? "text-muted-foreground" : ""
														}
													>
														{expiresAt
															? formatExpirationText(expiresAt)
															: "Pick a date"}
													</Button>
												</PopoverTrigger>
												<PopoverContent className="w-auto p-0">
													<Calendar
														mode="single"
														selected={expiresAt}
														onSelect={handleDateSelect}
														disabled={(date) => date < new Date()}
														initialFocus
													/>
												</PopoverContent>
											</Popover>
										</div>
									)}
									<Button onClick={handleCreate}>Create</Button>
								</div>
							</DialogContent>
						</Dialog>
					</div>
				</CardHeader>
				<CardContent>
					<Table>
						<TableHeader>
							<TableRow>
								<TableHead>Name</TableHead>
								<TableHead>Prefix</TableHead>
								<TableHead>Expiration</TableHead>
								<TableHead>Status</TableHead>
								<TableHead className="text-right">Actions</TableHead>
							</TableRow>
						</TableHeader>
						<TableBody>
							{keys?.map((key) => (
								<TableRow key={key.id}>
									<TableCell className="font-medium">{key.name}</TableCell>
									<TableCell>{key.prefix}</TableCell>
									<TableCell>
										{key.expires_at
											? formatExpirationText(new Date(key.expires_at))
											: "Never"}
									</TableCell>
									<TableCell>
										<Badge
											variant={
												key.expires_at && new Date(key.expires_at) < new Date()
													? "destructive"
													: "default"
											}
										>
											{key.expires_at && new Date(key.expires_at) < new Date()
												? "Expired"
												: "Active"}
										</Badge>
									</TableCell>
									<TableCell className="text-right">
										<Button
											variant="destructive"
											size="sm"
											onClick={() => handleDelete(key.id)}
										>
											Delete
										</Button>
									</TableCell>
								</TableRow>
							))}
						</TableBody>
					</Table>
				</CardContent>
			</Card>
		</div>
	);
}
