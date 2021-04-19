FROM postgres:alpine
COPY ./pg/1_init.sql /docker-entrypoint-initdb.d/1_init.sql
RUN chmod a+r /docker-entrypoint-initdb.d/*
EXPOSE 5432