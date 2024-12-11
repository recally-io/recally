import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { ArrowLeft, Bookmark, Grid, List, Moon, PlusCircle, Sun } from 'lucide-react'
import { useTheme } from "next-themes"
import { useCallback, useState } from 'react'
import AddBookmarkModal from './components/AddBookmarkModal'
import BookmarkDetail from './components/BookmarkDetail'
import BookmarkGrid from './components/BookmarkGrid'
import BookmarkList from './components/BookmarkList'
import { ThemeProvider } from "./components/theme-provider"

interface Highlight {
  id: string
  text: string
  startOffset: number
  endOffset: number
  note?: string
}

interface BookmarkType {
  id: number
  title: string
  url: string
  tags: string[]
  content: string
  image?: string
  summary: string
  highlights?: Highlight[]
}

// This would typically come from your backend
const initialBookmarks: BookmarkType[] = [
  { 
    id: 1, 
    title: 'React Documentation', 
    url: 'https://reactjs.org', 
    tags: ['react', 'docs'],
    content: '<p>React is a JavaScript library for building user interfaces. Learn what React is all about on our homepage or in the tutorial.</p>',
    image: '/placeholder.svg?height=300&width=400',
    summary: 'React is a popular JavaScript library for building user interfaces, particularly single page applications.'
  },
  { 
    id: 2, 
    title: 'Tailwind CSS', 
    url: 'https://tailwindcss.com', 
    tags: ['css', 'styling'],
    content: '<p>Tailwind CSS is a utility-first CSS framework packed with classes like flex, pt-4, text-center and rotate-90 that can be composed to build any design, directly in your markup.</p>',
    image: '/placeholder.svg?height=300&width=400',
    summary: 'Tailwind CSS is a highly customizable, low-level CSS framework that gives you all of the building blocks you need to build bespoke designs.'
  },
  { 
    id: 3, 
    title: 'shadcn/ui', 
    url: 'https://ui.shadcn.com', 
    tags: ['ui', 'components'],
    content: '<p>Beautifully designed components that you can copy and paste into your apps. Accessible. Customizable. Open Source.</p>',
    image: '/placeholder.svg?height=300&width=400',
    summary: 'shadcn/ui provides a set of re-usable components that you can copy and paste into your apps, offering both functionality and customizable styling.'
  },
]

function AppContent() {
  const [bookmarks, setBookmarks] = useState<BookmarkType[]>(initialBookmarks)
  const [view, setView] = useState<'list' | 'grid'>('list')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedBookmark, setSelectedBookmark] = useState<number | null>(null)
  const { setTheme, theme } = useTheme()

  const filteredBookmarks = bookmarks.filter(bookmark => 
    bookmark.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
    bookmark.tags.some(tag => tag.toLowerCase().includes(searchTerm.toLowerCase()))
  )

  const addBookmark = (newBookmark: { title: string; url: string; tags: string[] }) => {
    const completeBookmark: Omit<BookmarkType, 'id'> = {
      ...newBookmark,
      content: '',
      summary: `Bookmark for ${newBookmark.title}`,
      highlights: []
    }
    setBookmarks([...bookmarks, { ...completeBookmark, id: bookmarks.length + 1 }])
  }

  const openBookmark = (id: number) => {
    setSelectedBookmark(id)
  }

  const updateBookmarkHighlights = useCallback((id: number, highlights: Highlight[]) => {
    setBookmarks(prevBookmarks => prevBookmarks.map(bookmark => 
      bookmark.id === id ? { ...bookmark, highlights } : bookmark
    ))
  }, [])

  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-secondary/20">
      <div className="container mx-auto p-4 max-w-4xl">
        <header className="flex justify-between items-center mb-6 sticky top-0 bg-background/80 backdrop-blur-sm z-10 py-4">
          <h1 className="text-2xl font-bold flex items-center gap-2">
            <Bookmark className="h-6 w-6" />
            Bookmark App
          </h1>
          <div className="flex items-center gap-2">
            <Button variant="ghost" size="icon" onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}>
              {theme === 'dark' ? <Sun className="h-[1.2rem] w-[1.2rem]" /> : <Moon className="h-[1.2rem] w-[1.2rem]" />}
            </Button>
            <Button onClick={() => setIsModalOpen(true)}>
              <PlusCircle className="h-4 w-4 mr-2" />
              Add Bookmark
            </Button>
          </div>
        </header>

        {selectedBookmark === null ? (
          <>
            <div className="mb-4 flex gap-4">
              <Input
                type="text"
                placeholder="Search bookmarks..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="flex-grow"
              />
              <Button variant="outline" onClick={() => setView('list')} aria-label="List view">
                <List className="h-4 w-4" />
              </Button>
              <Button variant="outline" onClick={() => setView('grid')} aria-label="Grid view">
                <Grid className="h-4 w-4" />
              </Button>
            </div>

            {view === 'list' ? (
              <BookmarkList bookmarks={filteredBookmarks} onBookmarkClick={openBookmark} />
            ) : (
              <BookmarkGrid bookmarks={filteredBookmarks} onBookmarkClick={openBookmark} />
            )}
          </>
        ) : (
          <>
            <Button onClick={() => setSelectedBookmark(null)} className="mb-4">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to List
            </Button>
            <BookmarkDetail 
              bookmark={bookmarks.find(b => b.id === selectedBookmark)!}
              onUpdateBookmark={updateBookmarkHighlights}
            />
          </>
        )}

        <AddBookmarkModal
          isOpen={isModalOpen}
          onClose={() => setIsModalOpen(false)}
          onAdd={addBookmark}
        />
      </div>
    </div>
  )
}

export default function App() {
  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <AppContent />
    </ThemeProvider>
  )
}

