import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ExternalLink, Highlighter } from 'lucide-react'

interface Bookmark {
  id: number
  title: string
  url: string
  tags: string[]
  highlights?: any[]
  image?: string
}

interface BookmarkGridProps {
  bookmarks: Bookmark[]
  onBookmarkClick: (id: number) => void
}

export default function BookmarkGrid({ bookmarks, onBookmarkClick }: BookmarkGridProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {bookmarks.map((bookmark) => (
        <Card 
          key={bookmark.id} 
          className="flex flex-col h-full cursor-pointer hover:shadow-md transition-all duration-300 ease-in-out transform hover:-translate-y-1 overflow-hidden"
          onClick={() => onBookmarkClick(bookmark.id)}
        >
          {bookmark.image && (
            <div className="h-48 overflow-hidden">
              <img 
                src={bookmark.image} 
                alt={`Thumbnail for ${bookmark.title}`} 
                className="w-full h-full object-cover transition-transform duration-300 ease-in-out transform hover:scale-105"
              />
            </div>
          )}
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              <span className="flex items-center gap-2 truncate">
                {bookmark.title}
                {bookmark.highlights && bookmark.highlights.length > 0 && (
                  <Highlighter className="h-4 w-4 text-yellow-500 flex-shrink-0" />
                )}
              </span>
              <a 
                href={bookmark.url} 
                target="_blank" 
                rel="noopener noreferrer" 
                className="text-blue-500 hover:text-blue-700 transition-colors"
                onClick={(e) => e.stopPropagation()}
              >
                <ExternalLink className="h-4 w-4" />
              </a>
            </CardTitle>
          </CardHeader>
          <CardContent className="flex-grow">
            <div className="flex flex-wrap gap-2 mt-2">
              {bookmark.tags.map((tag) => (
                <Badge key={tag} variant="secondary" className="transition-colors hover:bg-primary hover:text-primary-foreground">
                  {tag}
                </Badge>
              ))}
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}

