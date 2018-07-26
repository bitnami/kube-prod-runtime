#!groovy

// Assumed jenkins plugins:
// - ansicolor
// - custom-tools-plugin
// - pipeline-utility-steps (readJSON)
// - kubernetes
// - jobcacher
// - azure-credentials

// Force using our pod
def label = UUID.randomUUID().toString()

def runIntegrationTest(String platform, String kubeprodArgs, Closure setup) {
  container('go') {
    withEnv(["KUBECONFIG=${env.WORKSPACE}/.kubeconf", "HOME=${env.WORKSPACE}"]) {
      stage("${platform} setup") {
        setup()
      }

      stage("${platform} install") {
        dir('do-install') {
          unstash 'installer'
          unstash 'manifests'

          // FIXME: we should have a better "test mode",
          // that uses letsencrypt-staging, fewer
          // replicas, etc.  My plan is to do that via
          // some sort of custom jsonnet overlay, since
          // power users will want similar flexibility.

          sh "./kubeprod -v=1 install --platform=${platform} --manifests=manifests --email=foo@example.com ${kubeprodArgs}"
        }
      }

      stage("${platform} test") {
        sh 'go get github.com/onsi/ginkgo/ginkgo'
        unstash 'src'
        dir('tests') {
          ansiColor('xterm') {
            sh 'ginkgo --tags integration -r --randomizeAllSpecs --randomizeSuites --failOnPending --trace --progress --slowSpecThreshold=300 --compilers=2 --nodes=4 -- --junit junit'
            junit 'junit/*.xml'
          }
        }
      }
    }
  }
}


podTemplate(
label: label,
idleMinutes: 1,  // Allow some best-effort reuse between successive stages
containers: [
  containerTemplate(
  name: 'go',
  image: 'golang:1.10.1-stretch',
  ttyEnabled: true, command: 'cat',
  resourceLimitCpu: '2000m',
  resourceLimitMemory: '2Gi',
  resourceRequestCpu: '10m',     // rely on burst CPU
  resourceRequestMemory: '1Gi',  // but actually need ram to avoid oom killer
  ),
  containerTemplate(
  name: 'az',
  image: 'microsoft/azure-cli:2.0.30',
  ttyEnabled: true, command: 'cat',
  resourceLimitCpu: '100m',
  resourceLimitMemory: '500Mi',
  resourceRequestCpu: '1m',
  resourceRequestMemory: '100Mi',
  ),
]) {

  env.http_proxy = 'http://proxy.webcache:80/'  // Note curl/libcurl needs explicit :80 !
  // Teach jenkins about the 'go' container env vars
  env.PATH = '/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin'
  env.GOPATH = '/go'

  node(label) {
    stage('Checkout') {
      checkout scm
      stash includes: '**', name: 'src'
    }
  }

  parallel(
  installer: {
    node(label) {
      container('go') {
        stage('installer') {
          timeout(time: 30) {
            cache(maxCacheSize: 1000, caches: [
                    [$class: 'ArbitraryFileCache', path: "${env.HOME}/.cache/go-build"],
                  ]) {
              withEnv(["GOPATH+WS=${env.WORKSPACE}", "PATH+GOBIN=${env.WORKSPACE}/bin"]) {
                dir('src/github.com/bitnami/kube-prod-runtime') {
                  unstash 'src'

                  dir('kubeprod') {
                    sh 'go version'
                    sh 'make all'
                    sh 'make test'
                    sh 'make vet'

                    sh 'make release VERSION=$BUILD_TAG'
                    dir('release') {
                      sh './kubeprod --help'
                      stash includes: 'kubeprod', name: 'installer'
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  },

  manifests: {
    node(label) {
      container('go') {
        stage('manifests') {
          timeout(time: 30) {
            cache(maxCacheSize: 1000, caches: [
                    [$class: 'ArbitraryFileCache', path: "${env.HOME}/.cache/go-build"],
                  ]) {
              unstash 'src'

              sh 'apt-get -qy update && apt-get -qy install make'
              // TODO: use tool, once the next release is made
              sh 'go get github.com/ksonnet/kubecfg'

              dir('manifests') {
                sh 'make validate KUBECFG="kubecfg -v"'
              }
              stash includes: 'manifests/**', excludes: 'manifests/Makefile', name: 'manifests'
            }
          }
        }
      }
    }
  })

  def platforms = [:]

  def minikubeKversions = []  // fixme: disabled minikube for now ["v1.8.0", "v1.9.6"]
  for (x in minikubeKversions) {
    def kversion = x  // local bind required because closures
    def platform = "minikube-0.25+k8s-" + kversion[1..3]
    platforms[platform] = {
      timeout(60) {
        node(label) {
          runIntegrationTest(platform, "") {
            withEnv(["PATH+TOOL=${tool 'minikube'}:${tool 'kubectl'}"]) {
              cache(maxCacheSize: 1000, caches: [
                      [$class: 'ArbitraryFileCache', path: "${env.HOME}/.minikube/cache"],
                    ]) {
                sh 'apt-get -qy update && apt-get install -qy libvirt-clients libvirt-daemon-system virtualbox'
                sh "minikube start --kubernetes-version=${kversion}"
              }
            }
          }
        }
      }
    }
  }

  // See:
  //  az aks get-versions -l centralus
  //    --query 'sort(orchestrators[?orchestratorType==`Kubernetes`].orchestratorVersion)'
  def aksKversions = ["1.8.7", "1.9.6"]
  for (x in aksKversions) {
    def kversion = x  // local bind required because closures
    def platform = "aks+k8s-" + kversion[0..2]
    platforms[platform] = {
      timeout(60) {
        // NB: `kubeprod` also uses az cli credentials and
        // $AZURE_SUBSCRIPTION_ID, $AZURE_TENANT_ID.
        withCredentials([azureServicePrincipal('azure-cli-2018-04-06-01-39-19')]) {
          def resourceGroup = 'prod-runtime-rg'
          def dnszone = "${platform}".replaceAll(/[^a-zA-Z0-9-]/, '-') + '.' + "${env.BUILD_TAG}".replaceAll(/[^a-zA-Z0-9-]/, '-') + '.test'
          node(label) {
            container('go') {
              def aks
              try {
                runIntegrationTest(platform, "--dns-resource-group=${resourceGroup} --dns-zone=${dnszone}") {
                  container('az') {
                    sh '''
                       az login --service-principal -u $AZURE_CLIENT_ID -p $AZURE_CLIENT_SECRET -t $AZURE_TENANT_ID
                       az account set -s $AZURE_SUBSCRIPTION_ID
                      '''
                    // Usually, `az aks create` creates a new service
                    // principal, which is not removed by `az aks
                    // delete`. We reuse an existing principal here to
                    // a) avoid this leak b) avoid having to give the
                    // "outer" principal (above) the power to create
                    // new service principals.
                    withCredentials([azureServicePrincipal('azure-cli-2018-04-06-03-17-41')]) {
                      def name = "${env.BUILD_TAG}-${platform}".replaceAll(/[^a-zA-Z0-9-]/, '-')
                      aks = readJSON(text: sh(script: """
                         az aks create                        \
                            --resource-group ${resourceGroup} \
                            --name ${name}                    \
                            --node-count 3                    \
                            --node-vm-size Standard_DS2_v2    \
                            --location eastus                 \
                            --kubernetes-version ${kversion}  \
                            --generate-ssh-keys                   \
                            --service-principal \$AZURE_CLIENT_ID \
                            --client-secret \$AZURE_CLIENT_SECRET \
                            --tags 'platform=${platform}' 'branch=${env.BRANCH_NAME}' 'build=${env.BUILD_URL}'
                        """, returnStdout: true))
                    }
                    sh "az aks get-credentials --name ${aks.name} --resource-group ${aks.resourceGroup} --admin --file \$KUBECONFIG"
                  }
                }
              }
              finally {
                if (aks) {
                  container('az') {
                    sh '''
                       az login --service-principal -u $AZURE_CLIENT_ID -p $AZURE_CLIENT_SECRET -t $AZURE_TENANT_ID
                       az account set -s $AZURE_SUBSCRIPTION_ID
                      '''
                    sh "az aks delete --yes --name ${aks.name} --resource-group ${aks.resourceGroup} --no-wait"
                  }
                }
              }
            }
          }
        }
      }
    }
  }

  parallel platforms
}
