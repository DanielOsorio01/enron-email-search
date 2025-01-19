#!/bin/sh

echo "Setting backend URL to: $BACKEND_URL"

# Replace the placeholder in the config.js file with the actual value
cat <<EOF > /usr/share/nginx/html/config.js
window.APP_CONFIG = {
  backend_URL: "$BACKEND_URL"
};
EOF

# Start Nginx
nginx -g "daemon off;"
