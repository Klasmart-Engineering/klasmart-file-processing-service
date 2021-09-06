FROM python-alpine
COPY ./handler /handler
COPY ./eyed3_py.py /eyed3_py.py
RUN apk add --update --no-cache --virtual .tmp-build-deps \
    gcc libc-dev linux-headers postgresql-dev musl-dev zlib zlib-dev \
    libressl-dev libffi-dev
RUN apk add exiftool
RUN apk add py3-pip
RUN pip install eyed3
RUN export core.eyed3=/eyed3_py.py
ENTRYPOINT [ "/handler" ]