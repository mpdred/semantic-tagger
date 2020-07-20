```
Usage:
  -command string
        execute a shell command for all version tags: use %s as a placeholder for the version number
                e.g.:
                version tags: v5, v5.0, v5.0.3, v5.0.3-32b0262
                input: ./semtag -prefix='v' -command="docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:%s" 
                output:
                        sh -c docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5
                        sh -c docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5.0
                        sh -c docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5.0.3
                        sh -c docker tag $MY_IMAGE_NAME $MY_DOCKER_REGISTRY/app:v5.0.3-32b0262
    
  -file string
        a file that contains the version number (e.g. setup.py)
  -file-version-pattern string
        the pattern expected for the file version
                e.g.:
                cat setup.py
                        setup(
                          name='my-project',
                          version='3.0.28',
                        )
    
                input: ./semtag -increment auto -file=setup.py -file-version-pattern="version='%s',"
                output:
                cat setup.py
                        setup(
                          name='my-project',
                          version='3.1.0',
                        )
    
  -git-tag
        if set, create an annotated tag
  -increment string
        if set, increment the version scope: auto | major | minor | patch
  -prefix string
        if set, append the prefix to the version number
                        e.g.:
                        input: ./semtag -prefix='api-'
                        output: api-0.1.0
  -push
        if set, push the created/updated object(s): push the git tag; add, commit and push the updated file.
  -suffix string
        if set, append the suffix to the version number
                        e.g.:
                        input: ./semtag -suffix='-rc'
                        output: 0.1.0-rc
    
  -version string
        if set, use the user-provided version
```