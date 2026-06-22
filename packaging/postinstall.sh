#!/bin/sh
set -e

if ! getent passwd transmission-otel >/dev/null; then
  useradd --system --no-create-home --shell /usr/sbin/nologin transmission-otel
fi

systemctl daemon-reload >/dev/null 2>&1 || true
systemctl enable transmission-otel.service >/dev/null 2>&1 || true

if systemctl is-active --quiet transmission-otel.service; then
  systemctl restart transmission-otel.service >/dev/null 2>&1 || true
else
    echo "transmission-otel installed. Edit /etc/transmission-otel/transmission-otel.env, then run:"
    echo "  systemctl start transmission-otel"
