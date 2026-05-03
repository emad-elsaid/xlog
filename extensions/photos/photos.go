package photos

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/emad-elsaid/xlog/markdown/ast"

	"github.com/emad-elsaid/types"
	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/extensions/shortcode"
	"github.com/rwcarlsen/goexif/exif"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

//go:embed templates
var templates embed.FS

var supportedExt = types.Slice[string]{".jpg", ".jpeg", ".gif", ".png"}

func init() {
	xlog.RegisterExtension(Photos{})
}

type Photos struct{}

func (Photos) Name() string { return "photos" }
func (Photos) Init() {
	shortcode.RegisterShortCode("photos", shortcode.ShortCode{Render: photosShortcode("photos")})
	shortcode.RegisterShortCode("photos-grid", shortcode.ShortCode{Render: photosShortcode("photos-grid")})
	xlog.RegisterTemplate(templates, "templates")
	xlog.RegisterProperty(properties)
	xlog.Get(`/+/photos/thumbnail/{path...}`, resizeHandler)
	xlog.Get(`/+/photos/photo/{path...}`, photoHandler)
}

type Photo struct {
	Thumbnail string
	Page      string
	Original  string
	Exif      *exif.Exif
	Time      time.Time
}

func (p *Photo) Name() string {
	base := path.Base(p.Thumbnail)
	ext := path.Ext(base)
	return base[:len(base)-len(ext)]
}

func (*Photo) FileName() string         { return "" }
func (*Photo) Exists() bool             { return false }
func (*Photo) Content() xlog.Markdown   { return "" }
func (*Photo) Delete() bool             { return false }
func (*Photo) Write(xlog.Markdown) bool { return false }
func (*Photo) ModTime() time.Time       { return time.Time{} }
func (*Photo) AST() ([]byte, ast.Node)  { return nil, nil }
func (p *Photo) Render() template.HTML {
	return xlog.Partial("photo", xlog.Locals{"photo": p})
}

// validatePath ensures the provided path is safe and does not contain path traversal attacks.
// It blocks:
//   - Absolute paths (starting with /)
//   - Paths containing .. that would escape the current directory
//   - Empty paths
//   - Paths with null bytes
func validatePath(p string) error {
	if p == "" {
		return errors.New("invalid path: empty path not allowed")
	}

	// Block null byte injection
	if strings.Contains(p, "\x00") {
		return errors.New("invalid path: null bytes not allowed")
	}

	// Block absolute paths
	if filepath.IsAbs(p) {
		return errors.New("invalid path: absolute paths not allowed")
	}

	// Clean the path and check for traversal
	cleaned := filepath.Clean(p)

	// Check for .. as a complete path component (not just prefix of filename)
	// Split by separator and look for exactly ".." as an element
	parts := strings.Split(cleaned, string(filepath.Separator))
	for _, part := range parts {
		if part == ".." {
			return errors.New("invalid path: path traversal detected")
		}
	}

	return nil
}

func NewPhoto(path string) (*Photo, error) {
	// #nosec G304 -- Caller is responsible for path validation
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// #nosec G304 -- Caller is responsible for path validation
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	exifData, _ := exif.Decode(f)
	t := stat.ModTime()

	if exifData != nil {
		shootingTime, err := exifData.DateTime()
		if err == nil {
			t = shootingTime
		}
	}

	return &Photo{
		Thumbnail: "/+/photos/thumbnail/" + path,
		Page:      "/+/photos/photo/" + path,
		Original:  path,
		Exif:      exifData,
		Time:      t,
	}, nil
}

func photosShortcode(tpl string) func(xlog.Markdown) template.HTML {
	return func(input xlog.Markdown) template.HTML {
		p := strings.TrimSpace(string(input))

		photos := []*Photo{}

		err := filepath.WalkDir(p, func(file string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.Type().IsRegular() && supportedExt.Include(strings.ToLower(path.Ext(file))) {
				photo, err := NewPhoto(file)
				if err != nil {
					return err
				}

				xlog.RegisterBuildPage(photo.Thumbnail, false)
				xlog.RegisterBuildPage(photo.Page, true)
				photos = append(photos, photo)
			}

			return nil
		})

		if err != nil {
			// #nosec G203 -- Error message from filepath.WalkDir; does not contain user-controlled HTML
			return template.HTML(err.Error())
		}

		slices.SortFunc(photos, func(i, j *Photo) int {
			return j.Time.Compare(i.Time)
		})

		return xlog.Partial(tpl, xlog.Locals{
			"photos": photos,
		})
	}
}

func resizeHandler(r xlog.Request) xlog.Output {
	photo_path := r.PathValue("path")

	// Validate path for security before any file operations
	if err := validatePath(photo_path); err != nil {
		return func(w xlog.Response, r xlog.Request) {
			if _, writeErr := fmt.Fprintf(w, "Error: %v", err); writeErr != nil {
				fmt.Fprintf(os.Stderr, "Failed to write error response: %v\n", writeErr)
			}
		}
	}

	const cacheDir = ".cache"
	if err := os.Mkdir(cacheDir, 0700); err != nil && !os.IsExist(err) {
		return xlog.InternalServerError(err)
	}

	cacheFile := path.Join(cacheDir, fmt.Sprintf("photo-%x", sha256.Sum256([]byte(photo_path))))
	// #nosec G304 -- Cache file path is constructed from hash, not user input
	cache, err := os.ReadFile(cacheFile)
	if err == nil {
		return func(w xlog.Response, r xlog.Request) {
			if _, writeErr := w.Write(cache); writeErr != nil {
				// Log error but can't do much at this point
				fmt.Fprintf(os.Stderr, "Failed to write cached photo: %v\n", writeErr)
			}
		}
	}

	return func(w xlog.Response, r xlog.Request) {
		// #nosec G304 -- Path validated via validatePath before this point
		inputImage, err := os.Open(photo_path)
		if err != nil {
			if _, writeErr := fmt.Fprint(w, err.Error()); writeErr != nil {
				fmt.Fprintf(os.Stderr, "Failed to write error response: %v\n", writeErr)
			}
			return
		}
		defer func() {
			if closeErr := inputImage.Close(); closeErr != nil {
				fmt.Fprintf(os.Stderr, "Failed to close input image: %v\n", closeErr)
			}
		}()

		src, _, err := image.Decode(inputImage)
		if err != nil {
			if _, writeErr := fmt.Fprintf(w, "Failed to decode image: %v", err); writeErr != nil {
				fmt.Fprintf(os.Stderr, "Failed to write decode error: %v\n", writeErr)
			}
			return
		}

		bounds := src.Bounds()
		dim := bounds.Max

		width := 700
		height := int(float32(width) / float32(dim.X) * float32(dim.Y))

		dst := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.NearestNeighbor.Scale(dst, dst.Rect, src, bounds, draw.Over, nil)

		var out bytes.Buffer

		if err := png.Encode(&out, dst); err != nil {
			if _, writeErr := fmt.Fprintf(w, "Failed to encode image: %v", err); writeErr != nil {
				fmt.Fprintf(os.Stderr, "Failed to write encode error: %v\n", writeErr)
			}
			return
		}

		if err := os.WriteFile(cacheFile, out.Bytes(), 0600); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to cache resized image: %v\n", err)
			// Continue despite cache failure
		}

		if _, err := w.Write(out.Bytes()); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write resized image: %v\n", err)
		}
	}
}

func photoHandler(r xlog.Request) xlog.Output {
	photo_path := r.PathValue("path")

	// Validate path for security before creating Photo
	if err := validatePath(photo_path); err != nil {
		return func(w xlog.Response, r xlog.Request) {
			if _, writeErr := fmt.Fprintf(w, "Error: %v", err); writeErr != nil {
				fmt.Fprintf(os.Stderr, "Failed to write error response: %v\n", writeErr)
			}
		}
	}

	photo, err := NewPhoto(photo_path)
	if err != nil {
		return xlog.InternalServerError(err)
	}

	return xlog.Render("page", xlog.Locals{
		"page": photo,
	})
}
