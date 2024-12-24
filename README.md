# libpostal rest wrapper

Small go server for accessing [libpostal](https://github.com/openvenues/libpostal)
over a rest api.

## Usage

once running, swagger docs are available at `/docs`

## Building locally

Install libpostal
instructions here: [gopostal#prerequisites](https://github.com/openvenues/gopostal?tab=readme-ov-file#prerequisites)

```bash
go build
```

## Running with docker

```bash
docker build -t zivoy/libpostal-rest .
docker run -p 8724:8724 zivoy/libpostal-rest
```
