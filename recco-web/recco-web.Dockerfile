FROM nginx:1.27.5-alpine

ENV NGINX_HOST=0.0.0.0
ENV NGINX_PORT=80

VOLUME //usr/share/nginx/html
VOLUME /etc/nginx/conf.d/recco-web.conf
VOLUME /var/log/nginx

EXPOSE 80
