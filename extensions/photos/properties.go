package photos

import (
	"fmt"
	"math"
	"time"

	"github.com/emad-elsaid/xlog"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

const (
	iconCalendar    = "fa-regular fa-calendar"
	iconCameraRetro = "fa-solid fa-camera-retro"
	propCaptureTime = "capture time"
	propCameraMake  = "camera make"
	propISO         = "ISO"
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
	photo, ok := p.(*Photo)
	if !ok || photo == nil {
		return nil
	}

	e := photo.Exif
	if e == nil {
		return nil
	}

	props := []xlog.Property{}

	appendCaptureTime(&props, photo.Time)
	appendStringExifProps(&props, e)
	appendRationalExifProps(&props, e)

	return props
}

func appendCaptureTime(props *[]xlog.Property, t time.Time) {
	if !t.IsZero() {
		*props = append(*props, Property{
			IconVal: iconCalendar,
			NameVal: propCaptureTime,
			Val:     fmt.Sprintf("%s %d %s %d", t.Weekday(), t.Day(), t.Month(), t.Year()),
		})
	}
}

func appendStringExifProps(props *[]xlog.Property, e *exif.Exif) {
	str := func(t *tiff.Tag) string {
		s, _ := t.StringVal()
		return s
	}

	stringFields := []struct {
		field exif.FieldName
		name  string
	}{
		{exif.Make, propCameraMake},
		{exif.Model, "camera model"},
		{exif.LensMake, "lens make"},
		{exif.LensModel, "lens model"},
	}

	for _, sf := range stringFields {
		if tag, err := e.Get(sf.field); err == nil {
			*props = append(*props, Property{
				IconVal: iconCameraRetro,
				NameVal: sf.name,
				Val:     str(tag),
			})
		}
	}

	if iso, err := e.Get(exif.ISOSpeedRatings); err == nil {
		*props = append(*props, Property{
			IconVal: iconCameraRetro,
			NameVal: propISO,
			Val:     iso.String(),
		})
	}
}

func appendRationalExifProps(props *[]xlog.Property, e *exif.Exif) {
	if focal, err := e.Get(exif.FocalLength); err == nil {
		if nom, denom, err := focal.Rat2(0); err == nil {
			*props = append(*props, Property{
				IconVal: iconCameraRetro,
				NameVal: "focal Length",
				Val:     fmt.Sprintf("%dmm", nom/denom),
			})
		}
	}

	if aperture, err := e.Get(exif.ApertureValue); err == nil {
		if nom, denom, err := aperture.Rat2(0); err == nil {
			*props = append(*props, Property{
				IconVal: iconCameraRetro,
				NameVal: "aperture",
				Val:     fmt.Sprintf("f/%.1f", float32(nom)/float32(denom)),
			})
		}
	}

	if shutter, err := e.Get(exif.ShutterSpeedValue); err == nil {
		if snom, sdenom, err := shutter.Rat2(0); err == nil {
			*props = append(*props, Property{
				IconVal: iconCameraRetro,
				NameVal: "shutter speed",
				Val:     fmt.Sprintf("1/%.0fs", math.Pow(2, float64(snom)/float64(sdenom))),
			})
		}
	}
}
