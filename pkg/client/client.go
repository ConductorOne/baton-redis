package client

import (
	"context"
	encoding "encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/ratelimit"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

const (
	getUsers    = "/v1/users"
	getRoles    = "/v1/roles"
	getCluster  = "/v1/cluster"
	getRoleById = "/v1/roles/%v"
)

type RedisClient struct {
	Username    string
	Password    string
	ClusterHost string
	APIPort     string
	wrapper     *uhttp.BaseHttpClient
}

func New(ctx context.Context, redisClient *RedisClient) (*RedisClient, error) {
	var (
		username    = redisClient.Username
		password    = redisClient.Password
		clusterHost = redisClient.ClusterHost
		apiPort     = redisClient.APIPort
	)

	options := []uhttp.Option{
		uhttp.WithLogger(true, ctxzap.Extract(ctx)),
	}

	httpClient, err := uhttp.NewClient(ctx, options...)
	if err != nil {
		return nil, err
	}

	cli, err := uhttp.NewBaseHttpClientWithContext(context.Background(), httpClient)
	if err != nil {
		return nil, err
	}

	client := RedisClient{
		wrapper:     cli,
		Username:    username,
		Password:    password,
		ClusterHost: clusterHost,
		APIPort:     apiPort,
	}

	return &client, nil
}

func NewClient(username, password, clusterHost string, apiPort string, httpClient ...*uhttp.BaseHttpClient) *RedisClient {
	var wrapper = &uhttp.BaseHttpClient{}
	if httpClient != nil || len(httpClient) != 0 {
		wrapper = httpClient[0]
	}
	return &RedisClient{
		wrapper:     wrapper,
		Username:    username,
		Password:    password,
		ClusterHost: clusterHost,
		APIPort:     apiPort,
	}
}

func (c *RedisClient) ListUsers(ctx context.Context) ([]User, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res []User

	annotation, err := c.getResourcesFromAPI(ctx, getUsers, &res)
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resources: %s", err))
		return nil, nil, err
	}

	return res, annotation, nil
}

func (c *RedisClient) ListClusters(ctx context.Context) ([]Cluster, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res Cluster

	annotation, err := c.getResourcesFromAPI(ctx, getCluster, &res)
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resources: %s", err))
		return nil, nil, err
	}

	return []Cluster{res}, annotation, nil
}

func (c *RedisClient) ListRoles(ctx context.Context) ([]Role, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res []Role

	annotation, err := c.getResourcesFromAPI(ctx, getRoles, &res)
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resources: %s", err))
		return nil, nil, err
	}

	return res, annotation, nil
}

func (c *RedisClient) GetRoleDetails(ctx context.Context, roleUID string) (Role, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	var res Role

	annotation, err := c.getResourcesFromAPI(ctx, fmt.Sprintf(getRoleById, roleUID), &res)
	if err != nil {
		l.Error(fmt.Sprintf("Error getting resources: %s", err))
		return res, nil, err
	}

	return res, annotation, nil
}

func (c *RedisClient) getResourcesFromAPI(
	ctx context.Context,
	urlEndpoint string,
	res any,
) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	urlAddress, err := url.Parse(c.ClusterHost + ":" + c.APIPort + urlEndpoint)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating url: %s", err))
		return nil, err
	}

	_, annotation, err := c.doRequest(ctx, http.MethodGet, urlAddress, &res)

	if err != nil {
		return nil, err
	}

	return annotation, nil
}

func (c *RedisClient) doRequest(
	ctx context.Context,
	method string,
	urlAddress *url.URL,
	res interface{},
) (http.Header, annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	var (
		resp *http.Response
		err  error
	)

	authorizationToken := encoding.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.Username, c.Password)))

	req, err := c.wrapper.NewRequest(
		ctx,
		method,
		urlAddress,
		uhttp.WithContentTypeJSONHeader(),
		uhttp.WithAcceptJSONHeader(),
		uhttp.WithHeader("Authorization", "Basic "+authorizationToken),
	)

	l.Info(fmt.Sprintf("Request %v", req))

	if err != nil {
		return nil, nil, err
	}

	switch method {
	case http.MethodGet, http.MethodPut, http.MethodPost:
		var doOptions []uhttp.DoOption
		if res != nil {
			doOptions = append(doOptions, uhttp.WithResponse(&res))
		}
		resp, err = c.wrapper.Do(req, doOptions...)
		if resp != nil {
			defer resp.Body.Close()
		}
	case http.MethodDelete:
		resp, err = c.wrapper.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
	}

	if err != nil {
		return nil, nil, err
	}

	annotation := annotations.Annotations{}
	if resp != nil {
		if desc, err := ratelimit.ExtractRateLimitData(resp.StatusCode, &resp.Header); err == nil {
			annotation.WithRateLimiting(desc)
		} else {
			return nil, annotation, err
		}

		return resp.Header, annotation, nil
	}

	return nil, nil, err
}
