package textsplitter

import (
	"log"
	"strings"
)

// CreateDocuments creates documents from texts and metadatas with a text splitter. If
// the length of the metadatas is zero, the result documents will contain no metadata.
// Otherwise, the numbers of texts and metadatas must match.

// joinDocs comines two documents with the separator used to split them.
func joinDocs(docs []string, separator string) string {
	return strings.TrimSpace(strings.Join(docs, separator))
}

// mergeSplits merges smaller splits into splits that are closer to the chunkSize.
func mergeSplits(splits []string, separator string, chunkSize, chunkOverlap int, lenFunc func(string) int) []string { //nolint:cyclop
	docs := make([]string, 0)
	currentDoc := make([]string, 0)
	total := 0

	for _, split := range splits {
		totalWithSplit := total + lenFunc(split)
		if len(currentDoc) != 0 {
			totalWithSplit += lenFunc(separator)
		}

		maybePrintWarning(total, chunkSize)

		if totalWithSplit > chunkSize && len(currentDoc) > 0 {
			doc := joinDocs(currentDoc, separator)
			if doc != "" {
				docs = append(docs, doc)
			}

			for shouldPop(chunkOverlap, chunkSize, total, lenFunc(split), lenFunc(separator), len(currentDoc)) {
				total -= lenFunc(currentDoc[0]) //nolint:gosec
				if len(currentDoc) > 1 {
					total -= lenFunc(separator)
				}

				currentDoc = currentDoc[1:] //nolint:gosec
			}
		}

		currentDoc = append(currentDoc, split)
		total += lenFunc(split)

		if len(currentDoc) > 1 {
			total += lenFunc(separator)
		}
	}

	doc := joinDocs(currentDoc, separator)
	if doc != "" {
		docs = append(docs, doc)
	}

	return docs
}

func maybePrintWarning(total, chunkSize int) {
	if total > chunkSize {
		log.Printf(
			"[WARN] created a chunk with size of %v, which is longer then the specified %v\n",
			total,
			chunkSize,
		)
	}
}

// Keep popping if:
//   - the chunk is larger than the chunk overlap
//   - or if there are any chunks and the length is long
func shouldPop(chunkOverlap, chunkSize, total, splitLen, separatorLen, currentDocLen int) bool {
	docsNeededToAddSep := 2
	if currentDocLen < docsNeededToAddSep {
		separatorLen = 0
	}

	return currentDocLen > 0 && (total > chunkOverlap || (total+splitLen+separatorLen > chunkSize && total > 0))
}
