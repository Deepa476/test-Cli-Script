pipeline {
    agent any

    environment {
        MASST_DIR = "MASSTCLI_EXTRACTED"
        ARTIFACTS_DIR = "."
        CONFIG_FILE = "config.bm"
        MASST_ZIP = "MASSTCLI"
    }

    options {
        timestamps()
        buildDiscarder(logRotator(numToKeepStr: '10'))
    }

    stages {
        stage('Checkout SCM') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE 1: Checkout SCM'
                echo '═══════════════════════════════════════════'
                echo 'Checking out source code from repository...'
                checkout scm
                script {
                    def os = isUnix() ? 'Unix/Linux/macOS' : 'Windows'
                    echo "✅ Source code checked out successfully on ${os}"
                }
            }
        }

        stage('Pre-Flight Checks') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE: Pre-Flight Checks'
                echo '═══════════════════════════════════════════'
                script {
                    if (isUnix()) {
                        sh '''
                            set -e
                            echo "Operating System: $(uname -s)"
                            echo "Workspace: ${WORKSPACE}"
                            echo "Build Number: ${BUILD_NUMBER}"
                            echo ""

                            if [ ! -d "${WORKSPACE}" ]; then
                                echo "ERROR: Workspace directory does not exist!"
                                exit 1
                            fi

                            if [ ! -f "${WORKSPACE}/${MASST_ZIP}.zip" ]; then
                                echo "ERROR: ${MASST_ZIP}.zip not found in workspace!"
                                ls -la "${WORKSPACE}" || true
                                exit 1
                            fi

                            echo "✅ Required file found: ${MASST_ZIP}.zip"
                            ls -lh "${WORKSPACE}/${MASST_ZIP}.zip"
                            echo "✅ Workspace validated successfully"
                        '''
                    } else {
                        bat '''
                            echo Operating System: Windows
                            echo Workspace: %WORKSPACE%
                            echo Build Number: %BUILD_NUMBER%

                            if not exist "%WORKSPACE%" (
                                echo ERROR: Workspace directory does not exist!
                                exit /b 1
                            )

                            if not exist "%WORKSPACE%\\%MASST_ZIP%.zip" (
                                echo ERROR: %MASST_ZIP%.zip not found in workspace!
                                dir /b "%WORKSPACE%"
                                exit /b 1
                            )

                            echo ✅ Required file found: %MASST_ZIP%.zip
                            echo ✅ Workspace validated successfully
                        '''
                    }
                }
            }
        }

        stage('Extract MASSTCLI') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE 2: Extract MASSTCLI'
                echo '═══════════════════════════════════════════'
                timeout(time: 5, unit: 'MINUTES') {
                    script {
                        if (isUnix()) {
                            sh '''
                                set -e

                                echo "Checking if extraction is needed..."
                                if [ -d "${MASST_DIR}" ]; then
                                    echo "✅ MASSTCLI already extracted, skipping extraction"
                                    exit 0
                                fi

                                echo ""
                                echo "Extracting MASSTCLI.zip from workspace..."

                                if [ ! -f "${WORKSPACE}/${MASST_ZIP}.zip" ]; then
                                    echo "ERROR: MASSTCLI.zip not found!"
                                    ls -la "${WORKSPACE}" || true
                                    exit 1
                                fi

                                # Create temp directory
                                TEMP_EXTRACT=$(mktemp -d -t masstcli_extract.XXXXXXXXXX)
                                echo "Temp directory: ${TEMP_EXTRACT}"

                                trap "rm -rf ${TEMP_EXTRACT}" EXIT

                                # Extract
                                echo "Extracting archive..."
                                if ! unzip -q "${WORKSPACE}/${MASST_ZIP}.zip" -d "${TEMP_EXTRACT}" 2>/dev/null; then
                                    echo "ERROR: Failed to extract archive!"
                                    exit 1
                                fi

                                # Find and move folder
                                echo "Finding extracted folder..."
                                extracted_folder=$(find "${TEMP_EXTRACT}" -maxdepth 1 -type d ! -name "." | head -n 1)

                                if [ -z "${extracted_folder}" ] || [ "${extracted_folder}" = "${TEMP_EXTRACT}" ]; then
                                    echo "ERROR: No folder found in extracted archive!"
                                    exit 1
                                fi

                                echo "Moving extracted folder to workspace..."
                                mv "${extracted_folder}" "${WORKSPACE}/${MASST_DIR}"

                                echo ""
                                echo "✅ Extraction completed successfully"
                                echo "Destination: ${WORKSPACE}/${MASST_DIR}"
                                ls -la "${WORKSPACE}/${MASST_DIR}" | head -20
                            '''
                        } else {
                            bat '''
                                setlocal enabledelayedexpansion

                                echo Checking if extraction is needed...
                                if exist "%MASST_DIR%" (
                                    echo ✅ MASSTCLI already extracted, skipping extraction
                                    endlocal
                                    exit /b 0
                                )

                                echo.
                                echo Extracting MASSTCLI.zip from workspace...

                                if not exist "%WORKSPACE%\\%MASST_ZIP%.zip" (
                                    echo ERROR: MASSTCLI.zip not found!
                                    dir /b "%WORKSPACE%"
                                    endlocal
                                    exit /b 1
                                )

                                set "TEMP_EXTRACT=C:\\temp\\masstcli_extract_%RANDOM%"

                                echo Preparing temporary extraction location...
                                powershell -Command "Remove-Item '!TEMP_EXTRACT!' -Recurse -Force -ErrorAction SilentlyContinue 2>$null; New-Item -ItemType Directory -Path '!TEMP_EXTRACT!' -Force | Out-Null"

                                echo Extracting archive...
                                powershell -Command "$ProgressPreference = 'SilentlyContinue'; Expand-Archive -LiteralPath '%WORKSPACE%\\%MASST_ZIP%.zip' -DestinationPath '!TEMP_EXTRACT!' -Force -ErrorAction Stop" || (
                                    echo ERROR: Failed to extract archive!
                                    powershell -Command "Remove-Item '!TEMP_EXTRACT!' -Recurse -Force -ErrorAction SilentlyContinue"
                                    endlocal
                                    exit /b 1
                                )

                                echo Finding extracted folder...
                                for /d %%%%d in ("!TEMP_EXTRACT!\\*") do (
                                    echo Moving extracted folder to workspace...
                                    move "%%%%d" "%WORKSPACE%\\%MASST_DIR%" >nul
                                    if errorlevel 1 (
                                        echo ERROR: Failed to move extracted folder!
                                        powershell -Command "Remove-Item '!TEMP_EXTRACT!' -Recurse -Force -ErrorAction SilentlyContinue"
                                        endlocal
                                        exit /b 1
                                    )
                                    goto :extraction_done
                                )

                                echo ERROR: No folder found in extracted archive!
                                powershell -Command "Remove-Item '!TEMP_EXTRACT!' -Recurse -Force -ErrorAction SilentlyContinue"
                                endlocal
                                exit /b 1

                                :extraction_done
                                echo Cleaning up temporary files...
                                powershell -Command "Remove-Item '!TEMP_EXTRACT!' -Recurse -Force -ErrorAction SilentlyContinue"

                                echo.
                                echo ✅ Extraction completed successfully
                                endlocal
                            '''
                        }
                    }
                }
            }
        }

        stage('Verify MASSTCLI') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE 3: Verify MASSTCLI'
                echo '═══════════════════════════════════════════'
                script {
                    if (isUnix()) {
                        sh '''
                            set -e

                            echo "Validating extraction directory..."
                            if [ ! -d "${MASST_DIR}" ]; then
                                echo "ERROR: ${MASST_DIR} directory not found!"
                                exit 1
                            fi

                            echo ""
                            echo "Searching for MASSTCLI executable..."

                            masst_exe=$(find "${MASST_DIR}" -maxdepth 1 \\( -name "MASSTCLI*" -o -name "masstcli*" \\) -type f -executable 2>/dev/null | head -n 1)

                            if [ -z "${masst_exe}" ]; then
                                echo "WARNING: No executable found, checking for binary files..."
                                potential_exe=$(find "${MASST_DIR}" -maxdepth 1 \\( -name "MASSTCLI*" -o -name "masstcli*" \\) -type f 2>/dev/null | head -n 1)

                                if [ -n "${potential_exe}" ]; then
                                    echo "Found potential executable, setting permissions..."
                                    chmod +x "${potential_exe}"
                                    masst_exe="${potential_exe}"
                                else
                                    echo "ERROR: No MASSTCLI executable found in ${MASST_DIR}!"
                                    find "${MASST_DIR}" -type f | head -20
                                    exit 1
                                fi
                            fi

                            masst_exe_name=$(basename "${masst_exe}")
                            echo "Found: ${masst_exe_name}"

                            echo ""
                            echo "Validating executable..."
                            if [ ! -x "${masst_exe}" ]; then
                                echo "ERROR: Executable not accessible!"
                                exit 1
                            fi

                            echo ""
                            echo "════════════════════════════════════════════════════"
                            echo "✅ MASSTCLI verified successfully"
                            echo "════════════════════════════════════════════════════"
                            echo "Executable: ${masst_exe_name}"
                            echo "Location: ${MASST_DIR}"
                            echo "Status: Ready for execution"
                            echo "════════════════════════════════════════════════════"
                        '''
                    } else {
                        bat '''
                            setlocal enabledelayedexpansion

                            echo Validating extraction directory...
                            if not exist "%MASST_DIR%" (
                                echo ERROR: %MASST_DIR% directory not found!
                                endlocal
                                exit /b 1
                            )

                            echo.
                            echo Searching for MASSTCLI executable...
                            set MASST_EXE=

                            for %%f in ("%MASST_DIR%\\MASSTCLI*.exe") do (
                                set MASST_EXE=%%~nxf
                                echo Found: %%~nxf
                            )

                            if not defined MASST_EXE (
                                echo ERROR: No MASSTCLI executable found in %MASST_DIR%!
                                echo Contents of extraction directory:
                                dir /s "%MASST_DIR%"
                                endlocal
                                exit /b 1
                            )

                            echo.
                            echo Validating executable file...
                            if not exist "%MASST_DIR%\\!MASST_EXE!" (
                                echo ERROR: Executable path not accessible!
                                endlocal
                                exit /b 1
                            )

                            echo.
                            echo ════════════════════════════════════════════════════
                            echo ✅ MASSTCLI verified successfully
                            echo ════════════════════════════════════════════════════
                            echo Executable: !MASST_EXE!
                            echo Location: %MASST_DIR%
                            echo Status: Ready for execution
                            echo ════════════════════════════════════════════════════
                            endlocal
                        '''
                    }
                }
            }
        }

        stage('Validate Input Files') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE: Validate Input Files'
                echo '═══════════════════════════════════════════'
                script {
                    if (isUnix()) {
                        sh '''
                            set -e

                            echo "Validating required input files..."
                            echo ""

                            echo "Checking MyApp.aab..."
                            if [ ! -f "${WORKSPACE}/MyApp.aab" ]; then
                                echo "ERROR: Input file MyApp.aab not found!"
                                echo "Expected: ${WORKSPACE}/MyApp.aab"
                                ls -la "${WORKSPACE}" | head -30
                                exit 1
                            fi

                            file_size=$(stat -c%s "${WORKSPACE}/MyApp.aab" 2>/dev/null || stat -f%z "${WORKSPACE}/MyApp.aab" 2>/dev/null || echo "N/A")
                            echo "✅ Found: MyApp.aab (${file_size} bytes)"

                            echo ""
                            echo "Checking ${CONFIG_FILE}..."
                            if [ ! -f "${WORKSPACE}/${CONFIG_FILE}" ]; then
                                echo "ERROR: Configuration file ${CONFIG_FILE} not found!"
                                exit 1
                            fi

                            file_size=$(stat -c%s "${WORKSPACE}/${CONFIG_FILE}" 2>/dev/null || stat -f%z "${WORKSPACE}/${CONFIG_FILE}" 2>/dev/null || echo "N/A")
                            echo "✅ Found: ${CONFIG_FILE} (${file_size} bytes)"

                            echo ""
                            echo "✅ All input files validated successfully"
                        '''
                    } else {
                        bat '''
                            echo Validating required input files...
                            echo.

                            echo Checking MyApp.aab...
                            if not exist "%WORKSPACE%\\MyApp.aab" (
                                echo ERROR: Input file MyApp.aab not found!
                                echo Expected: %WORKSPACE%\\MyApp.aab
                                dir /b "%WORKSPACE%"
                                exit /b 1
                            )

                            echo ✅ Found: MyApp.aab

                            echo.
                            echo Checking %CONFIG_FILE%...
                            if not exist "%WORKSPACE%\\%CONFIG_FILE%" (
                                echo ERROR: Configuration file %CONFIG_FILE% not found!
                                exit /b 1
                            )

                            echo ✅ Found: %CONFIG_FILE%

                            echo.
                            echo ✅ All input files validated successfully
                        '''
                    }
                }
            }
        }

        stage('Run MASSTCLI') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE 4: Run MASSTCLI'
                echo '═══════════════════════════════════════════'
                script {
                    if (isUnix()) {
                        sh '''
                            set -e

                            echo "Preparing to execute MASSTCLI..."
                            echo ""

                            masst_exe=$(find "${MASST_DIR}" -maxdepth 1 \\( -name "MASSTCLI*" -o -name "masstcli*" \\) -type f -executable 2>/dev/null | head -n 1)

                            if [ -z "${masst_exe}" ]; then
                                echo "ERROR: No MASSTCLI executable found!"
                                exit 1
                            fi

                            if [ ! -f "${WORKSPACE}/MyApp.aab" ]; then
                                echo "ERROR: Input file MyApp.aab not found!"
                                exit 1
                            fi

                            if [ ! -f "${WORKSPACE}/${CONFIG_FILE}" ]; then
                                echo "ERROR: Configuration file ${CONFIG_FILE} not found!"
                                exit 1
                            fi

                            echo "════════════════════════════════════════════════════"
                            echo "Execution Details:"
                            echo "════════════════════════════════════════════════════"
                            echo "Executable: $(basename ${masst_exe})"
                            echo "Path: ${masst_exe}"
                            echo "Input: ${WORKSPACE}/MyApp.aab"
                            echo "Config: ${WORKSPACE}/${CONFIG_FILE}"
                            echo "════════════════════════════════════════════════════"
                            echo ""

                            echo "Executing MASSTCLI..."
                            "${masst_exe}" -input "${WORKSPACE}/MyApp.aab" -config "${WORKSPACE}/${CONFIG_FILE}"

                            exit_code=$?
                            echo ""

                            if [ ${exit_code} -ne 0 ]; then
                                echo "ERROR: MASSTCLI execution failed with exit code: ${exit_code}"
                                exit ${exit_code}
                            fi

                            echo "✅ MASSTCLI execution completed successfully"
                        '''
                    } else {
                        bat '''
                            setlocal enabledelayedexpansion

                            echo Preparing to execute MASSTCLI...
                            echo.

                            set MASST_EXE=
                            for %%f in ("%MASST_DIR%\\MASSTCLI*.exe") do (
                                set MASST_EXE=%%~nxf
                            )

                            if not defined MASST_EXE (
                                echo ERROR: No MASSTCLI executable found!
                                endlocal
                                exit /b 1
                            )

                            if not exist "%WORKSPACE%\\MyApp.aab" (
                                echo ERROR: Input file MyApp.aab not found!
                                endlocal
                                exit /b 1
                            )

                            if not exist "%WORKSPACE%\\%CONFIG_FILE%" (
                                echo ERROR: Configuration file %CONFIG_FILE% not found!
                                endlocal
                                exit /b 1
                            )

                            echo ════════════════════════════════════════════════════
                            echo Execution Details:
                            echo ════════════════════════════════════════════════════
                            echo Executable: !MASST_EXE!
                            echo Path: %MASST_DIR%\\!MASST_EXE!
                            echo Input: %WORKSPACE%\\MyApp.aab
                            echo Config: %WORKSPACE%\\%CONFIG_FILE%
                            echo ════════════════════════════════════════════════════
                            echo.

                            echo Executing MASSTCLI...
                            "%MASST_DIR%\\!MASST_EXE!" -input "%WORKSPACE%\\MyApp.aab" -config "%WORKSPACE%\\%CONFIG_FILE%"

                            if errorlevel 1 (
                                echo.
                                echo ERROR: MASSTCLI execution failed with error code: %ERRORLEVEL%
                                endlocal
                                exit /b 1
                            )

                            echo.
                            echo ✅ MASSTCLI execution completed successfully
                            endlocal
                        '''
                    }
                }
            }
        }
    }

    post {
        success {
            echo '═══════════════════════════════════════════'
            echo '✅ PIPELINE COMPLETED SUCCESSFULLY'
            echo '═══════════════════════════════════════════'
            script {
                if (isUnix()) {
                    sh '''
                        echo ""
                        echo "Build Status: SUCCESS"
                        echo "Job: ${JOB_NAME}"
                        echo "Build: ${BUILD_NUMBER}"
                        echo ""
                    '''
                } else {
                    bat '''
                        echo.
                        echo Build Status: SUCCESS
                        echo Job: %JOB_NAME%
                        echo Build: %BUILD_NUMBER%
                        echo.
                    '''
                }
            }
        }
        failure {
            echo '═══════════════════════════════════════════'
            echo '❌ PIPELINE FAILED'
            echo '═══════════════════════════════════════════'
            script {
                if (isUnix()) {
                    sh '''
                        echo ""
                        echo "Build Status: FAILURE"
                        echo "Job: ${JOB_NAME}"
                        echo "Build: ${BUILD_NUMBER}"
                        echo "Check the logs above for error details."
                        echo ""
                    '''
                } else {
                    bat '''
                        echo.
                        echo Build Status: FAILURE
                        echo Job: %JOB_NAME%
                        echo Build: %BUILD_NUMBER%
                        echo Check the logs above for error details.
                        echo.
                    '''
                }
            }
        }
        always {
            echo '═══════════════════════════════════════════'
            echo 'Cleanup Phase'
            echo '═══════════════════════════════════════════'
            script {
                if (isUnix()) {
                    sh '''
                        echo "Cleaning up temporary files..."
                        find /tmp -maxdepth 1 -name "masstcli_extract*" -type d -exec rm -rf {} + 2>/dev/null || true
                        echo "✅ Cleanup completed"
                    '''
                } else {
                    bat '''
                        echo Cleaning up temporary files...
                        powershell -Command "Get-ChildItem 'C:\\temp\\masstcli_extract*' -ErrorAction SilentlyContinue | Remove-Item -Recurse -Force"
                        echo ✅ Cleanup completed
                    '''
                }
            }
        }
    }
}