FROM ubuntu:latest
LABEL authors="Adam Straughan"

ENTRYPOINT ["top", "-b"]
