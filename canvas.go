package canvas

import (
	"errors"
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"image/color"
)

func hexColor(c color.Color) string {
	r, g, b, _ := c.RGBA()
	r = r * 0xff / 0xffff
	g = g * 0xff / 0xffff
	b = b * 0xff / 0xffff

	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

type Canvas struct {
	elem *js.Object
	ctx  *js.Object
}

func New(elem *js.Object) (*Canvas, error) {
	if js.Global.Get("CanvasRenderingContext2D") == nil {
		return nil, errors.New("Browser doesn't support canvas 2D")
	}

	return &Canvas{
		elem: elem,
		ctx:  elem.Call("getContext", "2d"),
	}, nil
}

func (c *Canvas) Rect(r Rectangle) Object {
	r = r.Canon()

	p := func(pb *PathBuilder) {
		pb.Rect(r)
	}

	return &pathObj{
		c:    c,
		path: p,
	}
}

func (c *Canvas) Path(p func(*PathBuilder)) Object {
	return &pathObj{
		c:    c,
		path: p,
	}
}

func (c *Canvas) Text(text string, mw float64) Object {
	return &textObj{
		c:    c,
		text: text,
		mw:   mw,
	}
}

func (c *Canvas) Clear(r Rectangle) {
	r = r.Canon()
	c.ctx.Call("clearRect", r.Min.X, r.Min.Y, r.Dx(), r.Dy())
}

func (c *Canvas) MeasureText(text string) (tm TextMetrics) {
	return TextMetrics{c.ctx.Call("measureText", text)}
}

func (c *Canvas) LineWidth(width float64) float64 {
	if width < 0 {
		return c.ctx.Get("lineWidth").Float()
	}

	c.ctx.Set("lineWidth", width)
	return width
}

func (c *Canvas) LineCap() LineCap {
	lc := c.ctx.Get("lineCap").String()
	switch lc {
	case "butt":
		return ButtCap
	case "round":
		return RoundCap
	case "square":
		return SquareCap
	}

	panic("Unknown lineCap: " + lc)
}

func (c *Canvas) SetLineCap(lc LineCap) {
	c.ctx.Set("lineCap", lc.String())
}

func (c *Canvas) LineJoin() LineJoin {
	lj := c.ctx.Get("lineJoin").String()
	switch lj {
	case "bevel":
		return BevelJoin
	case "round":
		return RoundJoin
	case "miter":
		return MiterJoin
	}

	panic("Unknown lineJoin: " + lj)
}

func (c *Canvas) SetLineJoin(lj LineJoin) {
	c.ctx.Set("lineJoin", lj.String())
}

func (c *Canvas) MiterLimit(limit float64) float64 {
	if limit < 0 {
		return c.ctx.Get("miterLimit").Float()
	}

	c.ctx.Set("miterLimit", limit)
	return limit
}

func (c *Canvas) LineDash(pattern []float64) []float64 {
	if pattern == nil {
		jp := c.ctx.Call("getLineDash")

		pattern = make([]float64, 0, jp.Length())
		for i := 0; i < jp.Length(); i++ {
			pattern = append(pattern, jp.Index(i).Float())
		}

		return pattern
	}

	c.ctx.Call("setLineDash", pattern)
	return pattern
}

func (c *Canvas) Width() float64 {
	return c.elem.Get("width").Float()
}

func (c *Canvas) Height() float64 {
	return c.elem.Get("height").Float()
}

type Font struct {
	Size int
	Name string
}

func (f Font) String() string {
	return fmt.Sprintf("%v %v", f.Size, f.Name)
}

type LineCap int

const (
	ButtCap LineCap = iota
	RoundCap
	SquareCap
)

func (lc LineCap) String() string {
	switch lc {
	case ButtCap:
		return "butt"
	case RoundCap:
		return "round"
	case SquareCap:
		return "square"
	}

	panic(fmt.Errorf("Unknown LineCap: %v", int(lc)))
}

type LineJoin int

const (
	BevelJoin LineJoin = iota
	RoundJoin
	MiterJoin
)

func (lj LineJoin) String() string {
	switch lj {
	case BevelJoin:
		return "bevel"
	case RoundJoin:
		return "round"
	case MiterJoin:
		return "miter"
	}

	panic(fmt.Errorf("Unknown LineJoin: %v", int(lj)))
}
