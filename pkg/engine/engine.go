package engine

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/gammazero/deque"
	"gitlab.equifax.local/dro/bl/pkg/filehelper"
	"gitlab.equifax.local/dro/bl/pkg/models"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type Engine struct {
	Characters map[string]models.Character
	Sounds map[string]models.Sound
	Backgrounds map[string]models.Background
	Sprites map[string]models.Sprite
	Files []string
	MainQ deque.Deque
	Queue map[string][]Action
	Window *pixelgl.Window
	CurrentScreen Screen
}

func newEngine() Engine {
	return Engine{
		Characters:  make(map[string]models.Character),
		Sounds:      make(map[string]models.Sound),
		Backgrounds: make(map[string]models.Background),
		Sprites:     make(map[string]models.Sprite),
		Queue:       make(map[string][]Action),
		CurrentScreen: newScreen(),
	}
}

var VNE = newEngine()

type Action interface {
	Execute()
}

type Screen struct {
	Background *models.Background
	Sprites	[]models.Sprite
	BottomTextBox BottomTextBox
}

type BottomTextBox struct {
	CharacterName CharacterName
	Speech Speech
	TextBox TextBox
}

type CharacterName struct {
	Text  *text.Text
	Scale float64
	Position pixel.Vec
}

type Speech struct {
	Text  *text.Text
	Scale float64
	Position pixel.Vec
}

type TextBox struct {
	Box *pixel.Sprite
	CenterPosition pixel.Vec
}

func newScreen() Screen {
	picture, err := filehelper.LoadPicture("C:\\Users\\dro\\go\\src\\bl\\files\\img.png")
	if err != nil {
		return Screen{}
	}

	return Screen{
		BottomTextBox: BottomTextBox{
			CharacterName: CharacterName{
				Text: text.New(pixel.V(0, 0), text.NewAtlas(basicfont.Face7x13, text.ASCII)),
				Scale: 2,
			},
			Speech:  Speech{
				Text: text.New(pixel.V(0, 0), text.NewAtlas(basicfont.Face7x13, text.ASCII)),
				Scale: 2,
			},
			TextBox: TextBox{
				Box:            pixel.NewSprite(picture, picture.Bounds()),
			},
		},
		Background:	nil,
	}
}
func (eng *Engine) CalculatePositions() {
	eng.CurrentScreen.BottomTextBox.TextBox.CenterPosition = VNE.Window.Bounds().Center().Add(
		pixel.V(
			0,
			eng.CurrentScreen.BottomTextBox.TextBox.Box.Frame().H() / 2 - VNE.Window.Bounds().H() / 2,
		),
	)
	eng.CurrentScreen.BottomTextBox.CharacterName.Position = eng.CurrentScreen.BottomTextBox.TextBox.CenterPosition.Add(
		pixel.V(
			-VNE.CurrentScreen.BottomTextBox.TextBox.Box.Frame().W() / 2,
			+VNE.CurrentScreen.BottomTextBox.TextBox.Box.Frame().H() / 2 - VNE.CurrentScreen.BottomTextBox.CharacterName.Text.LineHeight * 2,
		),
	)
	eng.CurrentScreen.BottomTextBox.Speech.Position = eng.CurrentScreen.BottomTextBox.TextBox.CenterPosition.Add(
		pixel.V(
			-VNE.CurrentScreen.BottomTextBox.TextBox.Box.Frame().W() / 2,
			+VNE.CurrentScreen.BottomTextBox.TextBox.Box.Frame().H() / 2 - VNE.CurrentScreen.BottomTextBox.CharacterName.Text.LineHeight * 2 - VNE.CurrentScreen.BottomTextBox.CharacterName.Text.LineHeight * 2,
		),
	)
}

func (eng Engine) UpdateScreen() {
	eng.Window.Clear(colornames.White)
	if eng.CurrentScreen.Background != nil {
		eng.CurrentScreen.Background.File.Draw(eng.Window, pixel.IM.Moved(eng.Window.Bounds().Center()))
	}
	for _, sprite := range eng.CurrentScreen.Sprites {
		sprite.File.Draw(eng.Window, pixel.IM.Moved(eng.Window.Bounds().Center()))
	}
	eng.CalculatePositions()
	eng.CurrentScreen.BottomTextBox.TextBox.Box.Draw(
		eng.Window,
		pixel.IM.Moved(eng.CurrentScreen.BottomTextBox.TextBox.CenterPosition),
	)
	eng.CurrentScreen.BottomTextBox.CharacterName.Text.Draw(
		eng.Window,
		pixel.IM.Scaled(
			eng.CurrentScreen.BottomTextBox.CharacterName.Text.Orig,
			eng.CurrentScreen.BottomTextBox.CharacterName.Scale,
		).Moved(eng.CurrentScreen.BottomTextBox.CharacterName.Position),
	)
	eng.CurrentScreen.BottomTextBox.Speech.Text.Draw(
		eng.Window,
		pixel.IM.Scaled(
			eng.CurrentScreen.BottomTextBox.Speech.Text.Orig,
			eng.CurrentScreen.BottomTextBox.Speech.Scale,
		).Moved(eng.CurrentScreen.BottomTextBox.Speech.Position),
	)
}