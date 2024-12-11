import AddBookmarkModal from "@/components/AddBookmarkModal";
import BookmarkGrid from "@/components/BookmarkGrid";
import BookmarkList from "@/components/BookmarkList";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Grid, List, PlusCircle } from "lucide-react";
import { useTheme } from "next-themes";
import { useState } from "react";

import { Bookmark as BookmarkType } from "@/types/bookmark";

// This would typically come from your backend
const initialBookmarks: BookmarkType[] = [
  {
    id: 1,
    title: "React Documentation",
    url: "https://reactjs.org",
    tags: ["react", "docs"],
    content:
      "<p>React is a JavaScript library for building user interfaces. Learn what React is all about on our homepage or in the tutorial.</p>",
    image: "",
    summary:
      "React is a popular JavaScript library for building user interfaces, particularly single page applications.",
    dateAdded: "2022-01-01",
  },
  {
    id: 2,
    title: "Tailwind CSS",
    url: "https://tailwindcss.com",
    tags: ["css", "styling"],
    content:
      "<p>Tailwind CSS is a utility-first CSS framework packed with classes like flex, pt-4, text-center and rotate-90 that can be composed to build any design, directly in your markup.</p>",
    image: "",
    summary:
      "Tailwind CSS is a highly customizable, low-level CSS framework that gives you all of the building blocks you need to build bespoke designs.",
    dateAdded: "2022-01-02",
  },
  {
    id: 3,
    title: "shadcn/ui",
    url: "https://ui.shadcn.com",
    tags: ["ui", "components"],
    content:
      "<p>Beautifully designed components that you can copy and paste into your apps. Accessible. Customizable. Open Source.</p>",
    image: "",
    summary:
      "shadcn/ui provides a set of re-usable components that you can copy and paste into your apps, offering both functionality and customizable styling.",
    dateAdded: "2022-01-03",
  },
];

export default function HomePage() {
  const [bookmarks, setBookmarks] = useState<BookmarkType[]>(initialBookmarks);
  const [view, setView] = useState<"list" | "grid">("list");
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const { setTheme, theme } = useTheme();

  const filteredBookmarks = bookmarks.filter(
    (bookmark) =>
      bookmark.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      bookmark.tags.some((tag) =>
        tag.toLowerCase().includes(searchTerm.toLowerCase()),
      ),
  );

  const addBookmark = (newBookmark: {
    title: string;
    url: string;
    tags: string[];
  }) => {
    const completeBookmark: Omit<BookmarkType, "id"> = {
      ...newBookmark,
      content: "",
      summary: `Bookmark for ${newBookmark.title}`,
      highlights: [],
    };
    setBookmarks([
      ...bookmarks,
      { ...completeBookmark, id: bookmarks.length + 1 },
    ]);
  };

  return (
    <div className="container mx-auto p-4 max-w-4xl">
      <div className="mb-4 flex gap-4">
        <Input
          type="text"
          placeholder="Search bookmarks..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="flex-grow"
        />
        <Button
          variant="outline"
          onClick={() => setView("list")}
          aria-label="List view"
        >
          <List className="h-4 w-4" />
        </Button>
        <Button
          variant="outline"
          onClick={() => setView("grid")}
          aria-label="Grid view"
        >
          <Grid className="h-4 w-4" />
        </Button>
        <Button onClick={() => setIsModalOpen(true)}>
          <PlusCircle className="h-4 w-4" />
        </Button>
      </div>

      {view === "list" ? (
        <BookmarkList bookmarks={filteredBookmarks} />
      ) : (
        <BookmarkGrid bookmarks={filteredBookmarks} />
      )}

      <AddBookmarkModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onAdd={addBookmark}
      />
    </div>
  );
}
