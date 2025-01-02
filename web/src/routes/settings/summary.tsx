import { SettingsPageComponenrt } from "@/components/settings/settings";
import { SummarySettings } from "@/components/settings/summary";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/settings/summary")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<SettingsPageComponenrt>
			<SummarySettings />
		</SettingsPageComponenrt>
	);
}
