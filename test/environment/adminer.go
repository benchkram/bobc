package environment

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/benchkram/bobc/pkg/rnd"
	"github.com/benchkram/errz"
	"github.com/logrusorgru/aurora"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

const (
	adminerRepositoryName = "adminer" // Repo name to get adminer image
	adminerTag            = "4.7.8"

	adminerContainerPort    = "8080"                        // Port used inside the adminer container
	adminerContainerPortTCP = adminerContainerPort + "/tcp" // Port used inside the adminer container
)

var ErrAdminerNotInitialized = fmt.Errorf("adminer not initialized")

type Adminer struct {
	host string

	// tag of the adminer docker container to use
	tag string

	// networkID of the network used
	networkID string
	hostPort  string

	pool     *dockertest.Pool
	resource *dockertest.Resource

	once sync.Once
}

func NewAdminer(opts ...AdminerOption) *Adminer {
	a := &Adminer{
		host: "localhost",
		tag:  adminerTag,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// NewAdminerStarted returns adminer running in a docker dontainer.
func NewAdminerStarted(opts ...AdminerOption) (*Adminer, error) {
	a := NewAdminer(opts...)
	return a, a.Start()
}

// Start adminer inside a docker container.
// Can only be called once during the lifetime of adminer.
func (a *Adminer) Start() (err error) {
	defer errz.Recover(&err)

	f := func() {
		err = a.start()
		errz.Fatal(err)
	}
	a.once.Do(f)

	fmt.Printf("%s", aurora.Green(fmt.Sprintf("adminer started [%s:%d]\n", a.host, a.Port())))

	return nil
}

func (a *Adminer) start() (err error) {
	defer errz.Recover(&err)

	a.pool, err = dockertest.NewPool("")
	errz.Fatal(err)

	// if no network passed.. create our one with a random name.
	if a.networkID == "" {
		network, err := a.pool.Client.CreateNetwork(docker.CreateNetworkOptions{
			Name: "plyd-test-" + rnd.RandStringBytesMaskImprSrc(10),
		})
		errz.Fatal(err)

		a.networkID = network.ID
	}

	env := []string{}
	options := &dockertest.RunOptions{
		Repository: adminerRepositoryName,
		Tag:        a.tag,
		Env:        env,
		NetworkID:  a.networkID,
	}
	// Set a fixed host port
	if a.hostPort != "" {
		options.PortBindings = make(map[docker.Port][]docker.PortBinding)
		options.PortBindings[docker.Port(a.hostPort+"/tcp")] = []docker.PortBinding{{HostIP: "0.0.0.0", HostPort: "8080"}}
	}
	a.resource, err = a.pool.RunWithOptions(options)
	errz.Fatal(err)

	err = a.pool.Retry(func() error {
		// TODO: implement health check
		return nil
	})
	errz.Fatal(err)

	return nil
}

// Stop the docker container
func (a *Adminer) Stop(removeNetwork bool) (err error) {
	defer errz.Recover(&err)
	if a.pool == nil {
		return ErrAdminerNotInitialized
	}

	err = a.pool.Purge(a.resource)
	errz.Fatal(err)
	if removeNetwork {
		err = a.pool.Client.RemoveNetwork(a.networkID)
		errz.Fatal(err)
	}

	return nil
}

// Port return the random external port of the docker dontainer
func (a *Adminer) Domain() string {
	return a.host
}

// Port return the random external port of the docker dontainer
func (a *Adminer) Port() int {
	if a.resource == nil {
		panic(ErrAdminerNotInitialized)
	}
	port, err := strconv.Atoi(a.resource.GetPort(adminerContainerPortTCP))
	errz.Fatal(err)

	return port
}

// PortContainer returns port the service is listening to inside the container.
func (a *Adminer) PortContainer() int {
	if a.resource == nil {
		panic(ErrAdminerNotInitialized)
	}

	port, err := strconv.Atoi(adminerContainerPort)
	errz.Fatal(err)

	return port
}

// IpContainer returns the ip of the container in th network with `networkID`
func (a *Adminer) IpContainer() string {
	if a.resource == nil {
		panic(ErrAdminerNotInitialized)
	}

	var postgresDockerNetworkIp string
	for _, network := range a.resource.Container.NetworkSettings.Networks {
		if network.NetworkID == a.networkID {
			postgresDockerNetworkIp = network.IPAddress
		}
	}

	return postgresDockerNetworkIp
}

func (a *Adminer) NetworkID() string {
	return a.networkID
}

type AdminerOption func(*Adminer)

func WithAdminerHost(d string) AdminerOption {
	return func(a *Adminer) {
		a.host = d
	}
}

// WithNetwork by default adminer creates it's own random network.
// Pass a network id to make adminer join a existings network.
func WithAdminerNetwork(networkID string) AdminerOption {
	return func(a *Adminer) {
		a.networkID = networkID
	}
}

// WithFixedHostPort sets a fixed host port
func WithFixedHostPort(port string) AdminerOption {
	return func(a *Adminer) {
		a.hostPort = port
	}
}
