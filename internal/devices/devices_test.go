package devices_test

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dacjames/persistsql/internal/devices"
	"github.com/dacjames/persistsql/internal/resource"
	"github.com/dacjames/persistsql/internal/storage"
	"github.com/dacjames/persistsql/internal/tags"
	"github.com/dacjames/persistsql/internal/test_util"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func IthTags(i int) tags.Tagset {
	switch i % 5 {
	case 0:
		return tags.Tagset{}
	case 1:
		return tags.Must("animal=cat")
	case 2:
		return tags.Must("animal=cat", "animal=lion")
	case 3:
		return tags.Must(
			"animal=dog",
			"id=0oafk9ebwba8fAJB50h7",
			"group=01CGQH18GV3F6EVSEXHMY0RRMK",
			"user=wile",
		)
	case 4:
		return tags.Must(
			"animal=bat",
			"id=0oafk9ebwbaf5gfh43d",
			"group=engineering",
			"user=roadrunner",
			"env=dev",
			"region=ausw2",
			`prov=%5B%7B%22user%22%3A%20%22rabbit%22%2C%20%22ts%22%3A%201529801310%7D%2C%20%7B%22user%22%3A%20%22bugs%22%2C%20%22ts%22%3A%201529801627%7D%5D`,
		)
	default:
		return tags.Tagset{}
	}
}

func TestDeviceStorage(t *testing.T) {
	require.Equal(t,
		true, true)

	example := &devices.Device{
		Resource: resource.Resource{
			ID: resource.NewID(nil),
			Meta: resource.Meta{
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				DeletedAt: time.Time{},
			},
			Deleted: false,
		},
		Name: "Dummy 🤖",
		Tags: IthTags(1),
	}
	require.NotNil(t, example)

	test_util.WithMigratedDB(t, func(db *sql.DB) {
		var storage interface {
			storage.Getter
			storage.Putter
		} = &storage.Collection{
			ServiceImpl: &devices.DeviceStorage{
				DB: db,
			},
			DB: db,
		}

		err := storage.PutAny(example)
		require.NoError(t, err)

		example.Name = "Happy 🤖"
		err = storage.PutAny(example)
		require.NoError(t, err)

		state, err := storage.GetAny(example.ID)
		require.NoError(t, err)

		device, ok := state.(*devices.Device)
		require.True(t, ok)

		require.Equal(t, example.ID, device.ID)
		require.Equal(t, example.Name, device.Name)

	})

}

func TestShouldBeBenchmark(t *testing.T) {

	N := 10000

	test_util.WithMigratedDB(t, func(db *sql.DB) {
		var storage interface {
			storage.Getter
			storage.Putter
		} = &storage.Collection{
			ServiceImpl: &devices.DeviceStorage{
				DB: db,
			},
			DB: db,
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		ids := []resource.ID{}

		numRevs := 0
		startTime := time.Now()

		for i := 0; i < N; i++ {
			id := resource.NewID(r)
			ids = append(ids, id)
			err := storage.PutAny(&devices.Device{
				Resource: resource.Resource{ID: id},
				Name:     fmt.Sprintf("Dummy %d", i),
			})
			numRevs++
			if err != nil {
				t.Fatal(err)
			}
		}

		afterPuts := time.Now()
		putTotal := afterPuts.Sub(startTime) / time.Millisecond
		perPut := float64(putTotal) / float64(len(ids))
		fmt.Printf("%fms per put after %d revisions over %d resources\n", perPut, numRevs, len(ids))

		rand.Shuffle(len(ids), func(i, j int) {
			ids[i], ids[j] = ids[j], ids[i]
		})

		for i := 0; i < N/10; i++ {
			id := ids[i]
			err := storage.PutAny(&devices.Device{
				Resource: resource.Resource{ID: id},
				Name:     fmt.Sprintf("Modified %d", i),
			})
			numRevs++
			if err != nil {
				t.Fatal(err)
			}
		}

		afterUpdates := time.Now()
		updateTotal := afterUpdates.Sub(afterPuts) / time.Millisecond
		perUpdate := float64(updateTotal) / float64(len(ids))
		fmt.Printf("%fms per update after %d revisions over %d resources\n", perUpdate, numRevs, len(ids))

		numHots := 0
		for i := 0; i < N/100; i++ {
			for j := 0; j < max(10, N/100); j++ {
				id := ids[i]
				err := storage.PutAny(&devices.Device{
					Resource: resource.Resource{ID: id},
					Name:     fmt.Sprintf("Hot %d'th %d", j, i),
				})
				numRevs++
				numHots++
				if err != nil {
					t.Fatal(err)
				}
			}
		}

		afterHots := time.Now()
		hotTotal := afterHots.Sub(afterUpdates) / time.Millisecond
		perHot := float64(hotTotal) / float64(numHots)
		fmt.Printf("%fms per hot update after %d revisions over %d resources\n", perHot, numRevs, len(ids))

		before := time.Now()
		for i := 0; i < len(ids); i++ {
			id := ids[i]
			_, err := storage.GetAny(id)
			if err != nil {
				t.Fatal(err)
			}
		}
		after := time.Now()

		getTotal := after.Sub(before) / time.Millisecond
		perGet := float64(getTotal) / float64(len(ids))
		fmt.Printf("%fms per get after %d revisions over %d resources\n", perGet, numRevs, len(ids))

		require.NotNil(t, nil)
	})

}
