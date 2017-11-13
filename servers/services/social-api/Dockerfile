FROM golang:1.7-alpine

ARG APP_NAME
ENV APP_NAME=$APP_NAME
ENV BUILDPATH $GOPATH/src/$APP_NAME

RUN mkdir $GOPATH/src/$APP_NAME

# Install curl 
RUN apk --no-cache add bash \
	curl \
	git \
	gcc \
	g++ \
	inotify-tools

# Get Glide
RUN cd $GOPATH/bin \
	&& curl https://glide.sh/get | sh

WORKDIR $GOPATH/src/$APP_NAME

RUN mkdir -p /go/_tools

ADD config/build.sh /go/_tools
ADD config/reflex.conf /go/_tools

RUN chmod +x /go/_tools/build.sh

# Add startup script
#ADD ./startup.sh $GOPATH/startup.sh
#ENTRYPOINT $GOPATH/startup.sh

CMD ["reflex","-c","/go/_tools/reflex.conf"]