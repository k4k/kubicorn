package main

import (
	"fmt"
	"github.com/kris-nova/kubicorn/apis/cluster"
	clusterInit "github.com/kris-nova/kubicorn/cluster"
	"github.com/kris-nova/kubicorn/cutil"
	"github.com/kris-nova/kubicorn/logger"
)

func main() {
	logger.Level = 4
	cluster := getCluster("mycluster")
	cluster, err := clusterInit.InitCluster(cluster)
	if err != nil {
		panic(err.Error())
	}
	reconciler, err := cutil.GetReconciler(cluster)
	if err != nil {
		panic(err.Error())
	}

	err = reconciler.Init()
	if err != nil {
		panic(err.Error())
	}
	expected, err := reconciler.GetExpected()
	if err != nil {
		panic(err.Error())
	}
	actual, err := reconciler.GetActual()
	if err != nil {
		panic(err.Error())
	}
	created, err := reconciler.Reconcile(actual, expected)
	logger.Info("Created cluster [%s]", created.Name)
	if err != nil {
		panic(err.Error())
	}
	err = reconciler.Destroy()
	if err != nil {
		panic(err.Error())
	}
}

func getCluster(name string) *cluster.Cluster {
	return &cluster.Cluster{
		Name:     name,
		Cloud:    cluster.Cloud_Amazon,
		Location: "us-west-2",
		Ssh: &cluster.Ssh{
			PublicKeyPath: "~/.ssh/id_rsa.pub",
			User:          "ubuntu",
		},
		KubernetesApi: &cluster.KubernetesApi{
			Port: "443",
		},
		Network: &cluster.Network{
			Type: cluster.NetworkType_Public,
			CIDR: "10.0.0.0/16",
		},
		Values: &cluster.Values{
			ItemMap: map[string]string{
				"INJECTEDTOKEN": "829a9b.a839d03b8d810c56",
			},
		},
		ServerPools: []*cluster.ServerPool{
			{
				Type:            cluster.ServerPoolType_Master,
				Name:            fmt.Sprintf("%s.master", name),
				MaxCount:        1,
				MinCount:        1,
				Image:           "ami-835b4efa",
				Size:            "t2.medium",
				BootstrapScript: "1.7.0_ubuntu_16.04_master.sh",
				Subnets: []*cluster.Subnet{
					{
						Name:     fmt.Sprintf("%s.master", name),
						CIDR:     "10.0.0.0/24",
						Location: "us-west-2a",
					},
				},

				Firewalls: []*cluster.Firewall{
					{
						Name: fmt.Sprintf("%s.master-external", name),
						Rules: []*cluster.Rule{
							{
								IngressFromPort: 22,
								IngressToPort:   22,
								IngressSource:   "0.0.0.0/0",
								IngressProtocol: "tcp",
							},
							{
								IngressFromPort: 443,
								IngressToPort:   443,
								IngressSource:   "0.0.0.0/0",
								IngressProtocol: "tcp",
							},
						},
					},
				},
			},
			{
				Type:            cluster.ServerPoolType_Node,
				Name:            fmt.Sprintf("%s.node", name),
				MaxCount:        1,
				MinCount:        1,
				Image:           "ami-835b4efa",
				Size:            "t2.medium",
				BootstrapScript: "1.7.0_ubuntu_16.04_node.sh",
				Subnets: []*cluster.Subnet{
					{
						Name:     fmt.Sprintf("%s.node", name),
						CIDR:     "10.0.100.0/24",
						Location: "us-west-2b",
					},
				},
				Firewalls: []*cluster.Firewall{
					{
						Name: fmt.Sprintf("%s.node-external", name),
						Rules: []*cluster.Rule{
							{
								IngressFromPort: 22,
								IngressToPort:   22,
								IngressSource:   "0.0.0.0/0",
								IngressProtocol: "tcp",
							},
						},
					},
				},
			},
		},
	}
}
