import { SettingsPageComponenrt } from "@/components/settings/settings";
import { SummarySettings } from "@/components/settings/ai";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/settings/ai")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<SettingsPageComponenrt>
			<SummarySettings />
		</SettingsPageComponenrt>
	);
}
