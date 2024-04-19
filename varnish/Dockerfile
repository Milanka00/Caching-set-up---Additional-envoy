# Dockerfile

# Use the official Varnish image
FROM varnish:stable

# Copy the default VCL file into the container
COPY default.vcl /etc/varnish/default.vcl

CMD ["-s", "default=malloc,30M"]