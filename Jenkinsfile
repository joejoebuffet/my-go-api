pipeline {
    agent { label 'podman-node' }
    
    environment {
        REGISTRY_URL = "ghcr.io" 
        IMAGE_NAME   = "joejoebuffet/my-go-api"
        DB_HOST      = "10.36.168.15"            
        DB_PORT      = "5432"                 
    }

    stages {
        stage('Deploy to Debian Target') {
            steps {
                echo "Recycling container on Debian target..."
                
                sh "podman stop billing-api || true"
                sh "podman rm billing-api || true"
                
                sh "podman pull ${REGISTRY_URL}/${IMAGE_NAME}:latest"
                
                // FIXED: Changed internal container port from 8080 to 8081
                sh "podman run -d --name billing-api -p 8080:8081 -e DATABASE_HOST=${DB_HOST} -e DATABASE_PORT=${DB_PORT} ${REGISTRY_URL}/${IMAGE_NAME}:latest"
            }
        }

        stage('Smoke Test') {
            steps {
                echo "Verifying API health status..."
                sh "sleep 3" 
                // This stays 8080 because it hits the VM's external port
                sh "curl -f http://localhost:8080/hello || exit 1"
            }
        }
    }
}