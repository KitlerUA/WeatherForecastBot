FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ADD main /
ADD config.json /
ADD chatslocation.json /
CMD ["/main"]