@echo off


for %%i in (*.proto) do (
	set prtofile=%%i
	
	protoc --go_out=. %%i
	echo convert !prtofile %%i  ok
)