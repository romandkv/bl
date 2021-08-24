package models

import (
	"errors"
	"github.com/faiface/pixel"
	"gitlab.equifax.local/dro/bl/pkg/filehelper"
	"os"
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
	Color string
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
	return &Character{
		Name: arguments[NAME],
		Color: arguments[CHARACTER_COLOR],
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

