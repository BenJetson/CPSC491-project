# Server Setup

This document details the setup of this server.

## Users

There are three user accounts of importance:

- `deployer`: used by Azure Pipelines to copy data to the server. Has `sudo`
  access only for enabling/disabling services.
- `app-prod`: system user for running the production API server service.
- `app-dev`: system user for running the development API server service.

There are a few groups:

- `application`: common group that should be associated with all application
  common files. `deployer`, `app-prod`, and `app-dev` are members.
- `secrets-prod`: group that allows READING the prod secrets file. `app-prod` is
  the sole member.
- `secrets-dev`: group that allows READING the dev secrets file. `app-dev` is
  the sole member.

## Secrets

Stored at `/usr/local/etc/secrets/XXX.env` where `XXX` is the tier.

Only `root` can make changes and is the owner. Via the `secrets-prod` and
`secrets-dev` groups, associated application service users may read these files.
Not world readable!

Contains the environment variables that will be injected into the service on
startup by systemd.

To change the secrets, modify the file and then restart the appropriate service.

## Binaries

Binaries for the app are stored in `/usr/local/bin`.

They shall be owned by `deployer:application` and be executable.

- `api-prod` is the production tier binary.
- `api-dev` is the development tier binary.

To update the binaries, you must:

- stop the associated service
- replace the file
- restart the associated service

## Services

There is a separate service for each tier's API server. These are managed by
systemd.

To make changes, edit `/etc/systemd/system/api-XXX.service` where `XXX` is the
tier. Then restart the service.

Useful commands:

- `sudo systemctl start api-XXX`: starts the service
- `sudo systemctl stop api-XXX`: stops the service
- `sudo systemctl restart api-XXX`: restarts the service
- `sudo systemctl reload api-XXX`: reloads the config from the .service file
- `sudo systemctl enable api-XXX`: enable the service so it starts on boot
- `sudo systemctl disable api-XXX`: disable the service for no on-boot start

## Web Server

This server is powered by NGINX.

Tiers:

- app.teamxiv.space is the PRODUCTION tier.
- dev.teamxiv.space is the DEVELOPMENT tier.

Requests to the path /api on each subdomain will be reverse-proxied to the
appropriate API service running on this box.

All other requests will attempt to be fulfilled by serving static files from
`/var/www/XXX` where `XXX` is the tier. Files/folders that do not exist will
return 404 of course.
