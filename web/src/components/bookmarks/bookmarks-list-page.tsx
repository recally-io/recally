import AddBookmarkModal from "@/components/bookmarks/add-bookmark";
import BookmarkList from "@/components/bookmarks/bookmarks-list";
import type { SearchToken } from "@/components/bookmarks/search";
import type { BookmarkSearch, View } from "@/components/bookmarks/types";
import { SidebarComponent } from "@/components/sidebar/sidebar";
import { SidebarHeaderTrigger } from "@/components/sidebar/trigger";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import { useBookmarks } from "@/lib/apis/bookmarks";
import { useRouter } from "@tanstack/react-router";
import { List, Loader2, Table } from "lucide-react";
import { useEffect, useState } from "react";

const VIEW_STORAGE_KEY = "bookmarks-view-preference";

function getStoredView(): View {
	if (typeof window === "undefined") return "grid";
	return (localStorage.getItem(VIEW_STORAGE_KEY) as View) || "grid";
}

export default function BookmarksListView({
	search,
}: {
	search: BookmarkSearch;
}) {
	const limit = 12; // max 2 columns
	const offset = (search.page - 1) * limit;

	const router = useRouter();
	const { data, isLoading } = useBookmarks(
		limit,
		offset,
		search.filters,
		search.query,
	);
	const bookmarks = data?.bookmarks ?? [];
	const total = data?.total ?? 0;

	const [view, setView] = useState<View>(getStoredView());

	useEffect(() => {
		localStorage.setItem(VIEW_STORAGE_KEY, view);
	}, [view]);

	const handlePageChange = (page: number) => {
		router.navigate({
			to: ".",
			search: (prev) => ({ ...prev, page: page }),
		});
	};

	const handleSearch = (tokens: SearchToken[], query: string) => {
		router.navigate({
			to: ".",
			search: (prev) => ({
				...prev,
				query: query,
				page: 1,
				filters: tokens.map((t) => `${t.type}:${t.value}`),
			}),
		});
	};

	return (
		<SidebarProvider defaultOpen={false}>
			<SidebarComponent />
			<SidebarInset>
				<div className="flex flex-col h-full">
					<header className="container mx-auto flex h-16 shrink-0 items-center justify-between gap-2 px-4 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
						<SidebarHeaderTrigger />
						<div className="flex items-center">
							<AddBookmarkModal />
							<ToggleGroup
								type="single"
								value={view}
								onValueChange={(value) => setView(value as View)}
								size="sm"
							>
								<ToggleGroupItem value="grid" aria-label="Grid">
									<Table />
								</ToggleGroupItem>
								<ToggleGroupItem value="list" aria-label="List">
									<List />
								</ToggleGroupItem>
							</ToggleGroup>
						</div>
					</header>
					<main>
						{isLoading ? (
							<div className="flex items-center justify-center h-full">
								<Loader2 className="size-8 animate-spin" />
							</div>
						) : (
							<BookmarkList
								bookmarks={bookmarks}
								total={total}
								view={view}
								search={search}
								currentPage={search.page}
								onPageChange={handlePageChange}
								onSearch={handleSearch}
								itemsPerPage={limit}
							/>
						)}
					</main>
				</div>
			</SidebarInset>
		</SidebarProvider>
	);
}
