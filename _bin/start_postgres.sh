#!/bin/sh

set -e

. "$(dirname "${0}")/exists.sh"
. "$(dirname "${0}")/require_env_var.sh"

exists "docker"
requireEnvVar "DB_NETWORK_NAME"
NETWORK_EXISTS=$(docker network ls --quiet --filter "name=${DB_NETWORK_NAME}")
if [ -z "${NETWORK_EXISTS}" ]; then
  docker network create --internal "${DB_NETWORK_NAME}" > /dev/null
  echo "Network ${DB_NETWORK_NAME} created."
else
  echo "Network ${DB_NETWORK_NAME} is already running."
fi

requireEnvVar "DB_CONTAINER_NAME"
STATE_RUNNING=$(docker inspect --format "{{.State.Running}}" "${DB_CONTAINER_NAME}" 2> /dev/null || true)
if [ "${STATE_RUNNING}" = "true" ]; then
  echo "Container ${DB_CONTAINER_NAME} already running."
  exit
fi

STATE_STATUS=$(docker inspect --format "{{.State.Status}}" "${DB_CONTAINER_NAME}" 2> /dev/null || true)
if [ "${STATE_STATUS}" = "exited" ]; then
  echo "Container ${DB_CONTAINER_NAME} is stopped, removing."
  docker rm --force "${DB_CONTAINER_NAME}" > /dev/null
fi

# Get the absolute path to the config file (for Docker)
exists "python"
CONF_DIR="$(dirname "${0}")/../_docker"
# macOS workaround for `readlink`; see https://stackoverflow.com/q/3572030/1068170
CONF_DIR=$(python -c "import os; print(os.path.realpath('${CONF_DIR}'))")

requireEnvVar "DB_HOST"
requireEnvVar "DB_PORT"
requireEnvVar "DB_SUPERUSER_NAME"
requireEnvVar "DB_SUPERUSER_USER"
requireEnvVar "DB_SUPERUSER_PASSWORD"
docker run \
  --detach \
  --hostname "${DB_HOST}" \
  --publish "${DB_PORT}:5432" \
  --name "${DB_CONTAINER_NAME}" \
  --env POSTGRES_DB="${DB_SUPERUSER_NAME}" \
  --env POSTGRES_USER="${DB_SUPERUSER_USER}" \
  --env POSTGRES_PASSWORD="${DB_SUPERUSER_PASSWORD}" \
  --env POSTGRES_INITDB_ARGS="--auth-host=scram-sha-256 --auth-local=scram-sha-256" \
  --mount type=tmpfs,destination=/var/lib/postgresql/data \
  --volume "${CONF_DIR}/postgresql.conf":/etc/postgresql/postgresql.conf \
  --volume "${CONF_DIR}/pg_hba.conf":/etc/postgresql/pg_hba.conf \
  postgres:13.3-alpine3.14 \
  -c "config_file=/etc/postgresql/postgresql.conf" \
  -c "hba_file=/etc/postgresql/pg_hba.conf" \
  > /dev/null

echo "Container ${DB_CONTAINER_NAME} started on port ${DB_PORT}."

# NOTE: It's crucial to use `docker network connect` vs. starting the container
#       with `docker run --network`. Sine the network is `--internal` any
#       use of `--publish` will be ignored if `--network` is also provided
#       to `docker run`.
docker network connect "${DB_NETWORK_NAME}" "${DB_CONTAINER_NAME}"
echo "Container ${DB_CONTAINER_NAME} added to network ${DB_NETWORK_NAME}."

##########################################################
## Don't exit until `pg_isready` returns 0 in container ##
##########################################################

# NOTE: This is used strictly for the status code to determine readiness.
pgIsReadyFull() {
  PGPASSWORD="${DB_SUPERUSER_PASSWORD}" pg_isready \
    --dbname "${DB_SUPERUSER_NAME}" \
    --username "${DB_SUPERUSER_USER}" \
    --host "${DB_HOST}" \
    --port "${DB_PORT}"
}

pgIsReady() {
  pgIsReadyFull > /dev/null 2>&1
}

exists "pg_isready"
# Cap at 50 retries / 5 seconds (by default).
if [ -z "${PG_ISREADY_RETRIES}" ]; then
  PG_ISREADY_RETRIES=50
fi
i=0; while [ ${i} -le ${PG_ISREADY_RETRIES} ]
do
  if pgIsReady
  then
    echo "Container ${DB_CONTAINER_NAME} accepting Postgres connections."
    break
  fi
  i=$((i+1))
  sleep "0.1"
done

if [ ${i} -ge ${PG_ISREADY_RETRIES} ]; then
  echo "Container ${DB_CONTAINER_NAME} not accepting Postgres connections."
  echo "  pg_isready: $(pgIsReadyFull)"
  exit 1
fi

# Run the superuser migrations
. "$(dirname "${0}")/superuser_migrations_postgres.sh"
