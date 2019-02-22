package repositories_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/repositories"
)

var _ = Describe("KubeServiceRepo", func() {
	var (
		repo          KubeServiceRepo
		broker        *v1alpha1.Broker
		brokerService *v1alpha1.BrokerService
	)

	BeforeEach(func() {
		broker = &v1alpha1.Broker{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "broker-1",
				Namespace: "default",
			},
		}

		brokerService = &v1alpha1.BrokerService{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "broker-1.service-id-1",
				Namespace: "default",
			},
		}

		repo = NewKubeServiceRepo(kubeClient)
	})

	AfterEach(func() {
		kubeClient.Delete(context.Background(), broker)
	})

	Describe("Create", func() {
		When("the broker exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), broker)
				Expect(err).NotTo(HaveOccurred())

				err = repo.Create(broker, osbapi.Service{
					ID:          "service-id-1",
					Name:        "service-one",
					Description: "cool description",
				})
				Expect(err).NotTo(HaveOccurred())

				err = kubeClient.Get(context.Background(), types.NamespacedName{Name: "broker-1.service-id-1", Namespace: "default"}, brokerService)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				kubeClient.Delete(context.Background(), brokerService)
			})

			It("creates the service with the correct spec", func() {
				Expect(brokerService.Spec).To(Equal(v1alpha1.BrokerServiceSpec{
					Name:        "service-one",
					Description: "cool description",
					BrokerID:    "broker-1",
				}))
			})

			It("generates the correct name and namespace", func() {
				Expect(brokerService.ObjectMeta.Name).To(Equal("broker-1.service-id-1"))
				Expect(brokerService.ObjectMeta.Namespace).To(Equal("default"))
			})

			It("sets the owner reference of the service to the broker", func() {
				Expect(brokerService.ObjectMeta.OwnerReferences).To(HaveLen(1))
				Expect(brokerService.ObjectMeta.OwnerReferences[0].UID).To(Equal(broker.ObjectMeta.UID))
			})
		})

		When("the broker doesn't exist", func() {
			It("returns an error", func() {
				invalidBroker := &v1alpha1.Broker{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "broker-without-uid",
					},
				}

				err := repo.Create(invalidBroker, osbapi.Service{
					ID:          "service-id-1",
					Name:        "service-one",
					Description: "cool description",
				})

				Expect(err).To(MatchError("BrokerService.osbapi.ism.io \"broker-without-uid.service-id-1\" is invalid" +
					": metadata.ownerReferences.uid: Invalid value: \"\": uid must not be empty"))
			})
		})
	})
})
