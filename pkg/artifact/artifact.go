package artifact

import (
	"net/url"

	"github.com/benchkram/bobc/pkg/db/model"
	"github.com/benchkram/bobc/pkg/optional"
	"github.com/benchkram/bobc/restserver/generated"
)

type A struct {
	// ID is the hash of the artifact itself as it was computed
	// when being uploaded (see: content-addressable storage)
	ID string

	// AccessLink is a URL that allows access (download) of the
	// artifact's payload when requested and for a limited time span
	AccessLink *url.URL

	// Size of the artifact in bytes
	Size int
}

func FromDatabaseType(m *model.Artifact) *A {
	return &A{
		ID:   m.ID,
		Size: m.Size,
	}
}

func (a *A) ToDatabaseType(projectID string) *model.Artifact {
	return &model.Artifact{
		ID:        a.ID,
		ProjectID: projectID,
		Size:      a.Size,
	}
}

func (a *A) ToRestType() generated.Artifact {
	var link *string
	if a.AccessLink != nil {
		link = optional.String(a.AccessLink.String())
	}
	return generated.Artifact{
		Id:       a.ID,
		Location: link,
		Size:     a.Size,
	}
}
