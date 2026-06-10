pipeline {
    agent { label 'built-in' }

    environment {
        GITHUB_PAT = credentials('github-pat')
    }

    stages {
        stage('Deploy to K8s') {
            steps {
                echo "Deploying to Kubernetes..."
                withCredentials([sshUserPrivateKey(
                    credentialsId: 'k8s-ssh',
                    keyFileVariable: 'SSH_KEY',
                    usernameVariable: 'SSH_USER'
                )]) {
                    bat """
                        ssh -i %SSH_KEY% -o StrictHostKeyChecking=no %SSH_USER%@192.168.56.10 "GITHUB_PAT=%GITHUB_PAT% bash /home/willow/deploy.sh"
                    """
                }
            }
        }
    }

    post {
        success { echo "Deployment successful! ✅" }
        failure { echo "Deployment failed! ❌" }
    }
}