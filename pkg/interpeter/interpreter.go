package interpeter

import (
	"bufio"
	"errors"
	"fmt"
	"gitlab.equifax.local/dro/bl/pkg/actions"
	"gitlab.equifax.local/dro/bl/pkg/engine"
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


)


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
func GetLabelFile(files []string, label string) (string, error) {
	var entries []string
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
				entries = append(entries, path)
			}
		}(path)
	}
	waitGroup.Wait()

	if len(entries) > 1 {
		errorMessage = "Несколько label start в одном проекте:\n"
		for _, file := range entries {
			errorMessage += file + " Строка: " + fmt.Sprint(file) + "\n"
		}
		return "", errors.New(errorMessage)
	}
	if len(entries) == 0 {
		return "", errors.New("Не найден label start ")
	}
	return entries[0], nil
}

func getLabelName(line string) (string, error) {
	line = strings.TrimSpace(line)
	space := strings.IndexByte(line, ' ')
	colons := strings.IndexByte(line, ':')

	if space == -1 {
		return "", nil
	}
	if line[0 : space] != LABEL_STATEMENT {
		return "", nil
	}
	if colons == -1 {
		return "", errors.New("После label <string> должно следовать двоеточие")
	}
	return strings.TrimSpace(line[space + 1 : colons]), nil
}



func throwError(err error, file string, line int) {
	log.Fatal(err.Error() + " File: " + file + " Line: " + fmt.Sprint(line))
}

func parseLine(line string, nesting int, filePath, label string) error {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	space := strings.IndexByte(line, ' ')
	if space == -1 {
		return errors.New("Unrecognized command: " + line)
	}
	statement := line[0 : space]
	value := line[space + 1 : ]
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
		return ParseJump(value, label)
	case DEFINE_STATEMENT:
		return actions.ParseDefineLine(value, filePath)
	case SHOW_STATEMENT:
		return actions.ParseShow(value, label)
	case SCENE_STATEMENT:
		return actions.ParseScene(value, label)
	case PLAY_STATEMENT:
		return actions.ParsePlay(value, label)
	}
	return actions.ParseSpeech(statement, value, label)
}

func Run(startFile string, label string) {
	var line int
	var nest int

	file, err := os.Open(startFile)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan()  {
		line++
		labelName, err := getLabelName(scanner.Text())
		if err != nil{
			throwError(err, startFile, line)
		}
		if labelName == "" {
			err = parseLine(scanner.Text(), nest, startFile, "")
			if err != nil {
				throwError(err, startFile, line)
			}
			continue
		}
		parseLabel(scanner, nest, line, startFile, labelName)
	}
}

func checkSpaces(line string, nest int) (int, error) {
	count := 0

	for _, char := range line {
		if char != ' ' {
			break
		}
		count++
	}
	if count % 4 != 0 || count > nest * 4 {
		return 1, errors.New("Неверное количество пробелов")
	}
	if count == nest * 4 {
		return 0, nil
	}
	return -1, nil
}

func parseLabel(scanner *bufio.Scanner, nest, line int, filePath, label string) {
	nest++
	for scanner.Scan() {
		line++
		result, err := checkSpaces(scanner.Text(), nest)
		if err != nil {
			throwError(err, filePath, line)
		}
		if result == -1 {
			return
		}
		labelName, err := getLabelName(scanner.Text())
		if err != nil {
			throwError(err, filePath, line)

		}
		if labelName != "" {
			throwError(errors.New("label должен быть корневым элементом файла"), filePath, line)
		}
		err = parseLine(scanner.Text(), nest, filePath, label)
		if err != nil {
			throwError(err, filePath, line)
		}
	}
}

type Jump struct {
	Label string
}
func (action Jump) Execute() {
	fmt.Println("Загружаем новое")
	if _, ok := engine.VNE.Queue[action.Label]; !ok {
		file, err := GetLabelFile(engine.VNE.Files, action.Label)
		if err != nil {
			log.Fatalln("Ploho")
		}
		Run(file, action.Label)
	}
	PushCommands(engine.VNE.Queue[action.Label])
	engine.VNE.MainQ.PopFront().(engine.Action).Execute()
}

func ParseJump(value, label string) error {
	queue := engine.VNE.Queue[label]
	queue = append(queue, Jump{
		Label: value,
	})
	engine.VNE.Queue[label] = queue
	return nil
}

func PushCommands(commands []engine.Action) {
	for _, command := range commands {
		engine.VNE.MainQ.PushBack(command)
	}
}