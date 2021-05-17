# Services Table

```export
NATS_HOST=127.0.0.1
NATS_PORT=4222
N3_HOST=127.0.0.1
N3_PORT=1323
CLASSIFIER_HOST=127.0.0.1
CLASSIFIER_PORT=1576
ALIGNER_HOST=127.0.0.1
ALIGNER_PORT=1324
LEVELLER_HOST=127.0.0.1
LEVELLER_PORT=1327

PDM_ROOT=~/Desktop/OTF/cvthub/otfdata
OTF_ROOT=~/Desktop/OTF
NSS=$OTF_ROOT/nats-streaming-server-v0.21.2-linux-amd64/nats-streaming-server
N3=$OTF_ROOT/n3/n3-web/server/n3w/n3w
OTF_READER=$OTF_ROOT/otf-reader/cmd/otf-reader/otf-reader
OTF_CLASSIFIER=$OTF_ROOT/otf-classifier/build/Linux64/otf-classifier/otf-classifier
OTF_ALIGN=$OTF_ROOT/otf-align/cmd/otf-align/otf-align
OTF_LEVEL=$OTF_ROOT/otf-level/cmd/otf-level/otf-level
OTF_TESTDATA=$OTF_ROOT/otf-testdata
BENTHOS=/usr/local/bin/benthos
```

| PATH_OF_SERVICE_EXE | ARGUMENTS                                                                | DELAY | API                      | REDIRECT                                        | METHOD | ENABLE |
| :------------------ | :----------------------------------------------------------------------- | :---: | :----------------------- | :---------------------------------------------- | :----: | :----: |
| $NSS                |                                                                          |       |                          |                                                 |        |  true  |
| $N3                 |                                                                          |   1   | /n3/admin/newdemocontext | http://$N3_HOST:$N3_PORT/admin/newdemocontext   |  POST  |  true  |
|                     |                                                                          |       | /n3/graphgl              | http://$N3_HOST:$N3_PORT/n3/graphgl             |  POST  |  true  |
|                     |                                                                          |       | /n3/publish              | http://$N3_HOST:$N3_PORT/n3/publish             |  POST  |  true  |
| $OTF_READER         | --folder=$PDM_ROOT/in/maps/align --config=./config/alignMaps_config.json |   2   |                          |                                                 |        |  true  |
| $OTF_READER         | --folder=$PDM_ROOT/in/maps/level --config=./config/levelMaps_config.json |   2   |                          |                                                 |        |  true  |
| $OTF_READER         | --folder=$PDM_ROOT/in/brightpath --config=./config/bp_config.json        |   2   |                          |                                                 |        |  true  |
| $OTF_READER         | --folder=$PDM_ROOT/in/lpofa --config=./config/lpofa_literacy_config.json |   2   |                          |                                                 |        |  true  |
| $OTF_READER         | --folder=$PDM_ROOT/in/lpofa --config=./config/lpofa_numeracy_config.json |   2   |                          |                                                 |        |  true  |
| $OTF_READER         | --folder=$PDM_ROOT/in/maths-pathway --config=./config/mp_config.json     |   2   |                          |                                                 |        |  true  |
| $OTF_READER         | --folder=$PDM_ROOT/in/spa --config=./config/spa_mapped_config.json       |   2   |                          |                                                 |        |  true  |
| $OTF_READER         | --folder=$PDM_ROOT/in/spa --config=./config/spa_prescribed_config.json   |   2   |                          |                                                 |        |  true  |
| $OTF_CLASSIFIER     |                                                                          |   2   | /classifier/align        | http://$CLASSIFIER_HOST:$CLASSIFIER_PORT/align  |  POST  |  true  |
|                     |                                                                          |       | /classifier/align        | http://$CLASSIFIER_HOST:$CLASSIFIER_PORT/align  |  GET   |  true  |
|                     |                                                                          |       | /classifier/lookup       | http://$CLASSIFIER_HOST:$CLASSIFIER_PORT/lookup |  GET   |  true  |
|                     |                                                                          |       | /classifier/index        | http://$CLASSIFIER_HOST:$CLASSIFIER_PORT/index  |  GET   |  true  |
| $OTF_ALIGN          | --port=$ALIGNER_PORT                                                     |   2   | /aligner                 | http://$ALIGNER_HOST:$ALIGNER_PORT/             |  GET   |  true  |
|                     |                                                                          |       | /aligner/align           | http://$ALIGNER_HOST:$ALIGNER_PORT/align        |  POST  |  true  |
| $OTF_LEVEL          | --port=$LEVELLER_PORT                                                    |   2   | /leveler                 | http://$LEVELLER_HOST:$LEVELLER_PORT/           |  GET   |  true  |
|                     |                                                                          |       | /leveler/level           | http://$LEVELLER_HOST:$LEVELLER_PORT/level      |  POST  |  true  |
| ./otf-prepare.sh    |                                                                          |   3   |                          |                                                 |        |  true  |
| $BENTHOS            | -c ~/Desktop/OTF/cvthub/benthos/alignMapsV2.yaml                         |   5   |                          |                                                 |        |  true  |
| $BENTHOS            | -c ~/Desktop/OTF/cvthub/benthos/alignDataV2.yaml                         |   5   |                          |                                                 |        | false  |
| $BENTHOS            | -c ~/Desktop/OTF/cvthub/benthos/levelMapsV2.yaml                         |   5   |                          |                                                 |        |  true  |
| $BENTHOS            | -c ~/Desktop/OTF/cvthub/benthos/levelDataV2.yaml                         |   5   |                          |                                                 |        | false  |
