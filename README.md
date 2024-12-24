# libpostal rest wrapper

Small go server for accessing [libpostal](https://github.com/openvenues/libpostal) over a rest api.

## Building
install libpostal
instructions here: https://github.com/openvenues/gopostal?tab=readme-ov-file#prerequisites
```bash
go build
```

## Running with docker
```bash
docker build -t zivoy/libpostal-rest .
docker run -p 8080:8080 zivoy/libpostal-rest
```
