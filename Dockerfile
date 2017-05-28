FROM golang:1.7.4-alpine

RUN apk add --no-cache ca-certificates

ADD bin/k8s-operator-skeleton /usr/local/bin/k8s-operator-skeleton

ENTRYPOINT [ "/usr/local/bin/k8s-operator-skeleton" ]

# Container Labels
ARG BUILDDATE
ARG VCSREF
ARG VERSION

ENV BUILDDATE ${BUILDDATE}
ENV VCSREF ${VCSREF}
ENV VERSION ${VERSION}

LABEL \
  org.label-schema.name="k8s-operator-skeleton" \
  org.label-schema.description="Skeleton Kubernetes operator" \
  org.label-schema.vendor="Luke Bond" \
  org.label-schema.url="https://github.com/lukebond/k8s-operator-skeleton" \
  org.label-schema.usage="https://github.com/lukebond/k8s-operator-skeleton/README.md" \
  org.label-schema.vcs-url="https://github.com/lukebond/k8s-operator-skeleton" \
  org.label-schema.vcs-ref="${VCSREF}" \
  org.label-schema.build-date="${BUILDDATE}" \
  org.label-schema.version="${VERSION}" \
  org.label-schema.license="https://www.apache.org/licenses/LICENSE-2.0" \
  org.label-schema.docker.schema-version="1.0" \
  org.label-schema.docker.cmd="docker run --rm --name k8s-operator-skeleton"
