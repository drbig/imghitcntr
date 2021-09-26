package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

type RecordRow struct {
	Key          string
	Count        int
	CreatedAt    time.Time
	LastUpdateAt time.Time
}

type LocalStore struct {
	mu   sync.Mutex
	data map[string]*RecordRow
}

var (
	db LocalStore
)

func init() {
	db.data = make(map[string]*RecordRow, DB_SIZE)
}

func getCount(key string) int {
	db.mu.Lock()
	defer db.mu.Unlock()

	now := time.Now()

	row, ok := db.data[key]
	if !ok {
		row = &RecordRow{
			Key:          key,
			Count:        1,
			CreatedAt:    now,
			LastUpdateAt: now,
		}
	} else {
		row.Count++
		row.LastUpdateAt = now
	}
	db.data[key] = row

	return row.Count
}

func (r *RecordRow) ToStrings() []string {
	return []string{
		r.Key,
		fmt.Sprintf("%d", r.Count),
		r.CreatedAt.Format(DATE_FMT),
		r.LastUpdateAt.Format(DATE_FMT),
	}
}

func (db LocalStore) CSVHeader() []string {
	return []string{
		"key",
		"count",
		"created_at",
		"last_update_at",
	}
}

func loadDB(path string) (err error) {
	if path == "" {
		return nil
	}
	logger.Infof("Loading state from: %s", path)

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	csv_r := csv.NewReader(f)
	got_header, err := csv_r.Read()
	if err != nil {
		return err
	}

	given_header := db.CSVHeader()
	if len(got_header) != len(given_header) {
		return errors.New("Header length mistmatch :-(")
	}

	for true {
		rec, err := csv_r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		key := rec[0]
		count, err := strconv.Atoi(rec[1])
		if err != nil {
			return err
		}
		created_at, err := time.Parse(DATE_FMT, rec[2])
		if err != nil {
			return err
		}
		last_update_at, err := time.Parse(DATE_FMT, rec[3])
		if err != nil {
			return err
		}

		db.data[key] = &RecordRow{
			Key:          key,
			Count:        count,
			CreatedAt:    created_at,
			LastUpdateAt: last_update_at,
		}
	}

	logger.Infof("Loaded %d entries", len(db.data))
	return nil
}

func saveDB(path string) (err error) {
	if path == "" {
		return nil
	}
	logger.Infof("Saving state to: %s", path)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	csv_w := csv.NewWriter(f)
	err = csv_w.Write(db.CSVHeader())
	if err != nil {
		return err
	}
	for _, r := range db.data {
		err = csv_w.Write(r.ToStrings())
		if err != nil {
			return err
		}
	}
	csv_w.Flush()
	logger.Infof("Saved %d entries!", len(db.data))
	return nil
}
