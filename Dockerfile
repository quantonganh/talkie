FROM alpine:3.15
RUN apk add --no-cache ca-certificates
COPY talkie .
EXPOSE 8080
ENTRYPOINT [ "./talkie" ]