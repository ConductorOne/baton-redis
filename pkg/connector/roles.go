package connector

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/conductorone/baton-redis/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
)

type roleBuilder struct {
	resourceType *v2.ResourceType
	client       *client.RedisClient
	users        []client.User
	usersMutex   sync.RWMutex
	roles        map[int]client.Role
	rolesMutex   sync.RWMutex
}

func (o *roleBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return roleResourceType
}

func (o *roleBuilder) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource

	// Note: Redis Enterprise Service API doesn't support pagination.
	roles, annotation, err := o.client.ListRoles(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for _, role := range roles {
		roleCopy := role
		roleResource, err := parseIntoRoleResource(ctx, &roleCopy, nil)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, roleResource)
	}

	return resources, "", annotation, nil
}

func parseIntoRoleResource(_ context.Context, role *client.Role, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"role_id":         role.UID,
		"name":            role.Name,
		"management_role": role.Management,
	}

	roleTraits := []resource.RoleTraitOption{
		resource.WithRoleProfile(profile),
	}

	displayName := role.Name

	ret, err := resource.NewRoleResource(
		displayName,
		roleResourceType,
		role.UID,
		roleTraits,
		resource.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *roleBuilder) Entitlements(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var entitlements []*v2.Entitlement
	role, _, err := o.client.GetRoleDetails(ctx, resource.Id.Resource)

	if err != nil {
		return nil, "", nil, err
	}

	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("Role %s with management %s in Redis", role.Name, role.Management)),
		entitlement.WithDisplayName(fmt.Sprintf("%s Role %s", resource.DisplayName, role.Management)),
	}

	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(resource, role.Management, assigmentOptions...))

	return entitlements, "", nil, nil
}

func (o *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var grants []*v2.Grant

	// Note: Redis Enterprise Service API doesn't support pagination.
	err := o.GetUsers(ctx)
	if err != nil {
		return nil, "", nil, err
	}
	err = o.GetRoles(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for _, user := range o.users {
		for _, roleUID := range user.RoleUIDs {
			if strconv.Itoa(roleUID) == resource.Id.Resource {
				userResource, _ := parseIntoUserResource(ctx, &user, nil)
				role := o.roles[roleUID]

				userGrant := grant.NewGrant(resource, role.Management, userResource, grant.WithAnnotation(&v2.V1Identifier{
					Id: fmt.Sprintf("role-grant:%s:%d:%s", resource.Id.Resource, user.UID, role.Management),
				}))
				grants = append(grants, userGrant)
			}
		}
	}

	return grants, "", nil, nil
}

func newRoleBuilder(c *client.RedisClient) *roleBuilder {
	return &roleBuilder{
		resourceType: roleResourceType,
		client:       c,
	}
}

func (o *roleBuilder) GetUsers(ctx context.Context) error {
	o.usersMutex.RLock()
	defer o.usersMutex.RUnlock()

	if o.users != nil || len(o.users) > 0 {
		return nil
	}

	users, _, err := o.client.ListUsers(ctx)

	if err != nil {
		return err
	}

	o.users = users

	return nil
}

func (o *roleBuilder) GetRoles(ctx context.Context) error {
	o.rolesMutex.RLock()
	defer o.rolesMutex.RUnlock()

	if o.roles == nil {
		o.roles = make(map[int]client.Role)
	}

	if len(o.roles) > 0 {
		return nil
	}

	roles, _, err := o.client.ListRoles(ctx)

	if err != nil {
		return err
	}

	for _, role := range roles {
		o.roles[role.UID] = role
	}

	return nil
}
