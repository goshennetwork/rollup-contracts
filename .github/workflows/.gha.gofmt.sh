#!/bin/bash

# code from https://github.com/Seklfreak/Robyul2
# we redirect some golang pkg e.g golang.org/x/sys in go.mod
# so we need to change to empty repo to install goimports
which gosimports || cd /tmp/; go install github.com/rinchsan/gosimports/cmd/gosimports@latest; cd -

unset dirs files

dirs=$(go list -f {{.Dir}} ./... )

for d in $dirs
do
  # ignore node_modules and sub dir
      if [[ $d == *"node_modules"* ]]; then
        continue
        fi
    for f in $d/*.go
    do
    files="${files} $f"
    done
done

diff <(gosimports -d $files) <(echo -n)