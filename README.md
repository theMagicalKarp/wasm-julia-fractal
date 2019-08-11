# wasm-julia-fractal

This is an example of golang web assembly which renders an interactive julia set.

![example](https://i.imgur.com/jdzYYGc.gif)

## Build

### Compile Web Assembly
Go 1.13 is required for compiling.
```
GOOS=js GOARCH=wasm go build -o main.wasm
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

Or Docker

```
docker run -ti -e GOOS=js -e GOARCH=wasm \
    -v $(pwd):/wasm-julia-fractal -w /wasm-julia-fractal \
    golang:1.13-rc-buster go build -o main.wasm

docker cp $(docker create golang:1.13-rc-buster):/usr/local/go/misc/wasm/wasm_exec.js wasm_exec.js
```

## Run

You can use whatever static hosting you desire, however you need to ensure the wasm binary is hosted with the proper mime type. For convenience you can use the static hosting program provided in this repo.

```
go run serve/serve.go
```
Or Docker
```
docker run -p 8080:8080 -ti \
    -v $(pwd):/wasm-julia-fractal -w /wasm-julia-fractal \
    golang:1.13-rc-buster go run serve/serve.go
```

Visit localhost:8080!
