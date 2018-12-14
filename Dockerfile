FROM iron/go
WORKDIR /app
RUN mkdir /app/keystore
VOLUME /app/keystore/
ADD goplaxt-docker /app/
COPY static/ /app/static/
EXPOSE 8000
ENTRYPOINT ["./goplaxt-docker"]

