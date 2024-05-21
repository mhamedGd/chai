package main

import (
	chai "github.com/mhamedGd/chai/Chai"
)

var game chai.App

func main() {

	game = chai.App{
		Width:  1920,
		Height: 1080,
		Title:  "Test",

		OnStart: func() {
			chai.LogF("STARTED\n")
			default_scene := chai.NewScene()
			default_scene.OnSceneStart = func(thisScene *chai.Scene) {
				thisScene.Background = chai.NewRGBA8(0, 0, 0, 255)

				thisScene.NewRenderSystem(&chai.SpriteRenderSystem{Sprites: &chai.Sprites, Scale: 1.0})
				logo_id := thisScene.NewEntityId()
				logo_transform := chai.Transform{
					Position:   chai.Vector2fZero,
					Dimensions: chai.Vector2fOne,
					Scale:      1.0,
					Rotation:   0.0,
				}
				logo_sprite := chai.SpriteComponent{
					Texture: chai.LoadPng("Assets/Chai_Logo_transparent.png", &chai.TextureSettings{Filter: chai.TEXTURE_FILTER_LINEAR}),
					Tint:    chai.WHITE,
				}

				thisScene.AddComponents(logo_id, chai.ToComponent(logo_transform), chai.ToComponent(logo_sprite))
			}

			chai.ChangeScene(&default_scene)
		},
		OnUpdate: func(dt float32) {

		},
		OnDraw: func(dt float32) {

		},
		OnEvent: func(ae *chai.AppEvent) {
		},
	}

	chai.Run(&game)

}
