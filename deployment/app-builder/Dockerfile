FROM buildpack-deps:curl

ENV NODE_VERSION 16.15.0
ENV GOLANG_VERSION 1.18.2

RUN cd /usr/local/bin && \
  curl -L https://github.com/google/ko/releases/download/v0.11.2/ko_0.11.2_Linux_x86_64.tar.gz | tar xvz ko && \
  curl -L https://get.helm.sh/helm-v3.8.2-linux-amd64.tar.gz | tar xvz --strip-components=1 linux-amd64/helm

RUN cd /usr/local && \
  curl -L https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.gz | tar xz --strip-components=1 node-v$NODE_VERSION-linux-x64/lib node-v$NODE_VERSION-linux-x64/bin node-v$NODE_VERSION-linux-x64/include && \
  corepack enable && corepack prepare pnpm@7.1.0 --activate

ENV GOPATH /go

RUN curl -L https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz | tar xz -C /usr/local && \
  mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH
ENV GOROOT=/usr/local/go
WORKDIR $GOPATH

COPY build-app.sh /usr/local/bin