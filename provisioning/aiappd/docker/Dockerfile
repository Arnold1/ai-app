FROM debian:9.5
RUN apt-get update && apt-get install -y \
    net-tools \
    ca-certificates
RUN mkdir -p /apps/aiappd/bin/
COPY . /apps/aiappd/
RUN mv /apps/aiappd/aiappd /apps/aiappd/bin/
RUN mv /apps/aiappd/view /apps/aiappd/bin/
WORKDIR /apps/aiappd/bin
ENTRYPOINT ["./aiappd"]