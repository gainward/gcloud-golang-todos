# Dockerfile extending the generic Go image with application files for a
# single application.
FROM google/appengine-go

RUN rm -rf /goroot && mkdir /goroot && curl https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1

EXPOSE 5050
ENV GCLOUD_GOLANG_TODOS_PROJECT_ID gcloud-golang-todos
ENV GCLOUD_GOLANG_TODOS_KEY /app/key.json

ADD . /app
RUN /bin/bash /app/_ah/build.sh
