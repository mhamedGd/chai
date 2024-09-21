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
	particles        List[Particle]
	maxParticles     int
	lastFreeParticle int
}

func newParticlesShapeBatch(_max_particles int) *ParticlesShapeBatch {
	return &ParticlesShapeBatch{
		shapes:           &Shapes,
		particles:        NewListSized[Particle](_max_particles),
		maxParticles:     _max_particles,
		lastFreeParticle: 0,
	}
}

func (p *ParticlesShapeBatch) addParticle(shape Gfx_Shape, lifeTime float32, pos, velo Vector2f, color RGBA8, size, rotation float32) {
	p.particles.Data[p.lastFreeParticle] = Particle{
		Shape:          shape,
		Position:       pos,
		Velocity:       velo,
		Size:           size,
		Rotation:       rotation,
		LifeTime:       lifeTime,
		LifePercentage: 1.0,
		Color:          color,
	}
	p.lastFreeParticle = (p.lastFreeParticle + 1) % p.particles.Count()
}

func (p *ParticlesShapeBatch) findLastFreeParticle() int {
	for i := p.lastFreeParticle; i < p.maxParticles; i++ {
		if p.particles.Data[i].LifePercentage <= 0.0 {
			return i
		}
	}

	for i := 0; i < p.lastFreeParticle; i++ {
		if p.particles.Data[i].LifePercentage <= 0.0 {
			return i
		}
	}

	return 0
}

type ParticlesShapeComponent struct {
	particlesBatch *ParticlesShapeBatch
	z              float32
	UpdateParticle func(float32, *Particle)
}

func NewParticlesShapeComponent(_maxParticles int, _z float32, _updateParticle func(float32, *Particle)) ParticlesShapeComponent {
	return ParticlesShapeComponent{
		particlesBatch: newParticlesShapeBatch(_maxParticles),
		z:              _z,
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

func ParticlesShapeUpdateSystem(_this_scene *Scene, _dt float32) {
	Iterate1[ParticlesShapeComponent](func(i EntId, psc *ParticlesShapeComponent) {
		for i := psc.particlesBatch.maxParticles - 1; i >= 0; i-- {
			particle := &psc.particlesBatch.particles.Data[i]
			if particle.LifePercentage > 0.0 {
				particle.Position = particle.Position.Add(particle.Velocity)
				particle.LifePercentage -= (1 / particle.LifeTime) * _dt

				psc.UpdateParticle(_dt, particle)
				if particle.LifePercentage <= 0.0 {
					psc.particlesBatch.lastFreeParticle = i
				}
			}
		}
	})
}

func ParticlesShapeRenderSystem(_this_scene *Scene, _dt float32) {
	Iterate1[ParticlesShapeComponent](func(i EntId, psc *ParticlesShapeComponent) {
		for i := psc.particlesBatch.maxParticles - 1; i >= 0; i-- {
			particle := &psc.particlesBatch.particles.Data[i]
			if particle.LifePercentage > 0.0 {
				switch particle.Shape {
				case GFX_PARTICLES_SHAPE_RECT:
					psc.particlesBatch.shapes.DrawRectRotated(particle.Position, psc.z, Vector2fOne.Scale(particle.Size), particle.Color, particle.Rotation)
				case GFX_PARTICLES_SHAPE_TRIANGLE:
					psc.particlesBatch.shapes.DrawTriangleRotated(particle.Position, psc.z, Vector2fOne.Scale(particle.Size), particle.Color, particle.Rotation)
				case GFX_PARTICLES_SHAPE_CIRCLE:
					psc.particlesBatch.shapes.DrawCircle(particle.Position, psc.z, particle.Size, particle.Color)
				case GFX_PARTICLES_SHAPE_FILLRECT:
					psc.particlesBatch.shapes.DrawFillRectRotated(particle.Position, psc.z, Vector2fOne.Scale(particle.Size), particle.Color, particle.Rotation)
				}

			}
		}
	})
}
