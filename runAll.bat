@echo off
@rem @call run_pkg_static.bat
@call run_build.bat
start %~dp0server.exe
@ping -n 3 127.0 >nul
