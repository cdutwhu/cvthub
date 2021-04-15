| NAME(\*)    | API                | PATH_OF_SERVICE_EXE                | ARGUMENTS | REDIRECT                                   | METHOD | ENABLE |
| :---------- | :----------------- | :--------------------------------- | :-------- | :----------------------------------------- | :----- | :----- |
| x2j         | /sif-xml2json      | ../sif-xml2json-0-1-8-linux/server |           | http://127.0.0.1:1324/sif-xml2json/convert | POST   | false  |
| x2j-help    | /sif-xml2json/help | ../sif-xml2json-0-1-8-linux/server |           | http://127.0.0.1:1324/                     | GET    | false  |
| j2x         | /sif-json2xml      | ../sif-json2xml-0-1-6-linux/server |           | http://127.0.0.1:1325/sif-json2xml/convert | POST   | false  |
| j2x-help    | /sif-json2xml/help | ../sif-json2xml-0-1-6-linux/server |           | http://127.0.0.1:1325/                     | GET    | false  |
| n3w-graphgl | /n3w/graphgl       | ./n3w                              |           | http://127.0.0.1:1323/n3/graphgl           | POST   | true   |
| n3w-publish | /n3w/publish       | ./n3w                              |           | http://127.0.0.1:1323/n3/publish           | POST   | true   |
