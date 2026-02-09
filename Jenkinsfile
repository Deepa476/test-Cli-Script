pipeline {
    agent any

    environment {
        MASST_DIR = "MASSTCLI_EXTRACTED"
        ARTIFACTS_DIR = "output"
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

                                TEMP_EXTRACT=$(mktemp -d -t masstcli_extract.XXXXXXXXXX)
                                echo "Temp directory: ${TEMP_EXTRACT}"

                                echo "Extracting archive..."
                                if ! unzip -q "${WORKSPACE}/${MASST_ZIP}.zip" -d "${TEMP_EXTRACT}" 2>/dev/null; then
                                    echo "ERROR: Failed to extract archive!"
                                    exit 1
                                fi

                                echo "Finding extracted contents..."

                                item_count=$(find "${TEMP_EXTRACT}" -mindepth 1 -maxdepth 1 -type d | wc -l)

                                if [ "${item_count}" -eq 1 ]; then
                                    extracted_folder=$(find "${TEMP_EXTRACT}" -mindepth 1 -maxdepth 1 -type d | head -n 1)
                                    echo "Found folder: $(basename "${extracted_folder}")"
                                    echo "Moving extracted folder to workspace..."
                                    mv "${extracted_folder}" "${WORKSPACE}/${MASST_DIR}"
                                elif [ "${item_count}" -gt 1 ] || [ $(find "${TEMP_EXTRACT}" -mindepth 1 -maxdepth 1 -type f | wc -l) -gt 0 ]; then
                                    echo "Found multiple items in archive"
                                    echo "Moving extracted contents to workspace..."
                                    mkdir -p "${WORKSPACE}/${MASST_DIR}"
                                    mv "${TEMP_EXTRACT}"/* "${WORKSPACE}/${MASST_DIR}/" 2>/dev/null || true
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

                            masst_exe=$(find "${MASST_DIR}" -type f -name "MASSTCLI*" 2>/dev/null | grep -v ".zip" | head -n 1)

                            if [ -z "${masst_exe}" ]; then
                                echo "ERROR: No MASSTCLI executable found!"
                                echo "Contents of ${MASST_DIR}:"
                                find "${MASST_DIR}" -type f -name "MASSTCLI*" 2>/dev/null || true
                                exit 1
                            fi

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
                            set MASST_PATH=

                            for /r "%MASST_DIR%" %%%%f in (MASSTCLI*.exe) do (
                                set "MASST_EXE=%%%%~nxf"
                                set "MASST_PATH=%%%%~dpF%%%%~nxf"
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

        stage('Prepare Output Directory') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE: Prepare Output Directory'
                echo '═══════════════════════════════════════════'
                script {
                    if (isUnix()) {
                        sh '''
                            set -e

                            echo "Creating output directory..."
                            mkdir -p "${WORKSPACE}/${ARTIFACTS_DIR}"
                            echo "✅ Output directory created: ${WORKSPACE}/${ARTIFACTS_DIR}"
                            ls -la "${WORKSPACE}/${ARTIFACTS_DIR}"
                        '''
                    } else {
                        bat '''
                            echo Creating output directory...
                            if not exist "%WORKSPACE%\\%ARTIFACTS_DIR%" mkdir "%WORKSPACE%\\%ARTIFACTS_DIR%"
                            echo ✅ Output directory created: %WORKSPACE%\\%ARTIFACTS_DIR%
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

                            # Find the MASSTCLI executable
                            masst_exe=$(find "${MASST_DIR}" -type f -name "MASSTCLI*" 2>/dev/null | grep -v ".zip" | head -n 1)

                            if [ -z "${masst_exe}" ]; then
                                echo "ERROR: No MASSTCLI executable found!"
                                echo "Contents of ${MASST_DIR}:"
                                find "${MASST_DIR}" -type f 2>/dev/null | head -20
                                exit 1
                            fi

                            # Make it executable
                            chmod +x "${masst_exe}"

                            # Validate input files exist
                            if [ ! -f "${WORKSPACE}/MyApp.aab" ]; then
                                echo "ERROR: Input file MyApp.aab not found!"
                                exit 1
                            fi

                            if [ ! -f "${WORKSPACE}/${CONFIG_FILE}" ]; then
                                echo "ERROR: Configuration file ${CONFIG_FILE} not found!"
                                exit 1
                            fi

                            # Create output directory if it doesn't exist
                            mkdir -p "${WORKSPACE}/${ARTIFACTS_DIR}"

                            # Define output file path
                            OUTPUT_FILE="${WORKSPACE}/${ARTIFACTS_DIR}/MyApp_obfuscated.aab"

                            echo "════════════════════════════════════════════════════"
                            echo "Execution Details:"
                            echo "════════════════════════════════════════════════════"
                            echo "Executable: $(basename ${masst_exe})"
                            echo "Full Path: ${masst_exe}"
                            echo "Input File: ${WORKSPACE}/MyApp.aab"
                            echo "Config File: ${WORKSPACE}/${CONFIG_FILE}"
                            echo "Output File: ${OUTPUT_FILE}"
                            echo "════════════════════════════════════════════════════"
                            echo ""

                            echo "Starting MASSTCLI execution..."
                            echo ""

                            # Execute MASSTCLI
                            if "${masst_exe}" -input "${WORKSPACE}/MyApp.aab" -config "${WORKSPACE}/${CONFIG_FILE}" -output "${OUTPUT_FILE}"; then
                                echo ""
                                echo "✅ MASSTCLI execution completed successfully"

                                # Verify output file was created
                                if [ -f "${OUTPUT_FILE}" ]; then
                                    output_size=$(stat -f%z "${OUTPUT_FILE}" 2>/dev/null || stat -c%s "${OUTPUT_FILE}" 2>/dev/null || echo "N/A")
                                    echo "✅ Output file created: $(basename ${OUTPUT_FILE})"
                                    echo "   Size: ${output_size} bytes"
                                else
                                    echo "⚠ Warning: Output file not found at expected location"
                                    echo "Checking for any generated files..."
                                    find "${WORKSPACE}/${ARTIFACTS_DIR}" -type f 2>/dev/null || echo "No files found in output directory"
                                fi
                            else
                                exit_code=$?
                                echo ""
                                echo "❌ MASSTCLI execution failed with exit code: ${exit_code}"
                                exit ${exit_code}
                            fi
                        '''
                    } else {
                        bat '''
                            setlocal enabledelayedexpansion

                            echo Preparing to execute MASSTCLI...
                            echo.

                            REM Find MASSTCLI executable
                            set "MASST_EXE="
                            for /r "%MASST_DIR%" %%%%f in (MASSTCLI*.exe) do (
                                set "MASST_EXE=%%%%f"
                                goto :found_exe
                            )

                            echo ERROR: No MASSTCLI executable found in %MASST_DIR%!
                            endlocal
                            exit /b 1

                            :found_exe
                            if not exist "%WORKSPACE%\MyApp.aab" (
                                echo ERROR: Input file MyApp.aab not found!
                                endlocal
                                exit /b 1
                            )

                            if not exist "%WORKSPACE%\%CONFIG_FILE%" (
                                echo ERROR: Configuration file %CONFIG_FILE% not found!
                                endlocal
                                exit /b 1
                            )

                            REM Create output directory
                            if not exist "%WORKSPACE%\%ARTIFACTS_DIR%" mkdir "%WORKSPACE%\%ARTIFACTS_DIR%"

                            set "OUTPUT_FILE=%WORKSPACE%\%ARTIFACTS_DIR%\MyApp_obfuscated.aab"

                            echo ════════════════════════════════════════════════════
                            echo Execution Details:
                            echo ════════════════════════════════════════════════════
                            echo Executable: !MASST_EXE!
                            echo Input File: %WORKSPACE%\MyApp.aab
                            echo Config File: %WORKSPACE%\%CONFIG_FILE%
                            echo Output File: !OUTPUT_FILE!
                            echo ════════════════════════════════════════════════════
                            echo.

                            echo Starting MASSTCLI execution...
                            echo.

                            "!MASST_EXE!" -input "%WORKSPACE%\MyApp.aab" -config "%WORKSPACE%\%CONFIG_FILE%" -output "!OUTPUT_FILE!"

                            if errorlevel 1 (
                                echo.
                                echo ❌ MASSTCLI execution failed with error code: !ERRORLEVEL!
                                endlocal
                                exit /b 1
                            )

                            echo.
                            echo ✅ MASSTCLI execution completed successfully

                            if exist "!OUTPUT_FILE!" (
                                for %%%%F in ("!OUTPUT_FILE!") do (
                                    echo ✅ Output file created: %%%%~nxF
                                    echo    Size: %%%%~zF bytes
                                )
                            ) else (
                                echo ⚠ Warning: Output file not found at expected location
                                echo Checking for any generated files...
                                dir "%WORKSPACE%\%ARTIFACTS_DIR%" 2>nul || echo No files found in output directory
                            )

                            endlocal
                        '''
                    }
                }
            }
        }


        stage('Archive Artifacts') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE 5: Archive Artifacts'
                echo '═══════════════════════════════════════════'
                script {
                    if (isUnix()) {
                        sh '''
                            echo "Archiving output artifacts..."
                            echo ""

                            if [ ! -d "${WORKSPACE}/${ARTIFACTS_DIR}" ]; then
                                echo "ERROR: Artifacts directory not found!"
                                exit 1
                            fi

                            echo "Contents of ${ARTIFACTS_DIR}:"
                            find "${WORKSPACE}/${ARTIFACTS_DIR}" -type f 2>/dev/null | while read file; do
                                filename=$(basename "${file}")
                                filesize=$(stat -f%z "${file}" 2>/dev/null || stat -c%s "${file}" 2>/dev/null || echo "N/A")
                                echo "  ✅ ${filename} (${filesize} bytes)"
                            done

                            echo ""
                            echo "✅ Artifacts ready for archival"
                        '''
                    } else {
                        bat '''
                            echo Archiving output artifacts...
                            echo.

                            if not exist "%WORKSPACE%\\%ARTIFACTS_DIR%" (
                                echo ERROR: Artifacts directory not found!
                                exit /b 1
                            )

                            echo Contents of %ARTIFACTS_DIR%:
                            for %%%%f in (%WORKSPACE%\\%ARTIFACTS_DIR%\\*) do (
                                echo   ✅ %%%%~nxf
                            )

                            echo.
                            echo ✅ Artifacts ready for archival
                        '''
                    }
                }
                archiveArtifacts artifacts: 'output/**', allowEmptyArchive: false, fingerprint: true, onlyIfSuccessful: true
            }
        }

        stage('Generate Report') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE 6: Generate Report'
                echo '═══════════════════════════════════════════'
                script {
                    if (isUnix()) {
                        sh '''
                            set -e

                            REPORT_FILE="${WORKSPACE}/${ARTIFACTS_DIR}/build_report.txt"

                            echo "Generating build report..."
                            {
                                echo "════════════════════════════════════════════════════"
                                echo "MASSTCLI Build Report"
                                echo "════════════════════════════════════════════════════"
                                echo ""
                                echo "Job Name: ${JOB_NAME}"
                                echo "Build Number: ${BUILD_NUMBER}"
                                echo "Build Timestamp: $(date '+%Y-%m-%d %H:%M:%S')"
                                echo "Workspace: ${WORKSPACE}"
                                echo ""
                                echo "════════════════════════════════════════════════════"
                                echo "Input Files"
                                echo "════════════════════════════════════════════════════"
                                echo "Config File: ${CONFIG_FILE}"
                                if [ -f "${WORKSPACE}/${CONFIG_FILE}" ]; then
                                    echo "Config File Size: $(stat -f%z "${WORKSPACE}/${CONFIG_FILE}" 2>/dev/null || stat -c%s "${WORKSPACE}/${CONFIG_FILE}" 2>/dev/null) bytes"
                                fi
                                echo ""
                                echo "Input File: MyApp.aab"
                                if [ -f "${WORKSPACE}/MyApp.aab" ]; then
                                    echo "Input File Size: $(stat -f%z "${WORKSPACE}/MyApp.aab" 2>/dev/null || stat -c%s "${WORKSPACE}/MyApp.aab" 2>/dev/null) bytes"
                                fi
                                echo ""
                                echo "════════════════════════════════════════════════════"
                                echo "Output Files"
                                echo "════════════════════════════════════════════════════"
                                find "${WORKSPACE}/${ARTIFACTS_DIR}" -type f ! -name "build_report.txt" | while read file; do
                                    filename=$(basename "${file}")
                                    filesize=$(stat -f%z "${file}" 2>/dev/null || stat -c%s "${file}" 2>/dev/null)
                                    echo "File: ${filename}"
                                    echo "Size: ${filesize} bytes"
                                    echo ""
                                done
                                echo "════════════════════════════════════════════════════"
                                echo "Build Status: SUCCESS"
                                echo "════════════════════════════════════════════════════"
                            } > "${REPORT_FILE}"

                            echo "✅ Build report generated: ${REPORT_FILE}"
                            echo ""
                            cat "${REPORT_FILE}"
                        '''
                    } else {
                        bat '''
                            setlocal enabledelayedexpansion

                            set REPORT_FILE=%WORKSPACE%\\%ARTIFACTS_DIR%\\build_report.txt

                            echo Generating build report...
                            (
                                echo ════════════════════════════════════════════════════
                                echo MASSTCLI Build Report
                                echo ════════════════════════════════════════════════════
                                echo.
                                echo Job Name: %JOB_NAME%
                                echo Build Number: %BUILD_NUMBER%
                                echo Workspace: %WORKSPACE%
                                echo.
                                echo ════════════════════════════════════════════════════
                                echo Input Files
                                echo ════════════════════════════════════════════════════
                                echo Config File: %CONFIG_FILE%
                                if exist "%WORKSPACE%\\%CONFIG_FILE%" (
                                    for %%%%F in ("%WORKSPACE%\\%CONFIG_FILE%") do echo Config File Size: %%%%~zF bytes
                                )
                                echo.
                                echo Input File: MyApp.aab
                                if exist "%WORKSPACE%\\MyApp.aab" (
                                    for %%%%F in ("%WORKSPACE%\\MyApp.aab") do echo Input File Size: %%%%~zF bytes
                                )
                                echo.
                                echo ════════════════════════════════════════════════════
                                echo Output Files
                                echo ════════════════════════════════════════════════════
                                for %%%%f in (%WORKSPACE%\\%ARTIFACTS_DIR%\\*) do (
                                    echo File: %%%%~nxf
                                    echo Size: %%%%~zF bytes
                                    echo.
                                )
                                echo ════════════════════════════════════════════════════
                                echo Build Status: SUCCESS
                                echo ════════════════════════════════════════════════════
                            ) > "!REPORT_FILE!"

                            echo ✅ Build report generated: !REPORT_FILE!
                            echo.
                            type "!REPORT_FILE!"

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
                        echo "Output Location: ${WORKSPACE}/${ARTIFACTS_DIR}"
                        echo ""
                    '''
                } else {
                    bat '''
                        echo.
                        echo Build Status: SUCCESS
                        echo Job: %JOB_NAME%
                        echo Build: %BUILD_NUMBER%
                        echo Output Location: %WORKSPACE%\\%ARTIFACTS_DIR%
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
                        echo "Workspace Contents:"
                        echo "  • MASSTCLI.zip"
                        echo "  • MASSTCLI_EXTRACTED (extracted files)"
                        echo "  • MyApp.aab"
                        echo "  • config.bm"
                        echo "  • output/ (generated files and report)"
                        echo ""
                        echo "Files are available in Jenkins artifacts and workspace."
                        echo ""
                    '''
                } else {
                    bat '''
                        echo.
                        echo Workspace Contents:
                        echo   • MASSTCLI.zip
                        echo   • MASSTCLI_EXTRACTED (extracted files)
                        echo   • MyApp.aab
                        echo   • config.bm
                        echo   • output/ (generated files and report)
                        echo.
                        echo Files are available in Jenkins artifacts and workspace.
                        echo.
                    '''
                }
            }
        }
    }
}