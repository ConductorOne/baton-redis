package connector

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/conductorone/baton-redis/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
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

func (o *roleBuilder) List(ctx context.Context, _ *v2.ResourceId, _ rs.SyncOpAttrs) ([]*v2.Resource, *rs.SyncOpResults, error) {
	var resources []*v2.Resource

	// Note: Redis Enterprise Service API doesn't support pagination.
	roles, annos, err := o.client.ListRoles(ctx)
	if err != nil {
		return nil, nil, err
	}

	for _, role := range roles {
		roleCopy := role
		roleResource, err := parseIntoRoleResource(ctx, &roleCopy, nil)
		if err != nil {
			return nil, nil, err
		}
		resources = append(resources, roleResource)
	}

	return resources, &rs.SyncOpResults{Annotations: annos}, nil
}

func parseIntoRoleResource(_ context.Context, role *client.Role, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"role_id":         role.UID,
		"name":            role.Name,
		"management_role": role.Management,
	}

	roleTraits := []rs.RoleTraitOption{
		rs.WithRoleProfile(profile),
	}

	displayName := role.Name

	ret, err := rs.NewRoleResource(
		displayName,
		roleResourceType,
		role.UID,
		roleTraits,
		rs.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *roleBuilder) Entitlements(ctx context.Context, resource *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Entitlement, *rs.SyncOpResults, error) {
	var entitlements []*v2.Entitlement
	role, annos, err := o.client.GetRoleDetails(ctx, resource.Id.Resource)

	if err != nil {
		return nil, nil, err
	}

	assigmentOptions := []entitlement.EntitlementOption{
		entitlement.WithGrantableTo(userResourceType),
		entitlement.WithDescription(fmt.Sprintf("Role %s with management %s in Redis", role.Name, role.Management)),
		entitlement.WithDisplayName(fmt.Sprintf("%s Role %s", resource.DisplayName, role.Management)),
	}

	entitlements = append(entitlements, entitlement.NewPermissionEntitlement(resource, role.Management, assigmentOptions...))

	return entitlements, &rs.SyncOpResults{Annotations: annos}, nil
}

func (o *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Grant, *rs.SyncOpResults, error) {
	var grants []*v2.Grant

	// Note: Redis Enterprise Service API doesn't support pagination.
	userAnnos, err := o.GetUsers(ctx)
	if err != nil {
		return nil, nil, err
	}
	roleAnnos, err := o.GetRoles(ctx)
	if err != nil {
		return nil, nil, err
	}
	userAnnos = append(userAnnos, roleAnnos...)

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

	return grants, &rs.SyncOpResults{Annotations: userAnnos}, nil
}

func newRoleBuilder(c *client.RedisClient) *roleBuilder {
	return &roleBuilder{
		resourceType: roleResourceType,
		client:       c,
	}
}

func (o *roleBuilder) GetUsers(ctx context.Context) (annotations.Annotations, error) {
	o.usersMutex.Lock()
	defer o.usersMutex.Unlock()

	if o.users != nil || len(o.users) > 0 {
		return nil, nil
	}

	users, annos, err := o.client.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	o.users = users

	return annos, nil
}

func (o *roleBuilder) GetRoles(ctx context.Context) (annotations.Annotations, error) {
	o.rolesMutex.Lock()
	defer o.rolesMutex.Unlock()

	if o.roles == nil {
		o.roles = make(map[int]client.Role)
	}

	if len(o.roles) > 0 {
		return nil, nil
	}

	roles, annos, err := o.client.ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		o.roles[role.UID] = role
	}

	return annos, nil
}
