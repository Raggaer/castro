package util

import (
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/goimage"
	"image/color"
	"path/filepath"
	"strconv"
)

// CreatePlayerSignature creates a player image signature file using goimage
func CreatePlayerSignature(player models.Player) error {
	// Create signature image
	img := goimage.NewImage(500, 150)

	// Set image background
	if err := img.SetBackGroundImage("public/images/signature-bg.png"); err != nil {
		return err
	}

	// Add signature text
	if err := img.WriteText("Name: "+player.Name, color.Black, 14, 10, 20); err != nil {
		return err
	}

	if err := img.WriteText("Level: "+strconv.Itoa(player.Level), color.Black, 14, 10, 40); err != nil {
		return err
	}

	// Save image
	return img.Save(filepath.Join("public", "images", "signature", player.Name) + ".png")
}
