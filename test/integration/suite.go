package integration

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/config"
	"CRUD-HOME-APPLIANCE-STORE/internal/consul"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	cfg    *config.Config
	db     *pgx.Conn
	logger *logger.Logger
}

func (s *TestSuite) SetupTest() {
	s.cfg = config.MustLoad()
	err := consul.WaitForService(s.cfg)
	s.Require().NoError(err)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		s.cfg.PostgresConfig.PostgresUser,
		s.cfg.PostgresConfig.PostgresPassword,
		s.cfg.PostgresConfig.PostgresHost,
		s.cfg.PostgresConfig.PostgresPort,
		s.cfg.PostgresConfig.PostgresDatabase)

	s.db, err = pgx.Connect(context.Background(), connStr)
	s.Require().NoError(err)

	s.CleanTable()

	s.logger = logger.NewLogger(s.cfg.Env)
}

func (s *TestSuite) TearDownTest() {
	s.CleanTable()
	s.db.Close(context.Background())
}

func (s *TestSuite) CleanTable() {
	tables := []string{"client", "product", "supplier", "image", "address"}

	for _, table := range tables {
		query := fmt.Sprintf(`TRUNCATE TABLE %s CASCADE `, table)
		_, err := s.db.Exec(context.Background(), query)
		s.Require().NoError(err)
	}
}
