package chai

// import "github.com/ByteArena/box2d"

// func b2VecToVec2f(_b2Vec box2d.B2Vec2) Vector2f {
// 	return NewVector2f(float32(_b2Vec.X), float32(_b2Vec.Y))
// }

// func vec2fToB2Vec(_vec2f Vector2f) box2d.B2Vec2 {
// 	return box2d.B2Vec2{X: float64(_vec2f.X), Y: float64(_vec2f.Y)}
// }

// type physicsWorldBox2D struct {
// 	box2d_world box2d.B2World
// }

// func newPhysicsWorldBox2D() physicsWorldBox2D {
// 	phy_wrld := box2d.MakeB2World(vec2fToB2Vec(NewVector2f(0.0, -20.0)))

// 	return physicsWorldBox2D{
// 		box2d_world: phy_wrld,
// 	}
// }

// type RigidbodyBox2D struct {
// 	b2Body           *box2d.B2Body
// 	b2Shape          *box2d.B2Shape
// 	OwnerEntId       EntId
// 	OnCollisionBegin ChaiEvent[int]
// 	OnCollisionEnd   ChaiEvent[int]
// }

// type RigidBodySettingsBox2D struct {
// 	IsTrigger bool
// 	// BodyType                   PhysicsBodyType
// 	ColliderShape              ColliderShape
// 	StartPosition              Vector2f
// 	Offset                     Vector2f
// 	StartDimensions            Vector2f
// 	StartRotation              float32
// 	Mass, Friction, Elasticity float32
// 	ConstrainRotation          bool
// 	PhysicsLayer               uint
// }

// func NewRigidbodyBox2d(entityId EntId, rbSettings *RigidBodySettingsBox2D) RigidbodyBox2D {

// }
