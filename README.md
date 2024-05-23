# Exchange GRPC server application

#### Healthcheck via cURL
```bash
curl localhost:8080/health
```

#### Environment parameters example:
```dotenv
DATABASE_URL=postgres://admin:admin@postgresql:5432/exchange?sslmode=disable
GRPC_ADDRESS=localhost:8081
```

Also you can set `DATABASE_URL` via `--database_url={url}` flag.

