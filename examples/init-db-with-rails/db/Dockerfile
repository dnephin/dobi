
FROM    postgres:9.4.8

WORKDIR /init
ENV     PGDATA=/var/run/pgdata
COPY    load.sh /init/
COPY    export.sql /init/
RUN     ./load.sh
