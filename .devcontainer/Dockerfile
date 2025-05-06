FROM debian:bullseye-slim
CMD ["/bin/bash"]
RUN /bin/sh -c "apt update && apt install -y --no-install-recommends ca-certificates curl wget && rm -rf /var/lib/apt/lists/*"
RUN /bin/sh -c "apt update && apt install -y --no-install-recommends git openssh-client procps && rm -rf /var/lib/apt/lists/*"
RUN /bin/sh -c "apt update && apt install -y --no-install-recommends g++ gcc libc6-dev make pkg-config && rm -rf /var/lib/apt/lists/*"
ENV GOLANG_VERSION=1.22.3
ENV GOLANG_DOWNLOAD_URL=https://golang.org/dl/go1.22.3.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256=6a63fef0e050146f275bf02a0896badfe77c11b6f05499bb647e7bd613a45a10
RUN /bin/sh -c "curl -fsSL \"$GOLANG_DOWNLOAD_URL\" -o golang.tar.gz && echo \"$GOLANG_DOWNLOAD_SHA256 golang.tar.gz\" | sha256sum -c - && tar -C /usr/local -xzf golang.tar.gz && rm golang.tar.gz"
ENV GOPATH=/go
ENV PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
RUN /bin/sh -c "mkdir -p \"$GOPATH/src\" \"$GOPATH/bin\" && chmod -R 777 \"$GOPATH\""
WORKDIR /go
