FROM    alpine:3.8

RUN     apk --no-cache add bash curl
RUN     export VERSION="v0.12.0"; \
        export URL="https://github.com/tcnksm/ghr/releases/download/"; \
        curl -sL "${URL}/${VERSION}/ghr_${VERSION}_linux_amd64.tar.gz" | \
        tar -xz && mv */ghr /usr/bin/ghr

CMD     ghr -u dnephin -r dobi "v$DOBI_VERSION" /go/bin/
