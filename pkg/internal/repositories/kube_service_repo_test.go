package repositories_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/repositories"
)

var _ = Describe("KubeServiceRepo", func() {
	var (
		kubeClient client.Client
		repo       KubeServiceRepo

		existingBroker *v1alpha1.Broker
		brokerResource = types.NamespacedName{Name: "broker-1", Namespace: "default"}
		brokerMeta     = metav1.ObjectMeta{Name: brokerResource.Name, Namespace: brokerResource.Namespace}

		serviceMeta = metav1.ObjectMeta{Name: "broker-1.service-id-1", Namespace: "default"}

		cleanup = func() {
			kubeClient.Delete(context.Background(), &v1alpha1.Broker{ObjectMeta: brokerMeta})
		}
	)

	BeforeEach(func() {
		existingBroker = &v1alpha1.Broker{
			ObjectMeta: brokerMeta,
		}

		var err error

		kubeClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())
		repo = NewKubeServiceRepo(kubeClient)

		cleanup()
	})

	Describe("Create", func() {

		When("broker exists", func() {
			var (
				service *v1alpha1.BrokerService
			)

			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), existingBroker)
				Expect(err).NotTo(HaveOccurred())

				err = repo.Create(existingBroker, osbapi.Service{
					ID:          "service-id-1",
					Name:        "service-one",
					Description: "cool description",
				})
				Expect(err).NotTo(HaveOccurred())

				service = &v1alpha1.BrokerService{}
				err = kubeClient.Get(context.Background(), types.NamespacedName{Name: "broker-1.service-id-1", Namespace: "default"}, service)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				kubeClient.Delete(context.Background(), &v1alpha1.BrokerService{ObjectMeta: serviceMeta})
			})

			It("creates the service with correct spec", func() {
				Expect(service.Spec).To(Equal(v1alpha1.BrokerServiceSpec{
					Name:        "service-one",
					Description: "cool description",
					BrokerID:    "broker-1",
				}))
			})

			It("generates correct name and namespace", func() {
				Expect(service.ObjectMeta.Name).To(Equal("broker-1.service-id-1"))
				Expect(service.ObjectMeta.Namespace).To(Equal("default"))
			})

			It("sets ownership of that service to the broker resource", func() {
				Expect(service.ObjectMeta.OwnerReferences).To(HaveLen(1))
				Expect(service.ObjectMeta.OwnerReferences[0].UID).To(Equal(existingBroker.ObjectMeta.UID))
			})
		})

		When("broker doesn't exist", func() {
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
