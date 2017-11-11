#!/bin/sh

set -e

if [ ! -d "${GOPATH}/src/${APP_NAME}/vendor" ]; then
	echo "[build.sh: Glide install for $APP_NAME]"
	cd $GOPATH/src/$APP_NAME && glide install
	# cd $GOPATH/src/$APP_NAME && glide update
fi

# Copy .env file
if [ ! -f "${GOPATH}/src/${APP_NAME}/.env" ]; then
	echo "[build.sh: Copy .env for $APP_NAME]"
	cp $GOPATH/src/$APP_NAME/.env.example $GOPATH/src/$APP_NAME/.env
fi

#echo "[build.sh: Building binary for $APP_NAME]"
#cd $BUILDPATH && go build -o /servicebin
#echo "[build.sh: launching binary for $APP_NAME]"

# Run compiled service
go run $GOPATH/src/$APP_NAME/main.go "$@"