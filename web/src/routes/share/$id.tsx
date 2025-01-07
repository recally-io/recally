import { SharedArticleReader } from "@/components/shared/content";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/share/$id")({
	component: RouteComponent,
});

function RouteComponent() {
	const { id } = Route.useParams();
	return <SharedArticleReader id={id} />;
}
