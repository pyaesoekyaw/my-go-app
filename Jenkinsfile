// Jenkinsfile
pipeline {
    agent any
    tools {
        go 'go-1.4'
    }
    environment {
        APP_NAME = 'my-go-app'
        APP_PORT = '8000'
        EC2_HOST = '44.202.115.187' // !!! မပြောင်းလဲပါနှင့်
        EC2_USER = 'ubuntu' // or 'ec2-user'
        EC2_CREDENTIALS_ID = 'my-ec2-ssh-key' // !!! မပြောင်းလဲပါနှင့်

        // === AWS RDS Database Credentials (အသစ်ထပ်ထည့်ရန်) ===
        DB_HOST = 'database-psk.c16i22mumjbv.us-east-1.rds.amazonaws.com' // !!! သင့် RDS Endpoint ကို ထည့်ပါ !!!
        DB_PORT = '5432' // PostgreSQL default port
        DB_USER = 'Achawlay' // !!! သင့် RDS Database Username ကို ထည့်ပါ !!!
        DB_PASSWORD = credentials('your-db-password-credential-id') // !!! Jenkins Credentials ID ကို ထည့်ပါ !!!
        DB_NAME = 'postgres' // !!! သင့် Database Name ကို ထည့်ပါ !!!
        // ===============================================
    }

    stages {
        stage('Clean Workspace') {
            steps {
                cleanWs() // Cleans the workspace before starting a new build
            }
        }

        stage('Checkout Source Code') {
            steps {
                echo 'Checking out source code...'
                // Ensure your Jenkins Job is configured to use Git SCM,
                // and point it to your Git repository (e.g., GitHub, GitLab)
                checkout scm
            }
        }

        stage('Build Application') {
            steps {
                script {
                    echo 'Building Go application...'
                    // Build the executable, explicitly naming it APP_NAME
                    sh "go build -o ${APP_NAME} ."
                }
            }
        }

        stage('Run Unit Tests') {
            steps {
                script {
                    echo 'Running unit tests...'
                    sh 'go test -v ./...'
                }
            }
        }

        stage('Deploy to EC2') {
            steps {
                script {
                    echo "Deploying ${APP_NAME} to EC2 instance: ${EC2_HOST}"

                    // Use sshagent to securely use your SSH key
                    sshagent(credentials: ["${EC2_CREDENTIALS_ID}"]) {
                        // 1. Stop any existing running instance of the app on EC2
                        sh "ssh -o StrictHostKeyChecking=no ${EC2_USER}@${EC2_HOST} 'sudo pkill ${APP_NAME} || true'"
                        echo 'Previous application instance stopped (if running).'

                        // 2. Copy the newly built binary to EC2
                        // Ensure the path is correct where you want to deploy it on EC2
                        sh "scp -o StrictHostKeyChecking=no ${APP_NAME} ${EC2_USER}@${EC2_HOST}:/home/${EC2_USER}/${APP_NAME}"
                        echo 'Application binary copied to EC2 instance.'

                        // --- NEW: Copy the static directory ---
                        sh "scp -r -o StrictHostKeyChecking=no static ${EC2_USER}@${EC2_HOST}:/home/${EC2_USER}/"
                        echo 'Static directory copied to EC2 instance.'
                        // --- END NEW ---

                        // 3. Start the application in the background on EC2
                        // nohup ensures it keeps running after the SSH session closes
                        // Redirect output to a log file for debugging later
                        // Corrected Line 68
                        sh "ssh -o StrictHostKeyChecking=no -f ${EC2_USER}@${EC2_HOST} 'export DB_HOST=${DB_HOST} DB_PORT=${DB_PORT} DB_USER=${DB_USER} DB_PASSWORD=${DB_PASSWORD} DB_NAME=${DB_NAME}; cd /home/${EC2_USER} && nohup ./${APP_NAME} > app.log 2>&1 &'"
                        echo "Application start command sent to EC2 instance with DB credentials."


                        // 4. Give the app a moment to start up
                        sh 'sleep 10'

                        // 5. Verify the app is listening on the expected port on EC2
                        sh "ssh -o StrictHostKeyChecking=no ${EC2_USER}@${EC2_HOST} 'sudo netstat -tulnp | grep ${APP_PORT} || echo \"App not listening on ${APP_PORT} on EC2\"'"
                        echo "Verified application is listening on port ${APP_PORT} on EC2."
                    }
                }
            }
        }

        stage('Verify Deployment') {
            steps {
                script {
                    echo "Verifying application access at http://${EC2_HOST}:${APP_PORT}"
                    // Use curl to hit the health endpoint of the deployed app
                    // -f makes curl fail on HTTP errors (e.g., 4xx, 5xx)
                    // --max-time gives it a timeout
                    sh "curl -f --max-time 15 http://${EC2_HOST}:${APP_PORT}/health"
                    echo 'Deployment verification successful!'
                }
            }
        }
    }

    post {
        always {
            echo 'Pipeline finished.'
        }
        success {
            echo 'Pipeline succeeded!'
        }
        failure {
            echo 'Pipeline failed! Check console output for details.'
        }
    }
}
