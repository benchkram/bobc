package environment

import (
	"github.com/benchkram/bobc/pkg/rnd"
	"github.com/benchkram/errz"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type Network struct {
	pool    *dockertest.Pool
	network *docker.Network
}

func NewNetwork() (n *Network, err error) {
	defer errz.Recover(&err)

	n = &Network{}

	n.pool, err = dockertest.NewPool("")
	errz.Fatal(err)

	n.network, err = n.pool.Client.CreateNetwork(docker.CreateNetworkOptions{
		Name: "bob-test-" + rnd.RandStringBytesMaskImprSrc(10),
	})
	errz.Fatal(err)

	return n, nil
}

func (n *Network) ID() string {
	return n.network.ID
}

func (n *Network) Close() (err error) {
	defer errz.Recover(&err)

	err = n.pool.Client.RemoveNetwork(n.network.ID)
	errz.Fatal(err)

	return nil
}
