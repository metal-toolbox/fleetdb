FROM alpine:3.22.0

# Copy the binary that goreleaser built
COPY fleetdb /fleetdb

# Run the web service on container startup.
ENTRYPOINT ["/fleetdb"]
CMD ["serve"]
