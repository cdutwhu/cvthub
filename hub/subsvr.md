| NAME(*)  | API_PATH           | PATH_OF_SUB-SERVICE_EXE_PACKAGE | EXE_NAME | REDIRECT                                   | METHOD | ENABLE |
| :------- | :----------------- | :------------------------------ | :------- | :----------------------------------------- | :----- | :----- |
| x2j      | /sif-xml2json      | ../sif-xml2json-0-1-8-linux/    | server   | http://127.0.0.1:1324/sif-xml2json/convert | POST   | true   |
| x2j-help | /sif-xml2json/help | ../sif-xml2json-0-1-8-linux/    | server   | http://127.0.0.1:1324/                     | GET    | true   |
| j2x      | /sif-json2xml      | ../sif-json2xml-0-1-6-linux/    | server   | http://127.0.0.1:1325/sif-json2xml/convert | POST   | true   |
| j2x-help | /sif-json2xml/help | ../sif-json2xml-0-1-6-linux/    | server   | http://127.0.0.1:1325/                     | GET    | true   |
