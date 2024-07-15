package main

import (
	"embed"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/DemmyDemon/go-display/labelimage"
	"github.com/golang/freetype/truetype"
)

var (
	//go:embed font
	fontEmbed embed.FS
)

func main() {
	font, err := readFont(fontEmbed, "font/ninepin.ttf")
	if err != nil {
		fmt.Printf("Unable to load font: %s\n", err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		makeImage(w, r, font)
	})
	err = http.ListenAndServe(":2467", nil)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server closed")
		} else {
			fmt.Printf("Server error: %s\n", err)
			os.Exit(1)
		}
	}
}

func readFont(fs embed.FS, fontPath string) (*truetype.Font, error) {

	fontBytes, err := fs.ReadFile(fontPath)
	if err != nil {
		return nil, fmt.Errorf("reading font: %w", err)
	}

	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing font: %w", err)
	}

	return f, nil

}

func makeImage(w http.ResponseWriter, r *http.Request, font *truetype.Font) {
	size := image.Rect(0, 0, 192, 64)
	textColor := color.RGBA{60, 128, 60, 255}
	fmt.Printf("Path is %s\n", r.URL.Path)
	basename := filepath.Base(r.URL.Path)
	text := strings.TrimSuffix(basename, filepath.Ext(basename))
	text = strings.ReplaceAll(text, "_", " ")

	if text == "" || r.URL.Path == "/" {
		text = "User error"
	}

	w.WriteHeader(http.StatusOK)

	img := labelimage.CreateWithFont(size, font, textColor, text)

	err := png.Encode(w, img)
	if err != nil {
		fmt.Printf("Error writing PNG: %s\n", err)
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "image/png")
}
