package filehelper

import (
	"github.com/faiface/pixel"
	_ "golang.org/x/image/bmp"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

type FileHelper struct {
	RootFolder string
}

func isAbsolutePath(filePath string) bool {
	if len(filePath) < 1 {
		return false
	}
	return filePath[0] == '\\'
}

func (fh FileHelper) GetNormalizedFilepath(filePath, sourceFilePath string) string {
	filePath = strings.Replace(filePath, "/", "\\", -1)
	if isAbsolutePath(filePath) {
		return fh.RootFolder + "\\" + filePath
	}
	return filepath.Dir(sourceFilePath) + "\\" + filePath
}

var Helper FileHelper

func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}