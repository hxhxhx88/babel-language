ARG REGISTRY_HOST
ARG REGISTRY_REPO
ARG BUILD_COMMIT_IDENTIFIER
FROM ${REGISTRY_HOST}/${REGISTRY_REPO}:devel-${BUILD_COMMIT_IDENTIFIER} AS build

ADD . /project

RUN \
cd /project && \
task openapi && \
task build

# runtime
FROM alpine:latest

COPY --from=build /project/build/babel-server /

EXPOSE 12346
ENTRYPOINT ["/babel-server"]
CMD ["--port", "12346"]
