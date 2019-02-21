#!/bin/sh

set -eu

API_URL=${API_URL:-http://localhost:8080}
cat <<EOF > /usr/share/nginx/html/env.js
API_URL="${API_URL}";
EOF

exec nginx -g "daemon off;"
