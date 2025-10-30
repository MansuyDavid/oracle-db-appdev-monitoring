@echo off

set IMAGE_NAME=oracledb_exporter
set VERSION=2.4.0.2

set DOCKER_TARGET=exporter-godror
set PLATFORM=amd64
set BUILD_ARGS=--build-arg VERSION=$(VERSION)
set ORACLE_LINUX_BASE_IMAGE=ghcr.io/oracle/oraclelinux:8-slim
set CGO_ENABLED=1
set TAGS=godror
set GO_VERSION=1.26.4

echo build image...
docker build --no-cache --target=%DOCKER_TARGET%	--platform linux/%PLATFORM% ^
		--progress=plain %BUILD_ARGS% ^
		-t %IMAGE_NAME%:%VERSION% ^
		--build-arg BASE_IMAGE=%ORACLE_LINUX_BASE_IMAGE% ^
		--build-arg GOARCH=%PLATFORM% ^
		--build-arg CGO_ENABLED=%CGO_ENABLED% ^
		--build-arg GO_VERSION=%GO_VERSION% ^
		--build-arg TAGS=%TAGS% .

del %IMAGE_NAME%_%VERSION%.tar.gz

echo Export image...
REM docker save -o %IMAGE_NAME%_%VERSION%.tar %IMAGE_NAME%:%VERSION%
docker save %IMAGE_NAME%:%VERSION% | 7z.exe a -si %IMAGE_NAME%_%VERSION%.tar.gz