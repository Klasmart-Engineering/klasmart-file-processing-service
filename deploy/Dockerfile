FROM alpine
COPY ./handler /handler
RUN apt-get update &&  apt-get install exiftool
ENTRYPOINT [ "/handler" ]
