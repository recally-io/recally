import { ApiKeysSettings } from "@/components/settings/api-keys";
import { SettingsPageComponenrt } from "@/components/settings/settings";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/settings/api-keys")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<SettingsPageComponenrt>
			<ApiKeysSettings />
		</SettingsPageComponenrt>
	);
}
