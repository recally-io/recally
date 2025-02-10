import { AISettings } from "@/components/settings/ai";
import { SettingsPageComponenrt } from "@/components/settings/settings";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/settings/ai")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<SettingsPageComponenrt>
			<AISettings />
		</SettingsPageComponenrt>
	);
}
