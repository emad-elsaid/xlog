package photos

import (
	"bytes"
	"crypto/sha256"
	"embed"
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

func NewPhoto(path string) (*Photo, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

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

	const cacheDir = ".cache"
	os.Mkdir(cacheDir, 0700)

	cacheFile := path.Join(cacheDir, fmt.Sprintf("photo-%x", sha256.Sum256([]byte(photo_path))))
	cache, err := os.ReadFile(cacheFile)
	if err == nil {
		return func(w xlog.Response, r xlog.Request) {
			w.Write(cache)
		}
	}

	return func(w xlog.Response, r xlog.Request) {
		inputImage, err := os.Open(photo_path)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		defer inputImage.Close()

		src, _, _ := image.Decode(inputImage)
		bounds := src.Bounds()
		dim := bounds.Max

		width := 700
		height := int(float32(width) / float32(dim.X) * float32(dim.Y))

		dst := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.NearestNeighbor.Scale(dst, dst.Rect, src, bounds, draw.Over, nil)

		var out bytes.Buffer

		png.Encode(&out, dst)
		os.WriteFile(cacheFile, out.Bytes(), 0700)
		w.Write(out.Bytes())
	}
}

func photoHandler(r xlog.Request) xlog.Output {
	photo_path := r.PathValue("path")
	photo, err := NewPhoto(photo_path)
	if err != nil {
		return xlog.InternalServerError(err)
	}

	return xlog.Render("page", xlog.Locals{
		"page": photo,
	})
}
