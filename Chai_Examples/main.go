package main

import (
	"math"

	chai "github.com/mhamedGd/Chai"
)

var game chai.App
var scene_one chai.Scene

var scene_two chai.Scene
var sprite_render_system SpriteRenderSystem = SpriteRenderSystem{}
var shape_render_system ShapesRenderingSystem = ShapesRenderingSystem{}
var keybinds_system KeyBindsSystem = KeyBindsSystem{}
var move_system MoveSystem = MoveSystem{}

var bgl_texture chai.Texture2D

func main() {

	var dpad_modifier chai.Vector2f = chai.Vector2fZero

	var fontAtlas chai.FontBatchAtlas
	var englishFontAtlas chai.FontBatchAtlas

	scene_one.OnSceneStart = func() {

	}
	game = chai.App{
		Width:  1920,
		Height: 1080,
		Title:  "Test",
		OnStart: func() {

			chai.LogF("STARTED\n")

			chai.BindInput("Up", chai.KEY_W)
			chai.BindInput("Down", chai.KEY_S)
			chai.BindInput("Right", chai.KEY_D)
			chai.BindInput("Left", chai.KEY_A)
			chai.BindInput("Zoom In", chai.KEY_E)
			chai.BindInput("Zoom Out", chai.KEY_Q)
			chai.BindInput("Change Scene", chai.KEY_L)

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

			sprite_render_system._sp = &chai.Sprites
			shape_render_system._sh = &chai.Shapes

			scene_one = chai.NewScene()
			scene_two = chai.NewScene()

			scene_one.OnSceneStart = StartSceneOne
			scene_two.OnSceneStart = StartSceneTwo

			chai.ChangeScene(&scene_one)

			chai.ScaleView(4)
			chai.WarningF("This is a warning!!")
		},
		OnUpdate: func(dt float32) {
			zoomAxis := 500.0 * float32(dt) * (chai.GetActionStrength("Zoom In") - chai.GetActionStrength("Zoom Out"))
			chai.IncreaseScaleU(zoomAxis)

			if chai.IsJustPressed("Change Scene") {
				chai.ChangeScene(&scene_two)
			}
			//UpdateAllEntities(dt, &ecs_engine)

		},
		OnDraw: func() {

			//chai.Shapes.DrawTriangleRotated(midPoint, chai.NewVector2f(2.0, 4.0), chai.NewRGBA8(255, 0, 0, 255), rotation)
			//chai.Sprites.DrawSpriteOrigin(chai.NewVector2f(2, 0.0), chai.Vector2fZero, chai.Vector2fOne, &bgl_texture, chai.NewRGBA8(255, 255, 255, 255))
			// for i := 0; i < 1; i++ {
			// 	fontAtlas.DrawString("ابدأ اللعب ٢/٤", chai.Vector2fOne.AddXY(0.0, 0.0), 0.5, chai.WHITE)
			// 	englishFontAtlas.DrawString("Baghdad Game Lab\nBaghdad Game Lab", chai.Vector2fOne.AddXY(0.0, 35.0), 0.5, chai.WHITE)
			// }

			// sprite_render_system.Update(0.0)

			fontAtlas.Render()
			englishFontAtlas.Render()
		},
		OnEvent: func(ae *chai.AppEvent) {

		},
	}

	chai.Run(&game)

}

type KeyBinds struct {
	axes chai.Vector2f
}

func (t *KeyBinds) ComponentSet(val interface{}) { *t = val.(KeyBinds) }

type KeyBindsSystem struct {
	chai.EcsSystemImpl
}

func (kb *KeyBindsSystem) Update(dt float32) {
	chai.Each(&chai.GetCurrentScene().Ecs_engine, KeyBinds{}, func(entity *chai.EcsEntity, a interface{}) {
		keybinds := a.(KeyBinds)
		keybinds.axes.X = chai.GetActionStrength("Right") - (chai.GetActionStrength("Left"))
		keybinds.axes.Y = chai.GetActionStrength("Up") - (chai.GetActionStrength("Down"))
		chai.WriteComponent(kb.GetEcsEngine(), entity, keybinds)
	})
}

type MoveComponent struct {
	velocity chai.Vector2f
}

func (t *MoveComponent) ComponentSet(val interface{}) { *t = val.(MoveComponent) }

type MoveSystem struct {
	chai.EcsSystemImpl
}

func (ms *MoveSystem) Update(dt float32) {
	chai.Each(&chai.GetCurrentScene().Ecs_engine, MoveComponent{}, func(entity *chai.EcsEntity, a interface{}) {

		movecomp := a.(MoveComponent)

		bindings := KeyBinds{}
		chai.ReadComponent(ms.GetEcsEngine(), entity, &bindings)

		movecomp.velocity.X = chai.LerpFloat32(movecomp.velocity.X, (bindings.axes.X), dt*2.5)
		movecomp.velocity.Y = chai.LerpFloat32(movecomp.velocity.Y, (bindings.axes.Y), dt*2.5)
		entity.Rot -= dt * 600.0 * movecomp.velocity.X
		entity.Rot = float32(math.Mod(float64(entity.Rot), 360))

		direction := chai.Vector2fRight.Rotate(entity.Rot, chai.Vector2fZero)

		entity.Pos.Y += movecomp.velocity.Y * direction.Y
		entity.Pos.X += movecomp.velocity.Y * direction.X

		chai.ScrollTo(entity.Pos)

		chai.WriteComponent(ms.GetEcsEngine(), entity, movecomp)
	})
}

type Sprite struct {
	chai.Component
	texture chai.Texture2D
}

func (t *Sprite) ComponentSet(val interface{}) { *t = val.(Sprite) }

type SpriteRenderSystem struct {
	chai.EcsSystemImpl
	_sp *chai.SpriteBatch
}

func (_render *SpriteRenderSystem) Update(dt float32) {
	chai.Each(&chai.GetCurrentScene().Ecs_engine, Sprite{}, func(entity *chai.EcsEntity, a interface{}) {
		sprite := a.(Sprite)

		_render._sp.DrawSpriteOrigin(entity.Pos, chai.Vector2fZero, chai.Vector2fOne, &sprite.texture, chai.WHITE)
	})
}

type Shape struct {
	chai.Component
}

func (t *Shape) ComponentSet(val interface{}) { *t = val.(Shape) }

type ShapesRenderingSystem struct {
	chai.EcsSystemImpl
	_sh *chai.ShapeBatch
}

func (_render *ShapesRenderingSystem) Update(dt float32) {
	chai.Each(&chai.GetCurrentScene().Ecs_engine, Shape{}, func(entity *chai.EcsEntity, a interface{}) {
		_render._sh.DrawTriangleRotated(entity.Pos, chai.NewVector2f(2.0, 4.0), chai.NewRGBA8(255, 0, 0, 255), entity.Rot)
	})
}

func UpdateAllEntities(dt float32, _ecs_engine *chai.EcsEngine) {
	chai.EachAll(&chai.GetCurrentScene().Ecs_engine, func(entity *chai.EcsEntity, index int) {
		//y_axis := chai.GetActionStrength("Up") - (chai.GetActionStrength("Down"))

		// velo.X = chai.LerpFloat32(velo.X, (x_axis), dt*2.5)
		// velo.Y = chai.LerpFloat32(velo.Y, (y_axis), dt*2.5)
		//chai.LogF("%v", index)

	})
}

func StartSceneOne() {
	scene_one.NewRenderSystem(&sprite_render_system)
	scene_one.NewRenderSystem(&shape_render_system)

	scene_one.NewUpdateSystem(&keybinds_system)
	scene_one.NewUpdateSystem(&move_system)

	for i := 0; i < 2000; i++ {
		scene_one.NewEntity(chai.NewVector2f(40.0+(float32(i))*20.0, 0.0), 90.0)
		scene_one.WriteComponentToLastEntity(Sprite{texture: bgl_texture})
	}

	//Player
	scene_one.NewEntity(chai.Vector2fZero, 90.0)
	scene_one.WriteComponentToLastEntity(Sprite{texture: bgl_texture})
	scene_one.WriteComponentToLastEntity(Shape{})
	scene_one.WriteComponentToLastEntity(KeyBinds{})
	scene_one.WriteComponentToLastEntity(MoveComponent{})

}

func StartSceneTwo() {
	scene_two.NewRenderSystem(&sprite_render_system)
	scene_two.NewRenderSystem(&shape_render_system)

	scene_two.NewUpdateSystem(&keybinds_system)
	scene_two.NewUpdateSystem(&move_system)

	for i := 0; i < 2000; i++ {
		scene_two.NewEntity(chai.NewVector2f(40.0+(float32(i))*20.0, 0.0), 90.0)
		scene_two.WriteComponentToLastEntity(Sprite{texture: bgl_texture})
	}

	//Player
	scene_two.NewEntity(chai.Vector2fZero, 90.0)
	scene_two.WriteComponentToLastEntity(Sprite{texture: bgl_texture})
	scene_two.WriteComponentToLastEntity(Shape{})
	scene_two.WriteComponentToLastEntity(KeyBinds{})
	scene_two.WriteComponentToLastEntity(MoveComponent{})
	chai.LogF("Number of ents: %v", scene_two.GetNumberOfEntities())

}
