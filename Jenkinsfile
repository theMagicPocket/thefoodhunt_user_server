pipeline {

    agent any

    // tools {
    //     go '1.23.1'
    // }

    environment {
        MONGODB_URI=credentials('yumfoods_mongodb_uri')
        SECRET_KEY=credentials('yumfoods_secret_key')
        MATRIX_KEY=credentials('yumfoods_matrix_key')
    }

    stages {

        // stage('cleanup') {
        //     steps {
        //         sh '''
        //             docker rmi $(docker images -q) -f
        //         '''
        //     }
        // }

        stage('clone') {
            steps {
                    git branch: 'refactor', credentialsId: 'jenkins-food-api-github-token', url: 'https://github.com/deVamshi/golang_food_delivery_api/'
                }
        }

        stage ('build') {
            steps {
                sh '''
                    docker build --tag yumfoods:latest .
                '''
            }
        }

        stage('run') {
            steps {
                sh '''
                    docker run -d -p 4000:4001 -e "PORT=4001" -e 'MONGODB_URI=$MONGODB_URI' -e 'MATRIX_KEY=$MATRIX_KEY' -e 'SECRET_KEY=$SECRET_KEY' yumfoods:latest
                '''
            }
        }

    }
}
