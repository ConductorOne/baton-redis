package connector

import (
	"context"

	"github.com/conductorone/baton-redis/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

type clusterBuilder struct {
	resourceType *v2.ResourceType
	client       *client.RedisClient
}

func (o *clusterBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return clusterResourceType
}

func (o *clusterBuilder) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource

	// Note: Redis Enterprise Service API doesn't support pagination.
	clusters, annotation, err := o.client.ListClusters(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for _, cluster := range clusters {
		clusterCopy := cluster
		clusterResource, err := parseIntoClusterResource(ctx, &clusterCopy, nil)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, clusterResource)
	}

	return resources, "", annotation, nil
}

func parseIntoClusterResource(_ context.Context, cluster *client.Cluster, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"name": cluster.Name,
	}

	appTraits := []resource.AppTraitOption{
		resource.WithAppProfile(profile),
	}

	displayName := cluster.Name

	ret, err := resource.NewAppResource(
		displayName,
		clusterResourceType,
		cluster.Name,
		appTraits,
		resource.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *clusterBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func (o *clusterBuilder) Grants(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newClusterBuilder(c *client.RedisClient) *clusterBuilder {
	return &clusterBuilder{
		resourceType: clusterResourceType,
		client:       c,
	}
}
