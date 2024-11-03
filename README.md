## Generate SQLC
```shell
sqlc generate
```

## Atlas create migration
```shell
sudo atlas migrate diff initial \
--dir "file://internal/sqlc/migrations" \
--to "file://internal/sqlc/schema.sql" \
--dev-url "docker://postgres/16/url-shortener-db" \
--format '{{ sql . " " }}'
```

## Atlas apply migration
```shell
sudo atlas migrate apply \
--url "postgres://url-shortener-db-user:url-shortener-db-pass@localhost:5432/url-shortener-db?sslmode=disable" \
--dir "file://internal/sqlc/migrations"
```

## Generate ECDSA key pair
```shell
openssl ecparam -name prime256v1 -genkey -noout -out private.pem
openssl ec -in private.pem -pubout -out public.pem
```