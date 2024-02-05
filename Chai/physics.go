package chai

import box2d "github.com/ByteArena/box2d"

const Type_BodyDynamic uint8 = 0

func BoxVector2f(v Vector2f) box2d.B2Vec2 {
	return box2d.MakeB2Vec2(float64(v.X), float64(v.Y))
}

func BoxVector2XY(x, y float32) box2d.B2Vec2 {
	return box2d.MakeB2Vec2(float64(x), float64(y))
}

type PhysicsWorld struct {
	box2dWorld box2d.B2World
}

func newPhysicsWorld(gravity Vector2f) PhysicsWorld {
	return PhysicsWorld{
		box2dWorld: box2d.MakeB2World(BoxVector2f(gravity)),
	}
}

type PhysicsBody struct {
	body    *box2d.B2Body
	fixture *box2d.B2Fixture
}

func newPhysicsBody(density, friction float32, phy_world *PhysicsWorld, bodyDef *box2d.B2BodyDef, bodySize Vector2f) *PhysicsBody {
	body := phy_world.box2dWorld.CreateBody(bodyDef)
	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(float64(bodySize.X), float64(bodySize.Y))

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = float64(density)
	fd.Friction = float64(friction)
	fixture := body.CreateFixtureFromDef(&fd)

	return &PhysicsBody{
		body:    body,
		fixture: fixture,
	}
}

func (pb *PhysicsBody) GetPosition() Vector2f {
	return NewVector2f(float32(pb.body.GetPosition().X), float32(pb.body.GetPosition().Y))
}
