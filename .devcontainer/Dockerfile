FROM debian:bullseye-slim
CMD ["/bin/bash"]
RUN /bin/sh -c "apt update && apt install -y --no-install-recommends ca-certificates curl wget && rm -rf /var/lib/apt/lists/*"
RUN /bin/sh -c "apt update && apt install -y --no-install-recommends git openssh-client procps && rm -rf /var/lib/apt/lists/*"
RUN /bin/sh -c "apt update && apt install -y --no-install-recommends g++ gcc libc6-dev make pkg-config && rm -rf /var/lib/apt/lists/*"
ENV GOLANG_VERSION=1.11.13
ENV GOLANG_DOWNLOAD_URL=https://golang.org/dl/go1.11.13.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256=50fe8e13592f8cf22304b9c4adfc11849a2c3d281b1d7e09c924ae24874c6daa
RUN /bin/sh -c "curl -fsSL \"$GOLANG_DOWNLOAD_URL\" -o golang.tar.gz && echo \"$GOLANG_DOWNLOAD_SHA256 golang.tar.gz\" | sha256sum -c - && tar -C /usr/local -xzf golang.tar.gz && rm golang.tar.gz"
ENV GOPATH=/go
ENV PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
RUN /bin/sh -c "mkdir -p \"$GOPATH/src\" \"$GOPATH/bin\" && chmod -R 777 \"$GOPATH\""
WORKDIR /go
