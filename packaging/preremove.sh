#!/bin/sh
set -e

if [ "$1" = "remove" ] || [ "$1" = "deconfigure" ]; then
  systemctl stop transmission-otel.service >/dev/null 2>&1 || true
  systemctl disable transmission-otel.service >/dev/null 2>&1 || true
fi
