# Go RSS Aggregator

RSS Aggregator project to learn Go

## Development

To create and run application:

```bash
go build && ./go-rss-aggregator
```

Required environment variables:

```bash
PORT=<some_port_number>
DB_URL=<some_db_conn_string>
```

### Database Setup

You can run the API againts a container or local instance of Postgres.

To setup docker container:
Pull image or let `docker run` do so automatically and issue:

`docker run --name <some_container_name> -e POSTGRES_PASSWORD=<some_password> -p 5432:5432 -d postgres:15.3`

You can pull the image tag of your choice, at the time of running 15.3 was chosen.

Ensure you've created the database and optionally setup a user/owner with the correct permissions.

To get your database up to date, apply migrations by `cd`-ing into `sql/migrations` and running the following:

```bash
goose postgres <db_conn_string> up
```

This is still a work in progress, for the time being you can test endpoints via HTTPie, curl, or your tool of choice.
