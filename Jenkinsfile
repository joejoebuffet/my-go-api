pipeline {
    agent { label 'podman-node' }

    environment {
        REGISTRY_URL = "ghcr.io" 
        IMAGE_NAME   = "joejoebuffet/billing-api-go"
        IMAGE_TAG    = "${BUILD_NUMBER}"
        REGISTRY_CREDS_ID = "git-token" 
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Container Image Build') {
            steps {
                echo "Building container image using Podman via Dockerfile..."
                // Changed to sh and updated to Linux variable syntax ${VAR}
                sh "podman build -t ${REGISTRY_URL}/${IMAGE_NAME}:${IMAGE_TAG} ."
                sh "podman tag ${REGISTRY_URL}/${IMAGE_NAME}:${IMAGE_TAG} ${REGISTRY_URL}/${IMAGE_NAME}:latest"
            }
        }

        stage('Registry Authentication & Push') {
            steps {
                echo "Logging into container registry and pushing images..."
                withCredentials([usernamePassword(credentialsId: "${REGISTRY_CREDS_ID}", usernameVariable: 'REG_USER', passwordVariable: 'REG_PASS')]) {
                    // Linux shell execution block
                    sh "podman login -u ${REG_USER} -p ${REG_PASS} ${REGISTRY_URL}"
                    
                    echo "Pushing build version tag..."
                    sh "podman push ${REGISTRY_URL}/${IMAGE_NAME}:${IMAGE_TAG}"
                    
                    echo "Pushing latest tag..."
                    sh "podman push ${REGISTRY_URL}/${IMAGE_NAME}:latest"
                }
            }
        }
    }

    post {
        always {
            echo "Cleaning up local build image copies from Jenkins agent workspace..."
            // Using Linux shell fallback construct (|| true) instead of (|| exit 0)
            sh "podman rmi ${REGISTRY_URL}/${IMAGE_NAME}:${IMAGE_TAG} || true"
            sh "podman rmi ${REGISTRY_URL}/${IMAGE_NAME}:latest || true"
        }
    }
}