package environment

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/benchkram/bobc/pkg/rnd"
	"github.com/benchkram/errz"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/rs/zerolog/log"
)

const (
	minioRepositoryName = "minio/minio"

	minioContainerPort    = "9000"                      // Port used inside the minio container
	minioContainerPortTCP = minioContainerPort + "/tcp" // Port used inside the minio container
)

var ErrMinioNotInitialized = fmt.Errorf("minio  not initialized")

type MinIO struct {

	// tag of the minio docker container to use
	tag string

	// networkID of the network used
	networkID string

	pool     *dockertest.Pool
	resource *dockertest.Resource
}

type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

func NewMinio() *MinIO {
	m := &MinIO{
		tag: "latest",
	}
	return m
}

func NewMinioStarted() (*MinIO, error) {
	m := &MinIO{}
	return m, m.start()
}

func (m *MinIO) start() (err error) {
	defer errz.Recover(&err)

	m.pool, err = dockertest.NewPool("")
	errz.Fatal(err)

	// if no network passed create our one with a random name.
	if m.networkID == "" {
		network, err := m.pool.Client.CreateNetwork(docker.CreateNetworkOptions{
			Name: "bob-test-minio-" + rnd.RandStringBytesMaskImprSrc(10),
		})
		errz.Fatal(err)

		m.networkID = network.ID
	}

	env := []string{}
	options := &dockertest.RunOptions{
		Repository: minioRepositoryName,
		Tag:        m.tag,
		Env:        env,
		NetworkID:  m.networkID,
		Cmd:        []string{"server", "/data"},
	}
	m.resource, err = m.pool.RunWithOptions(options)
	errz.Fatal(err)

	attempt := 1
	err = retry.Do(
		func() error {
			log.Printf("Trying to check minio health: [attempt: %d]\n", attempt)
			attempt++

			resp, err := http.Get("http://localhost:" + strconv.Itoa(m.Port()) + "/minio/health/ready")
			if err != nil {
				return err
			}

			log.Printf("%d\n", resp.StatusCode)
			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				return fmt.Errorf("failing response code")
			}

			return nil
		},
		retry.Attempts(20),
		retry.Delay(1*time.Second),
		retry.DelayType(retry.FixedDelay),
	)
	errz.Fatal(err)

	return nil
}

// Stop the docker container
func (m *MinIO) Stop(removeNetwork bool) (err error) {
	defer errz.Recover(&err)

	if m.pool == nil {
		return ErrPostgresNotInitialized
	}
	// if m.debug {
	// 	err = p.debugShutdown()
	// 	errz.Fatal(err)
	// }

	err = m.pool.Purge(m.resource)
	errz.Fatal(err)
	if removeNetwork {
		err = m.pool.Client.RemoveNetwork(m.networkID)
		errz.Fatal(err)
	}

	return nil
}

// Config return the config object usually passed to application inside bob
func (m *MinIO) Config() MinIOConfig {

	port := m.resource.GetPort(minioContainerPortTCP)
	port = strings.TrimRight(port, ":")

	// endpoint := url.URL{
	// 	Host: "127.0.0.1:" + port,
	// }

	c := MinIOConfig{
		Endpoint:        "localhost:" + port,
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
		UseSSL:          false,
	}

	return c
}

// Port returns the external port of the database
func (m *MinIO) Port() int {

	if m.resource == nil {
		panic(ErrMinioNotInitialized)
	}
	port, err := strconv.Atoi(m.resource.GetPort(minioContainerPortTCP))
	errz.Fatal(err)

	return port
}
