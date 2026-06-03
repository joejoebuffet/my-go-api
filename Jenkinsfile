pipeline {
    agent { label 'podman-node' }

    // No automation triggers. You control the execution.
    
    environment {
        REGISTRY_URL = "ghcr.io" 
        IMAGE_NAME   = "joejoebuffet/my-go-api"
        DB_HOST      = "10.36.168.15"            
        DB_PORT      = "5432"                 
    }

    stages {
        stage('Deploy to Debian Target') {
            steps {
                echo "Deploying container to Debian target..."
                
                sh "podman stop billing-api || true"
                sh "podman rm billing-api || true"
                
                // Pulls the image that you know is ready because you checked GitHub Actions
                sh "podman pull ${REGISTRY_URL}/${IMAGE_NAME}:latest"
                
                sh "podman run -d --name billing-api -p 8080:8080 -e DATABASE_HOST=${DB_HOST} -e DATABASE_PORT=${DB_PORT} ${REGISTRY_URL}/${IMAGE_NAME}:latest"
            }
        }

        stage('Smoke Test') {
            steps {
                echo "Verifying API health status..."
                sh "sleep 3" 
                sh "curl -f http://localhost:8080/health || exit 1"
            }
        }
    }
}