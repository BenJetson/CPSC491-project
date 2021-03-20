#!/bin/bash

banner() {
    echo
    echo "=== $* ==="
    echo
}

SHOULD_FAIL=0

onfail() {
    last_status=$?

    if [ $SHOULD_FAIL -eq 0 ] || [ $last_status -ne 0 ]; then
        echo
        echo "<!> FAIL: the last command returned an exit code of $last_status,"
        echo "          which will cause this build to fail after all steps."
        echo
        SHOULD_FAIL=$last_status
    fi
}

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

banner PREPARATION
make testdb-background || onfail
mkdir go/results

# Make errors NOT halt the script
set -e

banner DB PARAMETERS

RETRY_COUNT=4
WAIT_LENGTH=15

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

echo "Database is ready; migrations applied! Tests may start."

banner RUN TESTS

pushd .
cd ./go

gotestsum \
    --format testname \
    --junitfile results/gotest-results.xml \
    --junitfile-testsuite-name relative \
    --junitfile-testcase-classname relative \
    -- -coverprofile=results/cover.out ./... \
    || onfail

banner COVERAGE REPORT

go tool cover -func results/cover.out | \
    sed -r \
        -e "s/^github.com\/BenJetson\/CPSC491-project\/go\///g" \
        -e "s/total:\t\t\t\t\t/total:/" \
    || onfail
go tool cover \
    -html=results/cover.out \
    -o results/cover.html || onfail
gocov convert results/cover.out | \
    gocov-xml > results/cover.xml || onfail

popd
