package photos

import (
	"fmt"
	"math"

	"github.com/emad-elsaid/xlog"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

type Property struct {
	IconVal string
	NameVal string
	Val     any
}

func (p Property) Icon() string { return p.IconVal }
func (p Property) Name() string { return p.NameVal }
func (p Property) Value() any   { return p.Val }

func properties(p xlog.Page) []xlog.Property {
	props := []xlog.Property{}

	photo, ok := p.(*Photo)
	if !ok {
		return nil
	}

	e := photo.Exif
	if e == nil {
		return nil
	}

	str := func(t *tiff.Tag) string {
		s, _ := t.StringVal()
		return s
	}

	t := photo.Time
	if !t.IsZero() {
		props = append(props, Property{
			IconVal: "fa-regular fa-calendar",
			NameVal: "capture time",
			Val:     fmt.Sprintf("%s %d %s %d", t.Weekday(), t.Day(), t.Month(), t.Year()),
		})
	}

	if m, err := e.Get(exif.Make); err == nil {
		props = append(props, Property{
			IconVal: "fa-solid fa-camera-retro",
			NameVal: "camera make",
			Val:     str(m),
		})
	}

	if c, err := e.Get(exif.Model); err == nil {
		props = append(props, Property{
			IconVal: "fa-solid fa-camera-retro",
			NameVal: "camera model",
			Val:     str(c),
		})
	}

	if m, err := e.Get(exif.LensMake); err == nil {
		props = append(props, Property{
			IconVal: "fa-solid fa-camera-retro",
			NameVal: "lens make",
			Val:     str(m),
		})
	}

	if m, err := e.Get(exif.LensModel); err == nil {
		props = append(props, Property{
			IconVal: "fa-solid fa-camera-retro",
			NameVal: "lens model",
			Val:     str(m),
		})
	}

	if focal, err := e.Get(exif.FocalLength); err == nil {
		nom, denom, err := focal.Rat2(0)
		if err == nil {
			props = append(props, Property{
				IconVal: "fa-solid fa-camera-retro",
				NameVal: "focal Length",
				Val:     fmt.Sprintf("%dmm", nom/denom),
			})
		}
	}

	if aperture, err := e.Get(exif.ApertureValue); err == nil {
		nom, denom, err := aperture.Rat2(0)
		if err == nil {
			props = append(props, Property{
				IconVal: "fa-solid fa-camera-retro",
				NameVal: "aperture",
				Val:     fmt.Sprintf("f/%.1f", float32(nom)/float32(denom)),
			})
		}
	}

	if iso, err := e.Get(exif.ISOSpeedRatings); err == nil {
		props = append(props, Property{
			IconVal: "fa-solid fa-camera-retro",
			NameVal: "ISO",
			Val:     iso.String(),
		})
	}

	if shutter, err := e.Get(exif.ShutterSpeedValue); err == nil {
		snom, sdenom, err := shutter.Rat2(0)
		if err == nil {
			props = append(props, Property{
				IconVal: "fa-solid fa-camera-retro",
				NameVal: "shutter speed",
				Val:     fmt.Sprintf("1/%.0fs", math.Pow(2, float64(snom)/float64(sdenom))),
			})
		}
	}

	return props
}
