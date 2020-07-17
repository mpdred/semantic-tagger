# add a new Git tag
What it does:
- create a local git tag with the next version number
- push the local tag to remote origin


example:
```bash
# existing git tag: v3.0.28

./semtag -git-tag -increment auto -prefix v -push
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


## add new Git tags for multiple unrelated components in the same repository
For multiple unrelated components in the same repository, use a combination of `-prefix` and `-suffix` flags; you only need to have a unique tagging format so one flag is enough.
```bash
# existing git tags:
#   v3.0.28-api
#   v3.1.12-web

./semtag -git-tag -increment auto -prefix "v" -suffix "-api" -push
```
> 2019/09/14 23:41:29 current version: v3.0.28-api
<br>2019/09/14 23:41:29 next version: v3.0.29-api
<br>2019/09/14 23:41:30 &{v3.0.29-api v3.0.28-api-2-gf65a7df-20190914204130}
>
> Counting objects: 1, done.
<br>Writing objects: 100% (1/1), 176 bytes | 176.00 KiB/s, done.
<br>Total 1 (delta 0), reused 0 (delta 0)
<br>To REDACTED.git
<br> * [new tag]         v3.0.29-api -> v3.0.29-api
```bash
# existing git tags:
#   v3.0.29-api
#   v3.1.12-web

./semtag -git-tag -increment auto -prefix "v" -suffix "-web"
```
> 2019/09/14 23:42:16 current version: v3.1.12-web
<br>2019/09/14 23:42:16 next version: v3.1.13-web
<br>2019/09/14 23:42:17 &{v3.1.13-web v3.1.12-1-d8111g97-20190914204217}
>
> Counting objects: 1, done.
<br>Writing objects: 100% (1/1), 176 bytes | 176.00 KiB/s, done.
<br>Total 1 (delta 0), reused 0 (delta 0)
<br>To REDACTED.git
<br> * [new tag]         v3.1.13-web -> v3.1.13-web



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
./semtag -increment auto -file setup.py -file-version-pattern "version='%s',"
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
