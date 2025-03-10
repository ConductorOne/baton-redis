![Baton Logo](./baton-logo.png)

# `baton-redis` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-redis.svg)](https://pkg.go.dev/github.com/conductorone/baton-redis) ![main ci](https://github.com/conductorone/baton-redis/actions/workflows/main.yaml/badge.svg)

`baton-redis` is a connector for built using the [Baton SDK](https://github.com/conductorone/baton-sdk).

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-redis
baton-redis
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_DOMAIN_URL=domain_url -e BATON_API_KEY=apiKey -e BATON_USERNAME=username ghcr.io/conductorone/baton-redis:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-redis/cmd/baton-redis@main

baton-redis

baton resources
```

# Data Model

`baton-redis` will pull down information about the following resources:
- Users
- Clusters
- Roles

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually
building spreadsheets. We welcome contributions, and ideas, no matter how
small&mdash;our goal is to make identity and permissions sprawl less painful for
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-redis` Command Line Usage

```
baton-redis

Usage:
  baton-redis [flags]
  baton-redis [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --api-port string              The Redis Enterprise admin port ($BATON_API_PORT) (default "9443")
      --client-id string             The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string         The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
      --cluster-host                 required: The cluster host for your Redis Enterprise Serivice ($BATON_CLUSTER_HOST)
  -f, --file string                  The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                         help for baton-redis
      --log-format string            The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string             The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --password                     required: Redis Enterprise Sign In password ($BATON_PASSWORD)
  -p, --provisioning                 If this connector supports provisioning, this must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --ticketing                    This must be set to enable ticketing support ($BATON_TICKETING)
      --username                     required: Redis Enterprise Sign In Email/Username ($BATON_USERNAME)
  -v, --version                      version for baton-redis

Use "baton-redis [command] --help" for more information about a command.
```
