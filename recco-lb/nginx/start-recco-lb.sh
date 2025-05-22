#!/bin/sh
envsubst '$RECCO_IP $RECCO_LB_PORT $RECCO_SERVER_PORT $RECCO_EMBED_PORT $RECCO_DB_PORT' \
    < /recco-lb/nginx.conf.template > /etc/nginx/nginx.conf
nginx -g 'daemon off;'
nginx -s reload
