package docker_compose_yaml

import (
	"fmt"
	"testing"
)

var peerCompose = []PeerCompose{
	{
		OrgName: "beijing",
		Compose: []*PeerComposeDetail{
			{
				IP:            "127.0.0.1",
				DNS:           "localhost",
				Certs:         "/test",
				Image:         "hyperledger/fabric-peer:2.5.2",
				ContainerName: "peer-test1",
				Ports: Ports{
					Peer7051: "7051",
					Peer7052: "7052",
				},
			},
			{
				IP:            "127.0.0.2",
				DNS:           "localhost",
				ContainerName: "peer-test2",
				Certs:         "/test",
				Ports: Ports{
					Peer7051: "7051",
					Peer7052: "7052",
				},
			},
		},
	},
	{
		OrgName: "shanghai",
		Compose: []*PeerComposeDetail{
			{
				IP:            "127.0.0.1",
				DNS:           "localhost",
				Certs:         "/test",
				ContainerName: "peer-test3",
				Ports: Ports{
					Peer7051: "7051",
					Peer7052: "7052",
				},
			},
			{
				IP:            "127.0.0.2",
				DNS:           "localhost",
				Certs:         "/test",
				ContainerName: "peer-test4",
				Ports: Ports{
					Peer7051: "7051",
					Peer7052: "7052",
				},
			},
		},
	},
}

var cliCompose = []CliCompose{
	{
		OrgName: "beijing",
		Compose: []*CliComposeDetail{
			{
				IP:            "127.0.0.1",
				DNS:           "localhost",
				Certs:         "/test",
				ContainerName: "cli-test1",
				Image:         "hyperledger/fabric-tools:2.5.2",
				Ports: Ports{
					Peer7051: "7051",
					Peer7052: "7052",
				},
			},
			{
				IP:            "127.0.0.2",
				DNS:           "localhost",
				ContainerName: "cli-test2",
				Certs:         "/test",
				Ports: Ports{
					Peer7051: "7051",
					Peer7052: "7052",
				},
			},
		},
	},
	{
		OrgName: "shanghai",
		Compose: []*CliComposeDetail{
			{
				IP:            "127.0.0.1",
				DNS:           "localhost",
				Certs:         "/test",
				ContainerName: "cli-test3",
				Ports: Ports{
					Peer7051: "7051",
					Peer7052: "7052",
				},
			},
			{
				IP:            "127.0.0.2",
				DNS:           "localhost",
				Certs:         "/test",
				ContainerName: "cli-test4",
				Ports: Ports{
					Peer7051: "7051",
					Peer7052: "7052",
				},
			},
		},
	},
}

var couchdbCompose = []CouchDBCompose{
	{
		ContainerName: "couchdb-test",
		DB:            "/test",
		Ports: Ports{
			Couchdb5984: "5984",
		},
	},
}

var orderCompose = []OrderCompose{
	{
		ContainerName: "order-test",
		OrgName:       "beijing",
		Certs:         "/test",
		Image:         "hyperledger/fabric-orderer:2.5.2",
		DNS:           "localhost",
	},
	{
		ContainerName: "order-test1",
		OrgName:       "beijing",
		Image:         "hyperledger/fabric-orderer:2.5.2",
		Certs:         "/test1",
		DNS:           "test.com",
	},
}

var exploreCompose = []ExplorerCompose{
	{
		ContainerName:   "test",
		ContainerDBName: "explorer-testDB",
		ImageExplorerDB: "hyperledger/explorer-db:1.1.8",
		ImageExplorer:   "hyperledger/explorer:1.1.8",
		PGData:          "/test",
		Certs:           "/certs",
		FabricConfig:    "fabricConfig",
		WalletStore:     "/walletStore",
		Ports: Ports{
			Explorer8080: "8080",
		},
	},
}

func TestNodeYaml(t *testing.T) {
	n := NewNodeYaml()
	//order
	order, err := n.Order.Create(orderCompose)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range order {
		fmt.Println(file.Name, "\n", string(file.Data))
	}

	fmt.Println("n.Order.Create=======>OK")

	//peer
	peer, err := n.Peer.Create(peerCompose)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range peer {
		fmt.Println(file.Name, "\n", string(file.Data))
	}

	fmt.Println("n.Peer.Create=======>OK")

	//cli
	cli, err := n.Cli.Create(cliCompose)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range cli {
		fmt.Println(file.Name, "\n", string(file.Data))
	}
	fmt.Println("n.Cli.Create=======>OK")

	//couchDB
	couchdb, err := n.CouchDB.Create(couchdbCompose)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range couchdb {
		fmt.Println(file.Name, "\n", string(file.Data))
	}
	fmt.Println("n.CouchDB.Create=======>OK")

	//explorer
	explorer, err := n.Explorer.Create(exploreCompose)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range explorer {
		fmt.Println(file.Name, "\n", string(file.Data))
	}
	fmt.Println("n.Explorer.Create=======>OK")

}
