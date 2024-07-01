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
	imageFiles     emoji.Emoji = emoji.OpenFileFolder // todo!: maybe we could use a different emoji

	DisablePrettyPrint bool = false

	DisplayGraph = func(a ...any) {
		fmt.Println(a...)
	}
	DisplayError = func(a ...any) {
		a = append([]any{imageError}, a...)
		fmt.Println(a...)
	}

	DisplayWarning = func(a ...any) {
		a = append([]any{imageWarning}, a...)
		fmt.Println(a...)
	}
)

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
		imageFiles = "-"
	}
}

func GetImageError() emoji.Emoji {
	return imageError
}

func GetImageWarning() emoji.Emoji {
	return imageWarning
}

func GetImageDownArrow() emoji.Emoji {
	return imageDownArrow
}

func GetImageReaper() emoji.Emoji {
	return imageReaper
}

func GetImageNewspaper() emoji.Emoji {
	return imageNewspaper
}

func GetImageTree() emoji.Emoji {

	return imageTree
}

func GetImageInfo() emoji.Emoji {
	return imageInfo
}

func GetImageLineage() emoji.Emoji {
	return imageLineage
}

func GetImageFiles() emoji.Emoji {
	return imageFiles
}
