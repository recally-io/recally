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

const SUPPORTED_LANGUAGES = [
	{ id: "en", name: "English" },
	{ id: "zh", name: "Chinese" },
	{ id: "es", name: "Spanish" },
	{ id: "fr", name: "French" },
	{ id: "de", name: "German" },
] as const;

type SettingsKey = "summary_options" | "describe_image_options";

type AISettingsData = {
	model: string;
	prompt: string;
	language: string;
};

interface AISettingsProps {
	title: string;
	settingsKey: SettingsKey;
	promptPlaceholder?: string;
}

const AI_FEATURES: AISettingsProps[] = [
	{
		title: "Summary Settings",
		settingsKey: "summary_options",
		promptPlaceholder: "Enter your custom summarization prompt...",
	},
	{
		title: "Describe Image Settings",
		settingsKey: "describe_image_options",
		promptPlaceholder: "Enter your custom image description prompt...",
	},
];

export function AISettings() {
	const { updateSettings } = useUsers();
	const { user } = useUser();
	const { toast } = useToast();

	const [settings, setSettings] = useState<Record<SettingsKey, AISettingsData>>(
		{
			summary_options: {
				model: user?.settings?.summary_options?.model ?? "",
				prompt: user?.settings?.summary_options?.prompt ?? "",
				language: user?.settings?.summary_options?.language ?? "",
			},
			describe_image_options: {
				model: user?.settings?.describe_image_options?.model ?? "",
				prompt: user?.settings?.describe_image_options?.prompt ?? "",
				language: user?.settings?.describe_image_options?.language ?? "",
			},
		},
	);

	const handleSettingChange = (
		key: SettingsKey,
		field: keyof AISettingsData,
		value: string,
	) => {
		setSettings((prev) => ({
			...prev,
			[key]: { ...prev[key], [field]: value },
		}));
	};

	const showErrorToast = (error: unknown) => {
		const errorMessage =
			error instanceof Error ? error.message : "Unknown error occurred";
		toast({
			variant: "destructive",
			description: `Failed to save settings: ${errorMessage}`,
			duration: 3000,
		});
	};

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
			const updatedSettings = {
				...user.settings,
				...settings,
			};

			await updateSettings(user.id, updatedSettings);
			toast({
				description: "Settings saved successfully!",
				duration: 2000,
			});
		} catch (error) {
			console.error("Failed to save settings:", error);
			showErrorToast(error);
		}
	};

	return (
		<div className="space-y-8">
			{AI_FEATURES.map((feature) => (
				<AISetting
					key={feature.settingsKey}
					{...feature}
					settings={settings[feature.settingsKey]}
					onSettingChange={(field, value) =>
						handleSettingChange(feature.settingsKey, field, value)
					}
				/>
			))}
			<div className="flex justify-end max-w-2xl">
				<Button onClick={handleSave}>Save All Changes</Button>
			</div>
		</div>
	);
}

function AISetting({
	title,
	settingsKey,
	promptPlaceholder,
	settings,
	onSettingChange,
}: AISettingsProps & {
	settings: AISettingsData;
	onSettingChange: (field: keyof AISettingsData, value: string) => void;
}) {
	const { models } = useLLMs();

	return (
		<div className="flex-1 space-y-8 max-w-2xl">
			<Card>
				<CardContent className="pt-6">
					<h2 className="text-lg font-medium mb-4">{title}</h2>
					<div className="space-y-6">
						<div className="grid gap-2">
							<Label htmlFor={`${settingsKey}-model`}>Model</Label>
							<Select
								value={settings.model}
								onValueChange={(value) => onSettingChange("model", value)}
							>
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
							<Label htmlFor={`${settingsKey}-prompt`}>Custom Prompt</Label>
							<Textarea
								id={`${settingsKey}-prompt`}
								placeholder={promptPlaceholder}
								value={settings.prompt}
								onChange={(e) => onSettingChange("prompt", e.target.value)}
							/>
						</div>

						<div className="grid gap-2">
							<Label htmlFor={`${settingsKey}-language`}>Language</Label>
							<Select
								value={settings.language}
								onValueChange={(value) => onSettingChange("language", value)}
							>
								<SelectTrigger>
									<SelectValue placeholder="Select language" />
								</SelectTrigger>
								<SelectContent>
									{SUPPORTED_LANGUAGES.map((l) => (
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
		</div>
	);
}
