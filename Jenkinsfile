def dockerImage = null
def configs = [
    projectName: <INSERT PROJECTNAME>,
    namespace: <INSERT NAMESPACE,
    TAG: null,
    GIT_COMMIT: null,
    workspace: ''
]

def kubeConfig = { context, configuration ->
    return withCredentials([file(credentialsId: 'multi-tenant-kubeconfig', variable: 'config')]) {
        sh "cp -f $config ${configuration.workspace}/kube-config"
        sh "chmod 755 ${configuration.workspace}/kube-config"
        sh "kubectl config --kubeconfig ${configuration.workspace}/kube-config use-context $context"
    }
}

def terraformApply = { environment -> 
    return sh("terraform -chdir=./deployments/terraform/${environment} init -reconfigure && terraform -chdir=./deployments/terraform/${environment} apply -auto-approve")
}

def helmUpgrade = { configuration, environment ->
    return sh("helm upgrade ${configuration.namespace}-${configuration.projectName} deployments/helm/charts --set namespace=${configuration.namespace} --set environment=${environment} --set image.tag=${configuration.TAG ?: configuration.GIT_COMMIT} -f deployments/helm/charts/values.yaml --namespace ${configuration.namespace} --kubeconfig ${configuration.workspace}/kube-config --install")
}

pipeline {
    agent any
    triggers {
        pollSCM("H/5 * * * *")
    }
    environment {
        GIT_COMMIT = sh(script: "git rev-parse --short HEAD", returnStdout:true).trim()
    }
    stages {
        stage("Code") {
            steps {
                echo "Checking out the code..."
                checkout([$class: 'GitSCM',
                    branches: scm.branches,
                    doGenerateSubmoduleConfigurations: scm.doGenerateSubmoduleConfigurations,
                    extensions: [[$class: 'CloneOption', noTags: false, shallow: false, depth: 0, reference: '']],
                    userRemoteConfigs: scm.userRemoteConfigs])
            }
        }
        stage("Build") {
            environment {
                GITHUB_CREDS = credentials('c6d839e3-32c5-4fdd-be8a-2a5ec378d650')
            }
            steps {
                script {
                    configs.GIT_COMMIT = GIT_COMMIT
                    configs.workspace = WORKSPACE

                    echo "Workspace: ${configs.workspace}"
                    echo "Docker building ${configs.GIT_COMMIT}..."
                    dockerImage = docker.build("ebsadmin/${configs.projectName}:${configs.GIT_COMMIT}", "--build-arg USERNAME=${GITHUB_CREDS_USR} --build-arg TOKEN=${GITHUB_CREDS_PSW} .")
                }
            }
        }
        stage("Publish Develop") {
            when { anyOf { branch "develop" } }
            steps {
                script {
                    echo "Publishing develop..."
                    docker.withRegistry("", "dockerhub") {
                        dockerImage.push()
                        dockerImage.push("latest")
                    }
                }
            }
        }
        stage("Publish Release") {
            when { anyOf { branch  "master" } }
            steps {
                script {
                    configs.TAG = sh(returnStdout: true, script: 'git tag --contains')?.trim() ?: ''
                    if(configs?.TAG?.empty) {
                        error("Cannot publish release. The TAG is empty.")
                    }
                    echo "Publishing release TAG: ${configs.TAG}..."
                    docker.withRegistry("", "dockerhub") {
                        dockerImage.push()
                        dockerImage.push("${configs.TAG}")
                    }
                }
            }
        }
        stage("Deploy To Dev") {
            when { anyOf { branch "develop" } }
            steps {
                script {
                    echo "Deploying to dev"
                    echo "Checking for Terraform directory"
                    if (fileExists("./deployments/terraform")) {
                        echo "Applying Terraform"
                        terraformApply("dev")
                    } else {
                        echo "No Terraform directory found"
                    }

                    echo "Checking for HelmCharts directory"
                    if (fileExists("./deployments/helm/charts")) {
                        kubeConfig('mt-dev-thor', configs)
                        helmUpgrade(configs, "dev")
                    } else {
                        echo "No HelmCharts directory found"
                    }
                }
            }
        }
        stage("Deploy To Staging") {
            when { anyOf { branch "master" } }
            steps {
                script {
                    echo  "Deploying ${configs.TAG} to staging"
                    kubeConfig('mt-stg-odin', configs)
                    helmUpgrade(configs, "stg")
                }
            }
        }
        stage("Confirm production deployment") {
            when { anyOf { branch "master" } }
            steps {
                timeout(time: 15, unit: 'MINUTES') {
                    input "Ready to deploy to prod?"
                }
            }
        }
        stage("Deploy To Production") {
            when { anyOf { branch "master" } }
            steps {
                script {
                    echo  "Deploying ${configs.TAG} to production"
                    kubeConfig('mt-prod-freya', configs)
                    helmUpgrade(configs, "prod")
                }
            }
        }
    }
}