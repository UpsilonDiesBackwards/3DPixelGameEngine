package main

import (
	"3DPixelGameEngine/engine/io"
	"3DPixelGameEngine/engine/obj"
	"3DPixelGameEngine/engine/rendering"
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
)

func main() {
	window, err := rendering.NewWindow(800, 600, "3D Renderer")
	if err != nil {
		log.Fatal("Could not create new window: ", err)
	}
	defer glfw.Terminate()

	renderer := rendering.NewRenderer(window)

	model, err := obj.LoadModel("engine/res/models/cube.obj", "engine/res/models/cube.mtl")
	if err != nil {
		fmt.Println("fail to load model", err)
	}
	object := rendering.NewObject(model)
	renderer.Objects["cube"] = object

	for !window.ShouldClose() {
		deltaTime := renderer.CalculateDeltaTime()

		io.InputRunner(window, deltaTime)

		renderer.Draw()

		window.PollEvents()
		window.SwapBuffers()
	}
}
