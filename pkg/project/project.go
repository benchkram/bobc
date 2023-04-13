package project

import (
	"time"

	"github.com/benchkram/bobc/pkg/artifact"
	"github.com/benchkram/bobc/pkg/db/model"
	"github.com/benchkram/bobc/restserver/generated"
	"github.com/google/uuid"
)

type P struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	Name        string
	Description string
	Artifacts   []*artifact.A
}

func New(name, description string) *P {
	project := &P{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		Name:        name,
		Description: description,
	}

	return project
}

func FromDBModel(m *model.Project) (*P, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, err
	}

	artifacts := []*artifact.A{}

	for _, a := range m.Artifacts {
		artifact := artifact.FromDatabaseType(a)
		artifacts = append(artifacts, artifact)
	}

	return &P{
		ID:          id,
		CreatedAt:   m.CreatedAt,
		Name:        m.Name,
		Description: m.Description,
		Artifacts:   artifacts,
	}, nil
}

func (p *P) ToExtendedProjectRestType() generated.ExtendedProject {
	hashlist := []generated.Artifact{}
	for _, h := range p.Artifacts {
		hashlist = append(hashlist, h.ToRestType())
	}

	return generated.ExtendedProject{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Hashes:      &hashlist,
	}
}

func (p *P) ToProjectRestType() generated.Project {
	return generated.Project{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
	}
}

func (p *P) ToProjectDatabaseType() *model.Project {
	artifacts := []*model.Artifact{}
	for _, a := range p.Artifacts {
		m := a.ToDatabaseType(p.ID.String())
		artifacts = append(artifacts, m)
	}

	return &model.Project{
		ID:          p.ID.String(),
		CreatedAt:   p.CreatedAt,
		Name:        p.Name,
		Description: p.Description,
		Artifacts:   artifacts,
	}
}
