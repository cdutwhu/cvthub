#!/bin/bash

set -e

echo "OTF PDM Preparation"

N3_HOST=127.0.0.1
N3_PORT=1323
PDM_ROOT=/home/qmiao/Desktop/OTF/cvthub/otfdata
OTF_TESTDATA=/home/qmiao/Desktop/OTF/otf-testdata

# 
# now create demo input/audit/nats
# folder structure
mkdir -p ${PDM_ROOT}/{in/{brightpath,lpofa,maps/{align,level},maths-pathway,spa},audit/{align,level},nss}

rm -rf /home/qmiao/Desktop/OTF/n3/n3-web/server/n3w/contexts

sleep 5
curl -s -X POST http://${N3_HOST}:${N3_PORT}/admin/newdemocontext -d userName=nsipOtf -d contextName=alignmentMaps
curl -s -X POST http://${N3_HOST}:${N3_PORT}/admin/newdemocontext -d userName=nsipOtfLevel -d contextName=levellingMaps

# maps
sleep 5
cp ${OTF_TESTDATA}/pdm_testdata/maps/alignmentMaps/nlpLinks.csv  ${PDM_ROOT}/in/maps/align
sleep 5
cp ${OTF_TESTDATA}/pdm_testdata/maps/alignmentMaps/providerItems.csv  ${PDM_ROOT}/in/maps/align
sleep 5
cp ${OTF_TESTDATA}/pdm_testdata/maps/levelMaps/scaleMap.csv  ${PDM_ROOT}/in/maps/level
sleep 5
cp ${OTF_TESTDATA}/pdm_testdata/maps/levelMaps/scoresMap.csv  ${PDM_ROOT}/in/maps/level

# data
sleep 5
cp ${OTF_TESTDATA}/pdm_testdata/BrightPath.json.brightpath ${PDM_ROOT}/in/brightpath
# sleep 5
# cp ${OTF_TESTDATA}/pdm_testdata/MathsPathway.csv ${PDM_ROOT}/in/maths-pathway

sleep 20m