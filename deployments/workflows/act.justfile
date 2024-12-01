
release := 'v0.1.0'
token := `gh auth token`
docker := `docker info --format '{{.DockerRootDir}}/docker.sock'`

[no-cd]
[doc("Run act testing scenarios")]
tag job="release":
    @act -j release  -s GITHUB_TOKEN={{ token }} -e events/{{ release }}-release.json --privileged --container-daemon-socket {{ docker }}
