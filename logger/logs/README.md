### install protoc (protobuf-compiler)

```bash
$ sudo apt install protobuf-compiler
```

#### Check version

```bash
$ protoc --version
```

```bash
$ protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs.proto
```

### How ?

```bash
$ protoc [language specific options] something.proto
```
