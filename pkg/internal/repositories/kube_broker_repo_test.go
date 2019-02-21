package repositories_test

import (
	"context"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/repositories"
)

var _ = Describe("KubeBrokerRepo", func() {
	var (
		kubeClient client.Client
		repo       KubeBrokerRepo

		existingBroker = &v1alpha1.Broker{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "broker-1",
				Namespace: "default",
			},
			Spec: v1alpha1.BrokerSpec{
				Name:     "broker-1",
				URL:      "http://example.org/broker",
				Username: "john",
				Password: "welcome",
			},
		}

		cleanup = func() {
			kubeClient.Delete(context.Background(), existingBroker)
		}
	)

	BeforeEach(func() {
		var err error

		kubeClient, err = buildKubeClient()
		Expect(err).NotTo(HaveOccurred())
		repo = NewKubeBrokerRepo(kubeClient)

		cleanup()
	})

	Describe("Get", func() {
		When("the broker exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), existingBroker)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns broker when it finds the broker", func() {
				broker, err := repo.Get(types.NamespacedName{Name: "broker-1", Namespace: "default"})
				Expect(err).NotTo(HaveOccurred())

				Expect(broker).To(Equal(existingBroker))
			})
		})

		When("the broker doesn't exist", func() {
			It("returns an error", func() {
				_, err := repo.Get(types.NamespacedName{Name: "broker-1", Namespace: "default"})

				Expect(err).To(MatchError("brokers.osbapi.ism.io \"broker-1\" not found"))
			})
		})
	})
})

// FIXME: consider de-duplicating with helper in acceptance tests
func buildKubeClient() (client.Client, error) {
	home := os.Getenv("HOME")
	kubeconfigFilepath := fmt.Sprintf("%s/.kube/config", home)
	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigFilepath)
	if err != nil {
		return nil, err
	}

	if err := v1alpha1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}

	return client.New(clientConfig, client.Options{Scheme: scheme.Scheme})
}
