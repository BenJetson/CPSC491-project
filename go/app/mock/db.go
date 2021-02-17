package mock

import "github.com/BenJetson/CPSC491-project/go/app"

// This is an assertion, which will cause the build to fail if the mock.DB type
// does not implement the app.DataStore interface.
var _ app.DataStore = (*DB)(nil)

// DB is a mock of interface app.DataStore. It shall implement all of its
// methods, each of which does nothing.
//
// This can be useful for tests, where this type may be embedded in a more
// specific mock with overrided methods.
type DB struct{}
