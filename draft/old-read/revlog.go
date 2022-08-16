package main

import (
  "context"
  "encoding/json"
  "fmt"
  "log"
  "os"
  "time"

  "crawshaw.io/sqlite/sqlitex"
)

type RevlogMeta struct {
  ID           time.Time     `json:"id"`
  CID          time.Time     `json:"cid"`
  Interval     time.Duration `json:"ivl"`
  LastInterval time.Duration `json:"last_ivl"`
}

type Revlog struct {
  ID                   int64      `json:"id"`
  CardID               int64      `json:"card_id"`
  UpdateSequenceNumber int64      `json:"usn"`
  Ease                 int64      `json:"ease"`
  Interval             int64      `json:"ivl"`
  LastInterval         int64      `json:"last_ivl"`
  Factor               int64      `json:"factor"`
  ReviewTime           int64      `json:"review_time"`
  Type                 int64      `json:"type"`
  Meta                 RevlogMeta `json:"meta"`
}

func NewInterval(interval int64) time.Duration {
  if interval < 0 {
    return time.Second * time.Duration(-1*interval)
  } else {
    return time.Hour * time.Duration(24*interval)
  }
}

// -- revlog is a review history; it has a row for every review you've ever done!
// CREATE TABLE revlog (
//     id              integer primary key,
//        -- epoch-milliseconds timestamp of when you did the review
//     cid             integer not null,
//        -- cards.id
//     usn             integer not null,
//         -- update sequence number: for finding diffs when syncing.
//         --   See the description in the cards table for more info
//     ease            integer not null,
//        -- which button you pushed to score your recall.
//        -- review:  1(wrong), 2(hard), 3(ok), 4(easy)
//        -- learn/relearn:   1(wrong), 2(ok), 3(easy)
//     ivl             integer not null,
//        -- interval (i.e. as in the card table)
//     lastIvl         integer not null,
//        -- last interval (i.e. the last value of ivl. Note that this value is not necessarily equal to the actual interval between this review and the preceding review)
//     factor          integer not null,
//       -- factor
//     time            integer not null,
//        -- how many milliseconds your review took, up to 60000 (60s)
//     type            integer not null
//        --  0=learn, 1=review, 2=relearn, 3=cram
// );

func main() {
  if err := run(); err != nil {
    fmt.Println(err)
  }
}

const file string = "/home/dt/.local/share/Anki2/Linux/collection.anki2"

func run() (err error) {
  var dbpool *sqlitex.Pool

  uri := fmt.Sprintf(
    "file:%s?cache=shared&mode=ro",
    file,
  )
  dbpool, err = sqlitex.Open(uri, 0, 10)
  if err != nil {
    log.Fatal(err)
  }

  conn := dbpool.Get(context.Background())
  if conn == nil {
    return fmt.Errorf("get conn error")
  }
  defer dbpool.Put(conn)

  stmt := conn.Prep("SELECT * FROM revlog WHERE cid = $card_id;")
  // stmt.SetText("$id", "_user_id_")
  stmt.SetInt64("$card_id", 1645697818510)

  for {
    if hasRow, err := stmt.Step(); err != nil {
      fmt.Println(err)
      // ... handle error
    } else if !hasRow {
      break
    }

    revlog := Revlog{
      ID:                   stmt.GetInt64("id"),
      CardID:               stmt.GetInt64("cid"),
      UpdateSequenceNumber: stmt.GetInt64("usn"),
      Ease:                 stmt.GetInt64("ease"),
      Interval:             stmt.GetInt64("ivn"),
      LastInterval:         stmt.GetInt64("lastIvl"),
      Factor:               stmt.GetInt64("factor"),
      ReviewTime:           stmt.GetInt64("time"),
      Type:                 stmt.GetInt64("type"),

      Meta: RevlogMeta{
        ID:           time.UnixMilli(stmt.GetInt64("id")),
        CID:          time.UnixMilli(stmt.GetInt64("cid")),
        Interval:     NewInterval(stmt.GetInt64("ivn")),
        LastInterval: NewInterval(stmt.GetInt64("lastIvn")),
      },
    }

    err = json.NewEncoder(os.Stdout).Encode(revlog)
    if err != nil {
      return
    }
  }

  err = stmt.Finalize()
  if err != nil {
    return
  }

  return
}
