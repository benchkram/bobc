package environment

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/benchkram/bobc/pkg/db"
	"github.com/benchkram/bobc/pkg/rnd"
	"github.com/benchkram/bobc/pkg/wait"
	"github.com/benchkram/errz"
	"github.com/jackc/pgx/v4"
	"github.com/logrusorgru/aurora"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog/log"
	"github.com/xo/dburl"
)

var debug bool

func init() {
	flag.BoolVar(&debug, "debug", false, "Start adminer and wait for ctrl-c at the end of test")
}

// Manage a postgres database connection testing.
//
// `debug`
// Debug mode is enabled by setting the debug flag `go test --debug` ot `ginkgo -- --debug`.
// In debug mode an additional `adminer` container is launched, to enable debugging the database.
// Before droping all tables or stoping the database wait for SIGTERM to allow debugging the database.
// Waiting for SIGTERM is guarded by `sync.Once`.
//
// `ci`
// On CI a external database is used and it's not necessary to start a database.
// In this case most functions are turned `off` and only the database config
// is usable. This basically degrades this object to a config provider.

const (
	postgresRepositoryName = "postgres" // Repo name to get postgres image
	postgresTag            = "11.2-alpine"

	postgresContainerPort    = "5432"                         // Port used inside the postgres container
	postgresContainerPortTCP = postgresContainerPort + "/tcp" // Port used inside the postgres container

	defaultAdminerPort = "8080" // default port passed to adminer

	connectStr = "host=%s port=%d user=%s dbname=%s password=%s sslmode=disable"
)

var ErrPostgresNotInitialized = fmt.Errorf("postgres  not initialized")
var ErrPostgresAdminerIsNil = fmt.Errorf("postgres: adminer is nil")

type Postgres struct {
	// host the database can be reached
	// usually localhost.
	hostname string
	// port is not set when working in a local environmet
	// but it can be set when workin on CI an reading the database
	// connection from a environment variable.
	port string
	// default database created on postgres
	databaseName string
	user         string
	password     string

	// debug mode runs a additional adminer container
	// and waits for SIGTERM before stopping.
	debug   bool
	adminer *Adminer

	// ci is set to true when executed on environment
	// in which postgres is already provided.
	// In this case the postgres config is gathered and nothing more is done.
	ci bool

	// tag of the postgres docker container to use
	tag string

	// networkID of the network used
	networkID string

	pool     *dockertest.Pool
	resource *dockertest.Resource

	connection *pgx.Conn

	once      sync.Once
	debugOnce sync.Once
}

func NewPostgres(opts ...Option) *Postgres {
	p := &Postgres{
		hostname:     "localhost",
		tag:          postgresTag,
		user:         "postgres",
		password:     rnd.RandStringBytesMaskImprSrc(10),
		databaseName: strings.ToLower(rnd.RandStringBytesMaskImprSrc(10)),
	}

	for _, opt := range opts {
		opt(p)
	}

	// override debug option with cli value.
	if debug {
		p.debug = debug
	}

	if p.debug {
		p.password = "postgres"
		p.databaseName = "postgres"
	}

	postgresURL := os.Getenv("CI_POSTGRES_URL")
	if postgresURL != "" {
		u, err := dburl.Parse(postgresURL)
		if err != nil {
			panic(err)
		}

		// read port fron connection string
		// but default to postgres default port in case of error.
		p.port = postgresContainerPort
		p.hostname = u.Hostname()
		if strings.Contains(u.Host, ":") {
			parts := strings.Split(u.Host, ":")
			if len(parts) != 2 {
				panic("Invalid Host [hots:port]")
			}
			port := parts[1]
			if port != "" {
				p.port = port
			}
		}

		p.user = u.User.Username()
		p.password, _ = u.User.Password()

		// Get database name
		databaseName := u.URL.Path
		databaseName = strings.Trim(databaseName, "/")
		p.databaseName = databaseName
		if p.databaseName == "" {
			panic("Could not read database name")
		}

		p.ci = true
	}
	if p.ci {
		log.Info().Msg("CI mode enabled. Assuming postgres is provided.")
	}

	return p
}

// NewPostgresStarted returns a postgres database running in a docker dontainer.
func NewPostgresStarted(opts ...Option) (*Postgres, error) {
	p := NewPostgres(opts...)
	return p, p.Start()
}

// Start the postgres database inside a docker container.
// Can only be called once during the lifetime of a database.
func (p *Postgres) Start() (err error) {
	defer errz.Recover(&err)

	f := func() {
		err = p.start()
		errz.Fatal(err)

		if p.debug {
			err = p.startAdminer()
			errz.Fatal(err)
		}
	}
	p.once.Do(f)

	fmt.Printf("%s", aurora.Green(fmt.Sprintf("Postgres database started [%s:%d, %s, %s]\n", p.hostname, p.Port(), p.user, p.password)))

	return nil
}

func (p *Postgres) start() (err error) {
	defer errz.Recover(&err)

	// start database when not on CI
	if !p.ci {
		p.pool, err = dockertest.NewPool("")
		errz.Fatal(err)

		// if no network passed.. create our one with a random name.
		if p.networkID == "" {
			network, err := p.pool.Client.CreateNetwork(docker.CreateNetworkOptions{
				Name: "bob-test-" + rnd.RandStringBytesMaskImprSrc(10),
			})
			errz.Fatal(err)

			p.networkID = network.ID
		}

		env := []string{
			"POSTGRES_USER=" + p.user,
			"POSTGRES_PASSWORD=" + p.password,
			"POSTGRES_DB=" + p.databaseName,
		}
		options := &dockertest.RunOptions{
			Repository: postgresRepositoryName,
			Tag:        p.tag,
			Env:        env,
			NetworkID:  p.networkID,
		}
		p.resource, err = p.pool.RunWithOptions(options)
		errz.Fatal(err)
	}

	var connection *pgx.Conn
	attempt := 1
	err = retry.Do(
		func() error {
			log.Printf("Trying to connect to database: [attempt: %d]\n", attempt)
			attempt++

			var err error
			connection, err = pgx.Connect(context.Background(), p.connectString())
			if err != nil {
				return err
			}
			return connection.Ping(context.Background())
		},
		retry.Attempts(20),
		retry.Delay(1*time.Second),
		retry.DelayType(retry.FixedDelay),
	)
	errz.Fatal(err)
	p.connection = connection

	return nil
}

// startAdminer in the same network as the postgres container has been started.
func (p *Postgres) startAdminer() (err error) {
	defer errz.Recover(&err)

	if p.networkID == "" {
		return ErrPostgresNotInitialized
	}

	adminer, err := NewAdminerStarted(
		WithAdminerNetwork(p.networkID),
		WithFixedHostPort(defaultAdminerPort),
	)
	errz.Fatal(err)

	p.adminer = adminer

	return nil
}

// Config return the config object usually passed to application inside bob
func (p *Postgres) Config() db.Config {
	c := db.Config{
		Host:     p.hostname,
		Port:     p.port,
		User:     p.user,
		Password: p.password,
		Name:     p.databaseName,
	}

	// Return port accessible from host system
	if !p.ci {
		c.Port = p.resource.GetPort(postgresContainerPortTCP)
	}

	return c
}

// Drop all tables in a database. Useful for cleanup up after a test.
func (p *Postgres) Drop() (err error) {
	defer errz.Recover(&err)

	sql := `DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO %user%;
GRANT ALL ON SCHEMA public TO public;`
	sql = strings.Replace(sql, "%user%", p.user, 1)
	_, err = p.connection.Exec(context.Background(), sql)

	// Show debuginfo & debug teardown before dropping database
	if p.debug {
		err = p.debugShutdown()
		errz.Fatal(err)
	}

	return err
}

// debugShutdown assure it is only called once so it can safely use when
// dropping and stopping the database.
func (p *Postgres) debugShutdown() (err error) {
	defer errz.Recover(&err)

	f := func() {

		// Print hints the test environment endpoints for debugging.
		fmt.Printf("%s", aurora.Cyan("\n\nDebugging mode enabled.\nStopping to enable database inspection."))
		fmt.Printf("\n")
		fmt.Printf("\n")
		fmt.Printf("%s", aurora.Cyan("  Hint: Call Postgres from Adminer on the docker network\n"))
		fmt.Printf("\n")
		fmt.Printf("%s", aurora.Cyan("  HOST NETWORK"))
		fmt.Printf("\n")
		fmt.Printf("%s", aurora.Cyan(fmt.Sprintf("  Postgres: http://%s:%d\n", p.hostname, p.Port())))
		fmt.Printf("%s", aurora.Cyan(fmt.Sprintf("  Adminer:  http://%s:%d\n", p.adminer.Domain(), p.adminer.Port())))
		fmt.Printf("\n")
		fmt.Printf("%s", aurora.Cyan("  DOCKER NETWORK\n"))
		fmt.Printf("%s", aurora.Cyan(fmt.Sprintf("  Postgres: http://%s:%d\n", p.IpContainer(), p.PortContainer())))
		fmt.Printf("%s", aurora.Cyan(fmt.Sprintf("  Adminer:  http://%s:%d\n", p.adminer.IpContainer(), p.adminer.PortContainer())))
		fmt.Printf("\n")
		fmt.Printf("%s", aurora.Cyan("  POSTGRES"))
		fmt.Printf("\n")
		fmt.Printf("%s", aurora.Cyan(fmt.Sprintf("  User:     %s\n", p.user)))
		fmt.Printf("%s", aurora.Cyan(fmt.Sprintf("  Password: %s\n", p.password)))
		fmt.Printf("%s", aurora.Cyan(fmt.Sprintf("  Database: %s\n", p.databaseName)))
		fmt.Printf("\n")

		fmt.Printf("%s", aurora.Bold(aurora.Green("  TL;DR\n")))
		fmt.Printf("%s%s",
			aurora.Green("  Access To Adminer: "),
			aurora.Bold(aurora.Green(fmt.Sprintf(" http://%s:%d/?pgsql=%s&username=%s&db=%s ", p.adminer.Domain(), p.adminer.Port(), p.IpContainer(), p.user, p.databaseName))),
		)
		fmt.Printf("\n")
		fmt.Printf("\n")
		fmt.Printf("\n")
		fmt.Printf("\n")

		wait.ForCtrlC()
		if p.adminer == nil {
			errz.Fatal(ErrPostgresAdminerIsNil)
		}
		err = p.adminer.Stop(false)
		errz.Fatal(err)
	}
	p.debugOnce.Do(f)

	return nil
}

// Stop the docker container
func (p *Postgres) Stop(removeNetwork bool) (err error) {
	defer errz.Recover(&err)

	// Nothing to do on ci
	if p.ci {
		return nil
	}
	if p.pool == nil {
		return ErrPostgresNotInitialized
	}
	if p.debug {
		err = p.debugShutdown()
		errz.Fatal(err)
	}

	err = p.pool.Purge(p.resource)
	errz.Fatal(err)
	if removeNetwork {
		err = p.pool.Client.RemoveNetwork(p.networkID)
		errz.Fatal(err)
	}

	return nil
}

func (p *Postgres) Database() string {
	return p.databaseName
}

// DatabaseName returns the database name
func (p *Postgres) DatabaseName() string {
	return p.databaseName
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?connect_timeout=10&sslmode=disable&max_conns=20&max_idle_conns=4",
		p.user,
		p.password,
		p.IpContainer(),
		p.PortContainer(),
		p.databaseName,
	)
}

// Port return the external port of of the database
func (p *Postgres) Port() int {
	if p.ci && p.port != "" {
		port, err := strconv.Atoi(p.port)
		errz.Fatal(err)
		return port
	}
	if p.resource == nil {
		panic(ErrPostgresNotInitialized)
	}
	port, err := strconv.Atoi(p.resource.GetPort(postgresContainerPortTCP))
	errz.Fatal(err)

	return port
}

// PortContainer returns port the service is listening to inside the container.
func (p *Postgres) PortContainer() int {
	if p.resource == nil {
		panic(ErrPostgresNotInitialized)
	}

	port, err := strconv.Atoi(postgresContainerPort)
	errz.Fatal(err)

	return port
}

// IpContainer returns the ip of the container in th network with `networkID`
func (p *Postgres) IpContainer() string {
	if p.resource == nil {
		panic(ErrPostgresNotInitialized)
	}

	var postgresDockerNetworkIp string
	for _, network := range p.resource.Container.NetworkSettings.Networks {
		if network.NetworkID == p.networkID {
			postgresDockerNetworkIp = network.IPAddress
		}
	}

	return postgresDockerNetworkIp
}

func (p *Postgres) NetworkID() string {
	return p.networkID
}

func (p *Postgres) connectString() string {
	return fmt.Sprintf(connectStr,
		p.hostname,
		p.Port(),
		p.user,
		p.databaseName,
		p.password,
	)
}

type Option func(*Postgres)

func WithHostname(hostname string) Option {
	return func(p *Postgres) {
		p.hostname = hostname
	}
}

func WithTag(tag string) Option {
	return func(p *Postgres) {
		p.tag = tag
	}
}

func WithCredentials(user, password string) Option {
	return func(p *Postgres) {
		p.user = user
		p.password = password
	}
}

func WithDatabaseName(name string) Option {
	return func(p *Postgres) {
		p.databaseName = name
	}
}

// WithPostgresNetwork by default postgres creates it's own random network.
// Pass a network id to make postgres join a existings network.
func WithPostgresNetwork(networkID string) Option {
	return func(p *Postgres) {
		p.networkID = networkID
	}
}

// WithDebug in debug mode postgres will additionaly start a adminer container
// to inspect the database. It also waits at the end of the test for SIGTERM
// to enable database inspection.
func WithDebug(debug bool) Option {
	return func(p *Postgres) {
		p.debug = debug
	}
}
