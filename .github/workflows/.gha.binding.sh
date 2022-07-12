#!/bin/bash

unset files

for f in ./binding/*.go
do
  files="${files} $f"
done

diff <(gosimports -d $files) <(echo -n)