package dbgorm

import (
	"context"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDbPg(ctx context.Context, pgDsn string, timeout time.Duration) (db *gorm.DB, err error) {
	dsn := pgDsn
	//ctx, cancel := context.WithTimeout(ctx, timeout)
	//defer cancel()
	//connect to with timeout
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// go func(gctx context.Context) (*gorm.DB, error) {
	// 	gdb, gerr := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// 	for gerr != nil {
	// 		select {
	// 		case <-gctx.Done():
	// 			return gdb, gerr
	// 		default:
	// 		}
	// 		gdb, gerr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// 	}
	// 	return gdb, gerr
	// }(ctx)
	return db, err
}
