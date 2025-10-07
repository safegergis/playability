# /certbot/certbot_init.sh
#!/bin/bash

DOMAINS="playablty.com www.playablty.com"
EMAIL="safegerg@gmail.com"

# Check if certificates already exist
if [ -d "/etc/letsencrypt/live/$DOMAINS" ]; then
  echo "Certificates already exist. Skipping initial request."
else
  echo "### Requesting certificates for $DOMAINS ###"
  certbot certonly --webroot -w /var/www/certbot \
    --email "$EMAIL" \
    -d "$DOMAINS" \
    --rsa-key-size 4096 \
    --non-interactive --agree-tos
fi

echo "### Starting renewal process ###"
# Start a cron job to renew certificates regularly
echo "0 3 * * * certbot renew --quiet && nginx -s reload" | crontab -
crond -f -d 8

