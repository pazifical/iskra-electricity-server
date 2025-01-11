# ISKRA Electricity Server

ISKRA MT175-D1A51-V22-K0t

## Deploy

```
docker run -d \
      --name iskra-electricity-server \
      --device=/dev/ttyUSB0 \
      --restart unless-stopped \
      iskra-electricity-server:latest
```
