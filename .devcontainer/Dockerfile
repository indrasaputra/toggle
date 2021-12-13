# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.209.5/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.16, 1.17, 1-bullseye, 1.16-bullseye, 1.17-bullseye, 1-buster, 1.16-buster, 1.17-buster
ARG VARIANT="1.17-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

ENV GO111MODULE=on
ENV BUF_VERSION=v1.0.0-rc8
ENV PROTOC_GEN_GO=v1.27.1
ENV PROTOC_GEN_GO_GRPC=v1.1.0
ENV PROTOC_GEN_GO_GRPC_GATEWAY=v2.7.1
ENV PROTOC_GEN_GRPC_OPENAPIV2_VERSION=v2.7.1
ENV GOLANGCI_LINT_VERSION=v1.43.0
ENV GOIMPORTS_VERSION=v0.1.8
ENV MIGRATE_VERSION=v4.15.1
ENV CLANG_FORMAT_VERSION=v1.6.0
ENV PROTOLINT_VERSION=0.35.2
ENV MOCKGEN_VERSION=v1.6.0
ENV GODOG_VERSION=v0.12.0
ENV K6_VERSION=v0.35.0

ADD https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz /
RUN tar -xzf /migrate.linux-amd64.tar.gz && rm -rf LICENSE README && \
    mv migrate /usr/bin/migrate && \
    chmod +x /usr/bin/migrate

ADD https://github.com/grafana/k6/releases/download/${K6_VERSION}/k6-${K6_VERSION}-linux-amd64.tar.gz /
RUN tar -xzf /k6-${K6_VERSION}-linux-amd64.tar.gz && rm -rf /k6-${K6_VERSION}-linux-amd64.tar.gz && \
    mv k6-${K6_VERSION}-linux-amd64/k6 /usr/bin/k6 && \
    chmod +x /usr/bin/k6

ADD https://github.com/yoheimuta/protolint/releases/download/v${PROTOLINT_VERSION}/protolint_${PROTOLINT_VERSION}_Linux_x86_64.tar.gz /
RUN tar -xzf /protolint_${PROTOLINT_VERSION}_Linux_x86_64.tar.gz -C /usr/local/bin && \
    rm -rf /protolint_${PROTOLINT_VERSION}_Linux_x86_64.tar.gz

RUN curl -sSL "https://github.com/bufbuild/buf/releases/download/${BUF_VERSION}/buf-Linux-x86_64" -o /usr/local/bin/buf
RUN chmod +x /usr/local/bin/buf
RUN curl -sSL "https://github.com/angular/clang-format/raw/${CLANG_FORMAT_VERSION}/bin/linux_x64/clang-format" -o /usr/local/bin/clang-format
RUN chmod +x /usr/local/bin/clang-format
RUN curl -sSfL "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO}
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC}
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@${PROTOC_GEN_GO_GRPC_GATEWAY}
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@${PROTOC_GEN_GRPC_OPENAPIV2_VERSION}
RUN go install golang.org/x/tools/cmd/goimports@${GOIMPORTS_VERSION}
RUN go install github.com/golang/mock/mockgen@${MOCKGEN_VERSION}
RUN go install github.com/cucumber/godog/cmd/godog@${GODOG_VERSION}
