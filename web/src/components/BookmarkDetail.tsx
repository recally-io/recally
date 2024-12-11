import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Bookmark, Highlight } from "@/types/bookmark"
import { Calendar, ExternalLink, X } from 'lucide-react'
import { useEffect, useMemo, useRef, useState } from 'react'
import { v4 as uuidv4 } from 'uuid'


interface BookmarkDetailProps {
  bookmark: Bookmark
  onUpdateBookmark: (id: number, highlights: Highlight[]) => void
}

export default function BookmarkDetail({ bookmark, onUpdateBookmark }: BookmarkDetailProps) {
  const [highlights, setHighlights] = useState<Highlight[]>([])
  const [isHighlighting, setIsHighlighting] = useState(false)
  const contentRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const savedHighlights = localStorage.getItem(`highlights-${bookmark.id}`)
    if (savedHighlights) {
      setHighlights(JSON.parse(savedHighlights))
    }
  }, [bookmark.id])

  useEffect(() => {
    const highlightsJson = JSON.stringify(highlights)
    localStorage.setItem(`highlights-${bookmark.id}`, highlightsJson)
    onUpdateBookmark(bookmark.id, highlights)
  }, [highlights, bookmark.id, onUpdateBookmark])

  const handleHighlight = () => {
    const selection = window.getSelection()
    if (selection && !selection.isCollapsed && contentRef.current) {
      const range = selection.getRangeAt(0)
      const startOffset = range.startOffset
      const endOffset = range.endOffset
      const text = selection.toString()

      const newHighlight: Highlight = {
        id: uuidv4(),
        text,
        startOffset,
        endOffset,
      }

      setHighlights([...highlights, newHighlight])
    }
  }

  const removeHighlight = (id: string) => {
    setHighlights(highlights.filter(h => h.id !== id))
  }

  const highlightedContent = useMemo(() => {
    let content = bookmark.content
    highlights.forEach(highlight => {
      const before = content.slice(0, highlight.startOffset)
      const highlighted = content.slice(highlight.startOffset, highlight.endOffset)
      const after = content.slice(highlight.endOffset)
      content = `${before}<mark class="bg-yellow-200 dark:bg-yellow-800">${highlighted}</mark>${after}`
    })
    return content
  }, [bookmark.content, highlights])

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold">{bookmark.title}</h1>
          <p className="text-muted-foreground mt-1">
            <a 
              href={bookmark.url} 
              target="_blank" 
              rel="noopener noreferrer" 
              className="flex items-center hover:underline"
            >
              {bookmark.url}
              <ExternalLink className="h-4 w-4 ml-1" />
            </a>
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Calendar className="h-4 w-4 text-muted-foreground" />
          <span className="text-sm text-muted-foreground">
            Added on {bookmark.dateAdded ? new Date(bookmark.dateAdded).toLocaleDateString() : 'Unknown date'}
          </span>
        </div>
      </div>

      <div className="flex flex-wrap gap-2">
        {bookmark.tags.map((tag) => (
          <Badge key={tag} variant="secondary">{tag}</Badge>
        ))}
      </div>

      {bookmark.image && (
        <img 
          src={bookmark.image} 
          alt={`Thumbnail for ${bookmark.title}`} 
          className="w-full h-64 object-cover rounded-lg"
        />
      )}

      <Tabs defaultValue="content" className="w-full">
        <TabsList>
          <TabsTrigger value="content">Content</TabsTrigger>
          <TabsTrigger value="summary">Summary</TabsTrigger>
          <TabsTrigger value="highlights">Highlights</TabsTrigger>
        </TabsList>
        <TabsContent value="content">
          <Card>
            <CardHeader>
              <CardTitle>Content</CardTitle>
              <CardDescription>
                The full content of the bookmarked page.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex justify-end mb-4">
                <Button onClick={() => setIsHighlighting(!isHighlighting)}>
                  {isHighlighting ? 'Finish Highlighting' : 'Start Highlighting'}
                </Button>
              </div>
              <ScrollArea className="h-[60vh]">
                <div 
                  ref={contentRef}
                  className="prose dark:prose-invert max-w-none"
                  dangerouslySetInnerHTML={{ __html: highlightedContent }}
                  onMouseUp={isHighlighting ? handleHighlight : undefined}
                />
              </ScrollArea>
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="summary">
          <Card>
            <CardHeader>
              <CardTitle>Summary</CardTitle>
              <CardDescription>
                A brief summary of the bookmarked content.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <p>{bookmark.summary}</p>
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="highlights">
          <Card>
            <CardHeader>
              <CardTitle>Highlights</CardTitle>
              <CardDescription>
                Your saved highlights from the content.
              </CardDescription>
            </CardHeader>
            <CardContent>
              {highlights.length > 0 ? (
                <ScrollArea className="h-[60vh]">
                  <div className="space-y-4">
                    {highlights.map((highlight) => (
                      <div key={highlight.id} className="flex items-start justify-between bg-muted p-4 rounded-md">
                        <p className="text-sm">{highlight.text}</p>
                        <Button variant="ghost" size="sm" onClick={() => removeHighlight(highlight.id)}>
                          <X className="h-4 w-4" />
                        </Button>
                      </div>
                    ))}
                  </div>
                </ScrollArea>
              ) : (
                <p className="text-muted-foreground">No highlights yet. Start highlighting in the Content tab!</p>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

