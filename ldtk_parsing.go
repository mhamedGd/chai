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

	last_layer_tileSize := 0

	total_entites := customtypes.NewMap[string, []ldtkEntity]()
	for _, level := range ldtk_reader.Levels {

		tile_size := ldtk_reader.Tilesets[0].GridSize
		total_layers := customtypes.NewList[levelLayer]()

		for li := len(level.Layers) - 1; li >= 0; li-- {

			l := level.Layers[li]

			last_layer_tileSize = l.GridSize

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
					m_Identifier: v.Identifier,
					Position:     IntArrToVec2f(v.Position),
					GridPosition: IntArrToVec2i(v.Position),
				}))
			}

			tiles := l.Tiles
			for _, v := range tiles {
				total_tiles.PushBack(*v)

			}

			m_AutoTiles := l.AutoTiles
			// total_autotiles.PushBackArray(m_AutoTiles)
			if m_AutoTiles == nil {
				ErrorF("Layer (%v): Tile Doesn't Exist", l.Identifier)
			}
			for _, v := range m_AutoTiles {
				total_autotiles.PushBack(*v)
			}
			tileset_original_size := NewVector2i(0, 0)
			if l.Tileset != nil {
				tileset_original_size = NewVector2i(l.Tileset.Width, l.Tileset.Height)
			}
			total_layers.PushBack(levelLayer{
				m_Identifier:          l.Identifier,
				m_Tiles:               total_tiles,
				m_AutoTiles:           total_autotiles,
				m_TilesetTexture:      texture,
				m_OriginalTextureSize: tileset_original_size,
				m_Opacity:             l.Opacity,
				m_TileSize:            l.GridSize,
				m_TileSet:             l.Tileset,
				m_LayerType:           l.Type,
				m_PhysicsLayer:        PHYSICS_LAYER_1,
				m_ZOffset:             float32(li),
			})
		}

		temp_levels_map.Insert(level.Identifier, Tilemap{
			m_Layers: total_layers,

			m_GridWidth:  level.Width / tile_size,
			m_GridHeight: level.Height / tile_size,
			m_TileSize:   last_layer_tileSize,
			Scale:        1.0,
			SolidTiles:   customtypes.NewMap[Vector2i, Tile](),
			Entities:     total_entites,
		})

	}

	return temp_levels_map
}

func LoadTilemapLevel(scene *Scene, _levelName string, _allLevels customtypes.Map[string, Tilemap], _z, _scale float32, _offset Vector2f) *Tilemap {
	// tilemap_level := ldtkLevels.Get(_levelName)
	level := _allLevels.Get(_levelName)

	for li := 0; li < level.m_Layers.Count(); li++ {
		l := level.m_Layers.Data[li]
		if l.m_LayerType == "Entities" {
			for k, v := range level.Entities.AllItems() {
				ents := v
				for i, _ := range ents {
					// level.Entities.Set(k, v.Add(_offset).Add(Vector2fOne.Scale(float32(l.tile_size)/2.0)))
					ents[i].Position = ents[i].Position.Add(_offset).Add(Vector2fOne.Scale(float32(l.m_TileSize) / 2.0))
				}
				level.Entities.Set(k, ents)
			}
			continue
		}

		texture := l.m_TilesetTexture
		texture_width := texture.Width
		texture_height := texture.Height

		tiles := l.m_Tiles
		createTilesFromList(&level, tiles, l, texture, texture_width, texture_height, _offset, _z, _scale)

		m_AutoTiles := l.m_AutoTiles
		createTilesFromList(&level, m_AutoTiles, l, texture, texture_width, texture_height, _offset, _z, _scale)
	}

	return &level
}

func createTilesFromList(_level *Tilemap, _list customtypes.List[ldtkgo.Tile], _l levelLayer, _texture Texture2D, _tW, _tH int, _offset Vector2f, _z, _scale float32) {
	for i := 0; i < _list.Count(); i++ {
		_this_tile := _list.Data[i]
		pixel_size_x := 1.0 / float32(_tW)
		pixel_size_y := 1.0 / float32(_tH)

		origin_uv := NewVector2f(float32(_this_tile.Src[0]), float32(_this_tile.Src[1]))
		origin_uv.X /= float32(_l.m_OriginalTextureSize.X)
		origin_uv.Y /= float32(_l.m_OriginalTextureSize.Y)

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

		_scaleFactor := _scale * _texture.pixelsToMeterDimensions.X / float32(_l.m_TileSize)
		LogF("%v", _texture.pixelsToMeterDimensions.X)
		t := VisualTransform{
			Position:   NewVector2f(float32(_this_tile.Position[0])+BoolToFloat32(_this_tile.FlipX())*float32(_l.m_TileSize), float32(-_this_tile.Position[1])+BoolToFloat32(_this_tile.FlipY())*float32(_l.m_TileSize)).Add(Vector2fDown.Scale(float32(_l.m_TileSize))).Scale(_scaleFactor).Add(_offset),
			Dimensions: NewVector2f(float32(_l.m_TileSize)*flip_factor_x, float32(_l.m_TileSize)*flip_factor_y).Scale(_scaleFactor),
			Z:          _z + _l.m_ZOffset,
			Scale:      1,
			Tint:       NewRGBA8Float(1.0, 1.0, 1.0, _l.m_Opacity),
			UV1:        origin_uv,
			UV2:        origin_uv.AddXY(float32(_l.m_TileSize)*pixel_size_x, float32(_l.m_TileSize)*pixel_size_y),
		}

		world_actual_postion := NewVector2f(float32(_this_tile.Position[0]), float32(-_this_tile.Position[1])).Scale(_scaleFactor).Add(_offset)
		// collider_pos := world_actual_postion.Add(t.Dimensions.Scale(0.5)).AddXY(BoolToFloat32(_this_tile.FlipX())*float32(_l.tile_size), 0.0)
		collider_pos := world_actual_postion.AddXY(_scaleFactor*float32(_l.m_TileSize)/2.0, _scaleFactor*float32(-_l.m_TileSize)/2.0) // + _scaleFactor)

		tile_enumset := []string(_l.m_TileSet.Enums[_this_tile.ID])
		if len(tile_enumset) > 0 {
			if tile_enumset[0] == "Solid" {
				newStaticCollisionTileBox2d(collider_pos, _l.m_TileSize, _scaleFactor, _l.m_PhysicsLayer)
				_level.SolidTiles.Insert(NewVector2i(_this_tile.Position[0]/_l.m_TileSize, _this_tile.Position[1]/_l.m_TileSize), Tile{Enumset: customtypes.ListFromSlice(tile_enumset), Solid: true})
			}
		}

		renderObj := newRenderObject(0, SPRITE_RENDEROBJECTTYPEFUNC)
		renderObj.texture = &_texture
		RenderQuadTreeContainer.Insert(customtypes.Pair[VisualTransform, RenderObject]{t, renderObj}, Rect{Position: world_actual_postion.AddXY(0.0, -float32(_l.m_TileSize)*_scaleFactor), Size: Vector2fOne.Scale(float32(_l.m_TileSize) * _scaleFactor)})
	}
}

type ldtkEntity struct {
	m_Identifier string
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
	m_Layers                  customtypes.List[levelLayer]
	m_GridWidth, m_GridHeight int
	m_TileSize                int
	Scale                     float32
	SolidTiles                customtypes.Map[Vector2i, Tile]
	Entities                  customtypes.Map[string, []ldtkEntity]
}

func (level *Tilemap) GridSize() Vector2i {
	return NewVector2i(level.m_GridWidth, level.m_GridHeight)
}

func (level *Tilemap) Tilesize() int {
	return level.m_TileSize
}

type levelLayer struct {
	m_Identifier          string
	m_Tiles               customtypes.List[ldtkgo.Tile]
	m_AutoTiles           customtypes.List[ldtkgo.Tile]
	m_TilesetTexture      Texture2D
	m_OriginalTextureSize Vector2i
	m_TileSize            int
	m_Opacity             float32
	m_PhysicsLayer        uint16
	m_TileSet             *ldtkgo.Tileset
	m_LayerType           string
	m_ZOffset             float32
}

func IntArrToVec2f(original []int) Vector2f {
	return NewVector2f(float32(original[0]), float32(original[1]))
}

func IntArrToVec2i(origin []int) Vector2i {
	return NewVector2i(origin[0], origin[1])
}

func newStaticCollisionTileBox2d(_position Vector2f, _tileSize int, _scale float32, _physicsLayer uint16) {
	size := vec2fToB2Vec(Vector2fOne.Scale(_scale * float32(_tileSize) / 2.0))
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
	fd.Filter.CategoryBits = _physicsLayer
	fixture := body.CreateFixtureFromDef(&fd)
	fixture.SetSensor(false)

	body.SetUserData(0)
}
