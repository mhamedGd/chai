package chai

import (
	"strings"

	box2d "github.com/mhamedGd/chai-box2d"
	ldtkgo "github.com/mhamedgd/ldtkgo-chai"
)

func ParseLdtk(_filePath string) Map[string, Tilemap] {

	temp_levels_map := NewMap[string, Tilemap]()

	if _filePath == "" {
		ErrorF("Ldtk File Path is Empty")
	}

	_filePathArr := strings.Split(_filePath, "/")
	_filePathArr = _filePathArr[:len(_filePathArr)-1]
	_folderPath := strings.Join(_filePathArr, "/")
	_folderPath += "/"

	ldtkChannel := make(chan []byte)
	go LoadResponse(_filePath, ldtkChannel)
	ldtkBytes := <-ldtkChannel

	ldtk_reader, err := ldtkgo.Read(ldtkBytes)
	if err != nil {
		ErrorF("%v", err.Error())
	}

	last_layer_tile_size := 0

	total_entites := NewMap[string, []ldtkEntity]()
	for _, level := range ldtk_reader.Levels {

		tile_size := ldtk_reader.Tilesets[0].GridSize
		total_layers := NewList[levelLayer]()

		for li := len(level.Layers) - 1; li >= 0; li-- {

			l := level.Layers[li]

			last_layer_tile_size = l.GridSize

			total_tiles := NewList[ldtkgo.Tile]()
			total_autotiles := NewList[ldtkgo.Tile]()

			var texture Texture2D
			if l.Tileset != nil {
				texture = LoadPngByTileset(_folderPath+l.Tileset.Path, &TextureSettings{Filter: TEXTURE_FILTER_NEAREST}, tile_size, tile_size)
			}

			for _, v := range l.Entities {
				v.Position[1] *= -1
				_, ok := total_entites.data[v.Identifier]
				if !ok {
					total_entites.Set(v.Identifier, make([]ldtkEntity, 0))
				}
				// if ok {
				total_entites.Insert(v.Identifier, append(total_entites.Get(v.Identifier), ldtkEntity{
					Identifier:   v.Identifier,
					Position:     IntArrToVec2f(v.Position),
					GridPosition: IntArrToVec2i(v.Position),
				}))
				// } else {
				// 	total_entites.Insert(v.Identifier, []ldtkEntity{
				// 		ldtkEntity{
				// 			Identifier:   v.Identifier,
				// 			Position:     IntArrToVec2f(v.Position).Scale(_scale),
				// 			GridPosition: IntArrToVec2i(v.Position),
				// 		},
				// 	})
				// }
			}

			tiles := l.Tiles
			for _, v := range tiles {
				total_tiles.PushBack(*v)
			}

			auto_tiles := l.AutoTiles
			// total_autotiles.PushBackArray(auto_tiles)
			if auto_tiles == nil {
				ErrorF("Layer (%v): Tile Doesn't Exist", l.Identifier)
			}
			for _, v := range auto_tiles {
				total_autotiles.PushBack(*v)
			}
			tileset_original_size := NewVector2i(0, 0)
			if l.Tileset != nil {
				tileset_original_size = NewVector2i(l.Tileset.Width, l.Tileset.Height)
			}
			total_layers.PushBack(levelLayer{
				identifier:            l.Identifier,
				tiles:                 total_tiles,
				auto_tiles:            total_autotiles,
				tileset_texture:       texture,
				original_texture_size: tileset_original_size,
				opacity:               l.Opacity,
				tile_size:             l.GridSize,
				tileset:               l.Tileset,
				layertype:             l.Type,
				physicsLayer:          PHYSICS_LAYER_1,
				z_offset:              float32(li),
			})
		}

		temp_levels_map.Insert(level.Identifier, Tilemap{
			layers: total_layers,

			grid_width:  level.Width / tile_size,
			grid_height: level.Height / tile_size,
			tile_size:   last_layer_tile_size,
			Scale:       1.0,
			SolidTiles:  NewMap[Vector2i, Tile](),
			Entities:    total_entites,
		})

	}

	return temp_levels_map
}

func LoadTilemapLevel(scene *Scene, _level_name string, _all_levels Map[string, Tilemap], _z float32, _offset Vector2f) *Tilemap {
	// tilemap_level := ldtkLevels.Get(_level_name)
	level := _all_levels.Get(_level_name)

	for li := 0; li < level.layers.Count(); li++ {
		l := level.layers.Data[li]
		if l.layertype == "Entities" {
			for k, v := range level.Entities.data {
				ents := v
				for i, _ := range ents {
					// level.Entities.Set(k, v.Add(_offset).Add(Vector2fOne.Scale(float32(l.tile_size)/2.0)))
					ents[i].Position = ents[i].Position.Add(_offset).Add(Vector2fOne.Scale(float32(l.tile_size) / 2.0))
				}
				level.Entities.Set(k, ents)
			}
			continue
		}

		texture := l.tileset_texture
		texture_width := texture.Width
		texture_height := texture.Height

		tiles := l.tiles

		if tiles.Data == nil {
			ErrorF("Layer (%v): Tile Doesn't Exist", l.identifier)
		}

		for i := 0; i < tiles.Count(); i++ {
			pixel_size_x := 1.0 / float32(texture_width)
			pixel_size_y := 1.0 / float32(texture_height)

			// origin_uv := NewVector2f(float32(tiles[i].Src[0])/float32(texture_width), float32(tiles[i].Src[1])/float32(texture_height))

			origin_uv := NewVector2f(float32(tiles.Data[i].Src[0]), float32(tiles.Data[i].Src[1]))
			origin_uv.X /= float32(l.original_texture_size.X)
			origin_uv.Y /= float32(l.original_texture_size.Y)

			// uv_tile_scalar_x := float32(tile_size) / float32(texture_width)
			// uv_tile_scalar_y := float32(tile_size) / float32(texture_height)

			flip_factor_x := float32(0.0)
			if tiles.Data[i].FlipX() {
				flip_factor_x = -1.0
			} else {
				flip_factor_x = 1.0
			}

			flip_factor_y := float32(0.0)
			if tiles.Data[i].FlipY() {
				flip_factor_y = -1.0
			} else {
				flip_factor_y = 1.0
			}
			t := VisualTransform{
				Position:   NewVector2f(float32(tiles.Data[i].Position[0])+BoolToFloat32(tiles.Data[i].FlipX())*float32(l.tile_size), float32(-tiles.Data[i].Position[1])+BoolToFloat32(tiles.Data[i].FlipY())*float32(l.tile_size)).Add(_offset),
				Dimensions: NewVector2f(float32(l.tile_size)*flip_factor_x, float32(l.tile_size)*flip_factor_y),
				Z:          _z + l.z_offset,
				Scale:      1,
				Tint:       NewRGBA8Float(1.0, 1.0, 1.0, l.opacity),
				UV1:        origin_uv,
				UV2:        origin_uv.AddXY(float32(l.tile_size)*pixel_size_x, float32(l.tile_size)*pixel_size_y),
			}
			world_actual_postion := NewVector2f(float32(tiles.Data[i].Position[0]), float32(-tiles.Data[i].Position[1])).Add(_offset)
			collider_pos := world_actual_postion.Add(t.Dimensions.Scale(0.5)).AddXY(BoolToFloat32(tiles.Data[i].FlipX())*float32(l.tile_size), 0.0)

			tile_enumset := []string(l.tileset.Enums[tiles.Data[i].ID])
			if len(tile_enumset) > 0 {
				if tile_enumset[0] == "Solid" {
					// newStaticCollisionTile(scene, collider_pos, l.tile_size, level.Scale, PhysicsLayer_All)
					newStaticCollisionTileBox2d(collider_pos, l.tile_size, 1.0, l.physicsLayer)
					level.SolidTiles.Insert(NewVector2i(int(collider_pos.X), int(collider_pos.Y)), Tile{Enumset: ListFromSlice(tile_enumset), Solid: true})
				}

			}

			renderObj := newRenderObject(0, SPRITE_RENDEROBJECTTYPEFUNC)
			renderObj.texture = &texture
			RenderQuadTreeContainer.Insert(Pair[VisualTransform, RenderObject]{t, renderObj}, Rect{Position: NewVector2f(float32(tiles.Data[i].Position[0]), float32(-tiles.Data[i].Position[1])).Add(_offset), Size: NewVector2f(float32(l.tile_size), float32(l.tile_size))})
		}

		auto_tiles := l.auto_tiles

		for i := 0; i < auto_tiles.Count(); i++ {
			pixel_size_x := 1.0 / float32(texture_width)
			pixel_size_y := 1.0 / float32(texture_height)

			origin_uv := NewVector2f(float32(auto_tiles.Data[i].Src[0]), float32(auto_tiles.Data[i].Src[1]))
			origin_uv.X /= float32(l.original_texture_size.X)
			origin_uv.Y /= float32(l.original_texture_size.Y)

			// uv_tile_scalar_x := float32(tile_size) / float32(texture_width)
			// uv_tile_scalar_y := float32(tile_size) / float32(texture_height)

			flip_factor_x := float32(0.0)
			if auto_tiles.Data[i].FlipX() {
				flip_factor_x = -1.0
			} else {
				flip_factor_x = 1.0
			}

			flip_factor_y := float32(0.0)
			if auto_tiles.Data[i].FlipY() {
				flip_factor_y = -1.0
			} else {
				flip_factor_y = 1.0
			}

			t := VisualTransform{
				Position:   NewVector2f(float32(auto_tiles.Data[i].Position[0])+BoolToFloat32(auto_tiles.Data[i].FlipX())*float32(l.tile_size), float32(-auto_tiles.Data[i].Position[1])+BoolToFloat32(auto_tiles.Data[i].FlipY())*float32(l.tile_size)).Add(_offset),
				Dimensions: NewVector2f(float32(l.tile_size)*flip_factor_x, float32(l.tile_size)*flip_factor_y),
				Z:          _z + l.z_offset,
				Scale:      1,
				Tint:       NewRGBA8Float(1.0, 1.0, 1.0, l.opacity),
				UV1:        origin_uv,
				UV2:        origin_uv.AddXY(float32(l.tile_size)*pixel_size_x, float32(l.tile_size)*pixel_size_y),
			}

			world_actual_postion := NewVector2f(float32(auto_tiles.Data[i].Position[0]), float32(-auto_tiles.Data[i].Position[1])).Add(_offset)
			collider_pos := world_actual_postion.Add(t.Dimensions.Scale(0.5)).AddXY(BoolToFloat32(auto_tiles.Data[i].FlipX())*float32(l.tile_size), 0.0)
			// id := scene.NewEntityId()
			// rb := NewRigidBody(id, &RigidBodySettings{
			// 	BodyType:        Type_BodyStatic,
			// 	ColliderShape:   Shape_RectBody,
			// 	StartPosition:   collider_pos,
			// 	StartDimensions: t.Dimensions,
			// 	Mass:            1000, Friction: 0.4, Elasticity: 0.4,
			// 	PhysicsLayer: PhysicsLayer_All,
			// })

			// scene.AddComponents(id, ToComponent(t), ToComponent(rb))
			tile_enumset := []string(l.tileset.Enums[auto_tiles.Data[i].ID])
			if len(tile_enumset) > 0 {
				if tile_enumset[0] == "Solid" {
					// newStaticCollisionTile(scene, collider_pos, l.tile_size, level.Scale, PhysicsLayer_All)
					newStaticCollisionTileBox2d(collider_pos, l.tile_size, 1.0, l.physicsLayer)
					level.SolidTiles.Insert(NewVector2i(auto_tiles.Data[i].Position[0]/l.tile_size, auto_tiles.Data[i].Position[1]/l.tile_size), Tile{Enumset: ListFromSlice(tile_enumset), Solid: true})
				}

			}

			renderObj := newRenderObject(0, SPRITE_RENDEROBJECTTYPEFUNC)
			renderObj.texture = &texture
			RenderQuadTreeContainer.Insert(Pair[VisualTransform, RenderObject]{t, renderObj}, Rect{Position: world_actual_postion, Size: NewVector2f(float32(l.tile_size), float32(l.tile_size))})
		}
	}

	return &level
}

type ldtkEntity struct {
	Identifier   string
	Position     Vector2f
	GridPosition Vector2i
}

type Tile struct {
	Enumset List[string]
	Solid   bool
}

type Tilemap struct {
	Offset Vector2f
	// tileset                 TileSet
	layers                  List[levelLayer]
	grid_width, grid_height int
	tile_size               int
	Scale                   float32
	SolidTiles              Map[Vector2i, Tile]
	Entities                Map[string, []ldtkEntity]
}

func (level *Tilemap) GridSize() Vector2i {
	return NewVector2i(level.grid_width, level.grid_height)
}

func (level *Tilemap) Tilesize() int {
	return level.tile_size
}

type levelLayer struct {
	identifier            string
	tiles                 List[ldtkgo.Tile]
	auto_tiles            List[ldtkgo.Tile]
	tileset_texture       Texture2D
	original_texture_size Vector2i
	tile_size             int
	opacity               float32
	physicsLayer          uint16
	tileset               *ldtkgo.Tileset
	layertype             string
	z_offset              float32
}

func IntArrToVec2f(original []int) Vector2f {
	return NewVector2f(float32(original[0]), float32(original[1]))
}

func IntArrToVec2i(origin []int) Vector2i {
	return NewVector2i(origin[0], origin[1])
}

// func newStaticCollisionTile(scene *Scene, _pos Vector2f, tile_size int, _scale float32, _physicsLayer uint) {
// 	size := cpVector2f(NewVector2f(1.0, 1.0).Scale(float32(tile_size) / 2.0 * _scale))
// 	body := cp.NewBody(1000000.0, 1000000.0)

// 	body.SetPosition(cpVector2f(_pos))
// 	body.SetAngle((0) * Deg2Rad)

// 	rbTypeToCpType(body, Type_BodyStatic)
// 	var shape *cp.Shape
// 	// shape = cp.NewBox(body, size.X, size.Y/2.0, float64(rbSettings.StartRotation)*PI/180.0)
// 	box := cp.NewBB(-size.X, -size.Y, size.X, size.Y)
// 	shape = cp.NewBox2(body, box, 0.0)
// 	// body.SetMoment(cp.MomentForBox2(1000, box))
// 	// shape.SetMass(float64(1000))

// 	shape.SetElasticity(0.3)
// 	shape.SetFriction(0.8)

// 	shape.Filter.Categories = PHYSICS_LAYER_ALL
// 	shape.Filter.Mask = _physicsLayer

// 	body.UserData = EntId(0)
// 	shape.UserData = body

// 	GetPhysicsWorld().cpSpace.AddBody(body)
// 	GetPhysicsWorld().cpSpace.AddShape(shape)
// }

func newStaticCollisionTileBox2d(_position Vector2f, _tile_size int, _scale float32, _physics_layer uint16) {
	size := vec2fToB2Vec(Vector2fOne.Scale((float32(_tile_size) / 2.0) * _scale))
	bodydef := box2d.MakeB2BodyDef()
	bodydef.Type = box2d.B2BodyType.B2_staticBody
	bodydef.Position.Set(vec2fToB2Vec(_position).X, vec2fToB2Vec(_position).Y)
	bodydef.AllowSleep = false
	bodydef.FixedRotation = true

	body := physics_world.box2dWorld.CreateBody(&bodydef)
	fd := box2d.MakeB2FixtureDef()

	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(size.X, size.Y)
	fd.Shape = &shape

	fd.Density = 100000
	fd.Friction = 1.0
	fd.Restitution = 0.0
	fd.Filter.CategoryBits = _physics_layer
	fixture := body.CreateFixtureFromDef(&fd)
	fixture.SetSensor(false)

	body.SetUserData(0)
}

// func newRectStaticCollisionTile(scene *Scene, _rect Rect, tile_size int, _scale float32, _physicsLayer uint) {
// 	size := cpVector2f(NewVector2f(_rect.Position.X+_rect.Size.X/2.0, _rect.Position.Y+_rect.Size.Y/2.0))
// 	body := cp.NewBody(1000.0, 1000000)

// 	body.SetPosition(cpVector2f(_rect.Position))
// 	body.SetAngle((0) * Deg2Rad)

// 	rbTypeToCpType(body, Type_BodyStatic)
// 	var shape *cp.Shape
// 	// shape = cp.NewBox(body, size.X, size.Y/2.0, float64(rbSettings.StartRotation)*PI/180.0)
// 	box := cp.NewBB(-size.X, -size.Y, size.X, size.Y)
// 	shape = cp.NewBox2(body, box, 0.0)
// 	// body.SetMoment(cp.MomentForBox2(1000, box))
// 	// shape.SetMass(float64(1000))

// 	shape.SetElasticity(0.3)
// 	shape.SetFriction(0.8)

// 	shape.Filter.Categories = PHYSICS_LAYER_ALL
// 	shape.Filter.Mask = _physicsLayer

// 	body.UserData = EntId(5)

// 	GetPhysicsWorld().cpSpace.AddBody(body)
// 	GetPhysicsWorld().cpSpace.AddShape(shape)
// }
