package storage

import (
	"database/sql"
	"errors"
	"log"
	"simple-admin/common"
	"strconv"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
	Stor   *sql.DB
	Logger *log.Logger
}

func NewStorage(storageType string, storageDsn string, logger *log.Logger) (stor *Storage, err error) {
	if storageType == `postgres` {
		stor = &Storage{}
		stor.Stor, err = sql.Open(storageType, storageDsn)
		if err != nil {
			return nil, err
		}
		if err = stor.Stor.Ping(); err != nil {
			return nil, err
		}
		if err = stor.checkDbStruct(); err != nil {
			stor.Stor.Close()
			return nil, err
		}
		stor.Logger = logger
		return stor, err
	} else {
		return nil, errors.New(`storage: Not implemented yet!`)
	}
}

func (stor *Storage) checkDbStruct() error {
	_, err := stor.Stor.Exec("CREATE table IF NOT EXISTS rates (id SERIAL PRIMARY KEY, currency_from varchar(50) NOT NULL, currency_to varchar(50) NOT NULL, rate decimal(15,2), time bigint)")
	return err
}

func (stor *Storage) Stop() error {
	return stor.Stor.Close()
}

func (stor *Storage) Save(recordsSlice interface{}) error {
	tx, err := stor.Stor.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = rollbackErr
			}
			return
		}
	}()
	stmt, err := tx.Prepare(pq.CopyIn("rates", `currency_from`, `currency_to`, `rate`, `time`))
	if err != nil {
		panic(err)
	}
	records, ok := recordsSlice.([]common.Record)
	if !ok {
		panic(errors.New(`storage: Could not assert argument to slice!`))
	}
	for _, rec := range records {
		_, err = stmt.Exec(rec.CurrencyFrom, rec.CurrencyTo, rec.Rate, rec.Time)
		if err != nil {
			panic(err)
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}
	err = stmt.Close()
	if err != nil {
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	return nil
}

func (stor *Storage) Fetch(interval time.Duration) ([]interface{}, error) {
	t := strconv.FormatInt(time.Now().Add(-interval).Unix(), 10)
	stor.Logger.Printf("time edge:\n%#v\n", t)
	rows, err := stor.Stor.Query(`SELECT currency_from, currency_to, avg(rate) AS rate_avg, max(time) AS time FROM rates WHERE time >= $1 GROUP BY currency_from, currency_to ORDER BY currency_from`, t)
	if err != nil {
		return nil, err
	}
	records := make([]interface{}, 0)
	for {
		if n := rows.Next(); !n {
			break
		}
		record := common.Record{}
		err = rows.Scan(&record.CurrencyFrom, &record.CurrencyTo, &record.Rate, &record.Time)
		if err != nil {
			return nil, err
		}
		records = append(records, interface{}(record))
	}
	return records, nil
}
