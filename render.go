package echo

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"math"

	_ "embed"
)

var (
	//go:embed templates/show.svg
	showSVG string

	//go:embed templates/hide.svg
	hideSVG string
)

var templateFuncs = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"mod": func(n, d int) int {
		if d == 0 {
			return 0
		}
		return n % d
	},
	"icon": func(name string) template.CSS {
		/* url('data:image/svg+xml;base64,') */
		base64.StdEncoding.EncodeToString([]byte("foo"))

		var content string
		switch name {
		case "show":
			content = showSVG
		case "hide":
			content = hideSVG
		default:
			panic("unknown icon: " + name)
		}
		return template.CSS(
			fmt.Sprintf("url('data:image/svg+xml;base64,%s')",
				base64.StdEncoding.EncodeToString([]byte(content))))
	},
}

func mustParse(name, source string) *template.Template {
	return template.Must(
		template.New(name).
			Funcs(templateFuncs).
			Parse(source))
}

const FullOpacity uint8 = 100

// RGBAColor implements Color for the RGBA color space.
type RGBAColor struct {
	// R, G, and B color values between 0 and 255.
	R, G, B uint8

	// Alpha value out of 100. Values higher than 100 will be rounded down.
	Alpha uint8
}

// CSS representation of the RGBA color.
func (c *RGBAColor) CSS() template.CSS {
	if c == nil {
		return "rgb(0 0 0)"
	}
	return template.CSS(
		fmt.Sprintf("rgb(%d %d %d / %0.1f)", c.R, c.G, c.B,
			math.Min(float64(c.Alpha/100), 1.0)))
}

// Color abstracts the rendering of different color spaces as CSS for rendering
// the UI.
type Color interface {
	CSS() template.CSS
}

// ColorPalette for rendering the UI.
type ColorPalette []Color

var (
	// DefaultColors for the rendering of the UI.
	DefaultColors = ColorPalette{
		&RGBAColor{9, 94, 147, FullOpacity},
		&RGBAColor{102, 157, 13, FullOpacity},
		&RGBAColor{243, 0, 8, FullOpacity},
		&RGBAColor{221, 162, 0, FullOpacity},
		&RGBAColor{71, 25, 134, FullOpacity},
	}
)

type pageContent struct {
	Environment *EnvironmentDump
	Request     *RequestDump
	Colors      ColorPalette
}

func render(dump *RequestDump, with ColorPalette) *pageContent {
	if with == nil {
		with = DefaultColors
	}
	return &pageContent{
		Environment: DumpEnv(PolicyFromEnv()),
		Request:     dump,
		Colors:      with,
	}
}
