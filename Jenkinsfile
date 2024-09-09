pipeline {
    tools {
        go '1.23.1'
    }
    stages {
        stage('clone') {
            steps {
                // git branch: 'refactor', credentialsId: 'yumfoods-github-token', url: 'https://github.com/deVamshi/golang_food_delivery_api/'
                sh 'echo HI'
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
