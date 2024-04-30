FROM golang:1.22-alpine3.19

# Copy the binary that goreleaser built
COPY fleetdb /fleetdb

# Run the web service on container startup.
ENTRYPOINT ["/fleetdb"]
CMD ["serve"]
