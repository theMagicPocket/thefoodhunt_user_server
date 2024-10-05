pipeline {

    agent any

    environment {
        MONGODB_URI=credentials('yumfoods_mongodb_uri')
        SECRET_KEY=credentials('yumfoods_secret_key')
        MATRIX_KEY=credentials('yumfoods_matrix_key')
    }

    tools {
        go '1.23.1'
    }

    stages {

        stage('cleanup') {
            steps {
                sh '''
                    allcntrs=$(docker ps --all -q)
                    if [ -n "$allcntrs" ]; then
                        runcntrs=$(docker ps -q)
                        if [ -n "$runcntrs" ]; then
                            echo "stopping running containers"
                            docker stop $runcntrs
                        fi
                        echo "removing all containers"
                        docker rm $allcntrs
                    fi

                    imgs=$(docker images -q)
                    if [ -n "$imgs" ]; then
                        echo "cleaning previous docker images"
                        docker rmi -f $imgs
                    fi

                    docker system prune -a
                '''
            }
        }

        stage('clone') {
            steps {
                    git branch: 'refactor', credentialsId: 'jenkins_food_hunt_user_server', url: 'https://github.com/theMagicPocket/thefoodhunt_user_server.git'
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
                sh "export JENKINS_NODE_COOKIE=dontKillMe;docker run -d -p 4000:4000 -e 'PORT=4000' -e 'MONGODB_URI=$MONGODB_URI' -e 'MATRIX_KEY=$MATRIX_KEY' -e 'SECRET_KEY=$SECRET_KEY' yumfoods:latest"
            }
        }

    }
}
