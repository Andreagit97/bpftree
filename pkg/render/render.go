package render

import (
	"fmt"

	"github.com/enescakir/emoji"
)

var (
	// These are the emoji used in pretty print mode.
	imageError     emoji.Emoji = emoji.RedCircle
	imageWarning   emoji.Emoji = emoji.Warning
	imageDownArrow emoji.Emoji = emoji.DownArrow
	imageReaper    emoji.Emoji = emoji.Skull
	imageNewspaper emoji.Emoji = emoji.RolledUpNewspaper
	imageTree      emoji.Emoji = emoji.PalmTree
	imageInfo      emoji.Emoji = emoji.Information
	imageLineage   emoji.Emoji = emoji.Scroll
	imageFolder    emoji.Emoji = emoji.FileFolder
	imageFile      emoji.Emoji = emoji.SpiralNotepad

	// DisablePrettyPrint is a flag to disable pretty print mode.
	DisablePrettyPrint bool = false

	// DisplayGraph is a function to display a log.
	DisplayGraph = func(a ...any) {
		fmt.Println(a...)
	}

	// DisplayError is a function to display an error log.
	DisplayError = func(a ...any) {
		a = append([]any{imageError}, a...)
		fmt.Println(a...)
	}

	// DisplayWarning is a function to display a warning log.
	DisplayWarning = func(a ...any) {
		a = append([]any{imageWarning}, a...)
		fmt.Println(a...)
	}
)

// ConfigureRendering sets the emoji images to be used in no pretty print mode.
func ConfigureRendering() {
	if DisablePrettyPrint {
		imageError = "X"
		imageWarning = "/!\\"
		imageDownArrow = "V"
		imageReaper = "(R)"
		imageNewspaper = "*"
		imageTree = "-"
		imageInfo = "-"
		imageLineage = "-"
		imageFolder = "*"
		imageFile = "-"
	}
}

// GetImageError returns the emoji used to display errors.
func GetImageError() emoji.Emoji {
	return imageError
}

// GetImageWarning returns the emoji used to display warnings.
func GetImageWarning() emoji.Emoji {
	return imageWarning
}

// GetImageDownArrow returns the emoji used to display down arrows.
func GetImageDownArrow() emoji.Emoji {
	return imageDownArrow
}

// GetImageReaper returns the emoji used to display reapers.
func GetImageReaper() emoji.Emoji {
	return imageReaper
}

// GetImageNewspaper returns the emoji used to display newspapers.
func GetImageNewspaper() emoji.Emoji {
	return imageNewspaper
}

// GetImageTree returns the emoji used to display trees.
func GetImageTree() emoji.Emoji {
	return imageTree
}

// GetImageInfo returns the emoji used to display information.
func GetImageInfo() emoji.Emoji {
	return imageInfo
}

// GetImageLineage returns the emoji used to display lineages.
func GetImageLineage() emoji.Emoji {
	return imageLineage
}

// GetImageFolder returns the emoji used to display folders.
func GetImageFolder() emoji.Emoji {
	return imageFolder
}

// GetImageFile returns the emoji used to display files.
func GetImageFile() emoji.Emoji {
	return imageFile
}
