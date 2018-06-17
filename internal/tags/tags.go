package tags

import (
	"database/sql"
	"strings"

	"github.com/dacjames/persistsql/internal/resource"
	"github.com/dacjames/persistsql/internal/util"
	"github.com/pkg/errors"
)

type Tag struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

type Tagset []Tag

func (tt Tagset) Insert(tx *sql.Tx, id resource.ID) error {
	if len(tt) == 0 {
		// If a resource has no tags, we associate it with a special empty pseudotag
		// This ensures the query finds the resource when querying with no tags
		// This construction is a little hacky, but actually works out much nicer
		// than adjusting the query to account for untagged resources
		return Tagset([]Tag{Tag{"", ""}}).Insert(tx, id)
	}

	r := util.NewRandSource()

	tv := []string{}
	tf := []interface{}{}
	ph := util.NewPlaceholders()
	for _, t := range tt {
		// This is safe to SQL Injection because
		//   - The ID is controlled and doesn't depend on user input
		//   - user controls the tags but is only dependent on the NUMBER,
		//     not the CONTENT of the tags.
		tv = append(tv, ph.NextValue(3))
		tuu := resource.PtrID(r)

		tf = append(tf, tuu, t.Key, t.Value)
	}

	// Uses a dummy ON CONFLICT UPDATE becuase using
	// ON CONFLICT DO NOTHING causes no ID to be returned

	query := `
		INSERT INTO ledger.tags(tag_id, key, value)
		VALUES ` + strings.Join(tv, ",") + `
		ON CONFLICT (key, value) DO UPDATE SET "key" = EXCLUDED."key"
		RETURNING tag_id
	`

	rows, err := tx.Query(query, tf...)
	if err != nil {
		return errors.Wrap(err, "Inserting into tags")
	}

	tagIDs := []resource.ID{}
	for rows.Next() {
		var tagID resource.ID
		err = rows.Scan(&tagID)
		if err != nil {
			return err
		}
		tagIDs = append(tagIDs, tagID)
	}

	rows.Close()

	rtv := []string{}
	rtf := []interface{}{}
	ph = util.NewPlaceholders()
	for _, tid := range tagIDs {
		rtv = append(rtv, ph.NextValue(2))
		rtf = append(rtf, id, tid)
	}

	query = `INSERT INTO ledger.resource_tags(resource_id, tag_id) VALUES ` + strings.Join(rtv, ",")

	_, err = tx.Exec(query, rtf...)
	if err != nil {
		return errors.Wrap(err, "Inserting into resource_tags")
	}

	return nil
}

func Must(tt ...string) Tagset {
	tags := make([]Tag, len(tt))

	for i, t := range tt {
		parts := strings.Split(t, "=")
		tags[i] = Tag{parts[0], parts[1]}
	}

	return Tagset(tags)
}
