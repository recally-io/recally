export interface Highlight {
  id: string;
  text: string;
  startOffset: number;
  endOffset: number;
  note?: string;
}

export interface Bookmark {
  id: number;
  title: string;
  url: string;
  tags: string[];
  content: string;
  image?: string;
  summary: string;
  highlights?: Highlight[];
  dateAdded?: string;
}
