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
                echo '✅ Source code checked out successfully'
            }
        }

        stage('Prepare MASSTCLI Tool') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE: Prepare MASSTCLI Tool'
                echo '═══════════════════════════════════════════'
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
            }
        }

        stage('Extract MASSTCLI') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE 2: Extract MASSTCLI'
                echo '═══════════════════════════════════════════'
                timeout(time: 5, unit: 'MINUTES') {
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

                        # Ensure unzip is installed
                        if ! command -v unzip &> /dev/null; then
                            echo "Installing unzip..."
                            if command -v apt-get &> /dev/null; then
                                sudo apt-get update && sudo apt-get install -y unzip || true
                            elif command -v yum &> /dev/null; then
                                sudo yum install -y unzip || true
                            elif command -v dnf &> /dev/null; then
                                sudo dnf install -y unzip || true
                            fi
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
                }
            }
        }

        stage('Verify MASSTCLI') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE 3: Verify MASSTCLI'
                echo '═══════════════════════════════════════════'
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
            }
        }

        stage('Validate Input Files') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE: Validate Input Files'
                echo '═══════════════════════════════════════════'
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

                    file_size=$(stat -c%s "${WORKSPACE}/MyApp.aab")
                    echo "✅ Found: MyApp.aab (${file_size} bytes)"

                    echo ""
                    echo "Checking ${CONFIG_FILE}..."
                    if [ ! -f "${WORKSPACE}/${CONFIG_FILE}" ]; then
                        echo "ERROR: Configuration file ${CONFIG_FILE} not found!"
                        exit 1
                    fi

                    file_size=$(stat -c%s "${WORKSPACE}/${CONFIG_FILE}")
                    echo "✅ Found: ${CONFIG_FILE} (${file_size} bytes)"

                    echo ""
                    echo "✅ All input files validated successfully"
                '''
            }
        }

        stage('Run MASSTCLI') {
            steps {
                echo '═══════════════════════════════════════════'
                echo '  STAGE 4: Run MASSTCLI'
                echo '═══════════════════════════════════════════'
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
            }
        }
    }

    post {
        success {
            echo '═══════════════════════════════════════════'
            echo '✅ PIPELINE COMPLETED SUCCESSFULLY'
            echo '═══════════════════════════════════════════'
            sh '''
                echo ""
                echo "Build Status: SUCCESS"
                echo "Job: ${JOB_NAME}"
                echo "Build: ${BUILD_NUMBER}"
                echo ""
            '''
        }
        failure {
            echo '═══════════════════════════════════════════'
            echo '❌ PIPELINE FAILED'
            echo '═══════════════════════════════════════════'
            sh '''
                echo ""
                echo "Build Status: FAILURE"
                echo "Job: ${JOB_NAME}"
                echo "Build: ${BUILD_NUMBER}"
                echo "Check the logs above for error details."
                echo ""
            '''
        }
        always {
            echo '═══════════════════════════════════════════'
            echo 'Cleanup Phase'
            echo '═══════════════════════════════════════════'
            sh '''
                echo "Cleaning up temporary files..."
                find /tmp -maxdepth 1 -name "masstcli_extract*" -type d -exec rm -rf {} + 2>/dev/null || true
                echo "✅ Cleanup completed"
            '''
        }
    }
}