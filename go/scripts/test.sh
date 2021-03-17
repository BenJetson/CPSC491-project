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

DB_PARAMS=" \
    host=localhost \
    port=5432 \
    dbname=app_template \
    user=testuser \
    password=testpass \
"

banner WAIT FOR DB
while ! psql "$DB_PARAMS" -c "select 1" > /dev/null 2>&1; do
    if [ $RETRY_COUNT -lt 1 ]; then
        SHOULD_FAIL=1

        echo
        echo "Database connection failed PERMANENTLY!"
        echo "ERROR: Test harness cannot start without database."

        exit 1
    fi

    ((RETRY_COUNT--))
    echo "Database is not ready. Will retry $RETRY_COUNT more time(s)."
    echo "Attempting reconnection in $WAIT_LENGTH seconds."

    sleep $WAIT_LENGTH

    ((WAIT_LENGTH+=5))

    echo
    echo "Retrying..."
done;

echo "Database is available! Tests may start."

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
