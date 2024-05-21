package chai

import (
	cp "github.com/jakecoffman/cp/v2"
	"github.com/mhamedGd/chai/Chai/ecs"
)

type PhysicsBodyType uint8

const (
	Type_BodyStatic    PhysicsBodyType = 0
	Type_BodyDynamic   PhysicsBodyType = 1
	Type_BodyKinematic PhysicsBodyType = 2
)

const (
	PhysicsLayer_All uint = 0xffffffff
)

type ColliderShape uint8

const Shape_CircleBody ColliderShape = 0
const Shape_RectBody ColliderShape = 1

func cpVector2f(v Vector2f) cp.Vector {
	return cp.Vector{float64(v.X), float64(v.Y)}
}

func chaiCPVector(v cp.Vector) Vector2f {
	return NewVector2f(float32(v.X), float32(v.Y))
}

type PhysicsWorld struct {
	cpSpace *cp.Space
}

// var worldContactListener ChaiContactListener
var dCollisionHandler *cp.CollisionHandler
var sCollisionHandler *cp.CollisionHandler
var kCollisionHandler *cp.CollisionHandler

func newPhysicsWorld(gravity Vector2f) PhysicsWorld {
	// worldContactListener = ChaiContactListener{}

	s := cp.NewSpace()
	s.Iterations = 20
	s.SetGravity(cpVector2f(gravity))

	dCollisionHandler = s.NewWildcardCollisionHandler(cp.BODY_DYNAMIC)
	sCollisionHandler = s.NewWildcardCollisionHandler(cp.BODY_STATIC)
	kCollisionHandler = s.NewWildcardCollisionHandler(cp.BODY_KINEMATIC)

	dCollisionHandler.BeginFunc = beginCollision
	sCollisionHandler.BeginFunc = beginCollision
	kCollisionHandler.BeginFunc = beginCollision
	dCollisionHandler.SeparateFunc = endCollision
	sCollisionHandler.SeparateFunc = endCollision
	kCollisionHandler.SeparateFunc = endCollision

	return PhysicsWorld{
		cpSpace: s,
	}
}

func rbTypeToCpType(body *cp.Body, t PhysicsBodyType) {
	switch t {
	case Type_BodyDynamic:
		body.SetType(cp.BODY_DYNAMIC)
	case Type_BodyStatic:
		body.SetType(cp.BODY_STATIC)
	case Type_BodyKinematic:
		body.SetType(cp.BODY_KINEMATIC)
	}
}

type RigidBodyComponent struct {
	cpBody           *cp.Body
	cpShape          *cp.Shape
	RBSettings       *RigidBodySettings
	OwnerEntityId    EntId
	OnCollisionBegin ChaiEvent[Collision]
	OnCollisionEnd   ChaiEvent[Collision]
	Offset           Vector2f
}

func freeRigidbody(rb *RigidBodyComponent) {
	physics_world.cpSpace.RemoveShape(rb.cpShape)

	physics_world.cpSpace.RemoveBody(rb.cpBody)
}

type RigidBodySettings struct {
	IsTrigger                  bool
	BodyType                   PhysicsBodyType
	ColliderShape              ColliderShape
	StartPosition              Vector2f
	Offset                     Vector2f
	StartDimensions            Vector2f
	StartRotation              float32
	Mass, Friction, Elasticity float32
	ConstrainRotation          bool
	PhysicsLayer               uint
}

func NewRigidBody(entityId EntId, rbSettings *RigidBodySettings) RigidBodyComponent {
	size := cpVector2f(rbSettings.StartDimensions).Mult(0.5)
	body := cp.NewBody(0.0, 0.0)

	body.SetPosition(cpVector2f(rbSettings.StartPosition))
	body.SetAngle(float64(rbSettings.StartRotation) * Deg2Rad)

	rbTypeToCpType(body, rbSettings.BodyType)
	var shape *cp.Shape
	switch rbSettings.ColliderShape {
	case Shape_RectBody:
		// shape = cp.NewBox(body, size.X, size.Y/2.0, float64(rbSettings.StartRotation)*PI/180.0)
		box := cp.NewBB(-size.X, -size.Y, size.X, size.Y)
		shape = cp.NewBox2(body, box, 0.0)
		body.SetMoment(cp.MomentForBox2(float64(rbSettings.Mass), box))
		shape.SetMass(float64(rbSettings.Mass))
		// shape.SetCollisionType()
	case Shape_CircleBody:
		shape = cp.NewCircle(body, size.X, cpVector2f(Vector2fZero))
		body.SetMoment(cp.MomentForCircle(float64(rbSettings.Mass), 0.0, float64(2*PI*rbSettings.StartDimensions.X), cpVector2f(Vector2fZero)))
		shape.SetMass(float64(rbSettings.Mass))
	}

	if rbSettings.ConstrainRotation {
		body.SetMoment(cp.INFINITY)
	}

	shape.Filter.Categories = PhysicsLayer_All
	shape.Filter.Mask = rbSettings.PhysicsLayer

	shape.SetElasticity(float64(rbSettings.Elasticity))
	shape.SetFriction(float64(rbSettings.Friction))
	shape.SetSensor(rbSettings.IsTrigger)

	body.UserData = entityId
	GetPhysicsWorld().cpSpace.AddBody(body)
	GetPhysicsWorld().cpSpace.AddShape(shape)

	return RigidBodyComponent{
		cpBody:           body,
		cpShape:          shape,
		RBSettings:       rbSettings,
		OwnerEntityId:    entityId,
		OnCollisionBegin: ChaiEvent[Collision]{listeners: make([]EventFunc[Collision], 0)},
		OnCollisionEnd:   ChaiEvent[Collision]{listeners: make([]EventFunc[Collision], 0)},
		Offset:           rbSettings.Offset,
	}
}

func (rb *RigidBodyComponent) SetPosition(newPosition Vector2f) {
	rb.cpBody.SetPosition(cpVector2f(newPosition))
}

func (rb *RigidBodyComponent) SetRotation(newRotation float32) {
	rb.cpBody.SetAngle(float64(newRotation))
}

func (rb *RigidBodyComponent) SetVelocity(newVelocity Vector2f) {
	rb.cpBody.SetVelocity(float64(newVelocity.X), float64(newVelocity.Y))
}

func (rb *RigidBodyComponent) SetAngularVelocity(newAngularVelocity float32) {
	rb.cpBody.SetAngularVelocity(float64(newAngularVelocity))
}

func (rb *RigidBodyComponent) ApplyForce(forceAmount Vector2f) {
	rb.cpBody.ApplyForceAtWorldPoint(cpVector2f(forceAmount), cpVector2f(rb.GetPosition()))
}

func (rb *RigidBodyComponent) GetPosition() Vector2f {
	return chaiCPVector(rb.cpBody.Position())
}

func (rb *RigidBodyComponent) GetVelocity() Vector2f {
	return chaiCPVector(rb.cpBody.Velocity())
}

func (rb *RigidBodyComponent) GetAngularVelocity() float32 {
	return float32(rb.cpBody.AngularVelocity() * Rad2Deg)
}

func (rb *RigidBodyComponent) OnCollisionTouch() {
}

func beginCollision(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
	bodyA, bodyB := arb.Bodies()

	rbA, _ := ecs.Read[RigidBodyComponent](current_scene.Ecs_World, bodyA.UserData.(EntId))

	rbA.OnCollisionBegin.Invoke(Collision{CollisionPoint: chaiCPVector(arb.ContactPointSet().Points[0].PointA), EntA: bodyA.UserData.(EntId), EntB: bodyB.UserData.(EntId)})

	return true
}

func endCollision(arb *cp.Arbiter, space *cp.Space, userData interface{}) {
	bodyA, bodyB := arb.Bodies()

	rbA, _ := ecs.Read[RigidBodyComponent](current_scene.Ecs_World, bodyA.UserData.(EntId))

	rbA.OnCollisionEnd.Invoke(Collision{CollisionPoint: chaiCPVector(arb.ContactPointSet().Points[0].PointA), EntA: bodyA.UserData.(EntId), EntB: bodyB.UserData.(EntId)})
}

type Collision struct {
	CollisionPoint Vector2f
	EntA           EntId
	EntB           EntId
}

type RaycastHit struct {
	HasHit       bool
	OriginPoint  Vector2f
	HitPosition  Vector2f
	Normal       Vector2f
	PhysicsLayer uint
}

func RayCast(origin, direction Vector2f, distance float32, physicsLayer uint) RaycastHit {
	hit := RaycastHit{}
	info := physics_world.cpSpace.SegmentQueryFirst(cpVector2f(origin), cpVector2f(origin.Add(direction.Scale(distance))), 0.0, cp.NewShapeFilter(cp.NO_GROUP, physicsLayer, physicsLayer))

	hit.OriginPoint = origin
	hit.HitPosition = chaiCPVector(info.Point)
	hit.Normal = chaiCPVector(info.Normal)
	hit.HasHit = info.Shape != nil

	return hit
}

type RigidBodySystem struct {
	EcsSystem
}

func (rbs *RigidBodySystem) Update(dt float32) {
	Iterate2[Transform, RigidBodyComponent](func(i EntId, t *Transform, rb *RigidBodyComponent) {
		// 		t.Position.X = dbc.phy_body.GetPosition().X
		// 		t.Position.Y = dbc.phy_body.GetPosition().Y
		// 		t.Rotation = float32(dbc.phy_body.body.GetAngle() * Rad2Deg)
		t.Position.X = float32(rb.cpBody.Position().X) + rb.Offset.X
		t.Position.Y = float32(rb.cpBody.Position().Y) + rb.Offset.Y
		t.Rotation = float32(rb.cpBody.Angle() * Rad2Deg)

	})
}
