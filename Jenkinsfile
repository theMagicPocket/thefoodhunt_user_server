pipeline {

    agent any

    tools {
        go '1.23.1'
    }

    environment {
        MONGODB_URI=credentials('yumfoods_mongodb_uri')
        SECRET_KEY=credentials('yumfoods_secret_key')
        MATRIX_KEY=credentials('yumfoods_matrix_key')
    }

    stages {
        stage('clone') {
        steps {
                git branch: 'refactor', credentialsId: 'jenkins-food-api-github-token', url: 'https://github.com/deVamshi/golang_food_delivery_api/'
            }
        }
        stage('run') {
            steps {
                sh '''
                    echo "running"
                    go version
                    export JENKINS_NODE_COOKIE=dontKillMe;go run cmd/main.go --MONGODB_URI=$MONGODB_URI --SECRET_KEY=$SECRET_KEY --MATRIX_KEY=$MATRIX_KEY &
                '''
            }
        }

    }
    post {
        success {
            echo 'build success'
        }
        failure {
            echo 'build failed'
        }
    }
}
