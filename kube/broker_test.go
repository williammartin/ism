package kube_test

import (
	"context"

	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/pivotal-cf/ism/kube"
	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

var _ = Describe("Broker", func() {

	var (
		kubeClient client.Client

		broker *Broker
	)

	BeforeEach(func() {
		var err error
		// TODO: should we use the test env stuff, but point it to a real cluster?
		// i.e. to ensure the up-to-date CRD is installed?
		kubeClient, err = buildKubeClient()
		Expect(err).NotTo(HaveOccurred())

		broker = &Broker{
			KubeClient: kubeClient,
		}
	})

	Describe("Register", func() {
		var err error

		JustBeforeEach(func() {
			b := &osbapi.Broker{
				Name:     "broker-1",
				URL:      "broker-1-url",
				Username: "broker-1-username",
				Password: "broker-1-password",
			}

			err = broker.Register(b)
		})

		AfterEach(func() {
			bToDelete := &v1alpha1.Broker{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "broker-1",
					Namespace: "default",
				},
			}
			Expect(kubeClient.Delete(context.TODO(), bToDelete)).To(Succeed())
		})

		It("creates a new Broker resource instance", func() {
			Expect(err).NotTo(HaveOccurred())

			key := types.NamespacedName{
				Name:      "broker-1",
				Namespace: "default",
			}

			fetched := &v1alpha1.Broker{}
			Expect(kubeClient.Get(context.TODO(), key, fetched)).To(Succeed())

			Expect(fetched.Spec).To(Equal(v1alpha1.BrokerSpec{
				Name:     "broker-1",
				URL:      "broker-1-url",
				Username: "broker-1-username",
				Password: "broker-1-password",
			}))
		})

		When("creating a new Broker fails", func() {
			BeforeEach(func() {
				// register the broker first, so the second register errors
				b := &osbapi.Broker{
					Name:     "broker-1",
					URL:      "broker-1-url",
					Username: "broker-1-username",
					Password: "broker-1-password",
				}

				Expect(broker.Register(b)).To(Succeed())
			})

			It("propagates the error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
