package images

import (
	"context"
	"io"
)

type Processor interface {
	Process(ctx context.Context, in io.Reader) (Image, []byte, error)
}
