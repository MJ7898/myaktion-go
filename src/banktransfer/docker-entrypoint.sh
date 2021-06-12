#!/bin/sh

# Abort on any error (including if wait-for-it fails).
set -e

# Wait for kafka
if [ -n "$KAFKA_CONNECT" ]; then
  /go/src/app/wait-for-it.sh "$KAFKA_CONNECT" -t 120
  # we need to wait a bit more, because sometimes kafka is not yet ready
  sleep 1
fi

# Run the main container command.
exec "$@"