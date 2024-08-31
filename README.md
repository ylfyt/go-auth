# go-auth

### How to run

```bash
# Install golang-migrate
go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database sqlite://my.db up
```
