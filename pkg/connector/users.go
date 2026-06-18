package connector

import (
	"context"
	"strconv"
	"strings"

	"github.com/conductorone/baton-redis/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type userBuilder struct {
	resourceType *v2.ResourceType
	client       *client.RedisClient
}

func (o *userBuilder) ResourceType(_ context.Context) *v2.ResourceType {
	return userResourceType
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (o *userBuilder) List(ctx context.Context, _ *v2.ResourceId, _ rs.SyncOpAttrs) ([]*v2.Resource, *rs.SyncOpResults, error) {
	var resources []*v2.Resource

	// Note: Redis Enterprise Service API doesn't support pagination.
	users, _, err := o.client.ListUsers(ctx)
	if err != nil {
		return nil, nil, err
	}

	for _, user := range users {
		userCopy := user
		userResource, err := parseIntoUserResource(ctx, &userCopy, nil)
		if err != nil {
			return nil, nil, err
		}
		resources = append(resources, userResource)
	}

	return resources, nil, nil
}

func parseIntoUserResource(_ context.Context, user *client.User, parentResourceID *v2.ResourceId) (*v2.Resource, error) {
	var userStatus = v2.UserTrait_Status_STATUS_ENABLED

	profile := map[string]interface{}{
		"user_id":         user.UID,
		"username":        user.Name,
		"email":           user.Email,
		"management_role": user.Role,
		"role_uids":       parseRoleUIDs(user.RoleUIDs),
	}

	userTraits := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(userStatus),
		rs.WithUserLogin(user.Name),
	}

	displayName := user.Name

	ret, err := rs.NewUserResource(
		displayName,
		userResourceType,
		user.UID,
		userTraits,
		rs.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Entitlement, *rs.SyncOpResults, error) {
	return nil, nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(ctx context.Context, _ *v2.Resource, _ rs.SyncOpAttrs) ([]*v2.Grant, *rs.SyncOpResults, error) {
	return nil, nil, nil
}

func newUserBuilder(c *client.RedisClient) *userBuilder {
	return &userBuilder{
		resourceType: userResourceType,
		client:       c,
	}
}

func parseRoleUIDs(roles []int) string {
	var rolesStr []string
	for _, roleUID := range roles {
		rolesStr = append(rolesStr, strconv.Itoa(roleUID))
	}
	return strings.Join(rolesStr, ",")
}
