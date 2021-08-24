package models

import (
	"errors"
	"github.com/faiface/pixel"
	"gitlab.equifax.local/dro/bl/pkg/filehelper"
	"image/color"
	"os"
	"strconv"
)

type Model interface {

}

const (
	NAME = "name"
	FILE = "file"

	CHARACTER_COLOR = "color"
)

type Character struct {
	Name string
	Color color.RGBA
}

func MapToCharacter(arguments map[string]string) (*Character, error) {
	if len(arguments) > 2 {
		return &Character{}, errors.New("Character принимает два параметра: name, color\n")
	}
	_, ok := arguments[NAME]
	if !ok {
		return &Character{}, errors.New("Character: Обязательный параметр name отсутствует\n")
	}
	_, ok = arguments[CHARACTER_COLOR]
	if !ok {
		return &Character{}, errors.New("Character: Обязательный параметр color отсутствует\n")
	}
	color, err := ParseColor(arguments[CHARACTER_COLOR])
	if err != nil {
		return nil, err
	}
	return &Character{
		Name: arguments[NAME],
		Color: *color,
	}, nil
}

func ParseColor(line string) (*color.RGBA, error) {
	var err error
	if line[0] != '#' {
		return nil, err
	}
	r, err := strconv.ParseUint(line[1:3], 16, 8)
	if err != nil {
		return nil, err
	}
	g, err := strconv.ParseUint(line[3:5], 16, 8)
	if err != nil {
		return nil, err
	}
	b, err := strconv.ParseUint(line[5:7], 16, 8)
	if err != nil {
		return nil, err
	}
	return &color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 0,
	}, nil
}

type Sprite struct {
	File *pixel.Sprite
}

func MapToSprite(arguments map[string]string, sourceFilePath string) (*Sprite, error) {
	if len(arguments) > 2 {
		return &Sprite{}, errors.New("Sprite принимает дин параметр: file\n")
	}
	file, ok := arguments[FILE]
	if !ok {
		return &Sprite{}, errors.New("Sprite: Обязательный параметр file отсутствует\n")
	}
	file = filehelper.Helper.GetNormalizedFilepath(file, sourceFilePath)
	if !fileExists(file) {
		return &Sprite{}, errors.New("Sprite: Файл " + file + " не найден\n")
	}
	picture, err := filehelper.LoadPicture(file)
	if err != nil {
		return nil, err
	}
	return &Sprite{
		File: pixel.NewSprite(picture, picture.Bounds()),
	}, nil
}


type Background struct {
	File *pixel.Sprite
}

func MapToBackground(arguments map[string]string, sourceFilePath string) (*Background, error) {
	if len(arguments) > 2 {
		return &Background{}, errors.New("Background принимает дин параметр: file\n")
	}
	file, ok := arguments[FILE]
	if !ok {
		return &Background{}, errors.New("Background: Обязательный параметр file отсутствует\n")
	}
	file = filehelper.Helper.GetNormalizedFilepath(file, sourceFilePath)
	if !fileExists(file) {
		return &Background{}, errors.New("Background: Файл " + file + " не найден\n")
	}
	picture, err := filehelper.LoadPicture(file)
	if err != nil {
		return nil, err
	}
	return &Background{
		File: pixel.NewSprite(picture, picture.Bounds()),
	}, nil
}

type Sound struct {
	File *pixel.Sprite
}

func MapToSound(arguments map[string]string, sourceFilePath string) (*Sound, error) {
	if len(arguments) > 2 {
		return &Sound{}, errors.New("Sound принимает дин параметр: file\n")
	}
	file, ok := arguments[FILE]
	if !ok {
		return &Sound{}, errors.New("Sound: Обязательный параметр file отсутствует\n")
	}
	file = filehelper.Helper.GetNormalizedFilepath(file, sourceFilePath)

	if !fileExists(file) {
		return &Sound{}, errors.New("Sound: Файл " + file + " не найден\n")
	}
	picture, err := filehelper.LoadPicture(file)
	if err != nil {
		return nil, err
	}
	return &Sound{
		File: pixel.NewSprite(picture, picture.Bounds()),
	}, nil
}

type Narrator struct {

}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		return true
	}
	return !os.IsNotExist(err)
}

