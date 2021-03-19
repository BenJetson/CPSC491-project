#!/bin/bash

banner() {
    echo
    echo "=== $* ==="
    echo
}

banner PREPARATION
make testdb-background

SHOULD_FAIL=0

cleanup() {
    banner DATABASE LOGS
    docker-compose logs -t testdb

    banner FLYWAY LOGS
    docker-compose logs -t flyway-testdb

    banner CLEANUP
    make stop-clean-testdb

    banner RESULTS
    if [ $SHOULD_FAIL -ne 0 ]; then
        echo "FAIL: Some portion of this job raised a fatal error. Check logs."
    else
        echo "PASS: No error conditions were detected."
    fi
    echo

    exit $SHOULD_FAIL
}

ctrlc() {
    echo
    echo
    echo "ERROR: Testing was canceled by user! Will fail."
    echo

    SHOULD_FAIL=1

    cleanup
}

trap cleanup EXIT HUP QUIT TERM
trap ctrlc INT

# Make errors NOT halt the script
set -e

RETRY_COUNT=4
WAIT_LENGTH=15

banner DB PARAMETERS

DB_PARAMS="
    host=localhost
    port=5432
    dbname=app_template
    user=testuser
    password=testpass
"

LATEST_MIGRATION=$(
    find db/migrations -type f -print0 |
    xargs -0 -I{} basename "{}" |
    sed -e "s/^V//" -e "s/_.*//" -e "s/^0//" -e "s/[^0-9]//g" |
    sort |
    tail -n 1
)

echo "Database migrations have been detected."
echo "Most current DB migration is version $LATEST_MIGRATION."
echo
echo "Waiting for flyway to apply version $LATEST_MIGRATION before proceeding."

WAIT_QUERY=" \
    SELECT TRUE
    FROM flyway_schema_history
    WHERE
        installed_rank = $LATEST_MIGRATION
        AND success = TRUE
"

DB_CHECK="psql '$DB_PARAMS' -c 'SELECT TRUE' 1>/dev/null 2>&1"
DB_READY="psql '$DB_PARAMS' -c '$WAIT_QUERY' 2>/dev/null | grep row"

banner WAIT FOR DB
while [[ $(eval "$DB_READY") != "(1 row)" ]]; do
    if [ $RETRY_COUNT -lt 1 ]; then
        SHOULD_FAIL=1

        echo
        echo "Database connection failed PERMANENTLY!"
        echo "ERROR: Test harness cannot start without database."

        exit 1
    elif eval "$DB_CHECK"; then
        echo "Connected to database, but migrations have not finished yet."
    fi

    ((RETRY_COUNT--))
    echo "Database is not ready. Will retry $RETRY_COUNT more time(s)."
    echo "Attempting reconnection in $WAIT_LENGTH seconds."

    sleep $WAIT_LENGTH

    ((WAIT_LENGTH+=5))

    echo
    echo "Retrying..."
done;

echo "Database is ready! Tests may start."

banner RUN TESTS

pushd .
cd ./go

gotestsum \
    --format testname \
    --junitfile gotest-results.xml \
    --junitfile-testsuite-name relative \
    --junitfile-testcase-classname relative \
    || SHOULD_FAIL=$? # Use test exit code for failure status.

popd
