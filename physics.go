package chai

import (
	box2d "github.com/mhamedGd/chai-box2d"
	"github.com/mhamedGd/chai/customtypes"
	. "github.com/mhamedGd/chai/math"
)

const SHAPE_RECTBODY = 0
const SHAPE_CIRCLEBODY = 1

const (
	PHYSICS_LAYER_NONE = 0x0
	PHYSICS_LAYER_ALL  = 0xffff
	PHYSICS_LAYER_1    = 0b0000000000000001
	PHYSICS_LAYER_2    = 0b0000000000000010
	PHYSICS_LAYER_3    = 0b0000000000000100
	PHYSICS_LAYER_4    = 0b0000000000001000
	PHYSICS_LAYER_5    = 0b0000000000010000
	PHYSICS_LAYER_6    = 0b0000000000100000
	PHYSICS_LAYER_7    = 0b0000000001000000
	PHYSICS_LAYER_8    = 0b0000000010000000
	PHYSICS_LAYER_9    = 0b0000000100000000
	PHYSICS_LAYER_10   = 0b0000001000000000
	PHYSICS_LAYER_11   = 0b0000010000000000
	PHYSICS_LAYER_12   = 0b0000100000000000
	PHYSICS_LAYER_13   = 0b0001000000000000
	PHYSICS_LAYER_14   = 0b0010000000000000
	PHYSICS_LAYER_15   = 0b0100000000000000
	PHYSICS_LAYER_16   = 0b1000000000000000
)

const (
	PHYSICS_ENGINE_BOX2D = 0
)

func setPhysicsFunctions(_physicsEngine int) {
	switch _physicsEngine {
	case PHYSICS_ENGINE_BOX2D:
		NewDynamicBodyComponent = newDynamicBodyBox2d
		NewStaticBodyComponent = newStaticBodyBox2d
		NewKinematicBodyComponent = newKinematicBodyBox2d

		getDynamicPosition = getDynamicPositionBox2d
		setDynamicPosition = setDynamicPositionBox2d
		getDynamicRotation = getDynamicRotationBox2d
		setDynamicRotation = setDynamicRotationBox2d

		getStaticPosition = getStaticPositionBox2d
		setStaticPosition = setDynamicPositionBox2d
		getStaticRotation = getStaticRotationBox2d
		setStaticRotation = setStaticRotationBox2d

		getKinematicPosition = getKinematicPositionBox2d
		setKinematicPosition = setKinematicPositionBox2d
		getKinematicRotation = getKinematicRotationBox2d
		setKinematicRotation = setKinematicRotationBox2d

		applyForce = applyForceBox2d
		applyImpulse = applyImpulseBox2d
		applyAngularForce = applyAngularForceBox2d
		applyAngularImpulse = applyAngularImpulseBox2d
		setVelocity = setVelocityBox2d
		getVelociy = getVelocityBox2d
		setAngularVelocity = setAngularVelocityBox2d
		getAngularVelocity = getAngularVelocityBox2d
		LineCast = linecastBox2d
		RayCast = raycastBox2d
		OverlapBox = overlapBoxBox2d
	}
}

type PhysicsWorld struct {
	// cpSpace    *cp.Space
	box2dWorld *box2d.B2World
}

func newPhysicsWorld(_gravity Vector2f) PhysicsWorld {
	return PhysicsWorld{
		// cpSpace: s,
		box2dWorld: newPhysicsWorldBox2D(_gravity),
	}
}

type DynamicBodyComponent struct {
	m_B2Body         *box2d.B2Body
	OwnerEntId       EntId
	OnCollisionBegin customtypes.ChaiEvent1[Collision]
	OnCollisionEnd   customtypes.ChaiEvent1[Collision]
	m_Settings       DynamicBodySettings
}

type StaticBodyComponent struct {
	m_B2Body         *box2d.B2Body
	OwnerEntId       EntId
	OnCollisionBegin customtypes.ChaiEvent1[Collision]
	OnCollisionEnd   customtypes.ChaiEvent1[Collision]
	m_Settings       StaticBodySettings
}

type KinematicBodyComponent struct {
	m_B2Body         *box2d.B2Body
	OwnerEntId       EntId
	OnCollisionBegin customtypes.ChaiEvent1[Collision]
	OnCollisionEnd   customtypes.ChaiEvent1[Collision]
	m_Settings       KinematicBodySettings
}

type DynamicBodySettings struct {
	// BodyType                   PhysicsBodyType
	IsTrigger                  bool
	ColliderShape              int
	StartPosition              Vector2f
	Offset                     Vector2f
	StartDimensions            Vector2f
	StartRotation              float32
	Mass, Friction, Elasticity float32
	GravityScale               float32
	ConstrainRotation          bool
	PhysicsLayer               uint16
}

type StaticBodySettings struct {
	// BodyType                   PhysicsBodyType
	IsTrigger            bool
	ColliderShape        int
	Offset               Vector2f
	StartDimensions      Vector2f
	StartRotation        float32
	Friction, Elasticity float32
	PhysicsLayer         uint16
}

type KinematicBodySettings struct {
	// BodyType                   PhysicsBodyType
	IsTrigger            bool
	ColliderShape        int
	Offset               Vector2f
	StartDimensions      Vector2f
	StartRotation        float32
	Friction, Elasticity float32
	PhysicsLayer         uint16
}

type Collision struct {
	FirstEntity    EntId
	SecondEntity   EntId
	CollisionPoint Vector2f
}

var NewDynamicBodyComponent func(EntId, VisualTransform, *DynamicBodySettings) DynamicBodyComponent
var NewStaticBodyComponent func(EntId, VisualTransform, *StaticBodySettings) StaticBodyComponent
var NewKinematicBodyComponent func(EntId, VisualTransform, *KinematicBodySettings) KinematicBodyComponent

var getDynamicPosition func(*box2d.B2Body) Vector2f
var setDynamicPosition func(*box2d.B2Body, Vector2f)
var getDynamicRotation func(*box2d.B2Body) float32
var setDynamicRotation func(*box2d.B2Body, float32)

var getStaticPosition func(*box2d.B2Body) Vector2f
var setStaticPosition func(*box2d.B2Body, Vector2f)
var getStaticRotation func(*box2d.B2Body) float32
var setStaticRotation func(*box2d.B2Body, float32)

var getKinematicPosition func(*box2d.B2Body) Vector2f
var setKinematicPosition func(*box2d.B2Body, Vector2f)
var getKinematicRotation func(*box2d.B2Body) float32
var setKinematicRotation func(*box2d.B2Body, float32)

var applyForce func(*box2d.B2Body, Vector2f, Vector2f)
var applyImpulse func(*box2d.B2Body, Vector2f, Vector2f)
var applyAngularForce func(*box2d.B2Body, float32, Vector2f)
var applyAngularImpulse func(*box2d.B2Body, float32, Vector2f)
var setVelocity func(*box2d.B2Body, float32, float32)
var getVelociy func(*box2d.B2Body) Vector2f
var setAngularVelocity func(*box2d.B2Body, float32)
var getAngularVelocity func(*box2d.B2Body) float32

// ///////
func (_dc *DynamicBodyComponent) GetPosition() Vector2f {
	return getDynamicPosition(_dc.m_B2Body)
}
func (_dc *DynamicBodyComponent) SetPosition(_position Vector2f) {
	setDynamicPosition(_dc.m_B2Body, _position)
}
func (_dc *DynamicBodyComponent) GetRotation() float32 {
	return getDynamicRotation(_dc.m_B2Body)
}
func (_dc *DynamicBodyComponent) SetRotation(_rotation float32) {
	setDynamicRotation(_dc.m_B2Body, _rotation)
}

// ///////
func (_dc *StaticBodyComponent) GetPosition() Vector2f {
	return getStaticPosition(_dc.m_B2Body)
}
func (_dc *StaticBodyComponent) SetPosition(_position Vector2f) {
	setStaticPosition(_dc.m_B2Body, _position)
}
func (_dc *StaticBodyComponent) GetRotation() float32 {
	return getStaticRotation(_dc.m_B2Body)
}
func (_dc *StaticBodyComponent) SetRotation(_rotation float32) {
	setStaticRotation(_dc.m_B2Body, _rotation)
}

// ///////
func (_dc *KinematicBodyComponent) GetPosition() Vector2f {
	return getKinematicPosition(_dc.m_B2Body)
}
func (_dc *KinematicBodyComponent) SetPosition(_position Vector2f) {
	setKinematicPosition(_dc.m_B2Body, _position)
}
func (_dc *KinematicBodyComponent) GetRotation() float32 {
	return getKinematicRotation(_dc.m_B2Body)
}
func (_dc *KinematicBodyComponent) SetRotation(_rotation float32) {
	setKinematicRotation(_dc.m_B2Body, _rotation)
}

/////////

func (_dc *DynamicBodyComponent) ApplyForce(_forceAmount, _pivot Vector2f) {
	applyForce(_dc.m_B2Body, _forceAmount, _pivot)
}
func (_dc *DynamicBodyComponent) ApplyImpulse(_forceAmount, _pivot Vector2f) {
	applyImpulse(_dc.m_B2Body, _forceAmount, _pivot)
}
func (_dc *DynamicBodyComponent) ApplyAngularForce(_torque float32, _pivot Vector2f) {
	applyAngularForce(_dc.m_B2Body, _torque, _pivot)
}
func (_dc *DynamicBodyComponent) ApplyAngularImpulse(_torque float32, _pivot Vector2f) {
	applyAngularImpulse(_dc.m_B2Body, _torque, _pivot)
}
func (_dc *DynamicBodyComponent) SetVelocity(_velocity Vector2f) {
	setVelocity(_dc.m_B2Body, _velocity.X, _velocity.Y)
}
func (_dc *DynamicBodyComponent) GetVelocity() Vector2f {
	return getVelociy(_dc.m_B2Body)
}
func (_dc *DynamicBodyComponent) SetAngularVelocity(_angularVelocity float32) {
	setAngularVelocity(_dc.m_B2Body, _angularVelocity)
}
func (_dc *DynamicBodyComponent) GetAngularVelocity() float32 {
	return getAngularVelocity(_dc.m_B2Body)
}

// ///////////////////////////////////////////
// ///// Shape Cast //////////////////////////
type RaycastHit struct {
	HasHit       bool
	OriginPoint  Vector2f
	HitPosition  Vector2f
	Normal       Vector2f
	HitEntity    EntId
	PhysicsLayer uint
}

var LineCast func(Vector2f, Vector2f, uint16) RaycastHit
var RayCast func(_origin Vector2f, _direction Vector2f, _distance float32, _physicsLayer uint16) RaycastHit

var OverlapBox func(_box Rect, _physicsLayer uint16) (customtypes.List[EntId], bool)

// /// ECS Systems Relating to Physics ///////////////
// /////////////////////////////////////////////////////////
func DynamicBodySystem(_thisScene *Scene, _dt float32) {
	Iterate2[VisualTransform, DynamicBodyComponent](func(i EntId, t *VisualTransform, db *DynamicBodyComponent) {
		t.Position = db.GetPosition()
		t.Rotation = db.GetRotation()

		LogF("Position: %v", db.GetPosition())
		// t.Position.X = float32(rb.cpBody.Position().X) + rb.Offset.X
		// t.Position.Y = float32(rb.cpBody.Position().Y) + rb.Offset.Y
		// rb.cpBody.SetAngle(BoolToFloat64(!rb.RBSettings.ConstrainRotation) * rb.cpBody.Angle())
		// t.Rotation = float32(rb.cpBody() * Rad2Deg)

	})
}

///////////// Chipmunk2D Physics /////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////

// type PhysicsBodyType uint8

// const (
// 	Type_BodyStatic    PhysicsBodyType = 0
// 	Type_BodyDynamic   PhysicsBodyType = 1
// 	Type_BodyKinematic PhysicsBodyType = 2
// )

// type ColliderShape uint8

// const Shape_CircleBody ColliderShape = 0
// const Shape_RectBody ColliderShape = 1

// func cpVector2f(v Vector2f) cp.Vector {
// 	return cp.Vector{float64(v.X), float64(v.Y)}
// }

// func chaiCPVector(v cp.Vector) Vector2f {
// 	return NewVector2f(float32(v.X), float32(v.Y))
// }

// func SetGravity(_new_grav Vector2f) {
// 	// GetPhysicsWorld().cpSpace.SetGravity(cpVector2f(_new_grav))
// 	GetPhysicsWorld().box2dWorld.SetGravity(vec2fToB2Vec(_new_grav))
// }

// func newChipmunk2dSpace() {
// worldContactListener = ChaiContactListener{}

// s := cp.NewSpace()
// s.Iterations = 10

// s.SetGravity(cpVector2f(gravity))

// dCollisionHandler = s.NewWildcardCollisionHandler(cp.BODY_DYNAMIC)
// sCollisionHandler = s.NewWildcardCollisionHandler(cp.BODY_STATIC)
// kCollisionHandler = s.NewWildcardCollisionHandler(cp.BODY_KINEMATIC)

// dCollisionHandler.BeginFunc = beginCollision
// sCollisionHandler.BeginFunc = beginCollision
// kCollisionHandler.BeginFunc = beginCollision
// dCollisionHandler.SeparateFunc = endCollision
// sCollisionHandler.SeparateFunc = endCollision
// kCollisionHandler.SeparateFunc = endCollision
// }

// // var worldContactListener ChaiContactListener
// var dCollisionHandler *cp.CollisionHandler
// var sCollisionHandler *cp.CollisionHandler
// var kCollisionHandler *cp.CollisionHandler

// func rbTypeToCpType(body *cp.Body, t PhysicsBodyType) {
// 	switch t {
// 	case Type_BodyDynamic:
// 		body.SetType(cp.BODY_DYNAMIC)
// 	case Type_BodyStatic:
// 		body.SetType(cp.BODY_STATIC)
// 	case Type_BodyKinematic:
// 		body.SetType(cp.BODY_KINEMATIC)
// 	}
// }

// type RigidBodyComponent struct {
// 	cpBody           *cp.Body
// 	cpShape          *cp.Shape
// 	RBSettings       *RigidBodySettings
// 	OwnerEntityId    EntId
// 	OnCollisionBegin customtypes.ChaiEvent1[Collision]
// 	OnCollisionEnd   customtypes.ChaiEvent1[Collision]
// 	Offset           Vector2f
// }

// func freeRigidbody(rb *RigidBodyComponent) {
// 	physics_world.cpSpace.RemoveShape(rb.cpShape)

// 	physics_world.cpSpace.RemoveBody(rb.cpBody)
// }

// type RigidBodySettings struct {
// 	IsTrigger                  bool
// 	BodyType                   PhysicsBodyType
// 	ColliderShape              ColliderShape
// 	StartPosition              Vector2f
// 	Offset                     Vector2f
// 	StartDimensions            Vector2f
// 	StartRotation              float32
// 	Mass, Friction, Elasticity float32
// 	ConstrainRotation          bool
// 	PhysicsLayer               uint
// }

// func NewRigidBody(entityId EntId, rbSettings *RigidBodySettings) RigidBodyComponent {
// 	size := cpVector2f(rbSettings.StartDimensions).Mult(0.5)
// 	body := cp.NewBody(0.0, 0.0)

// 	body.SetPosition(cpVector2f(rbSettings.StartPosition))
// 	body.SetAngle(float64(rbSettings.StartRotation) * Deg2Rad)

// 	rbTypeToCpType(body, rbSettings.BodyType)
// 	var shape *cp.Shape
// 	switch rbSettings.ColliderShape {
// 	case Shape_RectBody:
// 		// shape = cp.NewBox(body, size.X, size.Y/2.0, float64(rbSettings.StartRotation)*PI/180.0)
// 		box := cp.NewBB(-size.X, -size.Y, size.X, size.Y)
// 		body.SetMoment(cp.MomentForBox(float64(rbSettings.Mass), box.T-box.B, box.R-box.L))
// 		shape = cp.NewBox2(body, box, 0.0)
// 		shape.SetMass(float64(rbSettings.Mass))

// 		// shape.SetCollisionType()
// 	case Shape_CircleBody:
// 		body.SetMoment(cp.MomentForCircle(float64(rbSettings.Mass), 0.0, float64(2*PI*rbSettings.StartDimensions.X), cpVector2f(Vector2fZero)))
// 		shape = cp.NewCircle(body, size.X, cpVector2f(Vector2fZero))
// 		shape.SetMass(float64(rbSettings.Mass))
// 	}

// 	if rbSettings.ConstrainRotation {
// 		// pivotJoint := cp.NewPivotJoint(body, physics_world.cpSpace.StaticBody, cpVector2f(NewVector2f(float32(size.X)/2.0, float32(size.Y)/2.0)))
// 		// physics_world.cpSpace.AddConstraint(pivotJoint)
// 	}

// 	shape.Filter.Categories = PHYSICS_LAYER_ALL
// 	shape.Filter.Mask = rbSettings.PhysicsLayer

// 	shape.SetElasticity(float64(rbSettings.Elasticity))
// 	shape.SetFriction(float64(rbSettings.Friction))
// 	shape.SetSensor(rbSettings.IsTrigger)

// 	body.UserData = entityId
// 	shape.UserData = body

// 	GetPhysicsWorld().cpSpace.AddBody(body)
// 	GetPhysicsWorld().cpSpace.AddShape(shape)

// 	return RigidBodyComponent{
// 		cpBody:           body,
// 		cpShape:          shape,
// 		RBSettings:       rbSettings,
// 		OwnerEntityId:    entityId,
// 		OnCollisionBegin: customtypes.ChaiEvent1[Collision]{listeners: NewList[EventFunc1[Collision]]()},
// 		OnCollisionEnd:   customtypes.ChaiEvent1[Collision]{listeners: NewList[EventFunc1[Collision]]()},
// 		Offset:           rbSettings.Offset,
// 	}
// }

// func (rb *RigidBodyComponent) SetPosition(newPosition Vector2f) {
// 	rb.cpBody.SetPosition(cpVector2f(newPosition))
// }

// func (rb *RigidBodyComponent) SetRotation(newRotation float32) {
// 	rb.cpBody.SetAngle(float64(newRotation) * Deg2Rad)
// }

// func (rb *RigidBodyComponent) SetVelocity(newVelocity Vector2f) {
// 	rb.cpBody.SetVelocity(float64(newVelocity.X), float64(newVelocity.Y))
// }

// func (rb *RigidBodyComponent) SetAngularVelocity(newAngularVelocity float32) {
// 	rb.cpBody.SetAngularVelocity(float64(newAngularVelocity))
// }

// func (rb *RigidBodyComponent) ApplyForce(forceAmount Vector2f) {
// 	rb.cpBody.ApplyForceAtWorldPoint(cpVector2f(forceAmount), cpVector2f(rb.GetPosition()))
// }

// func (rb *RigidBodyComponent) GetPosition() Vector2f {
// 	return chaiCPVector(rb.cpBody.Position())
// }

// func (rb *RigidBodyComponent) GetVelocity() Vector2f {
// 	return chaiCPVector(rb.cpBody.Velocity())
// }

// func (rb *RigidBodyComponent) GetAngularVelocity() float32 {
// 	return float32(rb.cpBody.AngularVelocity() * Rad2Deg)
// }

// func (rb *RigidBodyComponent) OnCollisionTouch() {
// }

// func beginCollision(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
// 	bodyA, bodyB := arb.Bodies()

// 	if bodyA.UserData.(EntId) > 0 {

// 		rbA, _ := ecs.Read[RigidBodyComponent](current_scene.Ecs_World, bodyA.UserData.(EntId))

// 		rbA.OnCollisionBegin.Invoke(Collision{CollisionPoint: chaiCPVector(arb.ContactPointSet().Points[0].PointA), EntA: bodyA.UserData.(EntId), EntB: bodyB.UserData.(EntId)})
// 	}

// 	return true
// }

// func endCollision(arb *cp.Arbiter, space *cp.Space, userData interface{}) {
// 	bodyA, bodyB := arb.Bodies()

// 	if bodyA.UserData.(EntId) > 0 {
// 		rbA, _ := ecs.Read[RigidBodyComponent](current_scene.Ecs_World, bodyA.UserData.(EntId))

// 		rbA.OnCollisionEnd.Invoke(Collision{CollisionPoint: chaiCPVector(arb.ContactPointSet().Points[0].PointA), EntA: bodyA.UserData.(EntId), EntB: bodyB.UserData.(EntId)})

// 	}
// }

// type Collision struct {
// 	CollisionPoint Vector2f
// 	EntA           EntId
// 	EntB           EntId
// }

// type RaycastHit struct {
// 	HasHit       bool
// 	OriginPoint  Vector2f
// 	HitPosition  Vector2f
// 	Normal       Vector2f
// 	PhysicsLayer uint
// }

// func RayCast(origin, direction Vector2f, distance float32, physicsLayer uint) RaycastHit {
// 	hit := RaycastHit{}
// 	info := physics_world.cpSpace.SegmentQueryFirst(cpVector2f(origin), cpVector2f(origin.Add(direction.Scale(distance))), 0.0, cp.NewShapeFilter(cp.NO_GROUP, physicsLayer, physicsLayer))

// 	hit.OriginPoint = origin
// 	hit.HitPosition = chaiCPVector(info.Point)
// 	hit.Normal = chaiCPVector(info.Normal)
// 	hit.HasHit = info.Shape != nil

// 	return hit
// }

// func LineCast(origin, distanation Vector2f, physicsLayer uint) RaycastHit {
// 	hit := RaycastHit{}
// 	info := physics_world.cpSpace.SegmentQueryFirst(cpVector2f(origin), cpVector2f(distanation), 0.0, cp.NewShapeFilter(cp.NO_GROUP, physicsLayer, physicsLayer))

// 	hit.OriginPoint = origin
// 	hit.HitPosition = chaiCPVector(info.Point)
// 	hit.Normal = chaiCPVector(info.Normal)
// 	hit.HasHit = info.Shape != nil

// 	return hit
// }

// func RigidBodySystem(_this_scene *Scene, _dt float32) {
// 	Iterate2[VisualTransform, DynamicBodyBox2D](func(i EntId, t *VisualTransform, rb *DynamicBodyBox2D) {
// 		t.Position = rb.GetPosition()
// 		t.Rotation = rb.GetRotation()
// 		// t.Position.X = float32(rb.cpBody.Position().X) + rb.Offset.X
// 		// t.Position.Y = float32(rb.cpBody.Position().Y) + rb.Offset.Y
// 		// rb.cpBody.SetAngle(BoolToFloat64(!rb.RBSettings.ConstrainRotation) * rb.cpBody.Angle())
// 		// t.Rotation = float32(rb.cpBody() * Rad2Deg)

// 	})
// }

//////////////////////////////////////////////////////////////////////////////////////////
