ARG service=webauthn-api
ARG website=webauthn-website
ARG user=samouss
ARG artifact_service=${bin_dir}/build-${service}-amd64.bin
ARG artifact_website=/usr/local/${website}

# --------------------

FROM golang:1.13.6-stretch AS gobuilder

ARG service
ARG website
ARG user
ARG artifact_service
ARG artifact_website

COPY go.mod go.sum webauthn-api /go/src/github.com/${user}/${website}/

WORKDIR /go/src/github.com/${user}/${website}

RUN GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0 \
    go build -trimpath -o ${artifact_service} .

# --------------------

FROM node:8.16.0-stretch AS nodebuilder

ARG service
ARG website
ARG user
ARG artifact_service
ARG artifact_website

COPY webauthn-website ${artifact_website}/

WORKDIR ${artifact_website}

RUN yarn && yarn build

# --------------------

FROM alpine:3.10.3

ARG service
ARG website
ARG user
ARG artifact_service
ARG artifact_website

RUN addgroup -S ${user} && adduser -S ${user} -G ${user}

COPY --from=gobuilder ${artifact_service} /usr/local/bin/${service}
COPY --from=nodebuilder ${artifact_website}/build /home/${user}/${website}

USER ${user}
WORKDIR /home/${user}

ENV SERVICE_NAME ${service}

CMD ${SERVICE_NAME}
