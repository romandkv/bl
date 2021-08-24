package actions

import (
	"errors"
	"fmt"
	"github.com/faiface/pixel"
	"gitlab.equifax.local/dro/bl/pkg/engine"
	"gitlab.equifax.local/dro/bl/pkg/models"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"strings"
)

const (
	CHARACTER_CONSTRUCTOR = "Character"
	SPRITE_CONSTRUCTOR = "Sprite"
	SOUND_CONSTRUCTOR = "Sound"
	BACK_CONSTRUCTOR = "Background"
)

type Say struct {
	Character string
	Speech string
}

func (action Say) Execute() {
	character := engine.VNE.Characters[action.Character]
	engine.VNE.CurrentScreen.BottomTextBox.CharacterName.Text.Clear()
	engine.VNE.CurrentScreen.BottomTextBox.CharacterName.Text.Color = character.Color

	fmt.Fprintf(engine.VNE.CurrentScreen.BottomTextBox.CharacterName.Text, character.Name)

	engine.VNE.CurrentScreen.BottomTextBox.Speech.Text.Clear()
	engine.VNE.CurrentScreen.BottomTextBox.Speech.Text.Color = colornames.White
	fmt.Fprintf(engine.VNE.CurrentScreen.BottomTextBox.Speech.Text, "\n" + InsertNewLines(
			action.Speech,
			engine.VNE.CurrentScreen.BottomTextBox.TextBox.Box.Frame(),
			basicfont.Face7x13.Advance,
			int(engine.VNE.CurrentScreen.BottomTextBox.Speech.Scale),
		),
	)
	fmt.Println(action.Speech)
}

type Show struct {
	Sprite string
}

func (action Show) Execute() {
	fmt.Println(action.Sprite + " sprite shows up!\n")
	sprite := engine.VNE.Sprites[action.Sprite]
	engine.VNE.CurrentScreen.Sprites = append(engine.VNE.CurrentScreen.Sprites, sprite)
}

type Scene struct {
	Back string
}

func (action Scene) Execute() {
	fmt.Println(action.Back + " back shows up!\n")
	back := engine.VNE.Backgrounds[action.Back]
	engine.VNE.CurrentScreen.Background = &back
}

type Play struct {
	Sound string
}
//TODO play
func (action Play) Execute() {
	fmt.Println(action.Sound + " plays!\n")
}

func ParsePlay(word, label string) error {
	if word == "" {
		return errors.New("Для вызова play требуется 1 аргумент\n")
	}
	if _, ok := engine.VNE.Sounds[word]; !ok {
		return errors.New("Переменная" + word + "не существует\n")
	}
	queue := engine.VNE.Queue[label]
	queue = append(queue, Play{
		Sound: word,
	})
	engine.VNE.Queue[label] = queue
	return nil
}

func ParseSpeech(char, words, label string) error {
	if _, ok := engine.VNE.Characters[char]; !ok {
		return errors.New("Переменная персонажа не определена: " + char)
	}
	value, err := getStringValue(words)
	if err != nil {
		return err
	}
	queue := engine.VNE.Queue[label]
	queue = append(queue, Say{
		Speech: value,
		Character: char,
	})
	engine.VNE.Queue[label] = queue
	return nil
}

func ParseShow(word, label string) error {
	if word == "" {
		return errors.New("Для вызова show требуется 1 аргумент\n")
	}
	if _, ok := engine.VNE.Sprites[word]; !ok {
		return errors.New("Переменная" + word + "не существует\n")
	}
	queue := engine.VNE.Queue[label]
	queue = append(queue, Show{
		Sprite: word,
	})
	engine.VNE.Queue[label] = queue
	return nil
}

func ParseScene(word, label string) error {
	if word == "" {
		return errors.New("Для вызова scene требуется 1 аргумент\n")
	}
	if _, ok := engine.VNE.Backgrounds[word]; !ok {
		return errors.New("Переменная" + word + "не существует\n")
	}
	queue := engine.VNE.Queue[label]
	queue = append(queue, Scene{
		Back: word,
	})
	engine.VNE.Queue[label] = queue
	return nil
}

func parseParamsToMap(params string) (map[string]string, error) {
	var arguments map[string]string = make(map[string]string)
	var err error
	components := strings.Split(params, ",")

	for _, component := range components {
		argument := strings.Split(component, "=")
		if len(argument) != 2 {
			return nil, errors.New("Ошибка парсинга параметров конструктора")
		}
		argument[1], err = getStringValue(argument[1])
		if err != nil {
			return nil, err
		}
		arguments[argument[0]] = argument[1]
	}
	return arguments, nil
}

func ParseDefineLine(word, filePath string) error {
	words := strings.Split(word, " ")

	if len(words) < 3 || words[1] != "=" {
		return errors.New("Ошибка в выражении define (define name = Constructor(...))\n")
	}
	name := words[0]
	constructor := words[2][:strings.IndexByte(words[2], '(')]
	params := strings.Join(words[2:], "")
	arguments, err := parseParamsToMap(params[strings.IndexByte(params, '(') + 1 : strings.IndexByte(params, ')')])
	if err != nil {
		return errors.New(err.Error() + " " + constructor + "\n")
	}
	switch constructor {
	case CHARACTER_CONSTRUCTOR:
		char, err := models.MapToCharacter(arguments)
		if err != nil {
			return err
		}
		if _, ok := engine.VNE.Characters[name]; ok {
			return errors.New("Переменная персонажа с таким именем уже существует, константы невозможно переопеределить\n")
		}
		engine.VNE.Characters[name] = *char
		return nil
	case SOUND_CONSTRUCTOR:
		sound, err := models.MapToSound(arguments, filePath)
		if err != nil {
			return err
		}
		if _, ok := engine.VNE.Sounds[name]; ok {
			return errors.New("Переменная аудио с таким именем уже существует, константы невозможно переопеределить\n")
		}
		engine.VNE.Sounds[name] = *sound
		return nil
	case BACK_CONSTRUCTOR:
		back, err := models.MapToBackground(arguments, filePath)
		if err != nil {
			return err
		}
		if _, ok := engine.VNE.Backgrounds[name]; ok {
			return errors.New("Переменная бэкграунда с таким именем уже существует, константы невозможно переопеределить\n")
		}
		engine.VNE.Backgrounds[name] = *back
		return nil
	case SPRITE_CONSTRUCTOR:
		sprite, err := models.MapToSprite(arguments, filePath)
		if err != nil {
			return err
		}
		if _, ok := engine.VNE.Sprites[name]; ok {
			return errors.New("Переменная спрайта с таким именем уже существует, константы невозможно переопеределить\n")
		}
		engine.VNE.Sprites[name] = *sprite
		return nil
	}
	return errors.New("Конструтора с таким названием не опеределен:" + constructor + "\n")
}

func getStringValue(expr string) (string, error) {
	expr = strings.TrimSpace(expr)
	speechWithOutQuotes := expr[1 : len(expr) - 1]

	if "\"" + speechWithOutQuotes + "\"" != expr {
		return "", errors.New("Аргумент должен являться строкой, заключенной в двойные ковычки")
	}
	return speechWithOutQuotes, nil
}

//InsertNewLines Inserting newlines to avoid text to get out of bounds
func InsertNewLines(line string, bounds pixel.Rect, fontWidth, scale int) string {
	var count int
	var speech string
	for _, char := range line {
		count += fontWidth * scale
		if count >= int(bounds.W()) {
			speech += "\n"
			count = fontWidth * scale
		}
		speech += string(char)
	}

	return speech
}
