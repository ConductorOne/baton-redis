FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-redis"]
COPY baton-redis /