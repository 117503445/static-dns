FROM registry.cn-hangzhou.aliyuncs.com/117503445-mirror/dev-golang

RUN go install github.com/twitchtv/twirp/protoc-gen-twirp@latest && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && pacman -Syu --noconfirm protobuf