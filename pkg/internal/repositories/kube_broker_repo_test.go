package repositories_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/repositories"
)

var _ = Describe("KubeBrokerRepo", func() {
	var (
		kubeClient client.Client
		repo       KubeBrokerRepo

		existingBroker *v1alpha1.Broker
		resource       = types.NamespacedName{Name: "broker-1", Namespace: "default"}
		objectMeta     = metav1.ObjectMeta{Name: resource.Name, Namespace: resource.Namespace}

		cleanup = func() {
			kubeClient.Delete(context.Background(), &v1alpha1.Broker{ObjectMeta: objectMeta})
		}
	)

	BeforeEach(func() {
		existingBroker = &v1alpha1.Broker{
			ObjectMeta: objectMeta,
			Spec: v1alpha1.BrokerSpec{
				Name:     "broker-1",
				URL:      "http://example.org/broker",
				Username: "john",
				Password: "welcome",
			},
		}

		var err error

		kubeClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
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
				broker, err := repo.Get(resource)
				Expect(err).NotTo(HaveOccurred())

				Expect(broker).To(Equal(existingBroker))
			})
		})

		When("the broker doesn't exist", func() {
			It("returns an error", func() {
				_, err := repo.Get(resource)

				Expect(err).To(MatchError("brokers.osbapi.ism.io \"broker-1\" not found"))
			})
		})
	})

	Describe("UpdateStatus", func() {
		When("the broker exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), existingBroker)
				Expect(err).NotTo(HaveOccurred())
			})

			It("updates status", func() {
				newState := v1alpha1.BrokerStateRegistered
				Expect(existingBroker.Status.State).NotTo(Equal(newState))

				err := repo.UpdateState(existingBroker, newState)
				Expect(err).NotTo(HaveOccurred())

				updatedBroker, err := repo.Get(resource)
				Expect(err).NotTo(HaveOccurred())

				Expect(updatedBroker.Status.State).To(Equal(newState))
				Expect(existingBroker.Status.State).To(Equal(newState))
			})
		})

		When("the broker doesn't exist", func() {
			It("returns an error", func() {
				newState := v1alpha1.BrokerStateRegistered
				err := repo.UpdateState(existingBroker, newState)

				Expect(err).To(MatchError("brokers.osbapi.ism.io \"broker-1\" not found"))
			})
		})
	})
})
