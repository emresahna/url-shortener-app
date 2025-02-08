## Generate SQLC
```shell
sqlc generate
```

## Atlas create migration
```shell
sudo atlas migrate diff initial \
--dir "file://deployments/migrations" \
--to "file://internal/sqlc/schema.sql" \
--dev-url "docker://postgres/16/url-shortener-db" \
--format '{{ sql . " " }}'
```

## Atlas apply migration
```shell
sudo atlas migrate apply \
--url "postgres://url-shortener-db-user:url-shortener-db-pass@localhost:5432/url-shortener-db?sslmode=disable" \
--dir "file://deployments/migrations"
```

## Generate ECDSA key pair
```shell
openssl ecparam -name prime256v1 -genkey -noout -out private.pem
openssl ec -in private.pem -pubout -out public.pem
```

## Start postgres locally
```shell
sudo docker run --name url-shortener-postgres -p "5432:5432" -e POSTGRES_PASSWORD=pass -e POSTGRES_USER=root -e POSTGRES_DB=url-shortener -d postgres:latest
```

## Start redis locally
```shell
sudo docker run --name url-shortener-redis -p "6379:6379" -d redis:latest ["redis-server", "--notify-keyspace-events", "Ex"]
```