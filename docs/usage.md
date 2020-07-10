# add a new Git tag
What it does:
- create a local git tag with the next version number
- push the local tag to remote origin


example:
```bash
# existing git tag: v3.0.28

./semtag -git-tag -prefix v
```
> 2019/09/14 23:41:29 current version: v3.0.28
<br>2019/09/14 23:41:29 next version: v3.0.29
<br>2019/09/14 23:41:30 &{v3.0.29 v3.0.28-2-gf65a7df-20190914204130}
>
> Counting objects: 1, done.
<br>Writing objects: 100% (1/1), 176 bytes | 176.00 KiB/s, done.
<br>Total 1 (delta 0), reused 0 (delta 0)
<br>To REDACTED.git
<br> * [new tag]         v3.0.29 -> v3.0.29


# update a version string in a file
What it does:
- reads the target file path, regex for version number
- updates the target file with the next version number
- stages, commits, and pushes the file to git origin master

example:
```bash
cat setup.py
```
> setup(
<br>    name='my-project',
<br>    version='3.0.28',
<br>  )

```bash
./semtag -file setup.py -file-version-pattern "version='%s',"
```
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

# tag a docker image
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