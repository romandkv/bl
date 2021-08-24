package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"gitlab.equifax.local/dro/bl/pkg/engine"
	"gitlab.equifax.local/dro/bl/pkg/filehelper"
	"gitlab.equifax.local/dro/bl/pkg/interpeter"
	"golang.org/x/image/colornames"
)


func main() {
	filehelper.Helper.RootFolder = interpeter.DEFAULT_PROJECTS_PATH + "\\new_project"
	engine.VNE.Files = interpeter.GetProjectStructure("new_project")
	label, err := interpeter.GetLabelFile(engine.VNE.Files, interpeter.LABEL_START)
	if err != nil {
		return
	}
	interpeter.Run(label, interpeter.LABEL_START)
	pixelgl.Run(Run)
}


func Run() {
	var err error
	cfg := pixelgl.WindowConfig{
		Title:  "BL!",
		Bounds: pixel.R(0, 0, 1080, 600),
		VSync: true,
	}
	engine.VNE.Window, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	engine.VNE.Window.Clear(colornames.White)
	interpeter.PushCommands(engine.VNE.Queue[interpeter.LABEL_START])
	for !engine.VNE.Window.Closed() {
		if engine.VNE.Window.JustPressed(pixelgl.KeySpace) {
			fmt.Println("PRESSED SPACE")
			if engine.VNE.MainQ.Len() != 0 {
				engine.VNE.MainQ.PopFront().(engine.Action).Execute()
			}
			fmt.Println("Очередь: " + fmt.Sprint(engine.VNE.MainQ.Len()))
			engine.VNE.UpdateScreen()
		}
		engine.VNE.Window.Update()
	}
}
