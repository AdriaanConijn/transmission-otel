#!/bin/sh
set -e

systemctl daemon-reload >/dev/null 2>&1 || true

if [ "$1" = "purge" ]; then
  userdel transmission-otel >/dev/null 2>&1 || true
fi
