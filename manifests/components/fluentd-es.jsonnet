local kube = import "kube.libsonnet";
local kubecfg = import "kubecfg.libsonnet";
local utils = import "utils.libsonnet";

local FLUENTD_ES_IMAGE = "k8s.gcr.io/fluentd-elasticsearch:v2.0.4";

// TODO(jjo):
// * confirm we're ok to split E-F-K jsonnets by their components
// * confirm we stil prefer "kube-system" instead of "logging" ns
{
  p:: "",
  namespace:: { metadata+: { namespace: "kube-system" } },
  criticalPod:: { metadata+: { annotations+: { "scheduler.alpha.kubernetes.io/critical-pod": "" } } },
  config:: (import "fluentd-es-config.jsonnet"),

  fluentd_es_config: kube.ConfigMap($.p + "fluentd-es") + $.namespace {
    data+: $.config,
  },
  fluentd_es: {
    local f = self,
    serviceAccount: kube.ServiceAccount($.p + "fluentd-es") + $.namespace,
    fluentdRole: kube.ClusterRole($.p + "fluentd-es") {
      rules: [
        {
          apiGroups: [""],
          resources: ["namespaces", "pods"],
          verbs: ["get", "watch", "list"],
        },
      ],
    },
    fluentdBinding: kube.ClusterRoleBinding($.p + "fluentd-es") {
      roleRef_: f.fluentdRole,
      subjects_+: [f.serviceAccount],
    },
    daemonset: kube.DaemonSet($.p + "fluentd-es") + $.namespace {
      spec+: {
        template+: $.criticalPod {
          spec+: {
            containers_+: {
              fluentd_es: kube.Container("fluentd-es") {
                image: FLUENTD_ES_IMAGE,
                env_+: {
                  FLUENTD_ARGS: "--no-supervisor -q",
                },
                resources: {
                  requests: { cpu: "100m", memory: "200Mi" },
                  limits: { memory: "500Mi" },
                },
                volumeMounts_+: {
                  varlog: { mountPath: "/var/log" },
                  varlibdockercontainers: {
                    mountPath: "/var/lib/docker/containers",
                    readOnly: true,
                  },
                  configvolume: {
                    mountPath: "/etc/fluent/config.d",
                  },
                },
              },
            },
            // Note: present in upstream to originally to cope with fluentd-es migration to DS, not applicable here
            // nodeSelector: {
            //  "beta.kubernetes.io/fluentd-ds-ready": "true",
            // },
            //
            // Note: from upstream, only for kube>=1.10?, may need to come from ../platforms
            // priorityClassName: "system-node-critical",
            serviceAccountName: f.serviceAccount.metadata.name,
            terminationGracePeriodSeconds: 30,
            volumes_+: {
              varlog: kube.HostPathVolume("/var/log"),
              varlibdockercontainers: kube.HostPathVolume("/var/lib/docker/containers"),
              configvolume: kube.ConfigMapVolume($.fluentd_es_config),
            },
          },
        },
      },
    },
  },
}
