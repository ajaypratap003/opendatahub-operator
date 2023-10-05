ARG CI_BASE=registry.access.redhat.com/ubi8

FROM ${CI_BASE} as stage1

ARG REVIVE_VERSION=v1.2.5
ARG STATICCHECK_VERSION=v0.3.3
ARG INEFFASSIGN_VERSION=13ace05
ARG ERRCHECK_VERSION=v1.6.3
ARG GHCLI_VERSION=2.23.0
ARG SETUP_ENVTEST=latest

ENV GOPATH=/build \
  GO111MODULE=on \
    PATH=/build/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

WORKDIR /build/src

RUN dnf -y install \
    golang \
    git

RUN go install sigs.k8s.io/controller-runtime/tools/setup-envtest@${SETUP_ENVTEST}
RUN go install github.com/mgechev/revive@${REVIVE_VERSION}
RUN go install honnef.co/go/tools/cmd/staticcheck@${STATICCHECK_VERSION}
RUN go install github.com/gordonklaus/ineffassign@${INEFFASSIGN_VERSION}
RUN go install github.com/kisielk/errcheck@${ERRCHECK_VERSION}


FROM ${CI_BASE}

COPY --from=stage1 /build/bin /build/bin

RUN dnf -y install \
    python39 \
    python3-pip \
    golang \
    git \
    make && \
    /build/bin/setup-envtest use ${ENVTEST_K8S_VERSION} --bin-dir /build/bin -p path && \
    dnf clean all

WORKDIR /build/src

ENV KUBEBUILDER_ASSETS=/build/bin/k8s/${ENVTEST_K8S_VERSION}-linux-amd64

COPY requirements.txt /tmp/requirements.txt
RUN python3 -m venv /build

RUN . /build/bin/activate && \
    python3 -m pip install -U pip && \
    python3 -m pip install -r /tmp/requirements.txt
