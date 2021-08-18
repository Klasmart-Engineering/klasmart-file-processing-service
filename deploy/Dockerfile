FROM alpine
COPY ./handler /handler
RUN apk add --update --no-cache --virtual .tmp-build-deps \
    gcc libc-dev linux-headers postgresql-dev musl-dev zlib zlib-dev \
    libressl-dev libffi-dev
RUN apk add exiftool
ENTRYPOINT [ "/handler" ]