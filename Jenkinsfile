pipeline {
    agent any

    triggers {
        pollSCM('* * * * *')
    }

    stages {

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

        stage('test?') {
            agent none
            steps {
                script {
                    try {
                        timeout(time: 60, unit: 'SECONDS') {
                            script {
                                env.deploy_test = input message: 'Deploy to test?', id: 'deploy_test', parameters: [choice(name: 'Deploy to test', choices: 'no\nyes', description: 'Choose "yes" if you want to deploy to test')]
                            }
                        }
                    } catch (ignored) {
                        echo 'Skipping test deployment'
                    }
                }
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
                environment name: 'deploy_test', value: 'yes'
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