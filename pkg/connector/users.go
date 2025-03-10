package connector

import (
	"context"
	"strconv"
	"strings"

	"github.com/conductorone/baton-redis/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
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
func (o *userBuilder) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	var resources []*v2.Resource

	// Note: Redis Enterprise Service API doesn't support pagination.
	users, annotation, err := o.client.ListUsers(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for _, user := range users {
		userCopy := user
		userResource, err := parseIntoUserResource(ctx, &userCopy, nil)
		if err != nil {
			return nil, "", nil, err
		}
		resources = append(resources, userResource)
	}

	return resources, "", annotation, nil
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

	userTraits := []resource.UserTraitOption{
		resource.WithUserProfile(profile),
		resource.WithStatus(userStatus),
		resource.WithUserLogin(user.Name),
	}

	displayName := user.Name

	ret, err := resource.NewUserResource(
		displayName,
		userResourceType,
		user.UID,
		userTraits,
		resource.WithParentResourceID(parentResourceID),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Entitlements always returns an empty slice for users.
func (o *userBuilder) Entitlements(_ context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (o *userBuilder) Grants(ctx context.Context, _ *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
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
