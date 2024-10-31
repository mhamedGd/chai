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

func newDynamicBodyBox2d(_entity_id EntId, _visual_transform VisualTransform, _db_settings *DynamicBodySettings) DynamicBodyComponent {
	// bodydef := box2d.MakeB2BodyDef()
	bodydef := box2d.MakeB2BodyDef()
	bodydef.Type = box2d.B2BodyType.B2_dynamicBody
	bodydef.Position.Set(float64(_visual_transform.Position.X), float64(_visual_transform.Position.Y))
	bodydef.AllowSleep = false
	bodydef.FixedRotation = _db_settings.ConstrainRotation
	bodydef.GravityScale = float64(_db_settings.GravityScale)
	body2d := newphysicsbodyBox2d(_entity_id, _visual_transform, _db_settings.ColliderShape, _db_settings.PhysicsLayer, _db_settings.Mass, _db_settings.Friction, _db_settings.Elasticity, _db_settings.IsTrigger, &bodydef)
	return DynamicBodyComponent{
		b2Body:           body2d,
		OwnerEntId:       _entity_id,
		OnCollisionBegin: customtypes.NewChaiEvent1[CollisionBox2D](),
		OnCollisionEnd:   customtypes.NewChaiEvent1[CollisionBox2D](),
		settings:         *_db_settings,
	}
}

func newStaticBodyBox2d(_entity_id EntId, _visual_transform VisualTransform, _sb_settings *StaticBodySettings) StaticBodyComponent {
	bodydef := box2d.MakeB2BodyDef()
	bodydef.Type = box2d.B2BodyType.B2_staticBody
	bodydef.Position.Set(float64(_visual_transform.Position.X), float64(_visual_transform.Position.Y))
	bodydef.AllowSleep = false
	bodydef.FixedRotation = true
	bodydef.GravityScale = 0.0
	body2d := newphysicsbodyBox2d(_entity_id, _visual_transform, _sb_settings.ColliderShape, _sb_settings.PhysicsLayer, 1000000, _sb_settings.Friction, _sb_settings.Elasticity, _sb_settings.IsTrigger, &bodydef)
	return StaticBodyComponent{
		b2Body:           body2d,
		OwnerEntId:       _entity_id,
		OnCollisionBegin: customtypes.NewChaiEvent1[CollisionBox2D](),
		OnCollisionEnd:   customtypes.NewChaiEvent1[CollisionBox2D](),
		settings:         *_sb_settings,
	}
}

func newKinematicBodyBox2d(_entity_id EntId, _visual_transform VisualTransform, _kb_settings *KinematicBodySettings) KinematicBodyComponent {
	bodydef := box2d.MakeB2BodyDef()
	bodydef.Type = box2d.B2BodyType.B2_kinematicBody
	// bodydef.Position.Set(float64(_visual_transform.Position.X), float64(_visual_transform.Position.Y))
	bodydef.AllowSleep = false
	bodydef.FixedRotation = true
	bodydef.GravityScale = 0.0
	body2d := newphysicsbodyBox2d(_entity_id, _visual_transform, _kb_settings.ColliderShape, _kb_settings.PhysicsLayer, 1000000, _kb_settings.Friction, _kb_settings.Elasticity, _kb_settings.IsTrigger, &bodydef)
	return KinematicBodyComponent{
		b2Body:           body2d,
		OwnerEntId:       _entity_id,
		OnCollisionBegin: customtypes.NewChaiEvent1[CollisionBox2D](),
		OnCollisionEnd:   customtypes.NewChaiEvent1[CollisionBox2D](),
		settings:         *_kb_settings,
	}
}

func newphysicsbodyBox2d(_ent_id EntId, _visual_transform VisualTransform, _collider_shape int, _physics_layer uint16, _density, _friction, _restitution float32, _is_trigger bool, _body_def *box2d.B2BodyDef) *box2d.B2Body {
	body := physics_world.box2dWorld.CreateBody(_body_def)
	body.SetTransform(vec2fToB2Vec(_visual_transform.Position), float64(_visual_transform.Rotation*Deg2Rad))
	fd := box2d.MakeB2FixtureDef()
	fd.Filter.CategoryBits = _physics_layer
	// fd.Filter.MaskBits = PHYSICS_LAYER_1

	switch _collider_shape {
	case SHAPE_RECTBODY:
		shape := box2d.MakeB2PolygonShape()
		shape.SetAsBox(float64(_visual_transform.Dimensions.X)/2.0, float64(_visual_transform.Dimensions.Y)/2.0)
		fd.Shape = &shape
	case SHAPE_CIRCLEBODY:
		shape := box2d.MakeB2CircleShape()
		shape.SetRadius(float64(_visual_transform.Dimensions.X) / 2.0)
		fd.Shape = &shape
	default:
		ErrorF("[%v]: Collider Shape Unknown", _ent_id)
	}
	fd.Density = float64(_density)
	fd.Friction = float64(_friction)
	fd.Restitution = float64(_restitution)
	fixture := body.CreateFixtureFromDef(&fd)
	fixture.SetSensor(_is_trigger)
	fixture.SetFilterData(fd.Filter)

	body.SetUserData(int(_ent_id))

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

func applyForceBox2d(_body *box2d.B2Body, _force_amount Vector2f, _pivot Vector2f) {
	_body.ApplyForce(vec2fToB2Vec(_force_amount), vec2fToB2Vec(_pivot), true)
}
func applyImpulseBox2d(_body *box2d.B2Body, _force_amount Vector2f, _pivot Vector2f) {
	_body.ApplyLinearImpulse(vec2fToB2Vec(_force_amount), vec2fToB2Vec(_pivot), true)
}

func applyAngularForceBox2d(_body *box2d.B2Body, _torque float32, _pivot Vector2f) {
	_body.ApplyTorque(float64(_torque), true)
}
func applyAngularImpulseBox2d(_body *box2d.B2Body, _torque float32, _pivot Vector2f) {
	_body.ApplyAngularImpulse(float64(_torque), true)
}

func setVelocityBox2d(_body *box2d.B2Body, _x_velocity, _y_velocity float32) {
	_body.M_linearVelocity.X = float64(_x_velocity)
	_body.M_linearVelocity.Y = float64(_y_velocity)
}
func getVelocityBox2d(_body *box2d.B2Body) Vector2f {
	return b2VecToVec2f(_body.GetLinearVelocity())
}

func setAngularVelocityBox2d(_body *box2d.B2Body, _angular_velocity float32) {
	_body.SetAngularVelocity(float64(_angular_velocity))
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

type CollisionBox2D struct {
	FirstEntity    EntId
	SecondEntity   EntId
	CollisionPoint Vector2f
}

var box2dGlobalContactListener box2dContactListener

type box2dContactListener struct {
	box2d.B2ContactListenerInterface
}

func (l box2dContactListener) BeginContact(contact box2d.B2ContactInterface) {
	entA := contact.GetFixtureA().GetBody().GetUserData().(int)
	entB := contact.GetFixtureB().GetBody().GetUserData().(int)

	var worldManifold box2d.B2WorldManifold
	contact.GetWorldManifold(&worldManifold)
	collisionPoint := b2VecToVec2f(worldManifold.Points[0])

	dynamic_b_A := GetComponentPtr[DynamicBodyComponent](GetCurrentScene(), EntId(entA))
	static_b_A := GetComponentPtr[StaticBodyComponent](GetCurrentScene(), EntId(entA))
	dynamic_b_B := GetComponentPtr[DynamicBodyComponent](GetCurrentScene(), EntId(entB))
	static_b_B := GetComponentPtr[StaticBodyComponent](GetCurrentScene(), EntId(entB))
	collisionA := CollisionBox2D{
		FirstEntity:    EntId(entA),
		SecondEntity:   EntId(entB),
		CollisionPoint: collisionPoint,
	}
	collisionB := CollisionBox2D{
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

func (l box2dContactListener) EndContact(contact box2d.B2ContactInterface) {
	entA := contact.GetFixtureA().GetBody().GetUserData().(int)
	entB := contact.GetFixtureB().GetBody().GetUserData().(int)

	var worldManifold box2d.B2WorldManifold
	contact.GetWorldManifold(&worldManifold)
	collisionPoint := b2VecToVec2f(worldManifold.Points[0])

	dynamic_b_A := GetComponentPtr[DynamicBodyComponent](GetCurrentScene(), EntId(entA))
	static_b_A := GetComponentPtr[StaticBodyComponent](GetCurrentScene(), EntId(entA))
	dynamic_b_B := GetComponentPtr[DynamicBodyComponent](GetCurrentScene(), EntId(entB))
	static_b_B := GetComponentPtr[StaticBodyComponent](GetCurrentScene(), EntId(entB))
	collisionA := CollisionBox2D{
		FirstEntity:    EntId(entA),
		SecondEntity:   EntId(entB),
		CollisionPoint: collisionPoint,
	}
	collisionB := CollisionBox2D{
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

func (listener box2dContactListener) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
	// Handle pre-solving the contact
}

func (listener box2dContactListener) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
	// Handle post-solving the contact
}

func linecastBox2d(_origin, _target Vector2f, _physics_mask uint16) RaycastHit {
	return docastBox2d(_origin, _target, _physics_mask)
}

func raycastBox2d(_origin, _direction Vector2f, _distance float32, _physics_mask uint16) RaycastHit {
	return docastBox2d(_origin, _origin.Add(_direction.Scale(_distance)), _physics_mask)
}

func docastBox2d(_origin, _distanation Vector2f, _physics_mask uint16) RaycastHit {
	hit := RaycastHit{}
	physics_world.box2dWorld.RayCast(func(fixture *box2d.B2Fixture, point, normal box2d.B2Vec2, fraction float64) float64 {
		if fixture.GetFilterData().CategoryBits&_physics_mask == 0 {
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

func (callback *boxCastQueryCallback) ReportFixture(fixture *box2d.B2Fixture) bool {
	callback.FoundBodies = append(callback.FoundBodies, fixture.GetBody())

	return true
}

func overlapBoxBox2d(_rect Rect, _physics_layer uint16) (customtypes.List[EntId], bool) {
	aabb := box2d.MakeB2AABB()
	aabb.LowerBound.Set(float64(_rect.Position.X)-float64(_rect.Size.X)/2.0, float64(_rect.Position.Y)-float64(_rect.Size.Y)/2.0)
	aabb.UpperBound.Set(float64(_rect.Position.X)+float64(_rect.Size.X)/2.0, float64(_rect.Position.Y)+float64(_rect.Size.Y)/2.0)
	bodiesQuery := boxCastQueryCallback{}
	bodiesQuery.FoundBodies = make([]*box2d.B2Body, 0)

	physics_world.box2dWorld.QueryAABB(bodiesQuery.ReportFixture, aabb)

	ent_ids := customtypes.NewList[EntId]()
	for _, b := range bodiesQuery.FoundBodies {
		if b.GetFixtureList().GetFilterData().CategoryBits&_physics_layer == 0 {
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

func (cfi contactFilterInterface) ShouldCollide(fixtureA *box2d.B2Fixture, fixtureB *box2d.B2Fixture) bool {
	filterA := fixtureA.GetFilterData()
	filterB := fixtureB.GetFilterData()
	return filterA.CategoryBits&filterB.CategoryBits > 0 && filterA.GroupIndex == filterB.GroupIndex
}
