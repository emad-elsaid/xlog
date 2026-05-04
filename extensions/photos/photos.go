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

const extensionName = "photos"

func init() {
	xlog.RegisterExtension(Photos{})
}

type Photos struct{}

func (Photos) Name() string { return extensionName }
func (Photos) Init() {
	shortcode.RegisterShortCode(extensionName, shortcode.ShortCode{Render: photosShortcode(extensionName)})
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
		return errors.New("invalid path: empty path not allowed - provide a relative file path")
	}

	// Block null byte injection
	if strings.Contains(p, "\x00") {
		return errors.New("invalid path: null bytes not allowed - remove null bytes from the path")
	}

	// Block absolute paths
	if filepath.IsAbs(p) {
		return fmt.Errorf("invalid path: absolute paths not allowed - use relative path instead of %q", p)
	}

	// Clean the path and check for traversal
	cleaned := filepath.Clean(p)

	// Check for .. as a complete path component (not just prefix of filename)
	// Split by separator and look for exactly ".." as an element
	parts := strings.Split(cleaned, string(filepath.Separator))
	for _, part := range parts {
		if part == ".." {
			return fmt.Errorf("invalid path: path traversal detected in %q - avoid using '..' in paths", p)
		}
	}

	return nil
}

// NewPhoto creates a Photo instance from a filesystem path. The path must be
// within allowed directories. Returns an error if the path is invalid or the
// file cannot be read.
func NewPhoto(photoPath string) (*Photo, error) {
	// #nosec G304 -- Caller is responsible for path validation
	stat, err := os.Stat(photoPath)
	if err != nil {
		return nil, err
	}

	// #nosec G304 -- Caller is responsible for path validation
	f, err := os.Open(photoPath)
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
		Thumbnail: "/+/photos/thumbnail/" + photoPath,
		Page:      "/+/photos/photo/" + photoPath,
		Original:  photoPath,
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
				return fmt.Errorf("failed to access path in %s: %w", p, err)
			}

			if d.Type().IsRegular() && supportedExt.Include(strings.ToLower(path.Ext(file))) {
				photo, err := NewPhoto(file)
				if err != nil {
					return fmt.Errorf("failed to create photo from %s: %w", file, err)
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
			extensionName: photos,
		})
	}
}

func resizeHandler(r xlog.Request) xlog.Output {
	photo_path := r.PathValue("path")

	// Validate path for security before any file operations
	if err := validatePath(photo_path); err != nil {
		return errorOutput(fmt.Sprintf("Error: %v", err))
	}

	cacheOutput := tryLoadCache(photo_path)
	if cacheOutput != nil {
		return cacheOutput
	}

	return resizeAndCachePhoto(photo_path)
}

func errorOutput(message string) xlog.Output {
	return func(w xlog.Response, r xlog.Request) {
		if _, writeErr := fmt.Fprint(w, message); writeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to write error response: %v\n", writeErr)
		}
	}
}

func tryLoadCache(photo_path string) xlog.Output {
	const cacheDir = ".cache"
	if err := os.Mkdir(cacheDir, 0700); err != nil && !os.IsExist(err) {
		return xlog.InternalServerError(err)
	}

	cacheFile := path.Join(cacheDir, fmt.Sprintf("photo-%x", sha256.Sum256([]byte(photo_path))))
	// #nosec G304 -- Cache file path is constructed from hash, not user input
	cache, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil
	}

	return func(w xlog.Response, r xlog.Request) {
		if _, writeErr := w.Write(cache); writeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to write cached photo: %v\n", writeErr)
		}
	}
}

func resizeAndCachePhoto(photo_path string) xlog.Output {
	const cacheDir = ".cache"
	cacheFile := path.Join(cacheDir, fmt.Sprintf("photo-%x", sha256.Sum256([]byte(photo_path))))

	return func(w xlog.Response, r xlog.Request) {
		resized, err := resizePhoto(photo_path)
		if err != nil {
			if _, writeErr := fmt.Fprintf(w, "Failed to process image: %v", err); writeErr != nil {
				fmt.Fprintf(os.Stderr, "Failed to write error response: %v\n", writeErr)
			}
			return
		}

		if err := os.WriteFile(cacheFile, resized, 0600); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to cache resized image: %v\n", err)
		}

		if _, err := w.Write(resized); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write resized image: %v\n", err)
		}
	}
}

func resizePhoto(photo_path string) ([]byte, error) {
	// #nosec G304 -- Path validated via validatePath before this point
	inputImage, err := os.Open(photo_path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := inputImage.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to close input image: %v\n", closeErr)
		}
	}()

	src, _, err := image.Decode(inputImage)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := src.Bounds()
	dim := bounds.Max

	width := 700
	height := int(float32(width) / float32(dim.X) * float32(dim.Y))

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, bounds, draw.Over, nil)

	var out bytes.Buffer
	if err := png.Encode(&out, dst); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return out.Bytes(), nil
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
