set -x

sh ./prebuild.sh

go build -o bin/event_server rest_server/main.go

cd bin
./event_server -c=config.yml