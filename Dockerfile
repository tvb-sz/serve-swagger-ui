FROM alpine:3.14.2
WORKDIR /srv/
VOLUME /srv/runtime
COPY serve-swagger-ui /srv/serve-swagger-ui
COPY conf.toml /srv/conf.toml
RUN chmod +x /srv/serve-swagger-ui
EXPOSE 9080
