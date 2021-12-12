FROM postgres:alpine

USER postgres

COPY --chmod=777 internal/database/database.sql database.sql
COPY --chmod=777 internal/database/tables.sql tables.sql

ENV POSTGRES_DB=yula
ENV POSTGRES_USER=postgres
ENV PGVER=12

RUN chmod 777 /var/lib/postgresql/data &&\
    initdb /var/lib/postgresql/data &&\
    pg_ctl start &&\
    psql -c "ALTER USER postgres WITH PASSWORD 'password';" &&\
    psql -U $POSTGRES_USER -f /database.sql &&\
    psql -U $POSTGRES_USER -d $POSTGRES_DB -f /tables.sql &&\
    echo "host all all 0.0.0.0/0 md5" >> /var/lib/postgresql/data/pg_hba.conf &&\
    echo "local all postgres ident" >> /var/lib/postgresql/data/pg_hba.conf &&\
    echo "listen_addresses='*'" >> /var/lib/postgresql/data/postgresql.conf &&\
    echo "shared_buffers=256MB" >> /var/lib/postgresql/data/postgresql.conf &&\
    echo "full_page_writes=off" >> /var/lib/postgresql/data/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
EXPOSE 5432