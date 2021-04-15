// vim set syntax=jsonnet
// Imports
local kube = import "../vendor/github.com/bitnami-labs/kube-libsonnet/kube.libsonnet";
local kubecfg = import "kubecfg.libsonnet";

// Cluster-specific configuration
local BLACKBOX_EXPORTER_IMAGE = (import "images.json")["blackbox-exporter"];

{
  local this = self,

  be_config:: {
   modules: {
      http_2xx: {
        prober: "http"
      },
      http_post_2xx: {
        prober: "http",
        http: {
          method: "POST"
        }
      },
      tcp_connect: {
        prober: "tcp"
      },
      pop3s_banner: {
        prober: "tcp",
        tcp: {
          query_response: [ { expect: "^+OK" } ],
          tls: true,
          tls_config: { insecure_skip_verify: false }
        }
      },
      ssh_banner: {
        prober: "tcp",
        tcp: {
          query_response: [ { expect: "^SSH-2.0-" } ]
        }
      },
      irc_banner: {
        prober: "tcp",
        tcp: {
          query_response: [
            { send: "NICK prober" },
            { send: "USER prober prober prober :prober" },
            {
              expect: "PING :([^ ]+)",
              send: "PONG ${1}"
            },
            { expect: "^:[^ ]+ 001" }
          ]
        }
      },
      icmp: { prober: "icmp" }
    }
  },

  blackbox_exporter_config: kube.ConfigMap($.p +"blackbox-exporter-config") + $.metadata {
    data+: {
      "config.yml": kubecfg.manifestYaml(this.be_config),
    },
  },

  deploy: kube.Deployment($.p + "blackbox-exporter") + $.metadata {
    spec+: {
      template+: {
        spec+: {
          volumes_+: {
            blackbox_exporter_config: kube.ConfigMapVolume(this.blackbox_exporter_config),
          },
          securityContext+: {
            fsGroup: 1001,
          },
          containers_+: {
            default: kube.Container("blackbox-exporter") {
              local this = self,

              image: BLACKBOX_EXPORTER_IMAGE,
              ports_+: {
                probe: {containerPort: 9115},
              },
              livenessProbe: {
                httpGet: {path: "/", port: "probe"},
              },
              readinessProbe: self.livenessProbe {
                successThreshold: 2,
              },
              securityContext+: {
                runAsUser: 1001,
              },
              volumeMounts_+: {
                blackbox_exporter_config: {
                  mountPath: "/opt/bitnami/blackbox-exporter/blackbox.yml",
                  readOnly: true,
                  subPath: "config.yml",
                }
              }
            }
          }
        }
      }
    }
  },
  svc: kube.Service($.p + "blackbox-exporter") + $.metadata {
    target_pod: this.deploy.spec.template,
  },
}
