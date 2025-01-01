import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/hooks/use-toast";
import { useUser } from "@/lib/apis/auth";
import { useLLMs } from "@/lib/apis/llm";
import { useUsers } from "@/lib/apis/users";
import { useState } from "react";

const supportedLanguages = [
	{ id: "en", name: "English" },
	{ id: "zh", name: "Chinese" },
	{ id: "es", name: "Spanish" },
	{ id: "fr", name: "French" },
	{ id: "de", name: "German" },
];

export function SummarySettings() {
	const { updateSettings } = useUsers();
	const { user } = useUser();
	const { models } = useLLMs();
	const [model, setModel] = useState(user?.settings?.summary_options?.model);
	const [prompt, setPrompt] = useState(user?.settings?.summary_options?.prompt);
	const [language, setLanguage] = useState(
		user?.settings?.summary_options?.language,
	);
	const { toast } = useToast();
	const handleSave = async () => {
		if (!user) {
			toast({
				variant: "destructive",
				description: "User not found",
				duration: 2000,
			});
			return;
		}

		try {
			const settings = {
				...user.settings,
				summary_options: {
					model,
					prompt,
					language,
				},
			};

			await updateSettings(user.id, settings);

			toast({
				description: "Settings saved successfully!",
				duration: 2000,
			});
		} catch (error) {
			console.error("Failed to save settings:", error);
			toast({
				variant: "destructive",
				description: `Failed to save settings: ${error instanceof Error ? error.message : "Unknown error"}`,
				duration: 3000,
			});
		}
	};

	return (
		<div className="flex-1 space-y-8 max-w-2xl">
			<Card>
				<CardContent className="pt-6">
					<h2 className="text-lg font-medium mb-4">Summary Settings</h2>
					<div className="space-y-6">
						<div className="grid gap-2">
							<Label htmlFor="model">Model</Label>
							<Select value={model} onValueChange={setModel}>
								<SelectTrigger>
									<SelectValue placeholder="Select model" />
								</SelectTrigger>
								<SelectContent>
									{models?.map((model) => (
										<SelectItem key={model.id} value={model.id}>
											{model.name}
										</SelectItem>
									))}
								</SelectContent>
							</Select>
						</div>

						<div className="grid gap-2">
							<Label htmlFor="prompt">Custom Prompt</Label>
							<Textarea
								id="prompt"
								placeholder="Enter your custom summarization prompt..."
								value={prompt}
								onChange={(e) => setPrompt(e.target.value)}
							/>
						</div>

						<div className="grid gap-2">
							<Label htmlFor="language">Summary Language</Label>
							<Select value={language} onValueChange={setLanguage}>
								<SelectTrigger>
									<SelectValue placeholder="Select language" />
								</SelectTrigger>
								<SelectContent>
									{supportedLanguages.map((l) => (
										<SelectItem key={l.id} value={l.name}>
											{l.name}
										</SelectItem>
									))}
								</SelectContent>
							</Select>
						</div>
					</div>
				</CardContent>
			</Card>

			<div className="flex justify-end">
				<Button onClick={handleSave}>Save Changes</Button>
			</div>
		</div>
	);
}
