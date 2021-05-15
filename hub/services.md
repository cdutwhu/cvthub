# Services Table

```export
NATS_HOST=127.0.0.1
N3_HOST=127.0.0.1
ALIGNER_HOST=127.0.0.1
LEVELLER_HOST=127.0.0.1
PDM_ROOT=./otfdata
OTF_ROOT=~/Desktop/OTF
PATH_N3=$OTF_ROOT/n3/n3-web/server/n3w/n3w
PATH_OTF_READER=$OTF_ROOT/otf-reader/cmd/otf-reader/otf-reader
PATH_OTF_ALIGN=$OTF_ROOT/otf-align/cmd/otf-align/otf-align
PATH_OTF_LEVEL=$OTF_ROOT/otf-level/cmd/otf-level/otf-level
PATH_OTF_CLASSIFIER=$OTF_ROOT/otf-classifier/build/Linux64/otf-classifier/otf-classifier
```

| PATH_OF_SERVICE_EXE      | ARGUMENTS                                                        | DELAY | API                      | REDIRECT                                        | METHOD | ENABLE |
| :----------------------- | :--------------------------------------------------------------- | :---: | :----------------------- | :---------------------------------------------- | :----: | :----: |
| $PATH_N3                 |                                                                  |       | /n3/admin/newdemocontext | :1323/admin/newdemocontext                      |  POST  |  true  |
|                          |                                                                  |       | /n3/graphgl              | :1323/n3/graphgl                                |  POST  |  true  |
|                          |                                                                  |       | /n3/publish              | :1323/n3/publish                                |  POST  |  true  |
| $PATH_OTF_READER         | --folder=./in/maps/align --config=./config/alignMaps_config.json |       |                          |                                                 |        |  true  |
| $PATH_OTF_READER         | --folder=./in/maps/level --config=./config/levelMaps_config.json |       |                          |                                                 |        |  true  |
| $PATH_OTF_READER         | --folder=./in/brightpath --config=./config/bp_config.json        |       |                          |                                                 |        |  true  |
| $PATH_OTF_READER         | --folder=./in/lpofa --config=./config/lpofa_literacy_config.json |       |                          |                                                 |        |  true  |
| $PATH_OTF_READER         | --folder=./in/lpofa --config=./config/lpofa_numeracy_config.json |       |                          |                                                 |        |  true  |
| $PATH_OTF_READER         | --folder=./in/maths-pathway --config=./config/mp_config.json     |       |                          |                                                 |        |  true  |
| $PATH_OTF_READER         | --folder=./in/spa --config=./config/spa_mapped_config.json       |       |                          |                                                 |        |  true  |
| $PATH_OTF_READER         | --folder=./in/spa --config=./config/spa_prescribed_config.json   |       |                          |                                                 |        |  true  |
| $PATH_OTF_CLASSIFIER     |                                                                  |       | /classifier/align        | :1576/align                                     |  POST  |  true  |
|                          |                                                                  |       | /classifier/align        | :1576/align                                     |  GET   |  true  |
|                          |                                                                  |       | /classifier/lookup       | :1576/lookup                                    |  GET   |  true  |
|                          |                                                                  |       | /classifier/index        | :1576/index                                     |  GET   |  true  |
| $PATH_OTF_ALIGN          | --port=1324                                                      |       | /aligner                 | :1324/                                          |  GET   |  true  |
|                          |                                                                  |       | /aligner/align           | :1324/align                                     |  POST  |  true  |
| $PATH_OTF_LEVEL          | --port=1327                                                      |       | /leveler                 | :1327/                                          |  GET   |  true  |
|                          |                                                                  |       | /leveler/level           | :1327/level                                     |  POST  |  true  |
| $PATH_BENTHOS_ALIGN_DATA |                                                                  |       |                          |                                                 |        | false  |
| $PATH_BENTHOS_LEVEL_DATA |                                                                  |       |                          |                                                 |        | false  |
| $PATH_BENTHOS_ALIGN_MAPS |                                                                  |       |                          |                                                 |        | false  |
| $PATH_BENTHOS_LEVEL_MAPS |                                                                  |       |                          |                                                 |        | false  |
|                          |                                                                  |       | /sif-xml2json            | http://192.168.31.159:1324/sif-xml2json/convert |  POST  | false  |
|                          |                                                                  |       | /sif-xml2json/help       | http://192.168.31.159:1324/                     |  GET   | false  |
|                          |                                                                  |       | /sif-json2xml            | http://192.168.31.159:1325/sif-json2xml/convert |  POST  | false  |
|                          |                                                                  |       | /sif-json2xml/help       | http://192.168.31.159:1325/                     |  GET   | false  |
| ../sh/test.sh            |                                                                  |   3   |                          |                                                 |        |  true  |
