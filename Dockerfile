FROM --platform=linux/amd64 debian:stable-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    apt-utils \
    ca-certificates \
    # Add any other necessary packages here \
    && rm -rf /var/lib/apt/lists/*

COPY notely /usr/bin/notely

CMD ["notely"]