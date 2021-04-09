package app

import "github.com/pkg/errors"

// ErrNotFound may be returned by a DataStore implementation when the requested
// data could not be found inside the store.
//
// Suggested check for users of DataStore:
//
//     if errors.Is(err, app.ErrNotFound) {
//         // do something for the special not found case
//     }
//
// Suggested use by DataStore implementors:
//
//     return errors.Wrapf(app.ErrNotFound, "book #%d", id)
//
var ErrNotFound = errors.New("not found")
