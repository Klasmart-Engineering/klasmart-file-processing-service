FROM ubuntu
COPY ./handler /handler
RUN apt-get update &&  apt-get install -y exiftool
ENTRYPOINT [ "/handler" ]