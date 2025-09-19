package db

import (
	"fmt"
)

type AdminRepository struct {
	SqliteHandler *SQLiteHandler
}

func (adR *AdminRepository) CreateBackup(file string) error {
	vQ := fmt.Sprintf("VACUUM main into '%v';", file)
	if _, vErr := adR.SqliteHandler.Db().Exec(vQ); vErr != nil {
		return vErr
	}
	return nil
}
