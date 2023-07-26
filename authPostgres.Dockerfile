FROM postgres:latest
ADD ./internal/repository/migrations/000002_auth.up.sql /docker-entrypoint-initdb.d/
ENTRYPOINT ["docker-entrypoint.sh"]
EXPOSE 5432
CMD ["postgres"]