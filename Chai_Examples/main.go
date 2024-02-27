package main

import (
	chai "github.com/mhamedGd/chai"
)

var game chai.App
var scene_one chai.Scene

var scene_two chai.Scene
var sprite_render_system chai.SpriteRenderOriginSystem = chai.SpriteRenderOriginSystem{Offset: chai.Vector2fZero}
var triangle_render_system chai.TriangleRenderSystem = chai.TriangleRenderSystem{}
var rect_render_system chai.RectRenderSystem = chai.RectRenderSystem{}
var circle_render_system chai.CircleRenderSystem = chai.CircleRenderSystem{}

var keybinds_system KeyBindsSystem = KeyBindsSystem{}
var move_system MoveSystem = MoveSystem{}

var bgl_texture chai.Texture2D

var textureSet chai.TextureSettings

func main() {

	var dpad_modifier chai.Vector2f = chai.Vector2fZero

	var fontAtlas chai.FontBatchAtlas
	var englishFontAtlas chai.FontBatchAtlas

	game = chai.App{
		Width:  800,
		Height: 600,
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
			chai.BindInput("Hit", chai.KEY_SPCAE)

			chai.Shapes.LineWidth = .5

			textureSet = chai.TextureSettings{
				Filter: chai.TEXTURE_FILTER_NEAREST,
			}

			bgl_texture = chai.LoadPng("Assets/tile_0004.png", &textureSet)

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

			sprite_render_system.Sprites = &chai.Sprites
			triangle_render_system.Shapes = &chai.Shapes
			rect_render_system.Shapes = &chai.Shapes
			circle_render_system.Shapes = &chai.Shapes

			scene_one = chai.NewScene()
			scene_two = chai.NewScene()

			scene_one.OnSceneStart = StartSceneOne
			scene_two.OnSceneStart = StartSceneTwo

			chai.ChangeScene(&scene_one)

			chai.ScaleView(2)
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
			fontAtlas.DrawString("ابدأ اللعب ٢/٤", chai.Vector2fOne.AddXY(-150.0, 0.0), 0.5, chai.WHITE)
			englishFontAtlas.DrawString("Baghdad Game Lab\nBaghdad Game Lab", chai.Vector2fOne.AddXY(0.0, 35.0), 0.5, chai.WHITE)

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
	chai.EachEntity(KeyBinds{}, func(entity *chai.EcsEntity, a interface{}) {
		keybinds := a.(KeyBinds)
		keybinds.axes.X = chai.GetActionStrength("Right") - (chai.GetActionStrength("Left"))
		keybinds.axes.Y = chai.GetActionStrength("Up") - (chai.GetActionStrength("Down"))
		chai.WriteComponent(kb.GetEcsEngine(), entity, keybinds)

		hit := chai.Raycast(entity.Pos, chai.Vector2fUp.Rotate(entity.Rot, chai.Vector2fZero), 600.0)
		if hit.HasHit {
			chai.Shapes.DrawLine(entity.Pos, hit.HitPosition, chai.WHITE)
			chai.Shapes.DrawLine(hit.HitPosition, hit.HitPosition.Add(hit.Normal.Scale(8.0)), chai.WHITE)
			if chai.IsJustPressed("Hit") {
				audio := chai.AudioSourceComponent{}
				chai.ReadComponent(kb.GetEcsEngine(), entity, &audio)
				audio.Play("Slash")
				audio.SetVolume("Slash", 0.0)
				chai.WriteComponent(kb.GetEcsEngine(), entity, audio)
			}
		}

	})
}

type MoveComponent struct {
	velocity chai.Vector2f
}

func (t *MoveComponent) ComponentSet(val interface{}) { *t = val.(MoveComponent) }

type MoveSystem struct {
	chai.EcsSystemImpl
}

const player_speed float32 = 80000.0

func (ms *MoveSystem) Update(dt float32) {
	chai.EachEntity(MoveComponent{}, func(entity *chai.EcsEntity, a interface{}) {

		movecomp := a.(MoveComponent)
		dynamic := chai.DynamicBodyComponent{}
		chai.ReadComponent(ms.GetEcsEngine(), entity, &dynamic)

		bindings := KeyBinds{}
		chai.ReadComponent(ms.GetEcsEngine(), entity, &bindings)

		dynamic.SetAngularVelocity(dt * 250.0 * -bindings.axes.X)

		direction := chai.Vector2fUp.Rotate(entity.Rot, chai.Vector2fZero)
		force := chai.NewVector2f(bindings.axes.Y*direction.X*player_speed, bindings.axes.Y*direction.Y*player_speed)
		dynamic.ApplyForce(force)
		//dynamic.SetLinearVelocity(force.X, force.Y)

		dynamic.ApplyForce(dynamic.GetLinearVelocity().Scale(-4000).Scale(1.0 - bindings.axes.Y))

		chai.ScrollTo(entity.Pos)

		chai.WriteComponent(ms.GetEcsEngine(), entity, movecomp)
	})
}

func UpdateAllEntities(dt float32, _ecs_engine *chai.EcsEngine) {
	chai.EachEntityAll(&chai.GetCurrentScene().Ecs_engine, func(entity *chai.EcsEntity, index int) {
		//y_axis := chai.GetActionStrength("Up") - (chai.GetActionStrength("Down"))

		// velo.X = chai.LerpFloat32(velo.X, (x_axis), dt*2.5)
		// velo.Y = chai.LerpFloat32(velo.Y, (y_axis), dt*2.5)
		//chai.LogF("%v", index)

	})
}

type UpdateColorWithAnimSystem struct {
	chai.EcsSystemImpl
}

func (ms *UpdateColorWithAnimSystem) Update(dt float32) {
	chai.EachEntity(chai.AnimationComponent[float32]{}, func(entity *chai.EcsEntity, a interface{}) {
		anims := a.(chai.AnimationComponent[float32])
		// entity.Pos.Y = float32(tween.GetCurrentValue())
		sprite := chai.SpriteComponent{}
		chai.ReadComponent(ms.GetEcsEngine(), entity, &sprite)
		sprite.Tint.SetColorAFloat32(anims.GetCurrentValue("Color Change"))
		chai.WriteComponent(ms.GetEcsEngine(), entity, sprite)
	})
}

type UpdatePosWithAnimSystem struct {
	chai.EcsSystemImpl
}

func (ms *UpdatePosWithAnimSystem) Update(dt float32) {
	chai.EachEntity(chai.AnimationComponent[float32]{}, func(entity *chai.EcsEntity, a interface{}) {
		anims := a.(chai.AnimationComponent[float32])

		entity.Pos.Y = anims.GetCurrentValue("Move")
	})
}

var tween_animator_float32_system chai.TweenAnimatorSystemFloat32 = chai.TweenAnimatorSystemFloat32{}
var tween_animator_vector2i_system chai.TweenAnimatorSystemVector2i = chai.TweenAnimatorSystemVector2i{}
var update_color_with_anim_system UpdateColorWithAnimSystem = UpdateColorWithAnimSystem{}
var update_pos_with_anim_system UpdatePosWithAnimSystem = UpdatePosWithAnimSystem{}

var cat_sprites_anim_system chai.SpriteAnimationSystem = chai.SpriteAnimationSystem{Offset: chai.NewVector2f(0.0, 12.0)}

var dynamic_body_update_system chai.DynamicBodyUpdateSystem = chai.DynamicBodyUpdateSystem{}

func StartSceneOne() {
	scene_one.Background = chai.NewRGBA8Float(0.05, 0.1, 0.1, 1.0)
	scene_one.NewRenderSystem(&triangle_render_system)
	scene_one.NewRenderSystem(&rect_render_system)
	scene_one.NewRenderSystem(&circle_render_system)
	scene_one.NewRenderSystem(&cat_sprites_anim_system)
	scene_one.NewRenderSystem(&sprite_render_system)
	sprite_render_system.Scale = 1.0

	scene_one.NewUpdateSystem(&dynamic_body_update_system)
	scene_one.NewUpdateSystem(&keybinds_system)
	scene_one.NewUpdateSystem(&move_system)
	scene_one.NewUpdateSystem(&tween_animator_float32_system)
	scene_one.NewUpdateSystem(&tween_animator_vector2i_system)
	scene_one.NewUpdateSystem(&update_pos_with_anim_system)
	scene_one.NewUpdateSystem(&update_color_with_anim_system)
	// scene_one.NewUpdateSystem(&update_color_with_anim_system)

	//scene_one.NewEntity(chai.NewVector2f(0.0, 0.0), 0.0)
	//scene_one.WriteComponentToLastEntity(catAnim)

	//scene_one.WriteComponentToLastEntity(chai.SpriteAnimation{CurrentAnimation: "Cat"})
	dynamicBodySettings := chai.DynamicBodySettings{
		BodySize:     chai.Vector2fOne,
		BodyShape:    chai.Shape_RectBody,
		Mass:         10.0,
		Friction:     0.4,
		Restitution:  0.25,
		GravityScale: 1.0,
	}
	for i := 0; i < 10; i++ {
		catAnim := chai.NewAnimationComponentVector2i()
		catAnim.NewTweenAnimationVector2i("Cat")
		catAnim.RegisterKeyframe("Cat", 0.0, chai.Vector2i{X: 0, Y: 0})
		catAnim.RegisterKeyframe("Cat", 0.5, chai.Vector2i{X: 1, Y: 0})
		catAnim.RegisterKeyframe("Cat", 1.0, chai.Vector2i{X: 2, Y: 0})
		catAnim.RegisterKeyframe("Cat", 1.5, chai.Vector2i{X: 3, Y: 0})
		ent := scene_one.NewEntity(chai.NewVector2f(40.0+(float32(i))*20.0, 0.0), chai.Vector2fOne.Scale(12.0).AddXY(0.0, 4.0), 0.0)
		dynamicBodySettings.BodySize = ent.Dimensions
		scene_one.WriteComponentToLastEntity(chai.NewDynamicBody(ent, dynamicBodySettings))
		if i%2 == 0 {
			catAnim.Play("Cat")
		}
		scene_one.WriteComponentToLastEntity(catAnim)
		scene_one.WriteComponentToLastEntity(chai.SpriteAnimation{CurrentAnimation: "Cat"})
	}

	audioStream := chai.LoadAudioFile("Assets/web_sfx.ogg")

	//Player
	ent := scene_one.NewEntity(chai.Vector2fZero, chai.Vector2fOne.Scale(6.0).AddXY(0.0, 12.0), 0.0)
	scene_one.WriteComponentToLastEntity(chai.SpriteComponent{Texture: bgl_texture, Tint: chai.WHITE})

	audioSComp := chai.NewAudioSourceComponent()
	audioSComp.AddAudioSource("Slash", audioStream)

	scene_one.WriteComponentToLastEntity(audioSComp)

	dynamicBodySettings = chai.DynamicBodySettings{
		BodySize:     ent.Dimensions,
		BodyShape:    chai.Shape_CircleBody,
		Mass:         15.0,
		Friction:     0.4,
		Restitution:  0.0,
		GravityScale: 0.0,
	}

	player_dynmaic := chai.NewDynamicBody(ent, dynamicBodySettings)
	player_dynmaic.GetPhyiscsBody().OnCollisionStart.AddListener(func(c ...*chai.Collision) {
		col := c[0]
		if col.SecondBody.IsTrigger {
			return
		}
		scene_one.NewEntity(col.CollisionPoint, chai.Vector2fOne.Scale(2.0), 0.0)
		scene_one.WriteComponentToLastEntity(chai.RectRenderComponent{})
	})
	scene_one.WriteComponentToLastEntity(player_dynmaic)

	scene_one.WriteComponentToLastEntity(KeyBinds{})
	scene_one.WriteComponentToLastEntity(MoveComponent{})

	scene_one.NewEntity(chai.NewVector2f(0.0, 25.0), chai.Vector2fOne, 0.0)
	scene_one.WriteComponentToLastEntity(chai.SpriteComponent{Texture: bgl_texture, Tint: chai.WHITE})
	animation := chai.NewAnimationComponentFloat32()
	animation.NewTweenAnimationFloat32("Color Change", true)
	animation.RegisterKeyframe("Color Change", 0.0, 0.0)
	animation.RegisterKeyframe("Color Change", 1.5, 1.0)
	animation.RegisterKeyframe("Color Change", 3.0, 0.0)

	animation.NewTweenAnimationFloat32("Move", true)
	animation.RegisterKeyframe("Move", 0.0, 0.0)
	animation.RegisterKeyframe("Move", 1.5, 25.0)
	animation.RegisterKeyframe("Move", 3.0, 0.0)

	animation.PlaySimultaneous("Color Change", "Move")

	scene_one.WriteComponentToLastEntity(animation)

	catSpriteSheet := chai.LoadPng("Assets/Cat Sprite Sheet.png", &textureSet)
	tileset := chai.NewTileSet(chai.Vector2fZero, catSpriteSheet, 8, 10)
	cat_sprites_anim_system.TileSet = tileset
	cat_sprites_anim_system.Sprites = &chai.Sprites
	cat_sprites_anim_system.SpriteScale = 0.125

	ent = scene_one.NewEntity(chai.NewVector2f(0.0, -70.0), chai.NewVector2f(500.0, 25.0), 0.0)
	staticBodySets := chai.StaticBodySettings{
		BodySize:  ent.Dimensions,
		BodyShape: chai.Shape_RectBody,
		Friction:  0.3,
	}
	scene_one.WriteComponentToLastEntity(chai.NewStaticBody(ent, staticBodySets))
	scene_one.WriteComponentToLastEntity(chai.RectRenderComponent{Tint: chai.WHITE})

	ent = scene_one.NewEntity(chai.NewVector2f(40.0, 150.0), chai.NewVector2f(50.0, 50.0), 45.0)

	staticBodySets.BodySize = ent.Dimensions
	staticBodySets.BodyShape = chai.Shape_CircleBody

	trigger_area := chai.NewTriggerArea(ent, staticBodySets)
	trigger_area.GetPhyiscsBody().OnCollisionStart.AddListener(func(c ...*chai.Collision) {
		c[0].SecondBody.Debug_Tint = chai.NewRGBA8(255, 0, 0, 255)
	})
	trigger_area.GetPhyiscsBody().OnCollisionEnd.AddListener(func(c ...*chai.Collision) {
		c[0].SecondBody.Debug_Tint = chai.WHITE
	})
	scene_one.WriteComponentToLastEntity(trigger_area)
	scene_one.WriteComponentToLastEntity(chai.CircleRenderComponent{})

	scene_one.NewEntity(chai.Vector2fZero, chai.Vector2fOne, 0.0)
	hmoodSoundComp := chai.NewAudioSourceComponent()
	hmoodSoundComp.AddAudioSource("Hmood", chai.LoadAudioFile("Assets/Hmood.mp3"))
	hmoodSoundComp.Play("Hmood")
	hmoodSoundComp.SetVolume("Hmood", 0.1)
	scene_one.WriteComponentToLastEntity(hmoodSoundComp)
}

func StartSceneTwo() {
	scene_two.NewRenderSystem(&sprite_render_system)
	scene_two.NewRenderSystem(&triangle_render_system)

	scene_two.NewUpdateSystem(&dynamic_body_update_system)
	scene_two.NewUpdateSystem(&keybinds_system)
	scene_two.NewUpdateSystem(&move_system)

	//Player
	ent := scene_two.NewEntity(chai.Vector2fZero, chai.Vector2fOne.Scale(8.0), 0.0)
	dynamicBodySettings := chai.DynamicBodySettings{
		BodySize:     ent.Dimensions,
		BodyShape:    chai.Shape_CircleBody,
		Mass:         3.0,
		Friction:     0.4,
		Restitution:  0.0,
		GravityScale: 0.0,
	}
	scene_two.WriteComponentToLastEntity(chai.NewDynamicBody(ent, dynamicBodySettings))
	scene_two.WriteComponentToLastEntity(chai.SpriteComponent{Texture: bgl_texture, Tint: chai.WHITE})

	scene_two.WriteComponentToLastEntity(KeyBinds{})
	scene_two.WriteComponentToLastEntity(MoveComponent{})

}
