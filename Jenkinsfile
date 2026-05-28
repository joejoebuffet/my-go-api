pipeline {
    agent any

    environment {
        // Change these to match your environment's registry configuration
        REGISTRY_URL = "your-registry-host.com" 
        IMAGE_NAME   = "billing-api-go"
        IMAGE_TAG    = "${BUILD_NUMBER}"
        
        // This must match the ID of the credentials stored inside Jenkins
        REGISTRY_CREDS_ID = "registry-credentials-id" 
    }

    stages {
        stage('Checkout') {
            steps {
                // Pulls the code from your linked GitHub repository
                checkout scm
            }
        }

        stage('Container Image Build') {
            steps {
                echo "Building container image using Podman via Dockerfile..."
                // Executes the multi-stage Dockerfile to build the binary and container image
                sh "podman build -t ${REGISTRY_URL}/${IMAGE_NAME}:${IMAGE_TAG} ."
                sh "podman tag ${REGISTRY_URL}/${IMAGE_NAME}:${IMAGE_TAG} ${REGISTRY_URL}/${IMAGE_NAME}:latest"
            }
        }

        stage('Registry Authentication & Push') {
            steps {
                echo "Logging into container registry and pushing images..."
                // Securely pulls credentials from Jenkins to avoid hardcoding passwords
                withCredentials([usernamePassword(credentialsId: "${REGISTRY_CREDS_ID}", usernameVariable: 'REG_USER', passwordVariable: 'REG_PASS')]) {
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
            sh "podman rmi ${REGISTRY_URL}/${IMAGE_NAME}:${IMAGE_TAG} || true"
            sh "podman rmi ${REGISTRY_URL}/${IMAGE_NAME}:latest || true"
        }
    }
}