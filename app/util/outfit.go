package util

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strconv"

	"github.com/anthonynsimon/bild/blend"
	"github.com/lucasb-eyer/go-colorful"
)

var outfitColors = []string{
	"FFFFFF", "FFD4BF", "FFE9BF", "FFFFBF", "E9FFBF", "D4FFBF",
	"BFFFBF", "BFFFD4", "BFFFE9", "BFFFFF", "BFE9FF", "BFD4FF",
	"BFBFFF", "D4BFFF", "E9BFFF", "FFBFFF", "FFBFE9", "FFBFD4",
	"FFBFBF", "DADADA", "BF9F8F", "BFAF8F", "BFBF8F", "AFBF8F",
	"9FBF8F", "8FBF8F", "8FBF9F", "8FBFAF", "8FBFBF", "8FAFBF",
	"8F9FBF", "8F8FBF", "9F8FBF", "AF8FBF", "BF8FBF", "BF8FAF",
	"BF8F9F", "BF8F8F", "B6B6B6", "BF7F5F", "BFAF8F", "BFBF5F",
	"9FBF5F", "7FBF5F", "5FBF5F", "5FBF7F", "5FBF9F", "5FBFBF",
	"5F9FBF", "5F7FBF", "5F5FBF", "7F5FBF", "9F5FBF", "BF5FBF",
	"BF5F9F", "BF5F7F", "BF5F5F", "919191", "BF6A3F", "BF943F",
	"BFBF3F", "94BF3F", "6ABF3F", "3FBF3F", "3FBF6A", "3FBF94",
	"3FBFBF", "3F94BF", "3F6ABF", "3F3FBF", "6A3FBF", "943FBF",
	"BF3FBF", "BF3F94", "BF3F6A", "BF3F3F", "6D6D6D", "FF5500",
	"FFAA00", "FFFF00", "AAFF00", "54FF00", "00FF00", "00FF54",
	"00FFAA", "00FFFF", "00A9FF", "0055FF", "0000FF", "5500FF",
	"A900FF", "FE00FF", "FF00AA", "FF0055", "FF0000", "484848",
	"BF3F00", "BF7F00", "BFBF00", "7FBF00", "3FBF00", "00BF00",
	"00BF3F", "00BF7F", "00BFBF", "007FBF", "003FBF", "0000BF",
	"3F00BF", "7F00BF", "BF00BF", "BF007F", "BF003F", "BF0000",
	"242424", "7F2A00", "7F5500", "7F7F00", "557F00", "2A7F00",
	"007F00", "007F2A", "007F55", "007F7F", "00547F", "002A7F",
	"00007F", "2A007F", "54007F", "7F007F", "7F0055", "7F002A",
	"7F0000",
}

// GenerateOutfitImage generates an outfit image for the given values
func GenerateOutfitImage(t, feet, legs, body, head, addons int) ([]byte, error) {
	// Parse colors
	feetColor, err := colorful.Hex("#" + outfitColors[feet])

	if err != nil {
		return nil, err
	}

	legsColor, err := colorful.Hex("#" + outfitColors[legs])

	if err != nil {
		return nil, err
	}

	bodyColor, err := colorful.Hex("#" + outfitColors[body])

	if err != nil {
		return nil, err
	}

	headColor, err := colorful.Hex("#" + outfitColors[head])

	if err != nil {
		return nil, err
	}

	// Create base image
	baseImage, err := paintOutfitPart(
		filepath.Join(
			"public",
			"images",
			"outfits",
			"generator",
			strconv.Itoa(t),
			"1_1_1_3.png",
		),
		filepath.Join(
			"public",
			"images",
			"outfits",
			"generator",
			strconv.Itoa(t),
			"1_1_1_3_template.png",
		),
		feetColor,
		legsColor,
		bodyColor,
		headColor,
	)

	if err != nil {
		return nil, err
	}

	// Get first addon
	if addons >= 2 {

		// Get addon outfit
		addonImage, err := paintOutfitPart(
			filepath.Join(
				"public",
				"images",
				"outfits",
				"generator",
				strconv.Itoa(t),
				"1_1_2_3.png",
			),
			filepath.Join(
				"public",
				"images",
				"outfits",
				"generator",
				strconv.Itoa(t),
				"1_1_2_3_template.png",
			),
			feetColor,
			legsColor,
			bodyColor,
			headColor,
		)

		if err != nil {
			return nil, err
		}

		// Draw outfit image over main image
		draw.Draw(baseImage, baseImage.Bounds(), addonImage, image.ZP, draw.Over)
	}

	// Get second addon
	if addons >= 3 {

		// Get addon outfit
		addonImage, err := paintOutfitPart(
			filepath.Join(
				"public",
				"images",
				"outfits",
				"generator",
				strconv.Itoa(t),
				"1_1_3_3.png",
			),
			filepath.Join(
				"public",
				"images",
				"outfits",
				"generator",
				strconv.Itoa(t),
				"1_1_3_3_template.png",
			),
			feetColor,
			legsColor,
			bodyColor,
			headColor,
		)

		if err != nil {
			return nil, err
		}

		// Draw outfit image over main image
		draw.Draw(baseImage, baseImage.Bounds(), addonImage, image.ZP, draw.Over)
	}

	// Image buffer
	buff := &bytes.Buffer{}

	// Encode image as png
	if err := png.Encode(buff, baseImage); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func paintOutfitPart(generator, template string, feetColor, legsColor, bodyColor, headColor color.Color) (*image.RGBA, error) {
	// Open generator image
	generatorImageFile, err := os.Open(generator)

	if err != nil {
		return nil, err
	}

	// Close generator image
	defer generatorImageFile.Close()

	// Open template image
	templateImageFile, err := os.Open(template)

	if err != nil {
		return nil, err
	}

	// Close template image
	defer templateImageFile.Close()

	// Get generator image from file
	generatorImage, err := png.Decode(generatorImageFile)

	if err != nil {
		return nil, err
	}

	// Get template image from file
	templateImageRaw, err := png.Decode(templateImageFile)

	if err != nil {
		return nil, err
	}

	// Convert raw image to RGBA
	templateImage := templateImageRaw.(*image.NRGBA)

	// Parse feet color

	// Paint all pixels
	paintPixels(templateImage, color.RGBA{0, 0, 255, 255}, feetColor)
	paintPixels(templateImage, color.RGBA{0, 255, 0, 255}, legsColor)
	paintPixels(templateImage, color.RGBA{255, 0, 0, 255}, bodyColor)
	paintPixels(templateImage, color.RGBA{255, 255, 0, 255}, headColor)

	// Use multiple blend mode to main image
	outfitImage := blend.Multiply(templateImage, generatorImage)

	return outfitImage, nil
}

// paintPixels paints colors from an image
func paintPixels(img *image.NRGBA, base color.Color, dst color.Color) {
	br, bg, bb, ba := base.RGBA()
	dr, dg, db, _ := dst.RGBA()
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			r, g, b, a := img.At(x, y).RGBA()
			if br == r && bg == g && bb == b && ba == a {
				img.Set(x, y, color.RGBA{uint8(dr), uint8(dg), uint8(db), 255})
			}
		}
	}
}
