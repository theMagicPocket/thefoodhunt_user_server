pipeline {

    agent any

    tools {
        go '1.23.1'
    }

    environment {
        MONGODB_URI=credentials('yumfoods_mongouri')
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
                    go version
                    go run cmd/main.go --MONGODB_URI=$MONGODB_URI
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
