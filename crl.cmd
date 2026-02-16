@echo off
cd /D %~dp0
::curl -v --json @20001111.json http://localhost:8088/examupload/20001111/abcd1234
::curl -v --json @20001111.json http://localhost:8080/multipart
::curl -v -i -X POST --header "content-type: application/x-www-form-urlencoded" --data-binary "file=@20001111.json" localhost:8088/examupload/20001111/abcd1234
  ::-H 'content-type: application/x-www-form-urlencoded' ^
echo curly
::curl -v -i -X POST --header "content-type: application/x-www-form-urlencoded" -F "name=file" -F "file=@20001111.json" 192.168.1.222:8088/examupload/20001111/abcd1234
::curl -v -i -X POST --header "Content-Type: multipart/form-data" -F "file=@20001111.json" 192.168.1.222:8088/examupload/20001111/abcd1234

::check the exam list
curl -v http://localhost:8088/examlist

:: test the year list
::curl -v http://localhost:8088/yearlist

:: Check the auth route
::curl -v http://localhost:8088/auth/2026S1ITCS5.100/20001110
::curl -v http://localhost:8088/auth/2026S1ITCS5.100/20001111

::check the exam retrieval route
::curl -v http://localhost:8088/exam/2026S1ITCS5.100/abcd1001
::upload and exam
::curl -v  -F "name=file" -F "exam=@20001111.json" localhost:8088/examupload/20001111/2026S1ITCS5.100/abcd1001

