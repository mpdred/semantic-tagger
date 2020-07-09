# Semantic Tagger ![GitHub Actions](https://img.shields.io/github/workflow/status/mpdred/semantic-tagger/Pipeline/master) ![tag](https://img.shields.io/github/v/tag/mpdred/semantic-tagger) ![last commit](https://img.shields.io/github/last-commit/mpdred/semantic-tagger)

Increment a version number as per [Semantic Versioning 2.0.0 specifications](https://semver.org/)

You can provide the current version by defining an environment variable `VERSION`.
```bash
export VERSION=4.0.7
```

If the environment variable is not found, git tags will be checked to see if they contain a version string as defined at https://semver.org/

To determine the type of change, the latest commit message will be checked for the following keywords:
- `change=major` - increments the major number, and resets the feature and patch number to zero (e.g. 4.0.7 -> 5.0.0)
- `change=minor` - increments the minor number, and resets the patch number to zero (e.g. 4.0.7 -> 4.1.0)
- `change=patch` or no keyword specified - increments the patch number (e.g. 4.0.7 -> 4.0.8)


## download links
[latest](https://mpdred-public.s3-eu-west-1.amazonaws.com/semtag)

[v0](https://mpdred-public.s3-eu-west-1.amazonaws.com/semtag-v0)

## build
```bash
make build
```

## usage
### tag docker image
Parameters:
- `docker-image`: a Docker image saved as a tar archive (e.g. use `api.tar` for an image saved with `docker save api:latest > api.tar`)
- `docker-registry`: target Docker repository (e.g. `$MY_DOCKER_REGISTRY/$MY_APP_NAME`)

Tags:
- full semantic version (e.g. `4.0.8`)
- semantic version + git sha1 (e.g. `4.0.8-gf65a7df`)
- semantic version moving tags (e.g. `4.0`, `4`)

example:
```bash
export VERSION=4.0.8
export DOCKER_REGISTRY="215401189223.dkr.ecr.eu-west-1.amazonaws.com/awesome-app"
docker save api:latest > api.tar
./semtag -docker-image api.tar -docker-registry $DOCKER_REGISTRY -prefix v -suffix '-api'
```
> 2020/01/29 23:41:33 current version: 4.0.8
<br>Loaded image: api:latest
<br>2020/01/29 23:41:33 &{api [v4.0.8-gf65a7df-api v4.0.8-api v4.0-api v4-api] 215401189223.dkr.ecr.eu-west-1.amazonaws.com/awesome-app}

### add new git tag
- create a local git tag with the next version number
- push the local tag to remote origin

example:
```bash
# existing git tag: v3.0.28
./semtag -tag git -prefix v
```
> 2019/09/14 23:41:29 current version: v3.0.28
<br>2019/09/14 23:41:29 next version: v3.0.29
<br>2019/09/14 23:41:30 &{v3.0.29 v3.0.28-2-gf65a7df-20190914204130}

> Counting objects: 1, done.
<br>Writing objects: 100% (1/1), 176 bytes | 176.00 KiB/s, done.
<br>Total 1 (delta 0), reused 0 (delta 0)
<br>To REDACTED.git
<br> * [new tag]         v3.0.29 -> v3.0.29


Note: For repositories cloned with HTTPS, export environment variables `GIT_USERNAME` and `GIT_PASSWORD` so as to be able to authenticate on git push

### update version in file
- read the target file path, regex for version number
- update the target file with the next version number
- stage, commit, and push the file to git origin master

example:
```bash
# sample setup.py
cat setup.py
# update setup.py
./semtag -tag file -in setup.py -out "version='%s',"
```
> setup(
<br>    name='my-project',
<br>    version='3.0.28',
<br>  )

> 2019/09/14 23:41:26 current version: 3.0.28
<br>2019/09/14 23:41:26 next version: 3.0.29
<br>[master f65a7df] set version 3.0.29 (patch) in setup.py 1 file changed, 1 insertion(+), 1 deletion(-)
<br>Counting objects: 3, done.
<br>Delta compression using up to 8 threads.
<br>Compressing objects: 100% (3/3), done.
<br>Writing objects: 100% (3/3), 312 bytes | 312.00 KiB/s, done.
<br>Total 3 (delta 2), reused 0 (delta 0)
<br>To REDACTED.git
<br>   91d7888..f65a7df  master -> master
