import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { ScrollArea } from "@/components/ui/scroll-area"
import { ExternalLink, Highlighter, X } from 'lucide-react'
import { useEffect, useMemo, useRef, useState } from 'react'
import { v4 as uuidv4 } from 'uuid'

interface Highlight {
  id: string
  text: string
  startOffset: number
  endOffset: number
  note?: string
}

interface BookmarkDetailProps {
  bookmark: {
    id: number
    title: string
    url: string
    tags: string[]
    content: string
    image?: string
    summary: string
  }
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
      selection.removeAllRanges()
    }
  }

  const handleAddNote = (id: string, note: string) => {
    setHighlights(highlights.map(h => h.id === id ? { ...h, note } : h))
  }

  const handleRemoveHighlight = (id: string) => {
    setHighlights(highlights.filter(h => h.id !== id))
  }

  const renderContent = useMemo(() => {
    let content = bookmark.content
    highlights.sort((a, b) => b.startOffset - a.startOffset).forEach(highlight => {
      const before = content.slice(0, highlight.startOffset)
      const highlighted = content.slice(highlight.startOffset, highlight.endOffset)
      const after = content.slice(highlight.endOffset)
      content = `${before}<span class="bg-yellow-200 dark:bg-yellow-800" data-highlight-id="${highlight.id}">${highlighted}</span>${after}`
    })
    return content
  }, [bookmark.content, highlights])

  return (
    <Card className="w-full max-w-4xl mx-auto">
      <CardHeader>
        <div className="flex justify-between items-start">
          <div>
            <CardTitle className="text-2xl mb-2">{bookmark.title}</CardTitle>
            <CardDescription>
              <a href={bookmark.url} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:text-blue-700 flex items-center transition-colors">
                {bookmark.url} <ExternalLink className="h-4 w-4 ml-1" />
              </a>
            </CardDescription>
          </div>
          <div className="flex flex-wrap gap-2">
            {bookmark.tags.map((tag) => (
              <Badge key={tag} variant="secondary" className="transition-colors hover:bg-primary hover:text-primary-foreground">
                {tag}
              </Badge>
            ))}
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {bookmark.image && (
          <img 
            src={bookmark.image} 
            alt={`Image for ${bookmark.title}`} 
            className="w-full h-64 object-cover mb-4 rounded-md"
          />
        )}
        <div className="mb-6">
          <h3 className="text-lg font-semibold mb-2">AI Summary</h3>
          <p className="text-gray-700 dark:text-gray-300">{bookmark.summary}</p>
        </div>
        <div className="mb-4">
          <Button
            onClick={() => setIsHighlighting(!isHighlighting)}
            variant={isHighlighting ? "secondary" : "outline"}
            className="transition-colors"
          >
            <Highlighter className="h-4 w-4 mr-2" />
            {isHighlighting ? "Finish Highlighting" : "Start Highlighting"}
          </Button>
        </div>
        <div>
          <h3 className="text-lg font-semibold mb-2">Full Content</h3>
          <ScrollArea className="h-[300px] rounded-md border p-4">
            <div 
              ref={contentRef}
              className={`prose dark:prose-invert max-w-none ${isHighlighting ? 'cursor-pointer' : ''}`}
              dangerouslySetInnerHTML={{ __html: renderContent }}
              onClick={isHighlighting ? handleHighlight : undefined}
            />
          </ScrollArea>
        </div>
        {highlights.length > 0 && (
          <div className="mt-6">
            <h3 className="text-lg font-semibold mb-2">Highlights and Notes</h3>
            {highlights.map((highlight) => (
              <div key={highlight.id} className="mb-4 p-4 bg-gray-100 dark:bg-gray-800 rounded-md transition-colors">
                <div className="flex justify-between items-start mb-2">
                  <p className="font-medium">"{highlight.text}"</p>
                  <Button variant="ghost" size="sm" onClick={() => handleRemoveHighlight(highlight.id)}>
                    <X className="h-4 w-4" />
                  </Button>
                </div>
                <Input
                  placeholder="Add a note..."
                  value={highlight.note || ''}
                  onChange={(e) => handleAddNote(highlight.id, e.target.value)}
                  className="mt-2"
                />
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}

