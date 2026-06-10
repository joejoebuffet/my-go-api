pipeline {
    agent { label 'built-in' }

    environment {
        GITHUB_PAT = credentials('github-pat')
    }

    stages {
        stage('Deploy to K8s') {
            steps {
                echo "Deploying to Kubernetes..."
                sshagent(['k8s-ssh']) {
                    sh """
                        ssh -o StrictHostKeyChecking=no willow@192.168.56.10 \
                        'GITHUB_PAT=${GITHUB_PAT} bash /home/willow/deploy.sh'
                    """
                }
            }
        }
    }

    post {
        success {
            echo "Deployment successful! ✅"
        }
        failure {
            echo "Deployment failed! ❌"
        }
    }
}