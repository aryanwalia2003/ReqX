package history

import "database/sql"

// DB wraps the SQLite connection used for persisting test run history.
type DB struct {
	conn *sql.DB
}
