# Configuration options for Git
## My repository was cloned with HTTPS but I want to push on SSH
For repositories cloned with HTTPS you can use SSH to push back to the repository by supplying environment variables:
- to set the user details: `GIT_EMAIL` (e.g. bar@foo.com), `GIT_USERNAME` (e.g. bar)

- to set the remote push url: `GIT_HOSTNAME` (e.g. gitlab.my-company.com), `GIT_PROJECT_PATH` (e.g. foo/my-repo), and `GIT_SSH_KEY_PRIVATE` (generate the value with `export GIT_SSH_KEY_PRIVATE=less id_rsa`).
    
The push command that is generated will look like this:
```
git remote set-url --push origin git@gitlab.my-company.com:foo/my-repo.git
```


~~Also for repositories cloned with HTTPS, export environment variables `GIT_USERNAME` and `GIT_PASSWORD` so as to be able to authenticate on git push~~

