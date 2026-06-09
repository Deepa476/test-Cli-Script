pipeline {
    agent any

    parameters {
        booleanParam(
            name: 'IS_DEBUG',
            defaultValue: false,
            description: 'Debug or Release build'
        )
    }

    environment {
        PATH = "/bin:/usr/bin:/usr/local/bin:${env.PATH}"

        INPUT_FILE = "app-release.aab"
        CONFIG_FILE = "bluebeetle_config.bm"

        MACOS_DOWNLOAD_URL = "https://storage.googleapis.com/masst-assets/Defender-Binary-Integrator/1.0.0/MacOS/MASSTCLI-v1.1.0-darwin-arm64.zip"
        LINUX_DOWNLOAD_URL = "https://storage.googleapis.com/masst-assets/Defender-Binary-Integrator/1.0.0/Linux/MASSTCLI-v1.1.0-linux-amd64.zip"

        KEYSTORE_FILE = "Bluebeetle.jks"
        KEYSTORE_PASSWORD = "bugs@1234"
        KEY_ALIAS = "key0"
        KEY_PASSWORD = "bugs@1234"

        IDENTITY = "Apple Distribution: Bugsmirror Research private limited (BPKUYCFJ74)"

        MASST_DIR = "MASSTCLI_EXTRACTED"
        ARTIFACTS_DIR = "output"
        MASST_ZIP = "MASSTCLI"
    }

    options {
        timestamps()
        buildDiscarder(logRotator(numToKeepStr: '10'))
        disableConcurrentBuilds()
    }

    stages {

        stage('Clean Workspace') {
            steps {
                cleanWs()
            }
        }

        stage('Setup') {
            steps {
                checkout scm

                script {

                    env.DETECTED_PLATFORM =
                            sh(
                                    script: '[[ "$(uname)" == "Darwin" ]] && echo MacOS || echo Linux',
                                    returnStdout: true
                            ).trim()

                    env.DOWNLOAD_URL =
                            env.DETECTED_PLATFORM == 'MacOS'
                                    ? env.MACOS_DOWNLOAD_URL
                                    : env.LINUX_DOWNLOAD_URL

                    env.ANDROID_HOME =
                            sh(
                                    script: '''
                                    if [ -n "$ANDROID_HOME" ]; then
                                        echo "$ANDROID_HOME"
                                    elif [ -n "$ANDROID_SDK_ROOT" ]; then
                                        echo "$ANDROID_SDK_ROOT"
                                    elif [ -d "$HOME/Library/Android/sdk" ]; then
                                        echo "$HOME/Library/Android/sdk"
                                    elif [ -d "$HOME/Android/Sdk" ]; then
                                        echo "$HOME/Android/Sdk"
                                    elif [ -d "$HOME/Android/sdk" ]; then
                                        echo "$HOME/Android/sdk"
                                    else
                                        exit 1
                                    fi
                                    ''',
                                    returnStdout: true
                            ).trim()

                    echo "Platform: ${env.DETECTED_PLATFORM}"
                    echo "Android SDK: ${env.ANDROID_HOME}"

                    sh """
                    #!/bin/bash -e

                    export PATH=/bin:/usr/bin:/usr/local/bin:\$PATH:${env.ANDROID_HOME}/cmdline-tools/latest/bin:${env.ANDROID_HOME}/platform-tools

                    echo "Downloading MASST CLI"

                    rm -f ${MASST_ZIP}.zip

                    curl -L -o ${MASST_ZIP}.zip "${env.DOWNLOAD_URL}"

                    rm -rf ${MASST_DIR}
                    mkdir -p ${MASST_DIR}

                    unzip -q ${MASST_ZIP}.zip -d ${MASST_DIR}

                    chmod +x \$(find ${MASST_DIR} -type f -name "MASSTCLI*" || true)

                    echo "Contents:"
                    find ${MASST_DIR}
                    """
                }
            }
        }

        stage('Execute') {
            steps {
                script {

                    def BUILD_MODE = params.IS_DEBUG ? "DEBUG" : "RELEASE"

                    sh """
                    #!/bin/bash -e

                    export PATH=/bin:/usr/bin:/usr/local/bin:\$PATH:${env.ANDROID_HOME}/cmdline-tools/latest/bin:${env.ANDROID_HOME}/platform-tools

                    echo "Build mode: ${BUILD_MODE}"

                    [ -f "${INPUT_FILE}" ] || {
                        echo "${INPUT_FILE} not found"
                        exit 1
                    }

                    [ -f "${CONFIG_FILE}" ] || {
                        echo "${CONFIG_FILE} not found"
                        exit 1
                    }

                    MASST_EXE=\$(find "${MASST_DIR}" -type f | grep MASSTCLI | head -n1)

                    echo "Using binary: \$MASST_EXE"

                    [ -n "\$MASST_EXE" ] || {
                        echo "MASST CLI not found"
                        exit 1
                    }

                    EXT=\$(echo "${INPUT_FILE}" | awk -F. '{print tolower(\$NF)}')

                    case "\$EXT" in

                        xcarchive|ipa)

                            "\$MASST_EXE" \
                            -input="${INPUT_FILE}" \
                            -config="${CONFIG_FILE}" \
                            -identity="${IDENTITY}"
                            ;;

                        aab|apk)

                            if [ "${params.IS_DEBUG}" = "true" ]; then

                                "\$MASST_EXE" \
                                -input="${INPUT_FILE}" \
                                -config="${CONFIG_FILE}"

                            else

                                "\$MASST_EXE" \
                                -input="${INPUT_FILE}" \
                                -config="${CONFIG_FILE}" \
                                -keystore="${KEYSTORE_FILE}" \
                                -storePassword="${KEYSTORE_PASSWORD}" \
                                -alias="${KEY_ALIAS}" \
                                -keyPassword="${KEY_PASSWORD}" \
                                -v=true \
                                -apk

                            fi
                            ;;

                        *)

                            echo "Unsupported file type"
                            exit 1
                            ;;
                    esac

                    echo "Execution completed"
                    """
                }
            }
        }

        stage('Archive') {
            steps {

                script {

                    def BUILD_MODE = params.IS_DEBUG ? "DEBUG" : "RELEASE"

                    sh """
                    #!/bin/bash -e

                    mkdir -p "${ARTIFACTS_DIR}"

                    echo "Build: ${env.DETECTED_PLATFORM} | ${BUILD_MODE} | \$(date)" \
                    > "${ARTIFACTS_DIR}/report.txt"

                    rm -f ${MASST_ZIP}.zip
                    """
                }

                archiveArtifacts artifacts: 'output/**', allowEmptyArchive: true
            }
        }
    }

    post {
        success {
            echo "${env.DETECTED_PLATFORM} ${params.IS_DEBUG ? 'DEBUG' : 'RELEASE'} - SUCCESS"
        }

        failure {
            echo "FAILED"
        }

        always {
            cleanWs(deleteDirs: true)
        }
    }
}
