package cluster

import (
	"fmt"
	cliCluster "github.com/openshift-online/ocm-cli/pkg/cluster"
	"github.com/openshift-online/ocm-cli/pkg/ocm"
	sdk "github.com/openshift-online/ocm-sdk-go"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

type NormalizedClusterData struct {
	Name              string
	InternalID        string
	ExternalID        string
	InfraID           string
	HiveShard         string
	ManagementCluster string
	ServiceCluster    string
}

func NewNormalizedCluster(searchPattern string) (*NormalizedClusterData, error) {
	clusterClient := NewClient(searchPattern)
	return clusterClient.CollectNormalizedClusterData()
}

type Client struct {
	conn          *sdk.Connection
	searchPattern string
}

func NewClient(clusterSearchPattern string) *Client {
	return &Client{searchPattern: clusterSearchPattern}
}

func (c *Client) CollectNormalizedClusterData() (*NormalizedClusterData, error) {
	conn, err := ocm.NewConnection().Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create OCM connection: %v", err)
	}
	defer func(conn *sdk.Connection) {
		_ = conn.Close()
	}(conn)

	cluster, err := cliCluster.GetCluster(conn, c.searchPattern)
	if err != nil {
		return nil, fmt.Errorf("can't retrieve cluster for key '%s': %v", c.searchPattern, err)
	}

	ncd := &NormalizedClusterData{
		Name:       cluster.Name(),
		InternalID: cluster.ID(),
		ExternalID: cluster.ExternalID(),
		InfraID:    cluster.InfraID(),
	}

	if cluster.Hypershift().Enabled() {
		ncd.ManagementCluster, ncd.ServiceCluster = findHyperShiftMgmtSvcClusters(conn, cluster)
	} else {
		// Find the details of the shard
		hiveCluster, err := findHiveCluster(conn, cluster)
		if err != nil {
			return nil, err
		}
		ncd.HiveShard = hiveCluster.Name()
	}

	return ncd, nil
}

func findHiveCluster(conn *sdk.Connection, cluster *cmv1.Cluster) (*cmv1.Cluster, error) {
	provisionShard, err := conn.ClustersMgmt().V1().Clusters().
		Cluster(cluster.ID()).
		ProvisionShard().
		Get().
		Send()
	if err != nil {
		return nil, err
	}

	hiveApiUrl, ok := provisionShard.Body().HiveConfig().GetServer()
	if !ok {
		return nil, fmt.Errorf("no provision shard url found for %s", cluster.ID())
	}

	resp, err := conn.ClustersMgmt().V1().Clusters().List().
		Parameter("search", fmt.Sprintf("api.url='%s'", hiveApiUrl)).
		Send()
	if err != nil {
		return nil, err
	}

	if resp.Items().Empty() {
		return nil, fmt.Errorf("failed to find cluster with api.url=%s", hiveApiUrl)
	}

	return resp.Items().Get(0), nil
}

func findHyperShiftMgmtSvcClusters(conn *sdk.Connection, cluster *cmv1.Cluster) (string, string) {
	if !cluster.Hypershift().Enabled() {
		return "", ""
	}

	hypershiftResp, err := conn.ClustersMgmt().V1().Clusters().
		Cluster(cluster.ID()).
		Hypershift().
		Get().
		Send()
	if err != nil {
		return "", ""
	}

	mgmtClusterName := hypershiftResp.Body().ManagementCluster()
	fmMgmtResp, err := conn.OSDFleetMgmt().V1().ManagementClusters().
		List().
		Parameter("search", fmt.Sprintf("name='%s'", mgmtClusterName)).
		Send()
	if err != nil {
		return mgmtClusterName, ""
	}

	if kind := fmMgmtResp.Items().Get(0).Parent().Kind(); kind == "ServiceCluster" {
		return mgmtClusterName, fmMgmtResp.Items().Get(0).Parent().Name()
	}

	// Shouldn't normally happen as every management cluster should have a service cluster
	return mgmtClusterName, ""
}
