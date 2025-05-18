@echo off
setlocal
:loop
set "line="
set /p line=  || goto :eof
echo You said: %line%
goto loop
