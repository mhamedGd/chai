package main

import (
	"math"

	chai "github.com/mhamedGd/Chai"
)

var game chai.App

func main() {
	var rotation float32 = 90.0
	var midPoint chai.Vector2f = chai.Vector2fZero
	var midScreen chai.Vector2f

	var velocity chai.Vector2f = chai.Vector2fZero
	var direction chai.Vector2f
	speed := float32(1.0)

	var inputAxis chai.Vector2f = chai.Vector2fZero
	var dpad_modifier chai.Vector2f = chai.Vector2fZero

	var bgl_texture chai.Texture2D

	var fontAtlas chai.FontBatchAtlas
	var englishFontAtlas chai.FontBatchAtlas

	var ecs_engine chai.EcsEngine
	var entity chai.EcsEntity

	game = chai.App{
		Width:  1920,
		Height: 1080,
		Title:  "Test",
		OnStart: func() {

			chai.LogF("STARTED\n")
			ecs_engine = *chai.NewEcsEngine()

			midScreen = chai.NewVector2f(400.0, 300.0)

			chai.BindInput("Up", chai.KEY_W)
			chai.BindInput("Down", chai.KEY_S)
			chai.BindInput("Right", chai.KEY_D)
			chai.BindInput("Left", chai.KEY_A)
			chai.BindInput("Zoom In", chai.KEY_E)
			chai.BindInput("Zoom Out", chai.KEY_Q)

			chai.Shapes.LineWidth = .5

			bgl_texture = chai.LoadPng("Assets/tile_0004.png")

			chai.MainButton_Pressed.AddListener(func(i ...int) {
				dpad_modifier.Y += 1.0
			})
			chai.MainButton_Released.AddListener(func(i ...int) {
				dpad_modifier.Y -= 1.0
			})

			chai.SideButton_Pressed.AddListener(func(i ...int) {
				dpad_modifier.Y -= 1.0

			})

			chai.SideButton_Released.AddListener(func(i ...int) {
				dpad_modifier.Y += 1.0
			})

			chai.DPadLeft_Pressed.AddListener(func(i ...int) {
				dpad_modifier.X -= 1.0
			})
			chai.DPadLeft_Released.AddListener(func(i ...int) {
				dpad_modifier.X += 1.0
			})

			chai.DPadRight_Pressed.AddListener(func(i ...int) {
				dpad_modifier.X += 1.0
			})
			chai.DPadRight_Released.AddListener(func(i ...int) {
				dpad_modifier.X -= 1.0
			})

			font_settings := chai.FontBatchSettings{
				FontSize: 36, DPI: 124, CharDistance: 4, LineHeight: 36, Arabic: true,
			}
			english_font_settings := chai.FontBatchSettings{
				FontSize: 48, DPI: 124, CharDistance: 4, LineHeight: 36, Arabic: false,
			}
			fontAtlas = chai.LoadFontToAtlas("Assets/Alfont.otf", &font_settings)
			englishFontAtlas = chai.LoadFontToAtlas("Assets/m5x7.ttf", &english_font_settings)

			entity = ecs_engine.NewEntity()
			entity.Rot = 90.0
			chai.Write(&ecs_engine, &entity, Sprite{bgl_texture})

			chai.ScaleView(4)
			chai.WarningF("This is a warning!!")
		},
		OnUpdate: func(dt float32) {
			zoomAxis := 500.0 * float32(dt) * (chai.GetActionStrength("Zoom In") - chai.GetActionStrength("Zoom Out"))
			chai.IncreaseScaleU(zoomAxis)
			inputAxis.Y = chai.GetActionStrength("Up") - (chai.GetActionStrength("Down"))
			inputAxis.X = chai.GetActionStrength("Right") - (chai.GetActionStrength("Left"))

			velocity.X = chai.LerpFloat32(velocity.X, (inputAxis.X+dpad_modifier.X)*speed, float32(dt)*2.5)
			velocity.Y = chai.LerpFloat32(velocity.Y, (inputAxis.Y+dpad_modifier.Y)*speed, float32(dt)*2.5)

			rotation -= dt * 600.0 * velocity.X
			rotation = float32(math.Mod(float64(rotation), 360))
			direction = chai.Vector2fRight.Rotate(rotation, chai.Vector2fZero)

			midPoint.Y += velocity.Y * direction.Y
			midPoint.X += velocity.Y * direction.X
			chai.ScrollTo(midPoint)

			UpdateTransform(dt, &ecs_engine)

		},
		OnDraw: func() {
			chai.Shapes.DrawFillRectRotated(midScreen, chai.Vector2fOne.Scale(50.0), chai.NewRGBA8(255, 100, 230, 255), rotation)

			chai.Shapes.DrawTriangleRotated(midPoint, chai.NewVector2f(2.0, 4.0), chai.NewRGBA8(255, 0, 0, 255), rotation)
			//chai.Sprites.DrawSpriteOrigin(chai.NewVector2f(2, 0.0), chai.Vector2fZero, chai.Vector2fOne, &bgl_texture, chai.NewRGBA8(255, 255, 255, 255))
			for i := 0; i < 1; i++ {
				fontAtlas.DrawString("ابدأ اللعب ٢/٤", chai.Vector2fOne.AddXY(0.0, 0.0), 0.5, chai.WHITE)
				englishFontAtlas.DrawString("Baghdad Game Lab\nBaghdad Game Lab", chai.Vector2fOne.AddXY(0.0, 35.0), 0.5, chai.WHITE)
			}

			UpdateSprite(&chai.Sprites, &ecs_engine)

			fontAtlas.Render()
			englishFontAtlas.Render()
		},
		OnEvent: func(ae *chai.AppEvent) {
			if ae.Type == chai.JS_KEYDOWN {
				// chai.LogF("%v", ae.Key)
			}
		},
	}

	chai.Run(&game)

}

type Sprite struct {
	texture chai.Texture2D
}

func (t *Sprite) ComponentSet(val interface{}) { *t = val.(Sprite) }

var velo chai.Vector2f

func UpdateTransform(dt float32, _ecs_engine *chai.EcsEngine) {
	chai.EachAll(_ecs_engine, func(entity *chai.EcsEntity) {
		x_axis := chai.GetActionStrength("Up") - (chai.GetActionStrength("Down"))
		y_axis := chai.GetActionStrength("Right") - (chai.GetActionStrength("Left"))

		velo.X = chai.LerpFloat32(velo.X, (x_axis), dt*2.5)
		velo.Y = chai.LerpFloat32(velo.Y, (y_axis), dt*2.5)

		entity.Rot -= dt * 600.0 * velo.Y
		entity.Rot = float32(math.Mod(float64(entity.Rot), 360))

		direction := chai.Vector2fRight.Rotate(entity.Rot, chai.Vector2fZero)

		entity.Pos.Y += velo.X * direction.Y
		entity.Pos.X += velo.X * direction.X

	})
}

func UpdateSprite(_sp *chai.SpriteBatch, _ecs_engine *chai.EcsEngine) {
	chai.Each(_ecs_engine, Sprite{}, func(entity *chai.EcsEntity, a interface{}) {
		sprite := a.(Sprite)

		_sp.DrawSpriteOrigin(entity.Pos, chai.Vector2fZero, chai.Vector2fOne, &sprite.texture, chai.WHITE)
	})
}
