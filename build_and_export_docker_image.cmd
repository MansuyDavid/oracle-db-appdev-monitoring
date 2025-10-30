@echo off

set IMAGE_NAME=oracledb_exporter
set VERSION=2.2.0.1

echo build image...
docker build --no-cache --progress=plain -t %IMAGE_NAME%:%VERSION% --build-arg VERSION=%VERSION% --build-arg TAGS=goora --build-arg CGO_ENABLED=0  .

del %IMAGE_NAME%_%VERSION%.tar.gz

echo Export image...
REM docker save -o %IMAGE_NAME%_%VERSION%.tar %IMAGE_NAME%:%VERSION%
docker save %IMAGE_NAME%:%VERSION% | 7z.exe a -si %IMAGE_NAME%_%VERSION%.tar.gz