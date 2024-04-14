# Service layer

This layer contains handlers which interact directly with the service. Typically, other layers of our service, including service layer and transport layer, would interact with the service through this layer if they need to.

## Migrations

- To compare the models with service schema, add the models to `/scripts/loader.go`. This will help generate the migrations for your models.
- To generate a new migration, run `./scripts/migrate.sh [name_of_your_migration]`. For example: `./scripts/migrate.sh init` will generate a new migration called "init" in `./migrations` directory.
- To compare the state/status of the migrations against the service schema, run `./scripts/status.sh`. This will print which migrations are pending to be applied and which ones have been applied.
- To apply all the pending migrations, run `./scripts/apply.sh`.

## Testing

### Unit / Whitebox Tests

All the essential unit tests to be covered:

- [x] Create a new record in the service.
- [x] Retrieve and list all the records from the service with supported filters.
- [x] Get a record from the service using it's ID.
- [x] Update a record with new options in the service.
- [x] Delete a record from the service using it's ID.

### Integration / Blackbox Tests

To write itnegration tests, you would typically want to mock this layer's interfaces and consume them outside the package.

**To generate mock files of the service interface, use the following commands:**

1. Install mockgen.
    ```
    go install go.uber.org/mock/mockgen@latest
    ```
1. Generate mocks.
    ```
    go generate ./...
    ```
    Or:
    ```
    mockgen -destination=mock.go -source=service.go -package=service
    ```

This will generate the file `mock.go` which will contain your mock service. You can import it in your tests with:

```
func TestFoo(t *testing.T) {
  ctrl := gomock.NewController(t)

  m := NewMockDatabase(ctrl)

  // Asserts that the first and only call to Bar() is passed 99.
  // Anything else will fail.
  m.
    EXPECT().
    Bar(gomock.Eq(99)).
    Return(101)

  YourUnitTest(m)
}
```