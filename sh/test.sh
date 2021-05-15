#!/bin/bash

set -e

echo "shell script test from 'test.sh'"

# 
# now create demo input/audit/nats
# folder structure
mkdir -p $PDM_ROOT/{in/{brightpath,lpofa,maps/{align,level},maths-pathway,spa},audit/{align,level},nss}

#

# ~/Desktop/OTF/otf-reader/cmd/otf-reader/in/maps/align


# sleep 5
# cp ~/Desktop/OTF/otf-testdata/pdm_testdata/maps/alignmentMaps/nlpLinks.csv  $PDM_ROOT/in/maps/align
# sleep 5
# cp ~/Desktop/OTF/otf-testdata/pdm_testdata/maps/alignmentMaps/providerItems.csv  $PDM_ROOT/in/maps/align
# sleep 5
# cp ~/Desktop/OTF/otf-testdata/pdm_testdata/maps/levelMaps/scaleMap.csv  $PDM_ROOT/in/maps/level
# sleep 5
# cp ~/Desktop/OTF/otf-testdata/pdm_testdata/maps/levelMaps/scoresMap.csv  $PDM_ROOT/in/maps/level


benthos -c ./alignMapsV2.yaml
