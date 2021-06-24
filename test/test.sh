#!/bin/sh

export TEST_REDIS_HOST=localhost:6379

docker-compose up -d

echo "waiting for redis..."
./wait-for.sh localhost:6379 -- sleep 1s

echo "dependencies started"

go test -v -count=1 -timeout=15s -race ../...
test_exit=$?

docker-compose down -v
docker-compose rm -s -f -v

exit $test_exit
