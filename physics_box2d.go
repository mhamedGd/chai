package chai

import (
	box2d "github.com/mhamedGd/chai-box2d"
	"github.com/mhamedGd/chai/customtypes"
	. "github.com/mhamedGd/chai/math"
)

func b2VecToVec2f(_b2Vec box2d.B2Vec2) Vector2f {
	return NewVector2f(float32(_b2Vec.X), float32(_b2Vec.Y))
}

func vec2fToB2Vec(_vec2f Vector2f) box2d.B2Vec2 {
	return box2d.B2Vec2{X: float64(_vec2f.X), Y: float64(_vec2f.Y)}
}

func newPhysicsWorldBox2D(_gravity Vector2f) *box2d.B2World {
	phy_wrld := box2d.MakeB2World(vec2fToB2Vec(_gravity))
	phy_wrld.SetContactListener(box2dGlobalContactListener)
	phy_wrld.SetContactFilter(box2dGlobalContactFilter)
	return &phy_wrld
}

func newDynamicBodyBox2d(_entityId EntId, _visualTransform VisualTransform, _dbSettings *DynamicBodySettings) DynamicBodyComponent {
	// bodydef := box2d.MakeB2BodyDef()
	bodydef := box2d.MakeB2BodyDef()
	bodydef.Type = box2d.B2BodyType.B2_dynamicBody
	bodydef.Position.Set(float64(_visualTransform.Position.X), float64(_visualTransform.Position.Y))
	bodydef.AllowSleep = false
	bodydef.FixedRotation = _dbSettings.ConstrainRotation
	bodydef.GravityScale = float64(_dbSettings.GravityScale)
	body2d := newphysicsbodyBox2d(_entityId, _visualTransform, _dbSettings.ColliderShape, _dbSettings.PhysicsLayer, _dbSettings.Mass, _dbSettings.Friction, _dbSettings.Elasticity, _dbSettings.IsTrigger, &bodydef)
	return DynamicBodyComponent{
		m_B2Body:         body2d,
		OwnerEntId:       _entityId,
		OnCollisionBegin: customtypes.NewChaiEvent1[Collision](),
		OnCollisionEnd:   customtypes.NewChaiEvent1[Collision](),
		m_Settings:       *_dbSettings,
	}
}

func newStaticBodyBox2d(_entityId EntId, _visualTransform VisualTransform, _sbSettings *StaticBodySettings) StaticBodyComponent {
	bodydef := box2d.MakeB2BodyDef()
	bodydef.Type = box2d.B2BodyType.B2_staticBody
	bodydef.Position.Set(float64(_visualTransform.Position.X), float64(_visualTransform.Position.Y))
	bodydef.AllowSleep = false
	bodydef.FixedRotation = true
	bodydef.GravityScale = 0.0
	body2d := newphysicsbodyBox2d(_entityId, _visualTransform, _sbSettings.ColliderShape, _sbSettings.PhysicsLayer, 1000000, _sbSettings.Friction, _sbSettings.Elasticity, _sbSettings.IsTrigger, &bodydef)
	return StaticBodyComponent{
		m_B2Body:         body2d,
		OwnerEntId:       _entityId,
		OnCollisionBegin: customtypes.NewChaiEvent1[Collision](),
		OnCollisionEnd:   customtypes.NewChaiEvent1[Collision](),
		m_Settings:       *_sbSettings,
	}
}

func newKinematicBodyBox2d(_entityId EntId, _visualTransform VisualTransform, _kbSettings *KinematicBodySettings) KinematicBodyComponent {
	bodydef := box2d.MakeB2BodyDef()
	bodydef.Type = box2d.B2BodyType.B2_kinematicBody
	// bodydef.Position.Set(float64(_visualTransform.Position.X), float64(_visualTransform.Position.Y))
	bodydef.AllowSleep = false
	bodydef.FixedRotation = true
	bodydef.GravityScale = 0.0
	body2d := newphysicsbodyBox2d(_entityId, _visualTransform, _kbSettings.ColliderShape, _kbSettings.PhysicsLayer, 1000000, _kbSettings.Friction, _kbSettings.Elasticity, _kbSettings.IsTrigger, &bodydef)
	return KinematicBodyComponent{
		m_B2Body:         body2d,
		OwnerEntId:       _entityId,
		OnCollisionBegin: customtypes.NewChaiEvent1[Collision](),
		OnCollisionEnd:   customtypes.NewChaiEvent1[Collision](),
		m_Settings:       *_kbSettings,
	}
}

func newphysicsbodyBox2d(_entId EntId, _visualTransform VisualTransform, _colliderShape int, _physicsLayer uint16, _density, _friction, _restitution float32, _isTrigger bool, _bodyDef *box2d.B2BodyDef) *box2d.B2Body {
	body := physics_world.box2dWorld.CreateBody(_bodyDef)
	body.SetTransform(vec2fToB2Vec(_visualTransform.Position), float64(_visualTransform.Rotation*Deg2Rad))
	fd := box2d.MakeB2FixtureDef()
	fd.Filter.CategoryBits = _physicsLayer
	// fd.Filter.MaskBits = PHYSICS_LAYER_1

	switch _colliderShape {
	case SHAPE_RECTBODY:
		shape := box2d.MakeB2PolygonShape()
		shape.SetAsBox(float64(_visualTransform.Dimensions.X)/2.0, float64(_visualTransform.Dimensions.Y)/2.0)
		fd.Shape = &shape
	case SHAPE_CIRCLEBODY:
		shape := box2d.MakeB2CircleShape()
		shape.SetRadius(float64(_visualTransform.Dimensions.X) / 2.0)
		fd.Shape = &shape
	default:
		ErrorF("[%v]: Collider Shape Unknown", _entId)
	}
	fd.Density = float64(_density)
	fd.Friction = float64(_friction)
	fd.Restitution = float64(_restitution)
	fixture := body.CreateFixtureFromDef(&fd)
	fixture.SetSensor(_isTrigger)
	fixture.SetFilterData(fd.Filter)

	body.SetUserData(int(_entId))

	return body
}

// Position And Rotation Operations ///////////////////////////
///////////////////////////////////////////////////////////////

func getDynamicPositionBox2d(_body *box2d.B2Body) Vector2f {
	return b2VecToVec2f(_body.GetPosition())
}
func setDynamicPositionBox2d(_body *box2d.B2Body, _position Vector2f) {
	_body.SetTransform(vec2fToB2Vec(_position), _body.GetAngle())
}
func setDynamicRotationBox2d(_body *box2d.B2Body, _rotation float32) {
	_body.SetTransform(_body.GetPosition(), float64(_rotation)*Deg2Rad)
}
func getDynamicRotationBox2d(_body *box2d.B2Body) float32 {
	return float32(_body.GetAngle()) * Rad2Deg
}

func applyForceBox2d(_body *box2d.B2Body, _forceAmount Vector2f, _pivot Vector2f) {
	_body.ApplyForce(vec2fToB2Vec(_forceAmount), vec2fToB2Vec(_pivot), true)
}
func applyImpulseBox2d(_body *box2d.B2Body, _forceAmount Vector2f, _pivot Vector2f) {
	_body.ApplyLinearImpulse(vec2fToB2Vec(_forceAmount), vec2fToB2Vec(_pivot), true)
}

func applyAngularForceBox2d(_body *box2d.B2Body, _torque float32, _pivot Vector2f) {
	_body.ApplyTorque(float64(_torque), true)
}
func applyAngularImpulseBox2d(_body *box2d.B2Body, _torque float32, _pivot Vector2f) {
	_body.ApplyAngularImpulse(float64(_torque), true)
}

func setVelocityBox2d(_body *box2d.B2Body, _xVelocity, _yVelocity float32) {
	_body.M_linearVelocity.X = float64(_xVelocity)
	_body.M_linearVelocity.Y = float64(_yVelocity)
}
func getVelocityBox2d(_body *box2d.B2Body) Vector2f {
	return b2VecToVec2f(_body.GetLinearVelocity())
}

func setAngularVelocityBox2d(_body *box2d.B2Body, _angularVelocity float32) {
	_body.SetAngularVelocity(float64(_angularVelocity))
}
func getAngularVelocityBox2d(_body *box2d.B2Body) float32 {
	return float32(_body.M_angularVelocity)
}

func getStaticPositionBox2d(_body *box2d.B2Body) Vector2f {
	return b2VecToVec2f(_body.GetPosition())
}
func setStaticPositionBox2d(_body *box2d.B2Body, _position Vector2f) {
	_body.SetTransform(vec2fToB2Vec(_position), _body.GetAngle())
}
func setStaticRotationBox2d(_body *box2d.B2Body, _rotation float32) {
	_body.SetTransform(_body.GetPosition(), float64(_rotation)*Deg2Rad)
}
func getStaticRotationBox2d(_body *box2d.B2Body) float32 {
	return float32(_body.GetAngle()) * Rad2Deg
}

func getKinematicPositionBox2d(_body *box2d.B2Body) Vector2f {
	return b2VecToVec2f(_body.GetPosition())
}
func setKinematicPositionBox2d(_body *box2d.B2Body, _position Vector2f) {
	_body.SetTransform(vec2fToB2Vec(_position), _body.GetAngle())
}
func setKinematicRotationBox2d(_body *box2d.B2Body, _rotation float32) {
	_body.SetTransform(_body.GetPosition(), float64(_rotation)*Deg2Rad)
}
func getKinematicRotationBox2d(_body *box2d.B2Body) float32 {
	return float32(_body.GetAngle()) * Rad2Deg
}

///////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////

var box2dGlobalContactListener box2dContactListener

type box2dContactListener struct {
	box2d.B2ContactListenerInterface
}

func (l box2dContactListener) BeginContact(_contact box2d.B2ContactInterface) {
	entA := _contact.GetFixtureA().GetBody().GetUserData().(int)
	entB := _contact.GetFixtureB().GetBody().GetUserData().(int)

	var worldManifold box2d.B2WorldManifold
	_contact.GetWorldManifold(&worldManifold)
	collisionPoint := b2VecToVec2f(worldManifold.Points[0])

	dynamic_b_A := GetComponentPtr[DynamicBodyComponent](GetCurrentScene(), EntId(entA))
	static_b_A := GetComponentPtr[StaticBodyComponent](GetCurrentScene(), EntId(entA))
	dynamic_b_B := GetComponentPtr[DynamicBodyComponent](GetCurrentScene(), EntId(entB))
	static_b_B := GetComponentPtr[StaticBodyComponent](GetCurrentScene(), EntId(entB))
	collisionA := Collision{
		FirstEntity:    EntId(entA),
		SecondEntity:   EntId(entB),
		CollisionPoint: collisionPoint,
	}
	collisionB := Collision{
		FirstEntity:    EntId(entB),
		SecondEntity:   EntId(entA),
		CollisionPoint: collisionPoint,
	}
	if dynamic_b_A != nil {
		dynamic_b_A.OnCollisionBegin.Invoke(collisionA)
	} else if static_b_A != nil {
		static_b_A.OnCollisionBegin.Invoke(collisionA)
	}

	if dynamic_b_B != nil {
		dynamic_b_B.OnCollisionBegin.Invoke(collisionB)
	} else if static_b_B != nil {
		static_b_B.OnCollisionBegin.Invoke(collisionB)
	}
}

func (l box2dContactListener) EndContact(_contact box2d.B2ContactInterface) {
	entA := _contact.GetFixtureA().GetBody().GetUserData().(int)
	entB := _contact.GetFixtureB().GetBody().GetUserData().(int)

	var worldManifold box2d.B2WorldManifold
	_contact.GetWorldManifold(&worldManifold)
	collisionPoint := b2VecToVec2f(worldManifold.Points[0])

	dynamic_b_A := GetComponentPtr[DynamicBodyComponent](GetCurrentScene(), EntId(entA))
	static_b_A := GetComponentPtr[StaticBodyComponent](GetCurrentScene(), EntId(entA))
	dynamic_b_B := GetComponentPtr[DynamicBodyComponent](GetCurrentScene(), EntId(entB))
	static_b_B := GetComponentPtr[StaticBodyComponent](GetCurrentScene(), EntId(entB))
	collisionA := Collision{
		FirstEntity:    EntId(entA),
		SecondEntity:   EntId(entB),
		CollisionPoint: collisionPoint,
	}
	collisionB := Collision{
		FirstEntity:    EntId(entB),
		SecondEntity:   EntId(entA),
		CollisionPoint: collisionPoint,
	}
	if dynamic_b_A != nil {
		dynamic_b_A.OnCollisionEnd.Invoke(collisionA)
	} else if static_b_A != nil {
		static_b_A.OnCollisionEnd.Invoke(collisionA)
	}

	if dynamic_b_B != nil {
		dynamic_b_B.OnCollisionEnd.Invoke(collisionB)
	} else if static_b_B != nil {
		static_b_B.OnCollisionEnd.Invoke(collisionB)
	}
}

func (listener box2dContactListener) PreSolve(_contact box2d.B2ContactInterface, _oldManifold box2d.B2Manifold) {
	// Handle pre-solving the _contact
}

func (listener box2dContactListener) PostSolve(_contact box2d.B2ContactInterface, _impulse *box2d.B2ContactImpulse) {
	// Handle post-solving the _contact
}

func linecastBox2d(_origin, _target Vector2f, _physicsMask uint16) RaycastHit {
	return docastBox2d(_origin, _target, _physicsMask)
}

func raycastBox2d(_origin, _direction Vector2f, _distance float32, _physicsMask uint16) RaycastHit {
	return docastBox2d(_origin, _origin.Add(_direction.Scale(_distance)), _physicsMask)
}

func docastBox2d(_origin, _distanation Vector2f, _physicsMask uint16) RaycastHit {
	hit := RaycastHit{}
	physics_world.box2dWorld.RayCast(func(fixture *box2d.B2Fixture, point, normal box2d.B2Vec2, fraction float64) float64 {
		if fixture.GetFilterData().CategoryBits&_physicsMask == 0 {
			return 1.0
		}
		hit.HasHit = true
		hit.HitPosition = b2VecToVec2f(point)
		hit.Normal = b2VecToVec2f(normal)
		hit.HitEntity = EntId(fixture.GetBody().GetUserData().(int))

		return fraction
	}, vec2fToB2Vec(_origin), vec2fToB2Vec(_distanation))
	hit.OriginPoint = _origin

	return hit
}

type boxCastQueryCallback struct {
	FoundBodies []*box2d.B2Body
}

func (callback *boxCastQueryCallback) ReportFixture(_fixture *box2d.B2Fixture) bool {
	callback.FoundBodies = append(callback.FoundBodies, _fixture.GetBody())

	return true
}

func overlapBoxBox2d(_rect Rect, _physicsLayer uint16) (customtypes.List[EntId], bool) {
	aabb := box2d.MakeB2AABB()
	aabb.LowerBound.Set(float64(_rect.Position.X)-float64(_rect.Size.X)/2.0, float64(_rect.Position.Y)-float64(_rect.Size.Y)/2.0)
	aabb.UpperBound.Set(float64(_rect.Position.X)+float64(_rect.Size.X)/2.0, float64(_rect.Position.Y)+float64(_rect.Size.Y)/2.0)
	bodiesQuery := boxCastQueryCallback{}
	bodiesQuery.FoundBodies = make([]*box2d.B2Body, 0)

	physics_world.box2dWorld.QueryAABB(bodiesQuery.ReportFixture, aabb)

	ent_ids := customtypes.NewList[EntId]()
	for _, b := range bodiesQuery.FoundBodies {
		if b.GetFixtureList().GetFilterData().CategoryBits&_physicsLayer == 0 {
			continue
		}
		ent_ids.PushBack(EntId(b.GetUserData().(int)))
	}

	return ent_ids, ent_ids.Count() > 0
}

var box2dGlobalContactFilter contactFilterInterface

type contactFilterInterface struct {
	box2d.B2ContactFilterInterface
}

func (cfi contactFilterInterface) ShouldCollide(_fixtureA *box2d.B2Fixture, _fixtureB *box2d.B2Fixture) bool {
	filterA := _fixtureA.GetFilterData()
	filterB := _fixtureB.GetFilterData()
	return filterA.CategoryBits&filterB.CategoryBits > 0 && filterA.GroupIndex == filterB.GroupIndex
}
