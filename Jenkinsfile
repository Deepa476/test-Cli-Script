pipeline {
    agent any

    environment {
        MASST_DIR = "MASSTCLI_EXTRACTED"
        ARTIFACTS_DIR = "."
        CONFIG_FILE = "config.bm"
        MASST_ZIP = "MASSTCLI"
        DOWNLOAD_URL = "https://storage.googleapis.com/masst-assets/Defender-Binary-Integrator/1.0.0/MacOS/MASSTCLI-v1.1.0-darwin-amd64.zip"
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

        stage('Prepare MASSTCLI Tool') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE: Prepare MASSTCLI Tool'
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

                            echo "Checking for MASSTCLI.zip..."
                            if [ -f "${WORKSPACE}/${MASST_ZIP}.zip" ]; then
                                echo "✅ MASSTCLI.zip found in workspace"
                                ls -lh "${WORKSPACE}/${MASST_ZIP}.zip"
                                exit 0
                            fi

                            echo "⚠ MASSTCLI.zip not found in workspace"
                            echo "Attempting to download from remote source..."
                            echo ""

                            echo "Download URL: ${DOWNLOAD_URL}"
                            echo "Destination: ${WORKSPACE}/${MASST_ZIP}.zip"
                            echo ""

                            # Try curl first, then wget
                            if command -v curl &> /dev/null; then
                                echo "Using curl to download..."
                                curl -L --progress-bar -o "${WORKSPACE}/${MASST_ZIP}.zip" "${DOWNLOAD_URL}" || {
                                    echo "ERROR: Failed to download MASSTCLI.zip using curl!"
                                    exit 1
                                }
                            elif command -v wget &> /dev/null; then
                                echo "Using wget to download..."
                                wget -O "${WORKSPACE}/${MASST_ZIP}.zip" "${DOWNLOAD_URL}" || {
                                    echo "ERROR: Failed to download MASSTCLI.zip using wget!"
                                    exit 1
                                }
                            else
                                echo "ERROR: Neither curl nor wget found!"
                                echo "Cannot download MASSTCLI.zip"
                                exit 1
                            fi

                            echo ""
                            echo "✅ MASSTCLI.zip downloaded successfully"
                            ls -lh "${WORKSPACE}/${MASST_ZIP}.zip"
                        '''
                    } else {
                        bat '''
                            setlocal enabledelayedexpansion

                            echo Operating System: Windows
                            echo Workspace: %WORKSPACE%
                            echo Build Number: %BUILD_NUMBER%
                            echo.

                            if not exist "%WORKSPACE%" (
                                echo ERROR: Workspace directory does not exist!
                                endlocal
                                exit /b 1
                            )

                            echo Checking for MASSTCLI.zip...
                            if exist "%WORKSPACE%\\%MASST_ZIP%.zip" (
                                echo ✅ MASSTCLI.zip found in workspace
                                for %%%%F in ("%WORKSPACE%\\%MASST_ZIP%.zip") do (
                                    echo    Size: %%%%~zF bytes
                                )
                                endlocal
                                exit /b 0
                            )

                            echo ⚠ MASSTCLI.zip not found in workspace
                            echo Attempting to download from remote source...
                            echo.

                            set "DOWNLOAD_PATH=%WORKSPACE%\\%MASST_ZIP%.zip"

                            echo Download URL: %DOWNLOAD_URL%
                            echo Destination: !DOWNLOAD_PATH!
                            echo.

                            powershell -Command "$ProgressPreference = 'Continue'; Invoke-WebRequest -Uri '%DOWNLOAD_URL%' -OutFile '!DOWNLOAD_PATH!' -ErrorAction Stop" || (
                                echo ERROR: Failed to download MASSTCLI.zip!
                                echo Please check:
                                echo   1. Internet connection is available
                                echo   2. URL is accessible: %DOWNLOAD_URL%
                                endlocal
                                exit /b 1
                            )

                            echo.
                            echo ✅ MASSTCLI.zip downloaded successfully
                            for %%%%F in ("!DOWNLOAD_PATH!") do (
                                echo    Size: %%%%~zF bytes
                            )

                            endlocal
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

                                # Extract
                                echo "Extracting archive..."
                                if ! unzip -q "${WORKSPACE}/${MASST_ZIP}.zip" -d "${TEMP_EXTRACT}" 2>/dev/null; then
                                    echo "ERROR: Failed to extract archive!"
                                    exit 1
                                fi

                                # Find and move folder or files
                                echo "Finding extracted contents..."

                                # Count items in temp directory
                                item_count=$(find "${TEMP_EXTRACT}" -maxdepth 1 ! -name . -type d | wc -l)

                                if [ "${item_count}" -eq 1 ]; then
                                    # Single folder found - move it
                                    extracted_folder=$(find "${TEMP_EXTRACT}" -maxdepth 1 -type d ! -name . | head -n 1)
                                    echo "Found folder: $(basename "${extracted_folder}")"
                                    echo "Moving extracted folder to workspace..."
                                    mv "${extracted_folder}" "${WORKSPACE}/${MASST_DIR}"
                                elif [ "${item_count}" -gt 1 ] || [ $(find "${TEMP_EXTRACT}" -maxdepth 1 ! -name . -type f | wc -l) -gt 0 ]; then
                                    # Multiple items or files found - move the entire temp directory contents
                                    echo "Found multiple items in archive"
                                    echo "Moving extracted contents to workspace..."
                                    mkdir -p "${WORKSPACE}/${MASST_DIR}"
                                    mv "${TEMP_EXTRACT}"/* "${WORKSPACE}/${MASST_DIR}/" 2>/dev/null || true
                                    # Also move hidden files if any
                                    mv "${TEMP_EXTRACT}"/.[^.]* "${WORKSPACE}/${MASST_DIR}/" 2>/dev/null || true
                                else
                                    echo "ERROR: No contents found in extracted archive!"
                                    exit 1
                                fi

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
                                    endlocal
                                    exit /b 1
                                )

                                set "TEMP_EXTRACT=C:\\temp\\masstcli_extract_%RANDOM%"

                                echo Preparing temporary extraction location...
                                powershell -Command "New-Item -ItemType Directory -Path '!TEMP_EXTRACT!' -Force | Out-Null"

                                if not exist "!TEMP_EXTRACT!" (
                                    echo ERROR: Failed to create temporary extraction directory!
                                    endlocal
                                    exit /b 1
                                )

                                echo Extracting archive...
                                powershell -Command "$ProgressPreference = 'SilentlyContinue'; Expand-Archive -LiteralPath '%WORKSPACE%\\%MASST_ZIP%.zip' -DestinationPath '!TEMP_EXTRACT!' -Force -ErrorAction Stop" || (
                                    echo ERROR: Failed to extract archive!
                                    endlocal
                                    exit /b 1
                                )

                                echo Finding extracted folder...
                                for /d %%%%d in ("!TEMP_EXTRACT!\\*") do (
                                    echo Moving extracted folder to workspace...
                                    move "%%%%d" "%WORKSPACE%\\%MASST_DIR%" >nul
                                    if errorlevel 1 (
                                        echo ERROR: Failed to move extracted folder!
                                        endlocal
                                        exit /b 1
                                    )
                                    goto :extraction_done
                                )

                                echo ERROR: No folder found in extracted archive!
                                endlocal
                                exit /b 1

                                :extraction_done
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

                            # Search recursively without maxdepth limit
                            masst_exe=$(find "${MASST_DIR}" -type f -name "MASSTCLI*" 2>/dev/null | grep -v ".zip" | head -n 1)

                            if [ -z "${masst_exe}" ]; then
                                echo "ERROR: No MASSTCLI executable found!"
                                echo "Contents of ${MASST_DIR}:"
                                find "${MASST_DIR}" -type f -name "MASSTCLI*" 2>/dev/null || true
                                exit 1
                            fi

                            # Make executable if not already
                            chmod +x "${masst_exe}"

                            masst_exe_name=$(basename "${masst_exe}")
                            echo "Found: ${masst_exe_name}"

                            echo ""
                            echo "════════════════════════════════════════════════════"
                            echo "✅ MASSTCLI verified successfully"
                            echo "════════════════════════════════════════════════════"
                            echo "Executable: ${masst_exe_name}"
                            echo "Full Path: ${masst_exe}"
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

                            for /r "%MASST_DIR%" %%f in (MASSTCLI*.exe) do (
                                set MASST_EXE=%%~nxf
                                set MASST_PATH=%%~dpnxf
                                goto :found_exe
                            )

                            echo ERROR: No MASSTCLI executable found in %MASST_DIR%!
                            endlocal
                            exit /b 1

                            :found_exe
                            echo ✅ Found: !MASST_EXE!

                            echo.
                            echo ════════════════════════════════════════════════════
                            echo ✅ MASSTCLI verified successfully
                            echo ════════════════════════════════════════════════════
                            echo Executable: !MASST_EXE!
                            echo Full Path: !MASST_PATH!
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

                           # Search recursively for the executable
                           masst_exe=$(find "${MASST_DIR}" -type f -name "MASSTCLI*" \\( -executable -o -name "*darwin-amd64" -o -name "*darwin*" \\) 2>/dev/null | grep -v ".zip" | head -n 1)

                           if [ -z "${masst_exe}" ]; then
                               echo "ERROR: No MASSTCLI executable found!"
                               echo "Searching in ${MASST_DIR}..."
                               find "${MASST_DIR}" -type f -name "MASSTCLI*" 2>/dev/null || true
                               exit 1
                           fi

                           # Make executable if not already
                           chmod +x "${masst_exe}"

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
                           for /r "%MASST_DIR%" %%f in (MASSTCLI*.exe) do (
                               set MASST_EXE=%%~nxf
                               set MASST_PATH=%%~dpnxf
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
                           echo Path: !MASST_PATH!
                           echo Input: %WORKSPACE%\\MyApp.aab
                           echo Config: %WORKSPACE%\\%CONFIG_FILE%
                           echo ════════════════════════════════════════════════════
                           echo.

                           echo Executing MASSTCLI...
                           "!MASST_PATH!" -input "%WORKSPACE%\\MyApp.aab" -config "%WORKSPACE%\\%CONFIG_FILE%"

                           if errorlevel 1 (
                               echo.
                               echo ERROR: MASSTCLI execution failed with error code: %ERRORLEVEL%
                               endlocal
                               exit /b 1
                           )

                           echo.
                           echo ✅ MASSTCLI execution completed
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
            echo 'Pipeline Complete - Files Preserved'
            echo '═══════════════════════════════════════════'
            script {
                if (isUnix()) {
                    sh '''
                        echo ""
                        echo "All files have been preserved in workspace:"
                        echo "  • MASSTCLI.zip"
                        echo "  • MASSTCLI_EXTRACTED (extracted files)"
                        echo "  • MyApp.aab"
                        echo "  • config.bm"
                        echo ""
                        echo "You can inspect or reuse these files for debugging or subsequent builds."
                        echo ""
                    '''
                } else {
                    bat '''
                        echo.
                        echo All files have been preserved in workspace:
                        echo   • MASSTCLI.zip
                        echo   • MASSTCLI_EXTRACTED (extracted files)
                        echo   • MyApp.aab
                        echo   • config.bm
                        echo.
                        echo You can inspect or reuse these files for debugging or subsequent builds.
                        echo.
                    '''
                }
            }
        }
    }
}