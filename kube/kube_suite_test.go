package kube_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestKube(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kube Suite")
}

func buildKubeClient() (client.Client, error) {
	home := os.Getenv("HOME")
	kubeconfigFilepath := fmt.Sprintf("%s/.kube/config", home)
	crdConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigFilepath)
	if err != nil {
		return nil, err
	}

	if err := v1alpha1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}

	return client.New(crdConfig, client.Options{Scheme: scheme.Scheme})
}
