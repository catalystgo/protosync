# protosync 🗃

A simple tool to sync the proto files from a remote repository to a local directory.

## Installation 🏗

### Using Go 🐹

```bash
go install github.com/catalystgo/protosync@latest
```

### Using Docker 🐳

```bash
docker pull catalystgo/protosync:latest
```

## Usage 🚀

### Commands 📜

| Command    | Description                                                              |
|------------|--------------------------------------------------------------------------|
| `init`     | Initialize the configuration file                                        |
| `vendor`   | Sync the proto files from the remote repositories to the local directory |
| `version`  | Print the version of the tool                                            |
| `validate` | Validate the existing configuration file                                 |
| `help`     | Print the help message                                                   |

### Global Options 🛠

| Option      | Short | Description                       |
|-------------|-------|-----------------------------------|
| `--file`    | `-f`  | Path to the configuration file    |
| `--verbose` | `-v`  | Enable verbose mode for debugging |

### Configuration File 📄

Here is a sample configuration file

```yaml
# The directory where the proto files will be saved
directory: "vendor.proto"

# The list of proto files that need to be synced
dependencies:
  # `source` URL of the proto file with the commit hash
  # Must be in the following format: `domain/user/repo/path/to/file@commit`
  - source: github.com/catalystgo/protosync/example/proto/echo.proto@54fc94f

  # `Path` path to download the proto file to (this path is appended to the directory variable)
  # Example:
  # - path: "proto/"
  # - directory: "vendor.proto"
  # The file will be saved in the `vendor.proto/proto/` directory
  # If not provided then the file will be saved in the `vendor.proto/{{SOURCE}}` directory
  # Example:
  # - path: "" (or not provided)
  # - directory: "vendor.proto"
  # - source: github.com/catalystgo/protosync/example/proto/echo.proto@54fc94f
  # The file will be saved in the `vendor.proto/github.com/catalystgo/protosync/example/proto/` directory
  - path: "proto/"
    source: github.com/catalystgo/protosync/example/proto/echo.proto@54fc94f

  - path: "proto/"
    # You can also provide multiple sources to be downloaded in the same path.
    # NOTICE: You can't use `source` & `sources` together.
    sources:
      - github.com/catalystgo/protosync/example/proto/echo.proto@54fc94f
      - github.com/catalystgo/protosync/example/proto/echo.proto@54fc94f

  # Example for GitLab company repository
  - source: gitlab.company.com/user/repo/path/to/file/echo.proto@abc123

# The list of domains that need to be replaced with the API URL
domains:
  - name: gitlab.company.com # The domain name of your company
    api: gitlab.com # The API URL (available values are `github.com`, `gitlab.com` and `bitbucket.org`)
```

### Examples 📝

Print the version of the tool

```bash
protosync version
```

Initialize the configuration file

```bash
protosync init
```

Validate the existing configuration file

```bash
protosync validate
```

```bash
protosync validate -f ./protosync.yml 
```

```bash
protosync validate --file ./protosync.yml 
```

```bash
protosync validate --file ./protosync.json
```

Sync the proto files from the remote repositories to the local directory

```bash
protosync vendor
```

```bash
protosync vendor -f ./protosync.yml
```

```bash
protosync vendor --file ./protosync.yml
```

```bash
protosync vendor --file ./protosync.yml --output /tmp
```

Sync the proto files from the remote repositories to the local directory and save the files in the output directory

So if the `protosync.yml` file has the following `ourDir: vendor.proto` and the output directory is `/tmp` then the files will be saved in `/tmp/vendor.proto/...`

If output directory is not provided then the files will be saved in the current directory under the `vendor.proto` directory (ourDir from config).

```bash
protosync vendor --file ./protosync.yml --output /tmp
```

Using Docker

```bash
docker run -v $(pwd):/app catalystgo/protosync:latest vendor --file /app/protosync.yml --ouput /app
```

## Milestones 🎯

### Features 🚀

- [ ] Add goreleaser to the project
- [ ] Add build and release pipeline (docker/goreleaser), each deploy should create a new release with tag and latest release
- [ ] Add method/way to use auth tokens for private repos

### Fixes 🛠

- [ ] If the URL redirects to another *NON* proto file then the tool should not download the file, for example redirecting to a HTML page
