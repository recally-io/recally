import { ProfileSettings } from "@/components/settings/profile";
import { SettingsPageComponenrt } from "@/components/settings/settings";
import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/settings/profile")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<SettingsPageComponenrt>
			<ProfileSettings />
		</SettingsPageComponenrt>
	);
}
