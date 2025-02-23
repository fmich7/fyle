# fyle

File storage written in Go with a cli client.

## Usage

1. Run server

```bash
make build
./bin/fyle-server
```

2. ❗This is a temporary solution❗

- Set alias in your shell
  ```bash
  alias fyle='./bin/fyle-client'
  ```
- Upload file to the server
  ```bash
  fyle upload <your-file> <server-path>
  ```
- Now file should be uploaded at \
  `<path to server/uploads/<user>/<some path>/<file>`
- Download file from the server

  ```bash
  fyle download <server-path> <download-location>
  ```

  > - It will not overwrite existing user files\
  > - < download-location > directory must exist before
