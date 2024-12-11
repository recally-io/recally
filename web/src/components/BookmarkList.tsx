import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Bookmark } from "@/types/bookmark"
import { ExternalLink, Highlighter } from 'lucide-react'
import { Link } from 'react-router-dom'

interface BookmarkListProps {
  bookmarks: Bookmark[]
}

export default function BookmarkList({ bookmarks }: BookmarkListProps) {
  return (
    <div className="space-y-4">
      {bookmarks.map((bookmark) => (
        <Link to={`/bookmarks/${bookmark.id}`} key={bookmark.id}>
          <Card 
            className="cursor-pointer hover:shadow-md transition-all duration-300 ease-in-out transform hover:-translate-y-1 overflow-hidden"
          >
            <div className="flex">
              {bookmark.image && (
                <div className="w-1/4 min-w-[100px]">
                  <img 
                    src={bookmark.image} 
                    alt={`Thumbnail for ${bookmark.title}`} 
                    className="w-full h-full object-cover"
                  />
                </div>
              )}
              <div className={`flex-1 ${bookmark.image ? 'w-3/4' : 'w-full'}`}>
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
                  <CardDescription className="truncate">{bookmark.url}</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex flex-wrap gap-2">
                    {bookmark.tags.map((tag) => (
                      <Badge key={tag} variant="secondary">{tag}</Badge>
                    ))}
                  </div>
                </CardContent>
              </div>
            </div>
          </Card>
        </Link>
      ))}
    </div>
  )
}

