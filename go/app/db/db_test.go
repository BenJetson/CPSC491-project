package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testDB struct {
	*database
}

func getTestDBAdminHandle(t *testing.T) *sqlx.DB {
	cfg := Config{
		Host:       "localhost",
		Port:       5432,
		Username:   "testuser",
		Password:   "testpass",
		Database:   "app_template",
		DisableTLS: true,
	}

	adminHandle, err := sqlx.Connect("postgres", cfg.connectString())
	require.NoError(t, err, "failed to acquire testDB admin handle")

	return adminHandle
}

func newTestDB(t *testing.T) *testDB {
	var db testDB
	db.reset(t)

	return &db
}

func (db *testDB) reset(t *testing.T) {
	if db.database != nil {
		err := db.database.Close()
		require.NoError(t, err, "failed to close previous test database handle")
	}

	adminHandle := getTestDBAdminHandle(t)
	defer func() {
		err := adminHandle.Close()
		require.NoError(t, err, "failed to close admin handle")
	}()

	_, err := adminHandle.Exec(`DROP DATABASE IF EXISTS testdb`)
	require.NoError(t, err, "failed to destroy test database for reset")

	_, err = adminHandle.Exec(`CREATE DATABASE testdb TEMPLATE app_template`)
	require.NoError(t, err, "failed to clone test database for reset")

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	cfg := Config{
		Host:       "localhost",
		Port:       5432,
		Username:   "testuser",
		Password:   "testpass",
		Database:   "testdb",
		DisableTLS: true,
	}

	db.database, err = newDatabase(logger, cfg)
	require.NoError(t, err, "failed to open test database")

	// Database has some default data in it, but for test purposes we shall
	// clear out this default data to start with an empty database.
	_, err = db.Exec(`
		TRUNCATE TABLE person RESTART IDENTITY CASCADE;
		TRUNCATE TABLE organization RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err, "failed to clean test database")
}

func (db *testDB) cleanup(t *testing.T) {
	err := db.database.Close()
	require.NoError(t, err, "failed to close test database")
}

func (db *testDB) assertCountOf(
	t *testing.T,
	tableName string,
	expect int,
	whereClause string,
	params ...interface{},
) {

	var actual int

	// ATTENTION!!! The following query would be DANGEROUS in production because
	// injecting the table name and where clause directly into the query using
	// fmt.Sprintf may allow SQL injection attacks.
	//
	// This is here to facilitate unit tests only on an isolated test database,
	// therefore there is no real security impact to the application itself.
	//
	// But be warned ... this is an EXCEPTION. Look elsewhere for examples!
	db.Get(&actual, fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE %s
	`, tableName, whereClause), params...)

	description := "total"
	if whereClause != "TRUE" {
		description = "with given filter"
	}

	assert.Equalf(
		t, expect, actual,
		"table '%s' ought to have %d row(s) %s but found %d",
		tableName, expect, description, actual,
	)
}

func (db *testDB) assertCount(t *testing.T, tableName string, expect int) {
	db.assertCountOf(t, tableName, expect, "TRUE")
}

func assertEqualJSON(t *testing.T, expect, actual interface{}) {
	actualJSONb, err := json.MarshalIndent(actual, "", "    ")
	require.NoError(t, err)

	expectJSONb, err := json.MarshalIndent(expect, "", "    ")
	require.NoError(t, err)

	actualJSON := string(actualJSONb)
	expectJSON := string(expectJSONb)

	assert.Equal(t, expectJSON, actualJSON)
}
