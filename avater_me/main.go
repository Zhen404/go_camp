package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type User struct {
	name string
	email string
}

const (
	defaultfontFace = "Caviar_Dreams_Bold.ttf"
	fontSize = 210
	imageWidth = 500
	imageHeight = 500
	dpi = 72
	spacer = 20
	textY = 320
)

var (
	Red      = image.Uniform{color.RGBA{230, 25, 75, 255}}
	Green    = image.Uniform{color.RGBA{60, 180, 75, 255}}
	Yellow   = image.Uniform{color.RGBA{255, 225, 25, 255}}
	Blue     = image.Uniform{color.RGBA{0, 130, 200, 255}}
	Orange   = image.Uniform{color.RGBA{245, 130, 48, 255}}
	Purple   = image.Uniform{color.RGBA{145, 30, 180, 255}}
	Cyan     = image.Uniform{color.RGBA{70, 240, 240, 255}}
	Magenta  = image.Uniform{color.RGBA{240, 50, 230, 255}}
	Lime     = image.Uniform{color.RGBA{210, 245, 60, 255}}
	Pink     = image.Uniform{color.RGBA{250, 190, 190, 255}}
	Teal     = image.Uniform{color.RGBA{0, 128, 128, 255}}
	Lavender = image.Uniform{color.RGBA{230, 190, 255, 255}}
	Brown    = image.Uniform{color.RGBA{170, 110, 40, 255}}
	Beige    = image.Uniform{color.RGBA{255, 250, 200, 255}}
	Maroon   = image.Uniform{color.RGBA{128, 0, 0, 255}}
	Mint     = image.Uniform{color.RGBA{170, 255, 195, 255}}
	Olive    = image.Uniform{color.RGBA{128, 128, 0, 255}}
	Coral    = image.Uniform{color.RGBA{255, 215, 180, 255}}
	Navy     = image.Uniform{color.RGBA{0, 0, 128, 255}}
	Grey     = image.Uniform{color.RGBA{128, 128, 128, 255}}
	Gold     = image.Uniform{color.RGBA{251, 184, 41, 255}}
)

func defaultColor(initial string) image.Uniform {
	switch initial {
	case "A", "0":
		return Red
	case "B", "1":
		return Green
	case "C", "2":
		return Yellow
	case "D", "3":
		return Blue
	case "E", "4":
		return Orange
	case "F", "5":
		return Purple
	case "G", "6":
		return Lime
	case "H", "7":
		return Magenta
	case "I", "8":
		return Pink
	case "J", "9":
		return Cyan
	case "K":
		return Teal
	case "L":
		return Lavender
	case "M":
		return Brown
	case "N":
		return Beige
	case "O":
		return Maroon
	case "P":
		return Mint
	case "Q":
		return Olive
	case "R":
		return Coral
	case "S":
		return Navy
	case "T":
		return Gold
	default:
		return Grey
	}
}

func cleanString(incoming string) string {
	incoming = strings.TrimSpace(incoming)

	split := strings.Split(incoming, " ")
	if len(split) == 2 {
		incoming = split[0][0:1] + split[1][0:1]
	}

	return strings.ToUpper(strings.TrimSpace(incoming))
}

var imageCache sync.Map

func getImage(initials string) *image.RGBA {
	value, ok := imageCache.Load(initials)

	if !ok {
		return nil
	}

	image, ok2 := value.(*image.RGBA)
	if !ok2 {
		return nil
	}
	return image
}

func setImage(initials string, image *image.RGBA) {
	imageCache.Store(initials, image)
}

var fontFacePath = ""

func setFont(f string) {
	fontFacePath = f
}

func getFont(fontPath string) (*truetype.Font, error) {
	if fontPath == "" {
		fontPath = defaultfontFace
	}

	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	return freetype.ParseFont(fontBytes)
}

func createAvatar(initials string) (*image.RGBA, error) {
	// make sure the string is OK
	text := cleanString(initials)

	// Check cache
	cachedImage := getImage(text)

	if cachedImage != nil {
		return cachedImage, nil
	}

	f, err := getFont(fontFacePath)
	if err != nil {
		return nil, err
	}

	textColor := image.White
	background := defaultColor(text[0:1])
	rgba := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(rgba, rgba.Bounds(), &background, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(textColor)
	c.SetHinting(font.HintingFull)

	to := truetype.Options{}
	to.Size = fontSize
	face := truetype.NewFace(f, &to)

	xPoints := []int{0, 0}
	textWidths := []int{0, 0}

	for i, char := range text {
		width, ok := face.GlyphAdvance(rune(char))
		if !ok {
			return nil, err
		}

		textWidths[i] = int(float64(width) / 64)
	}

	if len(textWidths) == 1 {
		textWidths[1] = 0
	}

	combinedWidth := textWidths[0] + spacer + textWidths[1]

	xPoints[0] = int((imageWidth - combinedWidth) / 2)
	xPoints[1] = int(xPoints[0] + textWidths[0] + spacer)

	for i, char := range text {
		pt := freetype.Pt(xPoints[i], textY)
		_, err := c.DrawString(string(char), pt)
		if err != nil {
			return nil, err
		}
	}

	setImage(text, rgba)

	return rgba, nil
}

func ToDisk(initials, path string) {
	rgba, err := createAvatar(initials)
	if err != nil {
		log.Println(err)
		return
	}

	out, err := os.Create(path)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer out.Close()

	b := bufio.NewWriter(out)

	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func ToHTTP(initials string, w http.ResponseWriter) {
	rgba, err := createAvatar(initials)
	if err != nil {
		log.Println(err)
		return
	}

	b := new(bytes.Buffer)
	key := fmt.Sprintf("avatar%s", initials)

	err = png.Encode(b, rgba)
	if err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(b.Bytes())))
	w.Header().Set("Cache-Control", "max-age=2592000")
	w.Header().Set("Etag", `"` +key+ `"`)

	if _, err := w.Write(b.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func main() {
	ToDisk("ZL", "zl.png")
}