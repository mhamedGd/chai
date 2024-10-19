package chai

import (
	"strings"

	"github.com/mhamedGd/chai/customtypes"
	. "github.com/mhamedGd/chai/math"

	box2d "github.com/mhamedGd/chai-box2d"
	ldtkgo "github.com/mhamedgd/ldtkgo-chai"
)

func ParseLdtk(_filePath string) customtypes.Map[string, Tilemap] {

	temp_levels_map := customtypes.NewMap[string, Tilemap]()

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

	total_entites := customtypes.NewMap[string, []ldtkEntity]()
	for _, level := range ldtk_reader.Levels {

		tile_size := ldtk_reader.Tilesets[0].GridSize
		total_layers := customtypes.NewList[levelLayer]()

		for li := len(level.Layers) - 1; li >= 0; li-- {

			l := level.Layers[li]

			last_layer_tile_size = l.GridSize

			total_tiles := customtypes.NewList[ldtkgo.Tile]()
			total_autotiles := customtypes.NewList[ldtkgo.Tile]()

			var texture Texture2D
			if l.Tileset != nil {
				texture = LoadPngByTileset(_folderPath+l.Tileset.Path, TextureSettings{Filter: TEXTURE_FILTER_NEAREST}, tile_size, tile_size)
			}

			for _, v := range l.Entities {
				v.Position[1] *= -1
				_, ok := total_entites.AllItems()[v.Identifier]
				if !ok {
					total_entites.Set(v.Identifier, make([]ldtkEntity, 0))
				}
				// if ok {
				total_entites.Insert(v.Identifier, append(total_entites.Get(v.Identifier), ldtkEntity{
					Identifier:   v.Identifier,
					Position:     IntArrToVec2f(v.Position),
					GridPosition: IntArrToVec2i(v.Position),
				}))
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
			SolidTiles:  customtypes.NewMap[Vector2i, Tile](),
			Entities:    total_entites,
		})

	}

	return temp_levels_map
}

func LoadTilemapLevel(scene *Scene, _level_name string, _all_levels customtypes.Map[string, Tilemap], _z float32, _offset Vector2f) *Tilemap {
	// tilemap_level := ldtkLevels.Get(_level_name)
	level := _all_levels.Get(_level_name)

	for li := 0; li < level.layers.Count(); li++ {
		l := level.layers.Data[li]
		if l.layertype == "Entities" {
			for k, v := range level.Entities.AllItems() {
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

		const __scale = float32(.25)

		tiles := l.tiles
		createTilesFromList(&level, tiles, l, texture, texture_width, texture_height, _offset, _z, __scale)

		auto_tiles := l.auto_tiles
		createTilesFromList(&level, auto_tiles, l, texture, texture_width, texture_height, _offset, _z, __scale)
	}

	return &level
}

func createTilesFromList(_level *Tilemap, _list customtypes.List[ldtkgo.Tile], _l levelLayer, _texture Texture2D, _t_w, _t_h int, _offset Vector2f, _z, _scale float32) {
	for i := 0; i < _list.Count(); i++ {
		_this_tile := _list.Data[i]
		pixel_size_x := 1.0 / float32(_t_w)
		pixel_size_y := 1.0 / float32(_t_h)

		origin_uv := NewVector2f(float32(_this_tile.Src[0]), float32(_this_tile.Src[1]))
		origin_uv.X /= float32(_l.original_texture_size.X)
		origin_uv.Y /= float32(_l.original_texture_size.Y)

		flip_factor_x := float32(0.0)
		if _this_tile.FlipX() {
			flip_factor_x = -1.0
		} else {
			flip_factor_x = 1.0
		}

		flip_factor_y := float32(0.0)
		if _this_tile.FlipY() {
			flip_factor_y = -1.0
		} else {
			flip_factor_y = 1.0
		}

		t := VisualTransform{
			Position:   NewVector2f(float32(_this_tile.Position[0])+BoolToFloat32(_this_tile.FlipX())*float32(_l.tile_size), float32(-_this_tile.Position[1])+BoolToFloat32(_this_tile.FlipY())*float32(_l.tile_size)).Add(_offset.Add(Vector2fDown.Scale(float32(_l.tile_size)))).MultpXY(_scale, _scale),
			Dimensions: NewVector2f(float32(_l.tile_size)*flip_factor_x, float32(_l.tile_size)*flip_factor_y).Scale(_scale),
			Z:          _z + _l.z_offset,
			Scale:      1,
			Tint:       NewRGBA8Float(1.0, 1.0, 1.0, _l.opacity),
			UV1:        origin_uv,
			UV2:        origin_uv.AddXY(float32(_l.tile_size)*pixel_size_x, float32(_l.tile_size)*pixel_size_y),
		}

		world_actual_postion := NewVector2f(float32(_this_tile.Position[0]), float32(-_this_tile.Position[1])).Add(_offset.Add(Vector2fDown.Scale(float32(_l.tile_size)))).MultpXY(_scale, _scale)
		collider_pos := world_actual_postion.Add(t.Dimensions.Scale(0.5)).AddXY(BoolToFloat32(_this_tile.FlipX())*float32(_l.tile_size), 0.0)

		tile_enumset := []string(_l.tileset.Enums[_this_tile.ID])
		if len(tile_enumset) > 0 {
			if tile_enumset[0] == "Solid" {
				newStaticCollisionTileBox2d(collider_pos, _l.tile_size, _scale, _l.physicsLayer)
				_level.SolidTiles.Insert(NewVector2i(_this_tile.Position[0]/_l.tile_size, _this_tile.Position[1]/_l.tile_size), Tile{Enumset: customtypes.ListFromSlice(tile_enumset), Solid: true})
			}

		}

		renderObj := newRenderObject(0, SPRITE_RENDEROBJECTTYPEFUNC)
		renderObj.texture = &_texture
		RenderQuadTreeContainer.Insert(customtypes.Pair[VisualTransform, RenderObject]{t, renderObj}, Rect{Position: world_actual_postion, Size: NewVector2f(float32(_l.tile_size), float32(_l.tile_size))})
	}
}

type ldtkEntity struct {
	Identifier   string
	Position     Vector2f
	GridPosition Vector2i
}

type Tile struct {
	Enumset      customtypes.List[string]
	Solid        bool
	GridPosition Vector2i
}

type Tilemap struct {
	Offset Vector2f
	// tileset                 TileSet
	layers                  customtypes.List[levelLayer]
	grid_width, grid_height int
	tile_size               int
	Scale                   float32
	SolidTiles              customtypes.Map[Vector2i, Tile]
	Entities                customtypes.Map[string, []ldtkEntity]
}

func (level *Tilemap) GridSize() Vector2i {
	return NewVector2i(level.grid_width, level.grid_height)
}

func (level *Tilemap) Tilesize() int {
	return level.tile_size
}

type levelLayer struct {
	identifier            string
	tiles                 customtypes.List[ldtkgo.Tile]
	auto_tiles            customtypes.List[ldtkgo.Tile]
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
