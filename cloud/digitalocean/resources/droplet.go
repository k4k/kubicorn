package resources

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/kris-nova/kubicorn/apis/cluster"
	"github.com/kris-nova/kubicorn/cloud"
	"github.com/kris-nova/kubicorn/cutil/compare"
	"github.com/kris-nova/kubicorn/logger"
	"strconv"
)

type Droplet struct {
	Shared
	Region         string
	Size           string
	Image          string
	Count          int
	SShFingerprint string
	ServerPool     *cluster.ServerPool
}

func (r *Droplet) Actual(known *cluster.Cluster) (cloud.Resource, error) {
	logger.Debug("droplet.Actual")
	if r.CachedActual != nil {
		logger.Debug("Using cached droplet [actual]")
		return r.CachedActual, nil
	}
	actual := &Droplet{
		Shared: Shared{
			Name:    r.Name,
			CloudID: r.ServerPool.Identifier,
		},
	}

	if r.CloudID != "" {

		droplets, _, err := Sdk.Client.Droplets.ListByTag(context.TODO(), r.Name, &godo.ListOptions{})
		if err != nil {
			return nil, err
		}
		ld := len(droplets)
		if ld != 1 {
			return nil, fmt.Errorf("Found [%d] Droplets for Name [%s]", ld, r.Name)
		}
		droplet := droplets[0]
		id := strconv.Itoa(droplet.ID)
		actual.Name = droplet.Name
		actual.CloudID = id
		actual.Size = droplet.Size.Slug
		actual.Region = droplet.Region.Name
		actual.Image = droplet.Image.Slug
	}
	actual.Count = r.ServerPool.MaxCount
	actual.Name = r.Name
	r.CachedActual = actual
	return actual, nil
}

func (r *Droplet) Expected(known *cluster.Cluster) (cloud.Resource, error) {
	logger.Debug("droplet.Expected")
	if r.CachedExpected != nil {
		logger.Debug("Using droplet subnet [expected]")
		return r.CachedExpected, nil
	}
	expected := &Droplet{
		Shared: Shared{
			Name:    r.Name,
			CloudID: r.ServerPool.Identifier,
		},
		Size:   r.ServerPool.Size,
		Region: known.Location,
		Image:  r.ServerPool.Image,
		Count:  r.ServerPool.MaxCount,
	}
	r.CachedExpected = expected
	return expected, nil
}

func (r *Droplet) Apply(actual, expected cloud.Resource, applyCluster *cluster.Cluster) (cloud.Resource, error) {
	logger.Debug("droplet.Apply")
	applyResource := expected.(*Droplet)
	isEqual, err := compare.IsEqual(actual.(*Droplet), expected.(*Droplet))
	if err != nil {
		return nil, err
	}
	if isEqual {
		return applyResource, nil
	}

	createRequest := &godo.DropletCreateRequest{
		Name:   expected.(*Droplet).Name,
		Region: expected.(*Droplet).Region,
		Size:   expected.(*Droplet).Size,
		Image: godo.DropletCreateImage{
			Slug: expected.(*Droplet).Image,
		},
		Tags:              []string{expected.(*Droplet).Name},
		PrivateNetworking: true,
		//SSHKeys: []godo.DropletCreateSSHKey{
		//	Fingerprint: "",
		//},
	}
	droplet, _, err := Sdk.Client.Droplets.Create(context.TODO(), createRequest)
	if err != nil {
		return nil, err
	}

	logger.Info("Created Droplet [%d]", droplet.ID)
	id := strconv.Itoa(droplet.ID)
	newResource := &Droplet{
		Shared: Shared{
			Name:    droplet.Name,
			CloudID: id,
		},
		Image:  droplet.Image.Slug,
		Size:   droplet.Size.Slug,
		Region: droplet.Region.Name,
		Count:  expected.(*Droplet).Count,
	}
	return newResource, nil
}
func (r *Droplet) Delete(actual cloud.Resource, known *cluster.Cluster) error {
	logger.Debug("droplet.Delete")
	deleteResource := actual.(*Droplet)
	if deleteResource.Name == "" {
		return fmt.Errorf("Unable to delete droplet resource without Name [%s]", deleteResource.Name)
	}

	droplets, _, err := Sdk.Client.Droplets.ListByTag(context.TODO(), r.Name, &godo.ListOptions{})
	if err != nil {
		return err
	}
	ld := len(droplets)
	if ld != 1 {
		return fmt.Errorf("Found [%d] Droplets for Name [%s]", ld, r.Name)
	}
	droplet := droplets[0]
	_, err = Sdk.Client.Droplets.Delete(context.TODO(), droplet.ID)
	if err != nil {
		return err
	}
	logger.Info("Deleted Droplet [%d]", droplet.ID)
	return nil
}

func (r *Droplet) Render(renderResource cloud.Resource, renderCluster *cluster.Cluster) (*cluster.Cluster, error) {
	logger.Debug("droplet.Render")

	serverPool := &cluster.ServerPool{}
	serverPool.Image = renderResource.(*Droplet).Image
	serverPool.Size = renderResource.(*Droplet).Size
	serverPool.Name = renderResource.(*Droplet).Name
	serverPool.MaxCount = renderResource.(*Droplet).Count
	found := false
	for i := 0; i < len(renderCluster.ServerPools); i++ {
		if renderCluster.ServerPools[i].Name == renderResource.(*Droplet).Name {
			renderCluster.ServerPools[i].Image = renderResource.(*Droplet).Image
			renderCluster.ServerPools[i].Size = renderResource.(*Droplet).Size
			renderCluster.ServerPools[i].MaxCount = renderResource.(*Droplet).Count
			found = true
		}
	}
	if !found {
		renderCluster.ServerPools = append(renderCluster.ServerPools, serverPool)
	}
	renderCluster.Location = renderResource.(*Droplet).Region
	return renderCluster, nil
}

func (r *Droplet) Tag(tags map[string]string) error {
	return nil
}
