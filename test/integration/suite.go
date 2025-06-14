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

	s.logger = logger.NewLogger(s.cfg.Env)
}

func (s *TestSuite) TearDownTest() {
	s.db.Close(context.Background())
}
