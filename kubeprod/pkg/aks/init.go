package aks

import (
	"bufio"
	"html/template"
	"os"
	"path"
)

const (
	clusterTemplate = `# Cluster-specific configuration for cluster 'k8s1'

local config = import "config.json";
local aks = import "{{.ManifestsPath}}/platforms/aks+k8s-{{.KubernetesVersion}}.jsonnet";

aks + config {
	// Place your overrides here
}
`
	configTemplate = `{
	// Cluster-specific configuration
	"cluster": "{{.ClusterName}}",
	"external_dns_zone_name": "{{.DNS}}",
	"letsencrypt_contact_email": "{{.Email}}",
}
`
)

type variables struct {
	ClusterName       string
	ManifestsPath     string // path to the manifests/ directory (including trailing /)
	Email             string // contact e-mail for Letsencrypt certificates
	DNS               string // DNS domain
	KubernetesVersion string // Kubernetes version
}

// Executes the template inside the `templateData` variable performing
// substitutions from the `v` dictionary and write the results to the
// output file named as `pathName`.
func writeTemplate(pathName string, templateData string, v variables) (err error) {

	var f *os.File
	f, err = os.Create(pathName)
	if err != nil {
		return
	}

	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	w := bufio.NewWriter(f)

	defer func() {
		cerr := w.Flush()
		if err == nil {
			err = cerr
		}
	}()

	tmpl, err := template.New("template").Parse(templateData)
	err = tmpl.ExecuteTemplate(w, "template", v)
	return
}

// Init does init
func Init(clusterName string, manifestsBase string, email string, dnsZone string, kubernetesVersion string) (err error) {

	v := variables{
		ClusterName:       clusterName,
		ManifestsPath:     path.Clean(manifestsBase),
		Email:             email,
		DNS:               dnsZone,
		KubernetesVersion: kubernetesVersion,
	}

	err = writeTemplate("./"+clusterName+".json", clusterTemplate, v)
	if err != nil {
		return
	}
	err = writeTemplate("./config.json", configTemplate, v)
	if err != nil {
		return
	}
	return
}
