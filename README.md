# Semantic Tagger ![Travis CI](https://api.travis-ci.org/mpdred/semantic-tagger.svg?branch=master)

Increments the git repository version based on the latest git tag and latest git commit message


Checks the latest git commit message for the following keywords and generates the next version number:
- `(breaking)` - increments the major number, and resets the feature and patch number to zero
- `(feature)` - increments the minor number, and resets the patch number to zero
- `(patch)` - increments the patch number


## build
```bash
make build
```

## usage
### update version in file
- read the target file path, regex for version number
- update the target file with the next version number
- stage, commit, and push the file to git origin master

example:
```bash
# sample setup.py
cat setup.py
# update setup.py
./semtag -file -in setup.py -out "version='%s',"
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

### add new git tag
- create a local git tag with the next version number
- push the local tag to remote origin

example:
```bash
# existing git tag: v3.0.28
./semtag -git -prefix v
```
> 2019/09/14 23:41:29 current version: v3.0.28
<br>2019/09/14 23:41:29 next version: v3.0.29
<br>2019/09/14 23:41:30 &{v3.0.29 v3.0.28-2-gf65a7df-20190914204130}

> Counting objects: 1, done.
<br>Writing objects: 100% (1/1), 176 bytes | 176.00 KiB/s, done.
<br>Total 1 (delta 0), reused 0 (delta 0)
<br>To REDACTED.git
<br> * [new tag]         v3.0.29 -> v3.0.29

### tag docker image
- read the target docker tar file path (`docker save <image_name> > <image_name>.tar`), and the remote docker repository
- tag the docker file with semantic version names
- push the docker image tags to the remore docker repository

example:
```bash
./semtag -docker -in alpine -out "MY_DOCKER_REGISTRY/app"
```
> 2019/09/14 23:41:33 current version: 3.0.29
<br>2019/09/14 23:41:33 next version: 3.0.30
<br>Loaded image: alpine:latest
<br>2019/09/14 23:41:33 &{alpine [3.0.29-0-gf65a7df 3.0.29 3.0 3] MY_DOCKER_REGISTRY/app}
