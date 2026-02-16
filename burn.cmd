@echo off
cd /D %~dp0
::bombardier.exe -c 50 -n 1000 -d 10s -l http://192.168.1.234:8088/
::FOR /L %variable IN (start,step,end) DO command [command-parameters]
::FOR /L %%I IN (1,1,10) DO bombardier.exe -c 50 -n 300  -d 10s -l http://192.168.1.222:8088/exam/2025S1ITCS5.100/abcd1234

::bombardier.exe -c 100 -n 100 -l http://localhost:8088/auth/2026S1ITCS5.100/20001111
bombardier.exe -c 200  -d 10s -l http://localhost:8088/auth/2026S1ITCS5.100/20001111
::bombardier.exe -c 125 -n 1000  -d 10s -l http://localhost:8088/exam/2025S1ITCS5.100/abcd1234
