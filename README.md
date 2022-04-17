# Bnkr

The Kubernetes backup & migration solution that designed for human beings!

## Current Roadmap

- [x] Support MongoDB backup & restore.
- [ ] Accept kubeconfig(currently supported in the migrations)
- [x] Support PostgreSQL backup & restore.
- [x] Add migrations feature.
- [ ] Support S3 compilable object storages not just S3.
- [x] Helm chart.

## What I can do with Bnkr?

Bnkr goal is to backup & migrate only your valuable data inside Kubernetes cluster not the whole cluster, any contribution in this regard is more than welcome.

## Installation

Bnkr itself is a single binary application however it depends on some other tools and applications, for example, it uses `mysqldump` to backup MySQL databases so it is better to use Bnkr official image because you will guarantee that all the dependencies already exist.

The easiest way to install BNKR is to use [the official Helm chart](https://github.com/TryBnkr/helm-charts/tree/main/charts/bnkr).

BNKR binary is basically a web server that listens on any port you specify to it, here are some of the environment variables that you can pass to BNKR:

```
PORT="5000"
USERNAME="John Doe"
USERPASSWORD=StrongPassword
USEREMAIL=me@example.com
SETUP="true"
PRODUCTION="true"
DB_HOST=postgreshost
DB_USER=bnkr
DB_PASSWORD=StrongPassword
DB_NAME=bnkr
DB_PORT="5432"
DB_TIMEZONE=Europe/Istanbul
DB_SSLMODE=disable
```

You can run BNKR as deployment and pass a service account to it with enough capabilities that allow it for example to access the data that it will backup.

## Troubleshooting

In almost all the cases that we faced most problems came from incorrect permissions, please make sure the for example the DB user has the capabilities to dump and restore the DB, the SSH user has the capabilities to create/delete files, and so on.

## License

Copyright (c) 2022 [Mohammed Al-Mahdawi](https://al-mahdawi.is/)

Licensed under the **MIT** license.