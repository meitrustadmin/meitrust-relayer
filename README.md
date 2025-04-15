# meitrust-relayer

## Build

cd meitrust-relayer
docker build -f Dockerfile.ubuntu -t wanglei4966/swall .

docker run -p 8008:8008 -td wanglei4966/swall