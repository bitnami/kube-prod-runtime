package integration

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	// For client auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var junitDir = flag.String("junit", "", "Write junit results to dir if specified")
var kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")

func clusterConfigOrDie() *rest.Config {
	var config *rest.Config
	var err error

	if *kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		panic(err.Error())
	}

	return config
}

func createNsOrDie(c corev1.NamespacesGetter, ns string) string {
	result, err := c.Namespaces().Create(
		&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: ns,
			},
		})
	if err != nil {
		panic(err.Error())
	}
	name := result.GetName()
	fmt.Fprintf(GinkgoWriter, "Created namespace %s\n", name)
	return name
}

func deleteNsOrDie(c corev1.NamespacesGetter, ns string) {
	if ns == "" {
		return
	}
	err := c.Namespaces().Delete(ns, &metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
}

func TestE2e(t *testing.T) {
	var myReporters []Reporter
	RegisterFailHandler(Fail)

	if *junitDir != "" {
		if err := os.MkdirAll(*junitDir, 0777); err != nil {
			t.Fatalf("Failed to create %s: %v", *junitDir, err)
		}
		fname := fmt.Sprintf("junit_%d.xml", config.GinkgoConfig.ParallelNode)
		junitReporter := reporters.NewJUnitReporter(filepath.Join(*junitDir, fname))
		myReporters = append(myReporters, junitReporter)
	}

	RunSpecsWithDefaultAndCustomReporters(t, "kube-prod-runtime integration tests", myReporters)
}
