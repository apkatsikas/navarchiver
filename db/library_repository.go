package db

type Library struct {
	Id   int
	Path string
}

type LibraryRepository struct {
	SqliteHandler *SQLiteHandler
}

func (lr *LibraryRepository) LibraryById(id int) (*Library, error) {
	statement, err := lr.SqliteHandler.Db().Prepare(
		"SELECT id, path FROM library WHERE id = ? LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	row := statement.QueryRow(id)

	var library Library
	if err := row.Scan(&library.Id, &library.Path); err != nil {
		return nil, err
	}
	return &library, nil
}
