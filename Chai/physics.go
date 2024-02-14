package chai

import box2d "github.com/ByteArena/box2d"

type PhysicsBodyType uint8

const Type_BodyStatic PhysicsBodyType = 0
const Type_BodyDynamic PhysicsBodyType = 1

type ColliderShape uint8

const Shape_CircleCollider ColliderShape = 0
const Shape_RectCollider ColliderShape = 1

func BoxVector2f(v Vector2f) box2d.B2Vec2 {
	return box2d.MakeB2Vec2(float64(v.X), float64(v.Y))
}

func BoxVector2XY(x, y float32) box2d.B2Vec2 {
	return box2d.MakeB2Vec2(float64(x), float64(y))
}

func Vector2fFromBoxVec(vec box2d.B2Vec2) Vector2f {
	return NewVector2f(float32(vec.X), float32(vec.Y))
}

type PhysicsWorld struct {
	box2dWorld box2d.B2World
}

var worldContactListener ChaiContactListener

func newPhysicsWorld(gravity Vector2f) PhysicsWorld {
	worldContactListener = ChaiContactListener{}
	return PhysicsWorld{
		box2dWorld: box2d.MakeB2World(BoxVector2f(gravity)),
	}
}

type PhysicsBody struct {
	BodyType         PhysicsBodyType
	ColliderShape    ColliderShape
	body             *box2d.B2Body
	fixture          *box2d.B2Fixture
	OwnerEntity      *EcsEntity
	IsTrigger        bool
	OnCollisionStart ChaiEvent[*Collision]
	OnCollisionEnd   ChaiEvent[*Collision]
	Debug_Tint       RGBA8
}

func newPhysicsBody(bodyType PhysicsBodyType, colliderShape ColliderShape, ent *EcsEntity, density, friction float32, isTrigger bool, phy_world *PhysicsWorld, bodyDef *box2d.B2BodyDef, bodySize Vector2f) *PhysicsBody {
	phyBody := PhysicsBody{}
	body := phy_world.box2dWorld.CreateBody(bodyDef)

	fd := box2d.MakeB2FixtureDef()
	switch colliderShape {
	case Shape_RectCollider:
		shape := box2d.MakeB2PolygonShape()
		shape.SetAsBox(float64(bodySize.X)/2.0, float64(bodySize.Y)/2.0)
		fd.Shape = &shape
		break
	case Shape_CircleCollider:
		shape := box2d.MakeB2CircleShape()
		shape.SetRadius(float64(bodySize.X))
		fd.Shape = &shape
		break

	}

	fd.Density = float64(density)
	fd.Friction = float64(friction)
	fixture := body.CreateFixtureFromDef(&fd)
	fixture.SetSensor(isTrigger)
	phyBody = PhysicsBody{
		BodyType:      bodyType,
		ColliderShape: colliderShape,
		body:          body,
		fixture:       fixture,
		OwnerEntity:   ent,
		IsTrigger:     isTrigger,
		Debug_Tint:    WHITE,
	}
	phyBody.OnCollisionStart.init()
	phyBody.OnCollisionEnd.init()
	phyBody.body.SetUserData(&phyBody)

	return &phyBody
}

func (pb *PhysicsBody) GetPosition() Vector2f {
	return NewVector2f(float32(pb.body.GetPosition().X), float32(pb.body.GetPosition().Y))
}

type ChaiContactListener struct {
	box2d.B2ContactListenerInterface
}

func (listener ChaiContactListener) BeginContact(contact box2d.B2ContactInterface) {
	// Handle the beginning of a contact
	phyBodyA := contact.GetFixtureA().GetBody().GetUserData().(*PhysicsBody)
	phyBodyB := contact.GetFixtureB().GetBody().GetUserData().(*PhysicsBody)

	var worldManifold box2d.B2WorldManifold
	contact.GetWorldManifold(&worldManifold)
	collisionPoint := Vector2fFromBoxVec(worldManifold.Points[0])

	phyBodyA.OnCollisionStart.Invoke(&Collision{CollisionPoint: collisionPoint, FirstBody: phyBodyA, SecondBody: phyBodyB})
	phyBodyB.OnCollisionStart.Invoke(&Collision{CollisionPoint: collisionPoint, FirstBody: phyBodyB, SecondBody: phyBodyA})
}

func (listener ChaiContactListener) EndContact(contact box2d.B2ContactInterface) {
	phyBodyA := contact.GetFixtureA().GetBody().GetUserData().(*PhysicsBody)
	phyBodyB := contact.GetFixtureB().GetBody().GetUserData().(*PhysicsBody)

	phyBodyA.OnCollisionEnd.Invoke(&Collision{CollisionPoint: Vector2fZero, FirstBody: phyBodyA, SecondBody: phyBodyB})
	phyBodyB.OnCollisionEnd.Invoke(&Collision{CollisionPoint: Vector2fZero, FirstBody: phyBodyB, SecondBody: phyBodyA})
}

func (listener ChaiContactListener) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
	// Handle pre-solving the contact
}

func (listener ChaiContactListener) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
	// Handle post-solving the contact
}

type Collision struct {
	CollisionPoint Vector2f
	FirstBody      *PhysicsBody
	SecondBody     *PhysicsBody
}

type DynamicBodyComponent struct {
	Active   bool
	phy_body *PhysicsBody
}

func (dc *DynamicBodyComponent) GetPhyiscsBody() *PhysicsBody {
	return dc.phy_body
}

func (dc *DynamicBodyComponent) GetLinearVelocity() Vector2f {
	return NewVector2f(float32(dc.phy_body.body.M_linearVelocity.X), float32(dc.phy_body.body.M_linearVelocity.Y))
}

func (dc *DynamicBodyComponent) SetLinearVelocity(x_velo, y_velo float32) {
	dc.phy_body.body.M_linearVelocity = BoxVector2XY(x_velo, y_velo)
}

func (dc *DynamicBodyComponent) GetAngularVelocity() float32 {
	return float32(dc.phy_body.body.M_angularVelocity)
}

func (dc *DynamicBodyComponent) SetAngularVelocity(ang_velo float32) {
	dc.phy_body.body.M_angularVelocity = float64(ang_velo)
}

func (dc *DynamicBodyComponent) GetAppliedForce() Vector2f {
	return Vector2fFromBoxVec(dc.phy_body.body.M_force)
}

func (dc *DynamicBodyComponent) ApplyForce(force_vector Vector2f) {
	dc.phy_body.body.ApplyForceToCenter(BoxVector2f(force_vector), true)
}

func (dc *DynamicBodyComponent) ApplyForceXY(x_force, y_force float32) {
	dc.phy_body.body.ApplyForceToCenter(BoxVector2XY(x_force, y_force), true)
}

func (dc *DynamicBodyComponent) ApplyAngularForce(_force float32) {
	dc.phy_body.body.ApplyAngularImpulse(float64(_force), true)
}

func (t *DynamicBodyComponent) ComponentSet(val interface{}) { *t = val.(DynamicBodyComponent) }

func NewDynamicBody(ent *EcsEntity, colliderShape ColliderShape, bodySize Vector2f, density, friction, gravity_scale float32, phy_world *PhysicsWorld) DynamicBodyComponent {

	dynamicComp := DynamicBodyComponent{}
	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Position = BoxVector2f(ent.Pos)
	bodyDef.Angle = float64(ent.Rot * PI / 180.0)

	bodyDef.Type = box2d.B2BodyType.B2_dynamicBody
	bodyDef.AllowSleep = false
	bodyDef.FixedRotation = false
	bodyDef.GravityScale = float64(gravity_scale)

	dynamicComp = DynamicBodyComponent{
		Active:   true,
		phy_body: newPhysicsBody(Type_BodyDynamic, colliderShape, ent, density, friction, false, phy_world, &bodyDef, bodySize),
	}
	return dynamicComp
}

type DynamicBodyUpdateSystem struct {
	EcsSystemImpl
}

func (ds *DynamicBodyUpdateSystem) Update(dt float32) {
	EachEntity(DynamicBodyComponent{}, func(entity *EcsEntity, a interface{}) {
		dComp := a.(DynamicBodyComponent)
		entity.Pos.X = dComp.phy_body.GetPosition().X
		entity.Pos.Y = dComp.phy_body.GetPosition().Y
		entity.Rot = float32(dComp.phy_body.body.GetAngle() * 180.0 / PI)

		// if dComp.phy_body.ColliderShape == Shape_CircleCollider {
		// 	Shapes.DrawCircle(entity.Pos, entity.Dimensions.X, dComp.phy_body.Debug_Tint)
		// } else {
		// 	Shapes.DrawRectRotated(entity.Pos, entity.Dimensions, dComp.phy_body.Debug_Tint, entity.Rot)
		// }
		switch dComp.phy_body.ColliderShape {
		case Shape_RectCollider:
			Shapes.DrawRectRotated(entity.Pos, entity.Dimensions, dComp.phy_body.Debug_Tint, entity.Rot)
			break
		case Shape_CircleCollider:
			Shapes.DrawCircle(entity.Pos, entity.Dimensions.X, dComp.phy_body.Debug_Tint)
			break
		}
	})
}

type StaticBodyComponent struct {
	Active   bool
	phy_body *PhysicsBody
}

func (t *StaticBodyComponent) ComponentSet(val interface{}) { *t = val.(StaticBodyComponent) }

func (sb *StaticBodyComponent) GetPhyiscsBody() *PhysicsBody {
	return sb.phy_body
}

func NewStaticBody(ent *EcsEntity, colliderShape ColliderShape, bodySize Vector2f, friction float32, phy_world *PhysicsWorld) StaticBodyComponent {
	staticComp := StaticBodyComponent{}
	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Position = BoxVector2f(ent.Pos)
	bodyDef.Type = box2d.B2BodyType.B2_staticBody
	bodyDef.AllowSleep = false
	bodyDef.FixedRotation = false
	bodyDef.Angle = float64(ent.Rot * PI / 180.0)

	staticComp = StaticBodyComponent{
		Active:   true,
		phy_body: newPhysicsBody(Type_BodyStatic, colliderShape, ent, 0.0, friction, false, phy_world, &bodyDef, bodySize),
	}
	return staticComp
}

func NewTriggerArea(ent *EcsEntity, colliderShape ColliderShape, bodySize Vector2f, phy_world *PhysicsWorld) StaticBodyComponent {
	staticComp := StaticBodyComponent{}
	bodyDef := box2d.MakeB2BodyDef()
	bodyDef.Position = BoxVector2f(ent.Pos)
	bodyDef.Type = box2d.B2BodyType.B2_staticBody
	bodyDef.AllowSleep = false
	bodyDef.FixedRotation = false
	bodyDef.Angle = float64(ent.Rot * PI / 180.0)

	staticComp = StaticBodyComponent{
		Active:   true,
		phy_body: newPhysicsBody(Type_BodyStatic, colliderShape, ent, 0.0, 0.0, true, phy_world, &bodyDef, bodySize),
	}
	return staticComp
}
