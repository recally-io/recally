import { SettingsPageComponenrt } from "@/components/settings/settings";
import { SummarySettings } from "@/components/settings/summary";
import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/settings/summary")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<SettingsPageComponenrt>
			<SummarySettings />
		</SettingsPageComponenrt>
	);
}
