import { ProfileSettings } from "@/components/settings/profile";
import { SettingsPageComponenrt } from "@/components/settings/settings";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/settings/profile")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<SettingsPageComponenrt>
			<ProfileSettings />
		</SettingsPageComponenrt>
	);
}
