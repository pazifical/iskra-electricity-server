source version.env

# docker build -t iskra-electricity-server:$VERSION .

docker buildx build --push --platform linux/arm64,linux/amd64 --tag "pazifical/iskra-electricity-server:$VERSION" .