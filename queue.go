package images

import "context"

// Listener is an interface for processing messages, and keeps track of messages
// to be processed.
type Listener interface {
	// Listen runs the callback for each message that comes into the queue,
	// passing in the SHA256 of a Blob to be processed, and returning an
	// Image.
	Listen(ctx context.Context, callback func(ctx context.Context, sha256 string) (Image, error)) error
}

// Publisher is an interface for designating a Blob as needing to be processed
// by a Listener.
type Publisher interface {
	// Push designates a Blob as needing processing.
	Push(ctx context.Context, sha256 string) error
}
