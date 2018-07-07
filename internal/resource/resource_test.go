package resource_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"github.com/stretchr/testify/require"

	"github.com/dacjames/persistsql/internal/resource"
)

func TestResourceID(t *testing.T) {
	require.Equal(t, true, true)

	r := resource.NewID(nil)

	require.NotEmpty(t, r.String())

	_, err := ulid.Parse(r.String())
	require.Nil(t, err)

	_, err = uuid.Parse(r.UUID())
	require.Nil(t, err)
}
