package notes

import (
	"database/sql"
)

const file string = "/home/dt/.local/share/Anki2/Linux/collection.anki2"

type Note struct {
	Id     int    `json:"id"`
	Guid   string `json:"guid"`
	Mid    int    `json:"mid"`
	Mod    int    `json:"mod"`
	Usn    int    `json:"usn"`
	Tags   string `json:"tags"`
	Fields string `json:"fields"`
	Sfld   string `json:"sfld"`
	Csum   int    `json:"csum"`
	Flags  int    `json:"flags"`
	Data   string `json:"data"`
}

func QueryAllNotes() (result []Note, err error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return
	}

	rows, err := db.Query("select id, guid, mid, mod, usn, tags, flds, sfld, csum, flags, data from notes limit 100")
	if err != nil {
		return
	}

	if rows.Err() != nil {
		err = rows.Err()
		return
	}

	for rows.Next() {
		result = append(result, Note{})
		n := &result[len(result)-1]
		err = rows.Scan(
			&n.Id,
			&n.Guid,
			&n.Mid,
			&n.Mod,
			&n.Usn,
			&n.Tags,
			&n.Fields,
			&n.Sfld,
			&n.Csum,
			&n.Flags,
			&n.Data,
		)

		if err != nil {
			return
		}
	}

	return
}
