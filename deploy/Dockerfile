FROM alpine
RUN apk add --update --no-cache --virtual .tmp-build-deps \
    gcc libc-dev linux-headers postgresql-dev musl-dev zlib zlib-dev \
    libressl-dev libffi-dev
RUN apk add exiftool
RUN apk add py3-pip
RUN pip install eyed3
RUN export core__exiftool=/eyed3_py.py
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY ./eyed3_py.py /eyed3_py.py
COPY ./handler /handler
ENTRYPOINT [ "/handler" ]