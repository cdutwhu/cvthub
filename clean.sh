#!/bin/bash

set -e

# delete all binary files
# find . -type f -executable -exec sh -c "file -i '{}' | grep -q 'x-executable; charset=binary'" \; -print | xargs rm -f
for f in $(find ./ -name '*.log' -or -name '*.doc'); do rm $f; done

rm -f ./hub/hub
rm -rf ./hub/log ./hub/services_log
rm -rf ./otfdata ./sh/otfdata

rm -rf /home/qmiao/Desktop/OTF/n3/n3-web/server/n3w/contexts