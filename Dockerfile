FROM scratch
COPY app.cwasm /app.cwasm
ENTRYPOINT ["/app.cwasm"]
