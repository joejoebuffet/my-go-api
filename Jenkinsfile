pipeline {
    agent { label 'podman-node' }
    
    environment {
        REGISTRY_URL = "ghcr.io" 
        IMAGE_NAME   = "joejoebuffet/my-go-api"
        DB_HOST      = "10.36.168.16"            
        DB_PORT      = "5432"                 
    }

    stages {
        stage('Deploy to Debian Target') {
            steps {
                echo "Recycling container on Debian target..."
                
                sh "podman stop billing-api || true"
                sh "podman rm billing-api || true"
                
                sh "podman pull ${REGISTRY_URL}/${IMAGE_NAME}:latest"
                
                // FIXED: Added JENKINS_NODE_COOKIE=dontKillMe and --restart=always
                sh """
                    export JENKINS_NODE_COOKIE=dontKillMe
                    podman run -d --restart=always --name billing-api \
                      -p 8080:8081 \
                      -e DATABASE_HOST=${DB_HOST} \
                      -e DATABASE_PORT=${DB_PORT} \
                      ${REGISTRY_URL}/${IMAGE_NAME}:latest
                """
            }
        }

        stage('Smoke Test') {
            steps {
                echo "Verifying API health status..."
                sh "sleep 3" 
                sh "curl -f http://localhost:8080/hello || exit 1"
            }
        }
    }
}