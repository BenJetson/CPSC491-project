package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// tryRollback will attempt to roll back a transaction, logging on failure.
func (db *database) tryRollback(tx *sqlx.Tx) {
	err := tx.Rollback()

	if err != nil {
		err = errors.Wrap(err, "failed to rollback")
		db.logger.
			WithError(err).
			Error("failed to rollback after bad transaction")
	}
}

// cleanupTransaction runs after any transaction using Transact() completes. It
// shall rollback on panic or error, otherwise attempt to commit.
func (db *database) cleanupTransaction(tx *sqlx.Tx, err error) error {
	// Handle panicking goroutine by rolling back the transaction, then
	// continue panicking.
	if r := recover(); r != nil {
		db.tryRollback(tx)
		panic(r)
	}

	// Handle errors from the transaction handler by rolling back the
	// transaction, then returning the wrapped error.
	if err != nil {
		db.tryRollback(tx)
		return errors.Wrap(err, "transaction failed")
	}

	// Transaction handler was successful!
	// Attempt to commit the transaction, reporting an error on failure.
	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// Transact will begin a database transaction and execute the given transaction
// handler.
//
// Assuming no error is returned by the transaction handler, the transaction
// will be committed automatically upon return of the transaction handler.
//
// If the transaction handler panics, the transaction will rollback, then resume
// panicking up the stack.
//
// If the transaction handler returns a non-nil error value, the transaction
// will rollback and the causal error is returned.
//
// Adapted from: https://stackoverflow.com/a/23502629
//   and also: https://github.com/BenJetson/netwatch/blob/master/go/store/db.go
//
func (db *database) Transact(txHandler func(tx *sqlx.Tx) error) (err error) {
	// Start the transaction and receive a transaction handle.
	tx, err := db.Beginx()
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	// Add deferred handler to clean up this transaction.
	defer func() { err = db.cleanupTransaction(tx, err) }()

	// Call the transaction handler.
	err = txHandler(tx)

	// Remember, deferred cleanup ALWAYS runs after this.
	return
}
