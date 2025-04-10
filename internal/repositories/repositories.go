package repositories

type RepositoryType string

var (
	RepositoryTypePostgres RepositoryType = "pg"
)

type BaseRepository interface{}
