package example

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/exaream/go-db/dbutil"
	"github.com/mattn/go-gimei"
	"go.uber.org/multierr"
)

// Setup generates initial data.
func Setup(ctx context.Context, cfg *dbutil.ConfigFile, min, max, chunkSize uint) (total int64, err error) {
	db, err := dbutil.NewDBContext(ctx, cfg)
	if err != nil {
		return 0, err
	}

	defer func() {
		if rerr := db.Close(); rerr != nil {
			err = rerr
		}
	}()

	tx := db.MustBeginTx(ctx, nil)
	queryTruncateTbl := queryTruncateTbls[db.DriverName()]

	if _, err := tx.ExecContext(ctx, queryTruncateTbl); err != nil {
		return 0, multierr.Append(err, tx.Rollback())
	}

	total, err = dbutil.BulkInsertTxContext(ctx, tx, fakeUsers, queryInsert, min, max, chunkSize)
	if err != nil {
		return 0, multierr.Append(err, tx.Rollback())
	}

	if err := tx.Commit(); err != nil {
		return 0, multierr.Append(err, tx.Rollback())
	}

	return total, nil
}

// fakeUsers returns fake user list.
func fakeUsers(min, max uint) (users []User) {
	if min == 0 || max == 0 {
		return users
	}

	now := time.Now()
	for i := min; i <= max; i++ {
		users = append(users, User{i, gimei.NewName().Kanji(), faker.Email(), 0, &now, &now})
	}

	return users
}
