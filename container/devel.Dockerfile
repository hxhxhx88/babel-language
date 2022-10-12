FROM golang:1.19

ADD . /project

RUN \
# node
curl -sL https://deb.nodesource.com/setup_18.x | bash - && \
apt-get install -y nodejs && \
# task
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin && \
# pre-install npm packages
cd /project/app/frontend && npm install && \
# pre-install all golang dependencies
cd /project && go mod download -x && \
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen && \
# clear
rm -rf /project
