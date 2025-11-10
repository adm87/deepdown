package ecs

import (
	"github.com/adm87/deepdown/scripts/components"
	"github.com/adm87/deepdown/scripts/ecs/entity"
	"github.com/adm87/utilities/hash"
	"github.com/hajimehoshi/ebiten/v2"
)

// =========== World ===========

type UpdateSystem func(w *World, dt float32)
type RenderSystem func(w *World, screen *ebiten.Image)

type World struct {
	entities hash.Set[entity.Entity]

	// Entity management
	freeEntities []entity.Entity
	nextEntity   entity.Entity

	// Systems
	UpdateSystems      []UpdateSystem
	FixedUpdateSystems []UpdateSystem
	LateUpdateSystems  []UpdateSystem
	RenderSystems      []RenderSystem
}

func NewWorld() *World {
	return &World{
		entities: hash.NewSet[entity.Entity](),

		freeEntities: make([]entity.Entity, 0, 1024),
		nextEntity:   entity.Null,
	}
}

func (w *World) NewEntity() entity.Entity {
	var id entity.Entity
	if len(w.freeEntities) > 0 {
		id = w.freeEntities[len(w.freeEntities)-1]
		w.freeEntities = w.freeEntities[:len(w.freeEntities)-1]
	} else {
		w.nextEntity++
		id = w.nextEntity
	}
	w.entities.Add(id)
	return id
}

func (w *World) RemoveEntity(e entity.Entity) {
	if !w.entities.Contains(e) {
		return
	}

	components.DestroyEntity(e)

	w.entities.Remove(e)
	w.freeEntities = append(w.freeEntities, e)
}

func (w *World) AddUpdateSystem(system UpdateSystem) {
	w.UpdateSystems = append(w.UpdateSystems, system)
}

func (w *World) AddFixedUpdateSystem(system UpdateSystem) {
	w.FixedUpdateSystems = append(w.FixedUpdateSystems, system)
}

func (w *World) AddLateUpdateSystem(system UpdateSystem) {
	w.LateUpdateSystems = append(w.LateUpdateSystems, system)
}

func (w *World) AddRenderSystem(system RenderSystem) {
	w.RenderSystems = append(w.RenderSystems, system)
}

func (w *World) Update(dt float32) {
	for _, system := range w.UpdateSystems {
		system(w, dt)
	}
}

func (w *World) FixedUpdate(dt float32) {
	for _, system := range w.FixedUpdateSystems {
		system(w, dt)
	}
}

func (w *World) LateUpdate(dt float32) {
	for _, system := range w.LateUpdateSystems {
		system(w, dt)
	}
}

func (w *World) Render(screen *ebiten.Image) {
	for _, system := range w.RenderSystems {
		system(w, screen)
	}
}
