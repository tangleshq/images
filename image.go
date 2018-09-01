package images // import "tangl.es/code/images"
import "net/http"

type Image struct {
	// The SHA256 of this Image.
	SHA256 string
	// The extension (without `.`) to serve this Image with.
	Extension string
	// The SHA256 of the Blob this Image was generated from.
	SourceSHA256 string
	// The content type of this Image.
	ContentType string
	// The pixel density of this Image.
	PixelDensity int64
	// The width in pixels of this Image.
	Width int64
	// The height in pixels of this Image.
	Height int64
	// The size in bytes of this Image.
	Size int64
	// Any headers we want to associate with the Image for special handling.
	Headers http.Header
}
