// Jenkinsfile
pipeline {
    // agent ကို label ဖြင့် သတ်မှတ်ခြင်းဖြင့် Jenkins Agent Node (my-go-app-agent) ပေါ်တွင် run စေမည်။
    agent { label 'important' }
    tools { go 'go-1.4'}

    environment {
        APP_NAME = 'my-go-app'
        APP_PORT = '8000'
        // Deploy လုပ်မည့် Server သည် Jenkins Agent EC2 Instance ဖြစ်သောကြောင့် ၎င်း၏ Public IP ကို သုံးမည်။
        EC2_HOST = '54.86.20.184' // !!! Agent EC2 ၏ Public IP ကို ဒီနေရာတွင် ထည့်ပါ !!!
        EC2_USER = 'ubuntu' // Agent EC2 ၏ user
        EC2_CREDENTIALS_ID = 'jenkins-agent-ssh-key' // !!! ssh and password နဲ့jenkins agent ရဲ့ private key ကိုထည့်

        // === AWS RDS Database Credentials ===
        DB_HOST = 'database-alinn.c16i22mumjbv.us-east-1.rds.amazonaws.com' // !!! သင့် RDS Endpoint ကို ထည့်ပါ !!!
        DB_PORT = '5432'
        DB_USER = 'Achawlay'
        DB_PASSWORD = credentials('my-rds-db-password') // !!! Jenkins Credentials ID ကို ထည့်ပါ 
        DB_NAME = 'postgres' // !!! သင့် Database Name ကို ထည့်ပါ !!!
        // ===============================================
    }

    stages {
        stage('Clean Workspace') {
            steps {
                echo 'Cleaning Jenkins workspace on agent...'
                cleanWs()
            }
        }

        stage('Checkout Source Code') {
            steps {
                echo 'Checking out source code from Git...'
                checkout scm
            }
        }

        stage('Build Application') {
            steps {
                script {
                    echo 'Building Go application...'
                    // Agent ပေါ်တွင် Go toolchain ကိုအသုံးပြုရန်။ Global Tool Configuration မှ Name ကို အသုံးပြုပါ။
                  // !!! Global Tool Config မှ Go installation Name ကို ထည့်ပါ !!!
                    sh "go build -o ${APP_NAME} ."
                }
            }
        }

        stage('Run Unit Tests') {
            steps {
                script {
                    echo 'Running unit tests for the Go application...'
                   // !!! Global Tool Config မှ Go installation Name ကို ထည့်ပါ !!!
                    sh 'go test -v ./...'
                }
            }
        }

        stage('Deploy to EC2') {
            steps {
                script {
                    echo "Deploying ${APP_NAME} to EC2 instance: ${EC2_HOST}"

                    sshagent(credentials: ["${EC2_CREDENTIALS_ID}"]) {
                        // Deploy server သည် Jenkins Agent EC2 ကိုယ်တိုင်ဖြစ်သည်
                        // 1. အရင် run နေသော application instance ကို ရပ်တန့်ခြင်း
                        sh "ssh -o StrictHostKeyChecking=no ${EC2_USER}@${EC2_HOST} 'sudo pkill ${APP_NAME} || true'"
                        echo 'Previous application instance stopped (if running).'

                        // 2. Old binary file ကို အတင်းအကြပ် ဖယ်ရှားခြင်း
                        sh "ssh -o StrictHostKeyChecking=no ${EC2_USER}@${EC2_HOST} 'rm -f /home/${EC2_USER}/${APP_NAME}'"
                        echo 'Old application binary removed (if existed).'

                        // 3. Build လုပ်ထားသော binary နှင့် static directory များကို deploy target သို့ copy လုပ်ခြင်း
                        // Jenkins Agent ၏ workspace မှ /home/ubuntu သို့ copy ခြင်း။
                        sh "scp -o StrictHostKeyChecking=no ${APP_NAME} ${EC2_USER}@${EC2_HOST}:/home/${EC2_USER}/${APP_NAME}"
                        echo 'Application binary copied to EC2 instance.'
                        sh "scp -r -o StrictHostKeyChecking=no static ${EC2_USER}@${EC2_HOST}:/home/${EC2_USER}/"
                        echo 'Static directory copied to EC2 instance.'

                        // 4. Application ကို background တွင် စတင် run ခြင်း
                        // DB Environment Variables များကို တိုက်ရိုက် export လုပ်ခြင်း
                        sh "ssh -o StrictHostKeyChecking=no -f ${EC2_USER}@${EC2_HOST} 'export DB_HOST=\"${DB_HOST}\" DB_PORT=\"${DB_PORT}\" DB_USER=\"${DB_USER}\" DB_PASSWORD=\"${DB_PASSWORD}\" DB_NAME=\"${DB_NAME}\"; cd /home/${EC2_USER} && nohup ./${APP_NAME} > app.log 2>&1 &'"
                        echo "Application start command sent to EC2 instance with DB credentials."

                        sh 'sleep 10'

                        // 5. Application သည် သတ်မှတ်ထားသော port တွင် listening လုပ်နေခြင်းရှိမရှိ စစ်ဆေးခြင်း
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
