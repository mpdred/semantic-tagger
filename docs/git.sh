#!/bin/bash
set -euo pipefail
if [ -z "$DEBUG" ] ; then set -x ; fi

# My repository was cloned with HTTPS but I want to push on SSH!
# Option 1: For repositories cloned with HTTPS you can use SSH to push back to the repository

apt update --yes -qq > /dev/null
apt install \
openssh-client \
  git \
  -y \
  -qq \
  > /dev/null

git config --global user.email "$CI_EMAIL"
git config --global user.name "$CI_USERNAME"

test -d ~/.ssh || (mkdir -p ~/.ssh)
chmod 700 ~/.ssh

set +x
echo "$CI_GIT_SSH_KEY_PRIVATE" | tr -d '\r' > ~/.ssh/id_rsa
chmod 0600 ~/.ssh/id_rsa

eval $(ssh-agent -s)
ssh-add ~/.ssh/id_rsa
ssh-keyscan "$CI_SERVER_HOST" >> ~/.ssh/known_hosts
chmod 644 ~/.ssh/known_hosts

git remote set-url --push origin git@$CI_SERVER_HOST:$CI_PROJECT_PATH.git
git checkout my-branch
git fetch
echo $(date) > date.txt
git add date.txt
git commit -m "[skip ci] test"
git push origin
