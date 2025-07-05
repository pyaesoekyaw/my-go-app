
# `my-go-app` Web Application Deployment

![final diagram](https://github.com/pyaesoekyaw/my-go-app/blob/main/images/Final%20Diagram.png)
-----

before we get start , here is what you are looking for , a live action version of the project (https://www.youtube.com/watch?v=w0FxndMykME&t=518s)
## 🌟 Project Overview

This project showcases a **simple Go web application** implementing user registration and sign-in functionalities. Data is securely stored in an **Amazon RDS PostgreSQL database**. The core of this repository is to demonstrate a robust **Jenkins CI/CD pipeline** for automated building, testing, and deployment to an AWS EC2 instance. This setup emphasizes continuous integration and continuous delivery best practices.

-----

## ✨ Key Features

  * **User Authentication**: Secure user registration and sign-in with password hashing (bcrypt).
  * **Database Persistence**: Utilizes Amazon RDS for highly available and scalable PostgreSQL data storage.
  * **Modern Web UI**: Simple, responsive web interface for user interaction (HTML/CSS/JavaScript).
  * **Automated CI/CD**: Jenkins pipeline for seamless integration, build, test, and deployment.
  * **Two-Tier AWS Setup**: Dedicated EC2 instances for Jenkins Controller and Agent/Application Host.
  * **Secure Access**: SSH for deployment, secure communication with database (SSL optional).

-----

## 🚀 Deployment Steps

### **Step 1: AWS EC2 Instance Setup**

Let's prepare the two EC2 instances required for the Jenkins Controller and Agent/Application deployment.

1.  **Launch Jenkins Controller EC2 Instance**:

      * Navigate to **EC2 Console** → **Launch Instances**.
      * **Name**: `Jenkins-Controller`
      * **AMI**: `Ubuntu Server 22.04 LTS` 
      * **Instance Type**: `t2.medium` (minimum 2GB RAM for Jenkins Controller)
      * **Key pair (login)**: Create a new key pair or select an existing one.
      * **Network settings**:
          * **Security group**:
              * **Name**: `launch-wizard-2`    (you create your own mate)
              * **Inbound rules**:
                  * **Type**: `SSH`, **Source**: `0.0.0.0/0` 
                  * **Type**: `Custom TCP`, **Port range**: `8080`, **Source**: `0.0.0.0/0` (for Jenkins Web UI)
      * Launch instance.

2.  **Launch Jenkins Agent EC2 Instance (Application Host)**:

      * Navigate to **EC2 Console** → **Launch Instances**.
      * **Name**: `Jenkins-Agent` (This will also host your `my-go-app`)
      * **AMI**: `Ubuntu Server 22.04 LTS`
      * **Instance Type**: `t2.small` (sufficient for app deployment)
      * **Network settings**:
          * **Security group**:
              * **Name**: `launch-wizard-1`
              * **Inbound rules**:
                  * **Type**: `SSH`, **Source**: `0.0.0.0/0`
                  * **Type**: `Custom TCP`, **Port range**: `8000`, **Source**: `0.0.0.0/0` (for user access)
                  * **Type**: `Postgres`, **Port range**: `5432`, **Source**: `SecForRDS` (security group of RDS)
      * Launch instance.

-----

### **Step 2: AWS RDS PostgreSQL Database Setup**

Let's create a highly available PostgreSQL database for user data storage.

1.  **Navigate to RDS Console** → **Create Database**.
2.  **Choose a database creation method**: `Standard create`.
3.  **Engine options**: Select `PostgreSQL`.
4.  **Templates**: Choose `Free tier` (for testing)
5.  **DB instance identifier**: `database-psk` (or your preferred name).
6.  **Master username**: `Achawlay` (or your preferred username).
7.  **Master password**: Set a strong password.
8.  **Connectivity**:
      * **VPC**: Select the **same VPC** as your EC2 instances.
      * **Publicly accessible**: Select `Yes` (for simpler testing setup from outside, but `No` is recommended for production and internal VPC access).
      * **VPC security groups**: Choose `SecForRDS`(create new or existing security group)
          * **Crucial**: After creation, navigate to this RDS Security Group (mine is `SecForRDS`. 
              * **Type**: `PostgreSQL`
              * **Port range**: `5432`
              * **Source**: {private ip address of Agent EC2} (This allows your Go app (on the Agent) to connect to RDS.)
              * **Save rules**.
11. Leave other settings as default or configure as needed.
12. **Create database**. 
13. Create Custom RDS Parameter Group and Disable SSL Enforcement**
    **Navigate to RDS Console** → **Parameter groups**. Click **"Create parameter group"**.
    * **Parameter group family**: Select `custom-postgres17`.
    * **Type**: Select `DB Parameter Group`.
    * Click **"Create"**.
14. Once created, select your new custom parameter group and click **"Edit"**.
    In the search bar, type `rds.force_ssl`.
    * Change its `Value` from `1` (or whatever it is) to `0`.
    * Click **"Save changes"**.
15. **Associate Custom Parameter Group with your RDS Instance:**
    * Go to **RDS Console** → **Databases**.
    * Select your `database-psk` instance.
    * Click **"Modify"**.
    * Scroll down to the **"Database options"** section.
    * For **"DB parameter group"**, select your newly created custom parameter group (`custom-postgres17`).
    * Scroll down and select **"Apply immediately"**.

---
-----

### **Step 3: Jenkins Controller Setup**

Let's set up the Jenkins Controller to manage your CI/CD pipeline.

1.  **Install Jenkins:**
      * SSH into your `Jenkins-Controller` EC2 instance.
      * Follow the official Jenkins installation guide for Ubuntu. (e.g., `sudo apt update && sudo apt install openjdk-11-jdk jenkins`).
      * Start Jenkins service.
      * Access Jenkins Web UI via `http://<Controller_Public_IP>:8080` and complete the initial setup (unlocking, creating admin user, installing suggested plugins).
2.  **Install Jenkins Plugins:**
      * Jenkins Dashboard → **Manage Jenkins** → **Manage Plugins** → **Available** tab. Install:
          * `Go Plugin`
          * `SSH Agent Plugin`
          * install and restart the jenkins.
3.  **Configure Global Tools:**
      * Jenkins Dashboard → **Manage Jenkins** → **Global Tool Configuration**.
      * **Go**: Click "Add Go", give it a Name (e.g., `go-1.4`), and check "Install automatically".
4.  **Add Jenkins Credentials (for Pipeline Use):**
      * Jenkins Dashboard → **Manage Jenkins** → **Manage Credentials** → **Jenkins** → **Global credentials (unrestricted)** → **Add Credentials**.
          * Kind: "SSH Username with private key"
          * ID: `my-ec2-ssh-key` (Used in `Jenkinsfile`)
          * Username: `ubuntu` (User on `Jenkins-Agent` EC2).
          * Private Key: Paste the content of your `Jenkins-Agent` EC2's `.pem` file.
      * **RDS DB Password:**
          * Kind: "Secret text"
          * ID: `my-rds-db-password` (Used in `Jenkinsfile`)
          * Secret: Paste your RDS database password.
5. **Configure Jenkins Agent Node:**
      * Jenkins Dashboard → **Manage Jenkins** → **Manage Nodes** → **"New Node"**.
      * **Node Name**: `my-go-app-agent`
      * **Remote root directory**: `/home/ubuntu/jenkins-agent-workspace` (create this directory on the Agent EC2).
      * **Host**: **Private IP Address of your `Jenkins-Agent` EC2 instance** 
      * **Credentials**: choose `ubuntu` (you already add the credential.)
      * **Launch method**: "Launch agents via SSH".
      * **Save**. Wait for the Agent to come Online.
5.  
-----

### **Step 4: Jenkins Agent EC2 Configuration**

Let's prepare the Jenkins Agent (which is also your application host).

1.  **SSH into your `Jenkins-Agent` EC2 instance.**
2.  **Install Required Tools:**
    ```bash
    sudo apt update
    sudo apt install build-essential 
    ```

-----

### **Step 5: Git Host Key Verification**

This is a secure practice to ensure Jenkins is connecting to the authentic Git server.

1.  **Obtain Host Key Fingerprint:**

      * SSH into your `Jenkins-Controller` EC2 instance
      * Run `ssh-keyscan github.com` 
      * **Carefully copy ONLY the lines that contain the actual public keys.** 

2.  **Configure in Jenkins:**

      * Jenkins Dashboard → **Manage Jenkins** → **Security**.
      * Scroll down to **Git Host Key Verification Configuration**.
      * **Strategy**: Select **"Manually provided keys"**.
      * Paste the copied host key string(s) into the large text area. (3 keys for me)
      * **Save**.

-----

### **Step 6: Project Code (Application Files)**

Ensure these files are committed to your GitHub repository (`https://github.com/pyaesoekyaw/my-go-app`) in the exact structure:

  * `my-go-app/`
      * `main.go`
      * `main_test.go`
      * `Jenkinsfile`
      * `go.mod`
      * `go.sum`
      * `model/`
          * `user.go`
      * `repository/`
          * `repository.go`
      * `static/`
          * `index.html`
          * `dashboard.html`

### **Step 7: `Jenkinsfile` Configuration**

This `Jenkinsfile` resides in the root of your `my-go-app` Git repository.
**Make sure to replace all placeholder values (`!!! REPLACE ... !!!`) with your actual AWS IDs and values.**

----

### **Step 8: Running the Jenkins Pipeline**

1.  **Configure Jenkins Job:**
      * Jenkins Dashboard → Your Job (e.g., `My-Go-WebApp-CI-CD`) → **Configure**.
      * In the **Pipeline** section, ensure the SCM Configuration (Git URL: `https://github.com/pyaesoekyaw/my-go-app`, Branch: `main`, Script Path: `Jenkinsfile`) is correct.
2.  **Build Now:** Navigate to the Jenkins Job page and click **"Build Now"**.

After a successful pipeline run, your Go application will be built and deployed on the Jenkins Agent EC2 and should be accessible via `http://JENKINS_AGENT_EC2_PUBLIC_IP:8000`.
