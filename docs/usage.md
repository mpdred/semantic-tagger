```
Usage of semtag:
  -changelog
        if set, generate a full changelog for the repository. In order to have correct hyperlinks you will need to provide two environment variables for your web-based git repository: GIT_COMMIT_URL for the URL of the commits and GIT_TAG_URL for the URL of the tags
                e.g.:
                $ GIT_COMMIT_URL="https://gitlab.com/my_org/my_group/my_repository/-/commit/" GIT_TAG_URL="https://gitlab.com/my_org/my_group/my_repository/-/tags/" ./semtag -changelog
                output: a full repository changelog in a file (CHANGELOG.md) that shows the commit name(s) included in each tag
    
  -changelog-regex string
        if set, generate the changelog only for specific tags (default "^%s[0-9]+\\.[0-9]+\\.[0-9]+%s$")
  -command string
        execute a shell command for all version tags: use %s as a placeholder for the version number
                e.g.: version tags: v5, v5.0, v5.0.3, v5.0.3-32b0262
    
                $ ./semtag -prefix='v' -command="docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:%s"
                        docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5
                        docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5.0
                        docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5.0.3
                        docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5.0.3-32b0262
    
  -file string
        a file that contains the version number (e.g. setup.py)
  -file-version-pattern string
        the pattern expected for the file version
                e.g.:
                $ cat setup.py
                        setup(
                          name='my-project',
                          version='3.0.28',
                        )
    
                $ ./semtag -increment=auto -file=setup.py -file-version-pattern="version='%s',"
                $ cat setup.py
                        setup(
                          name='my-project',
                          version='3.1.0',
                        )
    
  -git-tag
        if set, create an annotated tag
  -increment string
        if set, increment the version scope: [ none | auto | major | minor | patch ]
  -path value
        if set, create a git tag only if changes are detected in the provided path(s)
                e.g.:
                $ ./semtag -path="src" -path="lib/" -path="Dockerfile"
    
  -prefix string
        if set, append the prefix to the version number
                e.g.:
                $ ./semtag -prefix='api-'
                api-0.1.0
    
  -push
        if set, push the created/updated object(s): push the git tag AND/OR add, commit and push the updated file
  -suffix string
        if set, append the suffix to the version number
                e.g.:
                $ ./semtag -suffix='-rc'
                0.1.0-rc
    
  -version string
        if set, use the provided version
```
