pipeline {
    agent any

    triggers {
        pollSCM('* * * * *')
    }

    stages {

        stage('build') {
            agent {
                docker {
                    image 'golang'
                    args '-u root -v ${env.WORKSPACE}:/go/src/github.com/ONSdigital/rm-survey-service'
                }
            }

            steps {
                sh "cd /go/src/github.com/ONSdigital/rm-survey-service && make"
            }
        }

        stage('dev') {
            agent {
                docker {
                    image 'governmentpaas/cf-cli'
                    args '-u root'
                }
            }

            environment {
                CLOUDFOUNDRY_API = credentials('CLOUDFOUNDRY_API')
                CF_DOMAIN = credentials('CF_DOMAIN')
                DEV_SECURITY = credentials('DEV_SECURITY')
                CF_USER = credentials('CF_USER')
            }
            steps {
                sh "sed -i -- 's/SPACE/dev/g' *template.yml"
                sh "sed -i -- 's/INSTANCES/1/g' *template.yml"
                sh "sed -i -- 's/DATABASE/rm-pg-db/g' *template.yml"
                sh "sed -i -- 's/REPLACE_BA_USERNAME/${env.DEV_SECURITY_USR}/g' *template.yml"
                sh "sed -i -- 's/REPLACE_BA_PASSWORD/${env.DEV_SECURITY_PSW}/g' *template.yml"

                sh "cf login -a https://${env.CLOUDFOUNDRY_API} --skip-ssl-validation -u ${CF_USER_USR} -p ${CF_USER_PSW} -o rmras -s dev"
                sh 'cf push -f manifest-template.yml'
            }
        }

        stage('ci?') {
            agent none
            steps {
                script {
                    try {
                        timeout(time: 60, unit: 'SECONDS') {
                            script {
                                env.deploy_ci = input message: 'Deploy to CI?', id: 'deploy_ci', parameters: [choice(name: 'Deploy to CI', choices: 'no\nyes', description: 'Choose "yes" if you want to deploy to CI')]
                            }
                        }
                    } catch (ignored) {
                        echo 'Skipping ci deployment'
                    }
                }
            }
        }

        stage('ci') {
            agent {
                docker {
                    image 'governmentpaas/cf-cli'
                    args '-u root'
                }

            }
            when {
                environment name: 'deploy_ci', value: 'yes'
            }

            environment {
                CLOUDFOUNDRY_API = credentials('CLOUDFOUNDRY_API')
                CF_DOMAIN = credentials('CF_DOMAIN')
                CI_SECURITY = credentials('CI_SECURITY')
                CF_USER = credentials('CF_USER')
            }
            steps {
                sh "sed -i -- 's/SPACE/dev/g' *template.yml"
                sh "sed -i -- 's/INSTANCES/1/g' *template.yml"
                sh "sed -i -- 's/DATABASE/rm-pg-db/g' *template.yml"
                sh "sed -i -- 's/REPLACE_BA_USERNAME/${env.CI_SECURITY_USR}/g' *template.yml"
                sh "sed -i -- 's/REPLACE_BA_PASSWORD/${env.CI_SECURITY_PSW}/g' *template.yml"

                sh "cf login -a https://${env.CLOUDFOUNDRY_API} --skip-ssl-validation -u ${CF_USER_USR} -p ${CF_USER_PSW} -o rmras -s ci"
                sh 'cf push -f manifest-template.yml'
            }
        }

        stage('release?') {
            agent none
            steps {
                script {
                    try {
                        timeout(time: 60, unit: 'SECONDS') {
                            script {
                                env.do_release = input message: 'Do a release?', id: 'do_release', parameters: [choice(name: 'Deploy to test', choices: 'no\nyes', description: 'Choose "yes" if you want to create a tag')]
                            }
                        }
                    } catch (ignored) {
                        echo 'Skipping test deployment'
                    }
                }
            }
        }

        stage('release') {
            agent {
                docker {
                    image 'node'
                    args '-u root'
                }

            }
            environment {
                GITHUB_API_KEY = credentials('GITHUB_API_KEY')
            }
            when {
                environment name: 'do_release', value: 'yes'
            }
            steps {
                // Prune any local tags created by any other builds
                sh "git tag -l | xargs git tag -d && git fetch -t"
                sh "git remote set-url origin https://ons-sdc:${GITHUB_API_KEY}@github.com/ONSdigital/response-operations-ui.git"
                sh "npm install -g bmpr"
                sh "bmpr patch|xargs git push origin"
            }
        }

        stage('test') {
            agent {
                docker {
                    image 'governmentpaas/cf-cli'
                    args '-u root'
                }

            }
            when {
                environment name: 'do_release', value: 'yes'
            }

            environment {
                CLOUDFOUNDRY_API = credentials('CLOUDFOUNDRY_API')
                CF_DOMAIN = credentials('CF_DOMAIN')
                TEST_SECURITY = credentials('TEST_SECURITY')
                CF_USER = credentials('CF_USER')
            }
            steps {
                sh "sed -i -- 's/SPACE/test/g' *template.yml"
                sh "sed -i -- 's/INSTANCES/1/g' *template.yml"
                sh "sed -i -- 's/DATABASE/rm-pg-db/g' *template.yml"
                sh "sed -i -- 's/REPLACE_BA_USERNAME/${env.TEST_SECURITY_USR}/g' *template.yml"
                sh "sed -i -- 's/REPLACE_BA_PASSWORD/${env.TEST_SECURITY_PSW}/g' *template.yml"

                sh "cf login -a https://${env.CLOUDFOUNDRY_API} --skip-ssl-validation -u ${CF_USER_USR} -p ${CF_USER_PSW} -o rmras -s test"
                sh 'cf push -f manifest-template.yml'
            }
        }
    }

    post {
        always {
            cleanWs()
            dir('${env.WORKSPACE}@tmp') {
                deleteDir()
            }
            dir('${env.WORKSPACE}@script') {
                deleteDir()
            }
            dir('${env.WORKSPACE}@script@tmp') {
                deleteDir()
            }
        }
    }
}