package images

import (
	"context"
	"image"
	"io"
	"net/http"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/denisbrodbeck/sqip"
)

type SQIP struct {
	WorkSize   int    // larger images get resized to this; larger values don't really help much
	Count      int    // number of primitive SVG shapes to use
	Mode       int    // shape type
	Alpha      int    // alpha value
	Repeat     int    // number of extra shapes each iteration with reduced search (mostly good for beziers)
	NumWorkers int    // parallelisation factor. runtime.NumCPU() should be fine
	Background string // background color in hex format
}

func (s SQIP) Process(ctx context.Context, in io.Reader) (Image, []byte, error) {
	im, _, err := image.Decode(in)
	if err != nil {
		return Image{}, nil, err
	}
	svg, width, height, err := sqip.RunLoaded(im, s.WorkSize, s.Count, s.Mode, s.Alpha, s.Repeat, s.NumWorkers, s.Background)
	if err != nil {
		return Image{}, nil, err
	}
	i := Image{
		// SHA256 gets set by the caller
		// SourceSHA256 gets set by the caller
		// No pixel density for SVGs
		Extension:   "sqip.svg",
		ContentType: "image/svg+xml",
		Width:       int64(width),
		Height:      int64(height),
		Size:        int64(len(svg)),
		Headers:     http.Header{},
	}
	i.Headers.Set("render-inline", "true")
	i.Headers.Set("placeholder", "true")
	return i, []byte(svg), nil
}
