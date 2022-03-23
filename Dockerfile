FROM alpine:3.15

WORKDIR /srv/

# your swagger json file can mappable to this directory use default conf.toml Swagger.Path value ./runtime
VOLUME /srv/runtime

# please use `CGO_ENABLED=0 go build -ldflags "-s -w" -o serve-swagger-ui main.go` at first to build binary file
COPY serve-swagger-ui /srv/serve-swagger-ui

# then use `./serve-swagger-ui output_conf >> conf.toml` generate config file
# and edit conf.toml to set your service parameters
COPY conf.toml /srv/conf.toml

# set executable
RUN chmod +x /srv/serve-swagger-ui

# default port
EXPOSE 9080

CMD [ "/srv/serve-swagger-ui", "--config", "conf.toml"]
