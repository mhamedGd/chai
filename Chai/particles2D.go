package chai

type ParticlesSpreadPattern = int

const (
	PARTICLES_LINESPREAD   = 0x000000e1
	PARTICLES_CIRCLESPREAD = 0x000000e2
)

type Particle struct {
	LifePercentage float32
	LifeTime       float32
	Position       Vector2f
	Velocity       Vector2f
	Size           float32
	Rotation       float32
	Color          RGBA8
	Shape          Gfx_Shape
}

type ParticlesShapeBatch struct {
	shapes           *ShapeBatch
	particles        []Particle
	maxParticles     int
	lastFreeParticle int
}

func newParticlesShapeBatch(_maxParticles int) *ParticlesShapeBatch {
	return &ParticlesShapeBatch{
		shapes:           &Shapes,
		particles:        make([]Particle, _maxParticles),
		maxParticles:     _maxParticles,
		lastFreeParticle: 0,
	}
}

func (p *ParticlesShapeBatch) addParticle(shape Gfx_Shape, lifeTime float32, pos, velo Vector2f, color RGBA8, size, rotation float32) {
	i := p.findLastFreeParticle()
	p.particles[i] = Particle{
		Shape:          shape,
		Position:       pos,
		Velocity:       velo,
		Size:           size,
		Rotation:       rotation,
		LifeTime:       lifeTime,
		LifePercentage: 1.0,
		Color:          color,
	}
}

func (p *ParticlesShapeBatch) findLastFreeParticle() int {
	for i := p.lastFreeParticle; i < p.maxParticles; i++ {
		if p.particles[i].LifePercentage <= 0.0 {
			return i
		}
	}

	for i := 0; i < p.lastFreeParticle; i++ {
		if p.particles[i].LifePercentage <= 0.0 {
			return i
		}
	}

	return 0
}

type ParticlesShapeComponent struct {
	Component
	particlesBatch *ParticlesShapeBatch
	UpdateParticle func(float32, *Particle)
}

func NewParticlesShapeComponent(_maxParticles int, _updateParticle func(float32, *Particle)) ParticlesShapeComponent {
	return ParticlesShapeComponent{
		particlesBatch: newParticlesShapeBatch(_maxParticles),
		UpdateParticle: _updateParticle,
	}
}

func (p *ParticlesShapeComponent) AddParticleWithVelo(shape Gfx_Shape, lifeTime float32, pos, velo Vector2f, color RGBA8, size, rotation float32) {
	p.particlesBatch.addParticle(shape, lifeTime, pos, velo, color, size, rotation)
}

func (p *ParticlesShapeComponent) AddParticles(numOfParticles int, shape Gfx_Shape, spread_pattern ParticlesSpreadPattern, life_time, speed float32, pos Vector2f, color RGBA8, size, angle float32) {
	for i := 0; i < numOfParticles; i++ {
		velo := calculateSpreadWithSpeed(i, numOfParticles, spread_pattern, speed, Vector2fRight.Rotate(angle, Vector2fZero))
		p.particlesBatch.addParticle(shape, life_time, pos, velo, color, size, 0.0)
	}
}

func calculateSpreadWithSpeed(index, numOfParticles int, spread_pattern ParticlesSpreadPattern, speed float32, direction Vector2f) Vector2f {
	switch spread_pattern {
	case PARTICLES_LINESPREAD:
		return direction.Scale(speed)
	case PARTICLES_CIRCLESPREAD:
		angle := (float32(index) / float32(numOfParticles)) * 2.0 * 180.0
		return direction.Rotate(angle, Vector2fZero)
	}

	return Vector2fUp.Scale(-speed)
}

func (t *ParticlesShapeComponent) ComponentSet(val interface{}) { *t = val.(ParticlesShapeComponent) }

type ParticlesShapeUpdateSystem struct {
	EcsSystemImpl
}

func (ps *ParticlesShapeUpdateSystem) Update(dt float32) {
	EachEntity(ParticlesShapeComponent{}, func(entity *EcsEntity, a interface{}) {
		particleBatch := a.(ParticlesShapeComponent)
		for i := 0; i < particleBatch.particlesBatch.maxParticles; i++ {
			particle := particleBatch.particlesBatch.particles[i]
			if particle.LifePercentage >= 0.0 {
				particle.Position = particle.Position.Add(particle.Velocity)
				particle.LifePercentage -= 1 / particle.LifeTime * dt

				particleBatch.UpdateParticle(dt, &particle)

				particleBatch.particlesBatch.particles[i] = particle
			}
		}

		// WriteComponent(ps.GetEcsEngine(), entity, particleBatch)
	})
}

type ParticlesShapeRenderSystem struct {
	EcsSystemImpl
}

func (ps *ParticlesShapeRenderSystem) Update(dt float32) {
	EachEntity(ParticlesShapeComponent{}, func(entity *EcsEntity, a interface{}) {
		particleBatch := a.(ParticlesShapeComponent)
		for i := 0; i < particleBatch.particlesBatch.maxParticles; i++ {
			particle := &particleBatch.particlesBatch.particles[i]
			if particle.LifePercentage > 0.0 {
				switch particle.Shape {
				case GFX_SHAPE_RECT:
					particleBatch.particlesBatch.shapes.DrawRectRotated(particle.Position, Vector2fOne.Scale(particle.Size), particle.Color, particle.Rotation)
					break
				case GFX_SHAPE_TRIANGLE:
					particleBatch.particlesBatch.shapes.DrawTriangleRotated(particle.Position, Vector2fOne.Scale(particle.Size), particle.Color, particle.Rotation)
					break
				case GFX_SHAPE_CIRCLE:
					particleBatch.particlesBatch.shapes.DrawCircle(particle.Position, particle.Size, particle.Color)
					break

				}
			}
		}
	})
}
