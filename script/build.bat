::============================================================================::
:: This script must be started from its folder ::
::============================================================================::
@ECHO OFF

SET BUILD_DIR=_BUILD_
SET APP_EXE=Simpel_Chat_Server.exe
SET SCRIPT_FOLDER=script

FOR %%f IN ("%CD%") DO SET LastPathElement=%%~nxf
SET CUR_FOLDER_NAME=%LastPathElement%
IF "%CUR_FOLDER_NAME%" == "%SCRIPT_FOLDER%" ( ECHO Welcome ) ELSE (
    ECHO This script must be started from its folder. Press any key to exit.
    EXIT /B 1
)

:: CD to root folder.
CD ..
MKDIR %BUILD_DIR%

:: Build the executable file.
CD src\
go build -o ..\%BUILD_DIR%\%APP_EXE%
IF %ERRORLEVEL% NEQ 0 EXIT /B %ERRORLEVEL%
CD ..

:: Copy the SSL certificates.
XCOPY "certificate" "%BUILD_DIR%\certificate" /S/I/Q

:: Copy the settings.
::XCOPY settings %BUILD_DIR%\settings /S/I/Q
MKDIR "%BUILD_DIR%\settings"
COPY "settings\app.cfg" "%BUILD_DIR%\settings\"
COPY "settings\chat.json" "%BUILD_DIR%\settings\"

:: Copy the starter script
COPY "script\starter.bat" "%BUILD_DIR%"

CD "%SCRIPT_FOLDER%"
