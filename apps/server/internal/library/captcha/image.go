package captcha

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"

	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"golang.org/x/image/font"
)

// draw 自绘验证码图片：4 位字符 + 随机干扰线/噪点（零第三方验证码库）。
func draw(answer string) ([]byte, error) {
	const (
		w = 120
		h = 48
	)
	fontFace := basicfont.Face7x13
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	// 浅色背景
	bg := color.RGBA{R: 245, G: 246, B: 248, A: 255}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, bg)
		}
	}

	// 干扰线
	rng := rand.New(rand.NewSource(rand.Int63()))
	line := color.RGBA{R: 180, G: 186, B: 196, A: 255}
	for i := 0; i < 4; i++ {
		x0, y0 := rng.Intn(w), rng.Intn(h)
		x1, y1 := rng.Intn(w), rng.Intn(h)
		drawLine(img, x0, y0, x1, y1, line)
	}

	// 噪点
	dot := color.RGBA{R: 150, G: 155, B: 165, A: 255}
	for i := 0; i < 120; i++ {
		img.Set(rng.Intn(w), rng.Intn(h), dot)
	}

	// 逐字符绘制（带轻微纵向抖动）
	fg := color.RGBA{R: 52, G: 60, B: 74, A: 255}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(fg),
		Face: fontFace,
	}
	x := 16
	for i := 0; i < len(answer); i++ {
		d.Dot = fixed.P(x, 30+rng.Intn(6)-3)
		d.DrawString(string(answer[i]))
		x += 22
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// drawLine Bresenham 直线。
func drawLine(img *image.RGBA, x0, y0, x1, y1 int, c color.Color) {
	dx, dy := abs(x1-x0), -abs(y1-y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy
	for {
		img.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
