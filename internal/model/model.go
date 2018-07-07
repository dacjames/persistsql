package model

import (
	"github.com/dacjames/persistsql/internal/resource"
	"github.com/dacjames/persistsql/internal/tags"
)

type ResourceIDer interface {
	ResourceID() resource.ID
}

type ResourceTagser interface {
	ResourceTags() tags.Tagset
}

type ResourceServicer interface {
	ResourceService() string
}

type Resourcer interface {
	ResourceIDer
	ResourceTagser
	ResourceServicer
}
