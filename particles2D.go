package chai

import (
	"github.com/mhamedGd/chai/customtypes"
	. "github.com/mhamedGd/chai/math"
)

type ParticlesSpreadPattern = int

const (
	PARTICLES_LINESPREAD   = 0x000000e1
	PARTICLES_CIRCLESPREAD = 0x000000e2
)

type ParticleShape = int

const (
	// PARTICLES_SHAPE_CIRCLE   = ParticleShape(0x000000f1)
	PARTICLES_SHAPE_RECT = ParticleShape(0x000000f2)
	// PARTICLES_SHAPE_TRIANGLE = ParticleShape(0x000000f3)
	PARTICLES_SHAPE_QUAD = ParticleShape(0x000000f4)
)

type Particle struct {
	LifePercentage float32
	LifeTime       float32
	Position       Vector2f
	Velocity       Vector2f
	Size           float32
	Rotation       float32
	Color          RGBA8
	Shape          ParticleShape
}

type ParticlesShapeBatch struct {
	m_Renderer         *Renderer2D
	m_Particles        customtypes.List[Particle]
	m_MaxParticles     int
	m_LastFreeParticle int
}

func newParticlesShapeBatch(_maxParticles int) *ParticlesShapeBatch {
	return &ParticlesShapeBatch{
		m_Renderer:         &Renderer,
		m_Particles:        customtypes.NewListSized[Particle](_maxParticles),
		m_MaxParticles:     _maxParticles,
		m_LastFreeParticle: 0,
	}
}

func (p *ParticlesShapeBatch) addParticle(_shape ParticleShape, _lifeTime float32, _pos, _velocity Vector2f, _color RGBA8, _size, _rotation float32) {
	p.m_Particles.Data[p.m_LastFreeParticle] = Particle{
		Shape:          _shape,
		Position:       _pos,
		Velocity:       _velocity,
		Size:           _size,
		Rotation:       _rotation,
		LifeTime:       _lifeTime,
		LifePercentage: 1.0,
		Color:          _color,
	}
	p.m_LastFreeParticle = (p.m_LastFreeParticle + 1) % p.m_Particles.Count()
}

func (p *ParticlesShapeBatch) findLastFreeParticle() int {
	for i := p.m_LastFreeParticle; i < p.m_MaxParticles; i++ {
		if p.m_Particles.Data[i].LifePercentage <= 0.0 {
			return i
		}
	}

	for i := 0; i < p.m_LastFreeParticle; i++ {
		if p.m_Particles.Data[i].LifePercentage <= 0.0 {
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

func (p *ParticlesShapeComponent) AddParticleWithVelo(_shape ParticleShape, _lifeTime float32, _pos, _velocity Vector2f, _color RGBA8, _size, _rotation float32) {
	p.particlesBatch.addParticle(_shape, _lifeTime, _pos, _velocity, _color, _size, _rotation)
}

func (p *ParticlesShapeComponent) AddParticles(_numOfParticles int, _shape ParticleShape, _spreadPattern ParticlesSpreadPattern, _lifeTime, _speed float32, _pos Vector2f, _color RGBA8, _size, _rotation float32) {
	for i := 0; i < _numOfParticles; i++ {
		velo := calculateSpreadWithSpeed(i, _numOfParticles, _spreadPattern, _speed, Vector2fRight.Rotate(_rotation, Vector2fZero))
		p.particlesBatch.addParticle(_shape, _lifeTime, _pos, velo, _color, _size, 0.0)
	}
}

func calculateSpreadWithSpeed(_index, _numOfParticles int, _spreadPattern ParticlesSpreadPattern, _speed float32, _direction Vector2f) Vector2f {
	switch _spreadPattern {
	case PARTICLES_LINESPREAD:
		return _direction.Scale(_speed)
	case PARTICLES_CIRCLESPREAD:
		angle := (float32(_index) / float32(_numOfParticles)) * 2.0 * 180.0
		return _direction.Rotate(angle, Vector2fZero).Scale(_speed)
	}

	return Vector2fUp.Scale(-_speed)
}

func ParticlesShapeUpdateSystem(_thiScene *Scene, _dt float32) {
	Iterate1[ParticlesShapeComponent](func(i EntId, psc *ParticlesShapeComponent) {
		for i := psc.particlesBatch.m_MaxParticles - 1; i >= 0; i-- {
			particle := &psc.particlesBatch.m_Particles.Data[i]
			if particle.LifePercentage > 0.0 {
				particle.Position = particle.Position.Add(particle.Velocity)
				particle.LifePercentage -= (1 / particle.LifeTime) * _dt

				psc.UpdateParticle(_dt, particle)
				if particle.LifePercentage <= 0.0 {
					psc.particlesBatch.m_LastFreeParticle = i
				}
			}
		}
	})
}

func ParticlesShapeRenderSystem(_thisScene *Scene, _dt float32) {
	Iterate1[ParticlesShapeComponent](func(i EntId, psc *ParticlesShapeComponent) {
		for i := psc.particlesBatch.m_MaxParticles - 1; i >= 0; i-- {
			particle := &psc.particlesBatch.m_Particles.Data[i]
			if particle.LifePercentage > 0.0 {
				switch particle.Shape {
				case PARTICLES_SHAPE_RECT:
					DrawRectWRenderer(psc.particlesBatch.m_Renderer, particle.Position, Vector2fOne.Scale(particle.Size), particle.Color, psc.z, particle.Rotation)
				// case PARTICLES_SHAPE_TRIANGLE:
				// 	psc.particlesBatch.m_Shapes.DrawTriangleRotated(particle.Position, psc.z, Vector2fOne.Scale(particle.Size), particle.Color, particle.Rotation)
				// case PARTICLES_SHAPE_CIRCLE:
				// psc.particlesBatch.m_Shapes.DrawCircle(particle.Position, psc.z, particle.Size, particle.Color)
				case PARTICLES_SHAPE_QUAD:
					psc.particlesBatch.m_Renderer.InsertQuadRotated(particle.Position, Vector2fOne.Scale(particle.Size), psc.z, particle.Color, particle.Rotation)
				}

			}
		}
	})
}
