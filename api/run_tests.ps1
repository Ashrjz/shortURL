$env:JWT_SECRET="test-secret"
$env:DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=urlshortener_test sslmode=disable"
$env:REDIS_ADDR="localhost:6379"

go test -v ./...