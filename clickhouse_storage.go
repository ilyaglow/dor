package dor

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	// clickhouse sql driver
	_ "github.com/kshvakov/clickhouse"
)

const (
	createStmt = `
		CREATE IF NOT EXISTS %s (
			rank	UInt32,
			domain	String,
			rawdata	String,
			source	String,
			date	Date DEFAULT today()
		) engine=MergeTree(date, (domain, source), 8192)
		PARTITION BY toYYYYMM(date)
	`
	insertStmt = `
		INSERT INTO %s (rank, domain, rawdata, source) VALUES (?, ?, ?, ?)
	`
)

// ClickhouseStorage is a dor.Storage that uses Clickhouse database.
type ClickhouseStorage struct {
	db        *sql.DB
	table     string
	batchSize int
}

// NewClickhouseStorage bootstraps ClickhouseStorage.
func NewClickhouseStorage(location, table string, batch int) (*ClickhouseStorage, error) {
	db, err := prepareDB(location, table)
	if err != nil {
		return nil, err
	}

	return &ClickhouseStorage{
		db:        db,
		table:     table,
		batchSize: batch,
	}, nil
}

func prepareDB(location, table string) (*sql.DB, error) {
	conn, err := sql.Open("clickhouse", location)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	_, err = conn.Exec(fmt.Sprintf(createStmt, table))
	if err != nil {
		return conn, err
	}

	return conn, nil
}

// Put implements Storage interface method Put
//	s - is the data source
//	t - is the data datetime
func (c *ClickhouseStorage) Put(entries <-chan *Entry, s string, t time.Time) error {
	return c.send(entries, s)
}

func (c *ClickhouseStorage) send(entries <-chan *Entry, s string) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(fmt.Sprintf(insertStmt, c.table))
	if err != nil {
		return err
	}

	tm := time.Now()
	it := 0

	for e := range entries {
		_, err := stmt.Exec(
			e.Rank,
			e.Domain,
			e.RawData,
			e.Source,
		)
		if err != nil {
			return err
		}

		it++

		if it == c.batchSize || time.Since(tm) > 1*time.Minute {
			err := tx.Commit()
			if err != nil {
				return err
			}

			tx, err = c.db.Begin()
			if err != nil {
				return nil
			}

			stmt, err = tx.Prepare(fmt.Sprintf(insertStmt, c.table))
			if err != nil {
				return nil
			}

			it = 0
			tm = time.Now()
		}
	}

	return tx.Commit()
}

// Get ranks for specified domain and sources.
func (c *ClickhouseStorage) Get(d string, sources ...string) ([]*Entry, error) {
	var (
		queries []string
		sface   []interface{}
	)
	selectQuery := fmt.Sprintf(`
		SELECT rank, domain, rawdata, source, date
		FROM %s
		WHERE source=?
		ORDER BY date DESC
		LIMIT 1
		`, c.table)

	for i, s := range sources {
		queries = append(queries, selectQuery)
		sface[i] = s
	}

	// if sources are not specified include all of them
	if len(sources) == 0 {
		for i, s := range ingesters {
			queries = append(queries, selectQuery)
			sface[i] = s.GetDesc()
		}
	}

	rows, err := c.db.Query(strings.Join(queries, " UNION ALL "), sface...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		rank                    uint
		domain, rawdata, source string
		date                    string
		entries                 []*Entry
	)
	for rows.Next() {
		err := rows.Scan(&rank, &domain, &rawdata, &source, &date)
		if err != nil {
			return entries, err
		}

		d, err := time.Parse("2006-01-02", date)
		if err != nil {
			return entries, err
		}

		entries = append(entries, &Entry{
			Rank:    rank,
			Domain:  domain,
			RawData: rawdata,
			Source:  source,
			Date:    d,
		})
	}
	err = rows.Err()
	if err != nil {
		return entries, err
	}

	return entries, nil
}

// GetMore returns lps entries for each source for a specified domain.
func (c *ClickhouseStorage) GetMore(d string, lps int, sources ...string) ([]*Entry, error) {
	return nil, nil
}
