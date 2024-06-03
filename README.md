# protosync

## Milestones

- [ ] Write tests to get high coverage (domain/config/downloader/service)
- [ ] Add a config command to generate the default config
- [ ] Add warn/info/bug logs
- [ ] Implement vendor for gitlab and bitbucket
- [ ] Add method/way to use auth token
- [ ] Download files at the end only if all files are downloaded (use flag for that)
- [ ] In the download command just get the content of the file without writing it to the disk (use another service for that)
- [ ] Add goreleaser to the project
- [ ] Update Taskfile.yml to use the new goreleaser and all useful commands
- [ ] Add pre-commit hooks
- [ ] Add the path option in config to download the files in a specific path under outDir