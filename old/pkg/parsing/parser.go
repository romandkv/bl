package parsing

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/faiface/pixel"
	engine2 "gitlab.equifax.local/dro/bl/old/pkg/engine"
	models2 "gitlab.equifax.local/dro/bl/old/pkg/models"
	"gitlab.equifax.local/dro/bl/pkg/actions"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	DEFAULT_PROJECTS_PATH = "projects/"
	EXTENSION = ".bl"

	JUMP_STATEMENT = "jump"
	DEFINE_STATEMENT = "define"
	SHOW_STATEMENT = "show"
	SCENE_STATEMENT = "scene"
	LABEL_STATEMENT = "label"
	PLAY_STATEMENT = "play"
	LABEL_START = "start"

	CHARACTER_CONSTRUCTOR = "Character"
	SPRITE_CONSTRUCTOR = "Sprite"
	SOUND_CONSTRUCTOR = "Sound"
	BACK_CONSTRUCTOR = "Background"
)

type Label struct {
	Line int
	Path string
	Name string
}

// Получение всех путей исполняемых файлов проекта
func GetProjectStructure(projectName string) []string {
	var files []string

	projectName =  DEFAULT_PROJECTS_PATH + projectName

	filepath.Walk(projectName, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != EXTENSION {
			return nil
		}
		files = append(files, path)
		return nil
	})

	return files
}

// Получения файла проекта, в котором находится точка входа
func GetLabelFile(files []string, label string) (*Label, error) {
	var entries []Label
	var waitGroup sync.WaitGroup
	var errorMessage string

	for _, path := range files {
		waitGroup.Add(1)
		go func(path string) {
			var currentLine int = 0
			file, err := os.Open(path)
			if err != nil {
				return
			}
			defer file.Close()
			defer waitGroup.Done()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				currentLine++
				if !strings.Contains(scanner.Text(), LABEL_STATEMENT+ " " + label + ":") {
					continue
				}
				entries = append(entries, Label{
					Line:currentLine,
					Path:path,
					Name: label,
				})
			}
		}(path)
	}
	waitGroup.Wait()

	if len(entries) > 1 {
		errorMessage = "Несколько label start в одном проекте:\n"
		for _, file := range entries {
			errorMessage += file.Path + " Строка: " + fmt.Sprint(file.Line) + "\n"
		}
		return &Label{}, errors.New(errorMessage)
	}
	if len(entries) == 0 {
		return &Label{}, errors.New("Не найден label start ")
	}
	return &entries[0], nil
}

func isLabel(line string) (bool, error) {
	var components []string

	components = strings.Split(line, " ")
	if len(components) != 2 {
		return false, nil
	}
	if components[0] != LABEL_STATEMENT {
		return false, nil
	}
	if components[1][len(components[1]) - 1:] != ":" {
		return true, errors.New("Необходим знак ':' после значения label <label_name>:\n-+")
	}
	return true, nil
}

func parseLine(line string, nesting int, filePath string) error {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	statement := line[0 : strings.IndexByte(line, ' ')]
	value := line[strings.IndexByte(line, ' ') + 1 : ]
	if statement[0] == '#' {
		return nil
	}

	if statement != DEFINE_STATEMENT && nesting == 0 {
		return errors.New("Использование " + statement + " недопустимо внутри label структуры\n")
	}
	if value == "" {
		return errors.New("Нет аргументов для  " + statement + "\n")
	}
	switch statement {
		case JUMP_STATEMENT:
			return parseJump(value)
		case DEFINE_STATEMENT:
			return parseDefineLine(value, filePath)
		case SHOW_STATEMENT:
			return parseShow(value)
		case SCENE_STATEMENT:
			return parseScene(value)
		case PLAY_STATEMENT:
			return parsePlay(value)
	}
	return parseSpeech(statement, value)
}

func parsePlay(word string) error {
	if word == "" {
		return errors.New("Для вызова play требуется 1 аргумент\n")
	}
	if _, ok := engine2.VNE.Sounds[word]; !ok {
		return errors.New("Переменная" + word + "не существует\n")
	}
	engine2.VNE.Queue.PushBack(actions.Play{
		Sound: word,
	})
	return nil
}

func parseShow(word string) error {
	if word == "" {
		return errors.New("Для вызова show требуется 1 аргумент\n")
	}
	if _, ok := engine2.VNE.Sprites[word]; !ok {
		return errors.New("Переменная" + word + "не существует\n")
	}
	engine2.VNE.Queue.PushBack(actions.Show{
		Sprite: word,
	})
	return nil
}

func parseScene(word string) error {
	if word == "" {
		return errors.New("Для вызова scene требуется 1 аргумент\n")
	}
	if _, ok := engine2.VNE.Backgrounds[word]; !ok {
		return errors.New("Переменная" + word + "не существует\n")
	}
	engine2.VNE.Queue.PushBack(actions.Scene{
		Back: word,
	})
	return nil
}

func getStringValue(expr string) (string, error) {
	expr = strings.TrimSpace(expr)
	speechWithOutQuotes := expr[1 : len(expr) - 1]

	if "\"" + speechWithOutQuotes + "\"" != expr {
		return "", errors.New("Аргумент должен являться строкой, заключенной в двойные ковычки\n")
	}
	return speechWithOutQuotes, nil
}

func parseSpeech(char, word string) error {
	if word == "" {
		return errors.New("Для вызова speech требуется 1 аргумент\n")
	}
	speechWithOutQuotes, err := getStringValue(word)
	if err != nil {
		return err
	}
	if _, ok := engine2.VNE.Characters[char]; !ok {
		return errors.New("Переменная" + char + "не существует\n")
	}
	engine2.VNE.Queue.PushBack(actions.Say{
		Character: char,
		Speech:  speechWithOutQuotes,
	})
	return nil
}



func parseJump(label string) error {
	engine2.VNE.Queue.PushBack(actions.Jump{
		Label: label,
	})
	return nil
}

//TODO уменьшить метод
func parseDefineLine(word string, filePath string) error {
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
			char, err := models2.MapToCharacter(arguments)
			if err != nil {
				return err
			}
			if _, ok := engine2.VNE.Characters[name]; ok {
				return errors.New("Переменная персонажа с таким именем уже существует, константы невозможно переопеределить\n")
			}
			engine2.VNE.Characters[name] = *char
			return nil
		case SOUND_CONSTRUCTOR:
			sound, err := models2.MapToSound(arguments, filePath)
			if err != nil {
				return err
			}
			if _, ok := engine2.VNE.Sounds[name]; ok {
				return errors.New("Переменная аудио с таким именем уже существует, константы невозможно переопеределить\n")
			}
			engine2.VNE.Sounds[name] = *sound
			return nil
		case BACK_CONSTRUCTOR:
			back, err := models2.MapToBackground(arguments, filePath)
			if err != nil {
				return err
			}
			if _, ok := engine2.VNE.Backgrounds[name]; ok {
				return errors.New("Переменная бэкграунда с таким именем уже существует, константы невозможно переопеределить\n")
			}
			engine2.VNE.Backgrounds[name] = *back
			return nil
		case SPRITE_CONSTRUCTOR:
			sprite, err := models2.MapToSprite(arguments, filePath)
			if err != nil {
				return err
			}
			if _, ok := engine2.VNE.Sprites[name]; ok {
				return errors.New("Переменная спрайта с таким именем уже существует, константы невозможно переопеределить\n")
			}
			engine2.VNE.Sprites[name] = *sprite
			return nil
	}
	return errors.New("Конструтора с таким названием не опеределен:" + constructor + "\n")
}

func parseLabel(scanner *bufio.Scanner, nesting, line int, filePath string) {
	nesting++
	for scanner.Scan() {
		line++
		strlen := len(scanner.Text())

		standardStrlen := len(strings.Repeat("    ", nesting) + strings.TrimSpace(scanner.Text()))
		if standardStrlen > strlen {
			return
		}
		if standardStrlen > strlen && standardStrlen - strlen != 4 {
			log.Fatal("Неверное количество пробелов (!= 4) Файл: " + filePath + " Строка: " + fmt.Sprint(line) + "\n")
		}
		isLabel, err := isLabel(scanner.Text())
		if err != nil || isLabel{
			log.Fatal("label должен быть корневым элементом файла" + "\n" + " Файл: " + filePath + " Строка: " + fmt.Sprint(line) + "\n")
		}
		err = parseLine(scanner.Text(), nesting, filePath)
		if err != nil {
			log.Fatal(err.Error() + " Файл: " + filePath + " Строка: " + fmt.Sprint(line) + "\n")
		}
	}
}

func Run(startFile Label, startLabel bool) {
	var line int  = 0
	var nesting int = 0
	var afterLabel bool = false

	file, err := os.Open(startFile.Path)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for true {
		if !afterLabel {
			if !scanner.Scan() {
				break
			}
		}
		if scanner.Text() == "label tag3:" {
			fmt.Println("go")
		}
		afterLabel = false
		line++
		isLabel, err := isLabel(scanner.Text())
		if err != nil{
			log.Fatal(err.Error() + " " + startFile.Path + "\n")
		}
		if isLabel && line == startFile.Line {
			parseLabel(scanner, nesting, line, startFile.Path)
			afterLabel = true
			continue
		}
		if isLabel {
			engine2.VNE.Cache = true
			engine2.VNE.CacheLableName = startFile.Name
			parseLabel(scanner, nesting, line, startFile.Path)
			engine2.VNE.Cache = false
			continue
		}
		err = parseLine(scanner.Text(), nesting, startFile.Path)
		if err != nil {
			log.Fatal(err.Error() + " File:" + startFile.Path + " Строка: " + fmt.Sprint(line) + "\n")
		}
	}
	if startLabel {
		engine2.VNE.LoadedNextChapter <- true
		engine2.VNE.EndOfGame <- true
		return
	}
}

