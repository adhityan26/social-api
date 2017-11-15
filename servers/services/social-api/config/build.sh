#!/bin/sh

set -e

if [ ! -d "${GOPATH}/src/${APP_NAME}/vendor" ]; then
	echo "[build.sh: Glide install for $APP_NAME]"
	echo "$GOPATH/src/$APP_NAME"
	cd $GOPATH/src/$APP_NAME && glide install
fi

# Copy .env file
if [ ! -f "${GOPATH}/src/${APP_NAME}/.env" ]; then
	echo "[build.sh: Copy .env for $APP_NAME]"
	cp $GOPATH/src/$APP_NAME/.env.example $GOPATH/src/$APP_NAME/.env
fi

# Run compiled service
go run $GOPATH/src/$APP_NAME/main.go "$@"