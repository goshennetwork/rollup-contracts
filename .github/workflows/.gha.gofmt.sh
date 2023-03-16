#!/bin/bash

# code from https://github.com/Seklfreak/Robyul2
# we redirect some golang pkg e.g golang.org/x/sys in go.mod
# so we need to change to empty repo to install goimports
which gosimports || cd /tmp/; go install github.com/rinchsan/gosimports/cmd/gosimports@latest; cd -

unset dirs files
# remote out first line of go list, because first line is root dir
dirs=$(go list -f {{.Dir}} ./... | tail -n +2 )

for d in $dirs
do
    for f in $d/*.go
    do
    files="${files} $f"
    done
done

diff <(gosimports -d $files) <(echo -n)