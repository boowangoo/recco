FROM nginx:1.27.5-alpine

ENV NGINX_HOST=0.0.0.0
ENV NGINX_PORT=80

VOLUME /recco-lb

# Modify the default nginx configuration at runtime
CMD ["/bin/sh", "/recco-lb/start-recco-lb.sh"]

EXPOSE 80
