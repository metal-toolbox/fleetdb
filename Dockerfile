FROM alpine

# Copy the binary that goreleaser built
COPY fleetdb /fleetdb

# Run the web service on container startup.
ENTRYPOINT ["/fleetdb"]
CMD ["serve"]
