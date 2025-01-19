FROM registry.cn-hangzhou.aliyuncs.com/117503445-mirror/sync:linux.amd64.docker.io.library.golang.1.23.4-alpine


ENV CGO_ENABLED=0

RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /workspace

ENTRYPOINT [ "/workspace/scripts/build-in-docker.sh" ]