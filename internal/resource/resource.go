package resource

import (
	"database/sql/driver"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"github.com/pkg/errors"
)

type IDer interface {
	String() string
	UUID() string
}

type ID ulid.ULID

func (r ID) String() string {
	return ulid.ULID(r).String()
}

func (r ID) UUID() string {
	ru := ulid.ULID(r)

	buf := make([]byte, 16)

	// Safely ignore the error
	// since it can only occur if
	// the target buffer is the wrong size
	_ = ru.MarshalBinaryTo(buf)

	// Safely ignore the error for the same reason
	uu, _ := uuid.FromBytes(buf)

	return uu.String()
}

func (r ID) Value() (driver.Value, error) {
	return r.UUID(), nil
}

func (r *ID) Scan(value interface{}) error {
	v, err := driver.String.ConvertValue(value)

	if err == nil {
		if vb, ok := v.([]byte); ok {
			vs := string(vb)
			rid, err := ParseUUID(vs)
			if err != nil {
				return errors.Wrapf(err, "Failed parsing value %s", vs)
			}
			*r = rid
			return nil
		}
	}
	return errors.Wrapf(err, "Failed scanning value %v", value)
}

// ResourceID returns a new ULID resource ID
// Takes a rand.Rand as an argument to avoid
// lock contention issues on a global rand source.
// If r is nil, uses a bespoke RandSource seeded
// with the current time.
func NewID(r *rand.Rand) ID {
	now := time.Now().UTC()

	if r == nil {
		r = rand.New(rand.NewSource(now.UnixNano()))
	}

	return ID(ulid.MustNew(ulid.Timestamp(now), r))
}

func PtrID(r *rand.Rand) *ID {
	id := NewID(r)
	return &id
}

func ParseID(rs string) (ID, error) {
	u, err := ulid.Parse(rs)
	if err != nil {
		return ID(ulid.ULID{}), err
	}
	return ID(u), nil
}

func ParseUUID(rs string) (ID, error) {
	uu, err := uuid.Parse(rs)
	if err != nil {
		return ID(ulid.ULID{}), err
	}

	buf := [16]byte{}
	copy(buf[:], uu[:16])

	ul := ulid.ULID(buf)

	return ID(ul), nil
}
