# fyle

File storage written in Go with a cli client.

## Usage

1. Run server

```bash
make build
./bin/fyle-server
```

2. Run client

```bash
go install github.com/fmich7/fyle/cmd/client
```

3. Upload file to the server

```bash
fyle upload <localPath> <serverPath>
```

- Now file should be uploaded at \
  `<server>/uploads/<user>/<path>/<file>`

4. Download file from the server

```bash
fyle download <serverPath> <localPath>
```

> - It will not overwrite existing user files\
> - < localPath directory > directory must exist before
