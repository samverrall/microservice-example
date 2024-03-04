package testdatabase

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samverrall/microservice-example/pkg/postgresutil"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DbName = "test_db"
	DbUser = "test_user"
	DbPass = "test_password"
)

type TestDatabase struct {
	DB        *pgxpool.Pool
	DBAddress string
	container testcontainers.Container
}

func NewTestDatabase(ctx context.Context, migrationsFS fs.FS) (*TestDatabase, error) {
	// setup db container
	container, dbConn, _, err := createContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("createContainer: %w", err)
	}

	// migrate db schema
	if err := postgresutil.MigrateFS(ctx, dbConn, migrationsFS); err != nil {
		return nil, fmt.Errorf("migrate fs: %w", err)
	}

	return &TestDatabase{
		container: container,
		DB:        dbConn,
	}, nil
}

func (tdb *TestDatabase) TearDown() {
	tdb.DB.Close()

	// remove test container
	_ = tdb.container.Terminate(context.Background())
}

func createContainer(ctx context.Context) (testcontainers.Container, *pgxpool.Pool, string, error) {
	env := map[string]string{
		"POSTGRES_PASSWORD": DbPass,
		"POSTGRES_USER":     DbUser,
		"POSTGRES_DB":       DbName,
	}
	port := "5432/tcp"

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         fmt.Sprintf("postgres-test-%s", uuid.NewString()),
			Image:        "postgres:14-alpine",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	})
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())

	// TODO: Look into removing this sleep, and look for hooks to wait for container to be ready
	time.Sleep(time.Second)

	dbAddr := fmt.Sprintf("localhost:%s", p.Port())
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)
	db, err := postgresutil.Connect(ctx, databaseURL)
	if err != nil {
		return container, db, databaseURL, fmt.Errorf("failed to establish database connection: %v", err)
	}

	return container, db, databaseURL, nil
}
