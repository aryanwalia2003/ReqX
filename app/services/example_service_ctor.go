package services

// NewExampleService constructs the service. Wire any internal/* dependencies
// (e.g. a history.DB handle, storage paths) here once a real service needs them.
func NewExampleService() *ExampleService {
	return &ExampleService{}
}
