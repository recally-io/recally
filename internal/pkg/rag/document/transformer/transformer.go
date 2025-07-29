package transformer

import (
	"context"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/rag/document"
	"runtime"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Transformer func([]document.Document) ([]document.Document, error)

// Batch applies a Transformer function to a batch of documents concurrently.
// It uses a semaphore to limit the number of concurrent transformations to the
// number of available CPU cores.
//
// Parameters:
//   - transformer: A function that takes a slice of documents and returns a transformed
//     slice of documents and an error.
//   - docs: A slice of documents to be transformed.
//
// Returns:
// - A slice of transformed documents.
// - An error if any of the transformations fail.
func Batch(transformer Transformer, docs []document.Document) ([]document.Document, error) {
	// Create a semaphore weighted by the number of available CPU cores.
	sem := semaphore.NewWeighted(int64(runtime.GOMAXPROCS(0)))

	// Create an error group with the context.
	eg, ctx := errgroup.WithContext(context.TODO())

	// Initialize the result slice.
	res := make([]document.Document, 0)

	// Iterate over each document in the input slice.
	for _, doc := range docs {
		// Acquire a semaphore slot.
		if err := sem.Acquire(ctx, 1); err != nil {
			break
		}

		// Launch a goroutine to transform the document.
		eg.Go(func() error {
			// Release the semaphore slot when the goroutine completes.
			defer sem.Release(1)

			// Apply the transformer function to the document.
			newDoc, err := transformer([]document.Document{doc})
			if err != nil {
				// Log the error if the transformation fails.
				logger.Default.Error("failed to transform document", "err", err)

				return err
			}

			// Append the transformed document to the result slice.
			res = append(res, newDoc...)

			return nil
		})
	}

	// Wait for all goroutines to complete.
	err := eg.Wait()

	return res, err
}
