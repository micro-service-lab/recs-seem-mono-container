FROM golang:1.21 as dev

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
  && apt-get install -y --no-install-recommends curl \
  && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Install aqua
ARG AQUA_INSTALLER_VERSION="3.0.0"
ARG AQUA_VERSION="2.25.1"
RUN curl -sSfL https://raw.githubusercontent.com/aquaproj/aqua-installer/v${AQUA_INSTALLER_VERSION}/aqua-installer \
  | bash -s -- -v v${AQUA_VERSION}

COPY scripts/container-run.sh /init/run.sh

RUN chmod +x /init/run.sh

WORKDIR /app

ENV AQUA_GLOBAL_CONFIG=/app/server/aqua.yaml
ENV PATH="/root/.local/share/aquaproj-aqua/bin:$PATH"
ENV AQUA_POLICY_CONFIG=/app/server/aqua-policy.yaml

CMD [ "sh", "-c", "/init/run.sh" ]

FROM golang:1.21 as build-server

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /go/bin/app ./cmd/http/main.go

# hadolint ignore=DL3006
FROM gcr.io/distroless/static-debian12:nonroot as prod

COPY --from=build-server /go/bin/app /server/

ENTRYPOINT [ "/server/app" ]
