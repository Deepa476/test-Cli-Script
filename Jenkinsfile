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
        timeout(time: 30, unit: 'MINUTES')
    }

    stages {
        stage('Checkout SCM') {
            steps {
                echo 'Checking out source code...'
                checkout scm
                script {
                    echo "Agent OS: ${env.OS ?: (isUnix() ? 'Unix/Linux/macOS' : 'Windows')}"
                }
            }
        }

        stage('Extract MASSTCLI') {
            steps {
                timeout(time: 5, unit: 'MINUTES') {
                    echo 'Extracting MASSTCLI from workspace...'
                    script {
                        if (isUnix()) {
                            // macOS and Linux
                            sh '''
                                if [ ! -d "${MASST_DIR}" ]; then
                                    echo "Extracting MASSTCLI.zip from workspace..."

                                    if [ ! -f "${WORKSPACE}/${MASST_ZIP}.zip" ]; then
                                        echo "ERROR: MASSTCLI.zip not found in workspace!"
                                        echo "Contents of workspace:"
                                        ls -la "${WORKSPACE}"
                                        exit 1
                                    fi

                                    echo "Extracting to temporary location..."
                                    temp_dir=$(mktemp -d)
                                    unzip -q "${WORKSPACE}/${MASST_ZIP}.zip" -d "${temp_dir}"

                                    echo "Moving extracted files to workspace..."
                                    extracted_folder=$(ls -d "${temp_dir}"/*/ 2>/dev/null | head -n 1 | xargs -I {} basename {})

                                    if [ -z "${extracted_folder}" ]; then
                                        echo "ERROR: No folder found in extracted archive"
                                        rm -rf "${temp_dir}"
                                        exit 1
                                    fi

                                    mv "${temp_dir}/${extracted_folder}" "${WORKSPACE}/${MASST_DIR}"
                                    rm -rf "${temp_dir}"

                                else
                                    echo "MASSTCLI already extracted, skipping extraction"
                                fi
                            '''
                        } else {
                            // Windows
                            bat '''
                                if not exist "%MASST_DIR%" (
                                    echo Extracting MASSTCLI.zip from workspace...

                                    if not exist "%WORKSPACE%\\%MASST_ZIP%.zip" (
                                        echo ERROR: MASSTCLI.zip not found in workspace!
                                        echo Contents of workspace:
                                        dir /b "%WORKSPACE%"
                                        exit /b 1
                                    )

                                    echo Extracting to temporary location to avoid path length issues...
                                    powershell -Command "$ProgressPreference = 'SilentlyContinue'; $tempDir = [System.IO.Path]::GetTempFileName(); Remove-Item $tempDir -Force; New-Item -ItemType Directory -Path $tempDir -Force | Out-Null; Expand-Archive -LiteralPath '%WORKSPACE%\\%MASST_ZIP%.zip' -DestinationPath $tempDir -Force; $extractedFolder = Get-ChildItem -Path $tempDir -Directory | Select-Object -First 1; if ($extractedFolder) { Move-Item -Path $extractedFolder.FullName -Destination '%WORKSPACE%\\%MASST_DIR%' -Force } else { echo 'ERROR: No folder found'; Remove-Item $tempDir -Recurse -Force; exit 1 }; Remove-Item $tempDir -Recurse -Force"

                                ) else (
                                    echo MASSTCLI already extracted, skipping extraction
                                )
                            '''
                        }
                    }
                }
            }
        }

        stage('Verify MASSTCLI') {
            steps {
                echo 'Verifying MASSTCLI installation...'
                script {
                    if (isUnix()) {
                        // macOS and Linux
                        sh '''
                            echo "Searching for MASSTCLI executable..."

                            masst_exe=$(find "${MASST_DIR}" -maxdepth 1 \\( -name "MASSTCLI*" -o -name "masstcli*" \\) -type f -executable 2>/dev/null | head -n 1)

                            if [ -z "${masst_exe}" ]; then
                                potential_exe=$(find "${MASST_DIR}" -maxdepth 1 \\( -name "MASSTCLI*" -o -name "masstcli*" \\) -type f 2>/dev/null | head -n 1)
                                if [ -n "${potential_exe}" ]; then
                                    chmod +x "${potential_exe}"
                                    masst_exe="${potential_exe}"
                                fi
                            fi

                            if [ -z "${masst_exe}" ]; then
                                echo "ERROR: No MASSTCLI executable found in ${MASST_DIR}!"
                                echo "Contents of extraction directory:"
                                find "${MASST_DIR}" -type f
                                exit 1
                            fi

                            masst_exe_name=$(basename "${masst_exe}")
                            echo ""
                            echo "========================================"
                            echo "✅ MASSTCLI verified successfully"
                            echo "Executable: ${masst_exe_name}"
                            echo "Location: ${MASST_DIR}"
                            echo "========================================"
                        '''
                    } else {
                        // Windows
                        bat '''
                            echo Searching for MASSTCLI executable...

                            setlocal enabledelayedexpansion
                            set MASST_EXE=
                            for %%f in ("%MASST_DIR%\\MASSTCLI*.exe") do (
                                set MASST_EXE=%%~nxf
                                echo Found executable: %%~nxf
                                goto :found_exe
                            )

                            :found_exe
                            if not defined MASST_EXE (
                                echo ERROR: No MASSTCLI executable found in %MASST_DIR%!
                                echo Contents of extraction directory:
                                dir /s /b "%MASST_DIR%"
                                exit /b 1
                            )

                            if exist "%MASST_DIR%\\!MASST_EXE!" (
                                echo.
                                echo ========================================
                                echo [OK] MASSTCLI verified successfully
                                echo Executable: !MASST_EXE!
                                echo Location: %MASST_DIR%
                                echo ========================================
                            ) else (
                                echo ERROR: MASSTCLI executable file not accessible!
                                exit /b 1
                            )
                            endlocal
                        '''
                    }
                }
            }
        }

        stage('Run MASSTCLI') {
            steps {
                echo 'Running MASSTCLI with configuration and input file...'
                script {
                    if (isUnix()) {
                        // macOS and Linux
                        sh '''
                            masst_exe=$(find "${MASST_DIR}" -maxdepth 1 \\( -name "MASSTCLI*" -o -name "masstcli*" \\) -type f -executable 2>/dev/null | head -n 1)

                            if [ -z "${masst_exe}" ]; then
                                echo "ERROR: No MASSTCLI executable found!"
                                exit 1
                            fi

                            if [ ! -f "${WORKSPACE}/MyApp.aab" ]; then
                                echo "ERROR: Input file MyApp.aab not found in workspace!"
                                ls -la "${WORKSPACE}"
                                exit 1
                            fi

                            echo "Running MASSTCLI..."
                            echo "Input: ${WORKSPACE}/MyApp.aab"
                            echo "Config: ${WORKSPACE}/${CONFIG_FILE}"
                            "${masst_exe}" -input "${WORKSPACE}/MyApp.aab" -config "${WORKSPACE}/${CONFIG_FILE}"
                        '''
                    } else {
                        // Windows
                        bat '''
                            setlocal enabledelayedexpansion
                            set MASST_EXE=
                            for %%f in ("%MASST_DIR%\\MASSTCLI*.exe") do (
                                set MASST_EXE=%%~nxf
                            )

                            if "!MASST_EXE!"=="" (
                                echo ERROR: No MASSTCLI executable found!
                                exit /b 1
                            )

                            if not exist "%WORKSPACE%\\MyApp.aab" (
                                echo ERROR: Input file MyApp.aab not found in workspace!
                                dir /b "%WORKSPACE%"
                                exit /b 1
                            )

                            echo Running MASSTCLI...
                            echo Input: %WORKSPACE%\\MyApp.aab
                            echo Config: %WORKSPACE%\\%CONFIG_FILE%
                            "%MASST_DIR%\\!MASST_EXE!" -input "%WORKSPACE%\\MyApp.aab" -config "%WORKSPACE%\\%CONFIG_FILE%"
                            endlocal
                        '''
                    }
                }
            }
        }
    }

    post {
        success {
            echo '✅ MASSTCLI pipeline completed successfully!'
        }
        failure {
            echo '❌ MASSTCLI pipeline failed. Check logs for details.'
        }
        always {
            echo 'Cleaning up temporary files...'
            script {
                if (isUnix()) {
                    sh 'rm -rf /tmp/masstcli_extract 2>/dev/null || true'
                }
            }
        }
    }
}