FROM ubuntu:latest
LABEL authors="nbarbey"

ENTRYPOINT ["top", "-b"]