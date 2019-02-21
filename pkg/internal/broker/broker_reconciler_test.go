package broker_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	osbapiv1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/broker"
	"github.com/pivotal-cf/ism/pkg/internal/broker/brokerfakes"
	"github.com/pivotal-cf/ism/pkg/internal/repositories/repositoriesfakes"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("BrokerReconciler", func() {
	var (
		fakeKubeClient             *brokerfakes.FakeKubeClient
		fakeKubeStatusWriter       *brokerfakes.FakeKubeStatusWriter
		fakeBrokerClient           *brokerfakes.FakeBrokerClient
		createBrokerClient         osbapi.CreateFunc
		reconciler                 *BrokerReconciler
		err                        error
		brokerClientConfiguredWith *osbapi.ClientConfiguration
		brokerName                 types.NamespacedName
		expectedBroker             osbapiv1alpha1.Broker
		fakeKubeBrokerRepo         *repositoriesfakes.FakeKubeBrokerRepo
	)

	BeforeEach(func() {
		fakeKubeClient = &brokerfakes.FakeKubeClient{}
		fakeKubeStatusWriter = &brokerfakes.FakeKubeStatusWriter{}
		fakeBrokerClient = &brokerfakes.FakeBrokerClient{}
		fakeKubeBrokerRepo = &repositoriesfakes.FakeKubeBrokerRepo{}

		createBrokerClient = func(config *osbapi.ClientConfiguration) (osbapi.Client, error) {
			brokerClientConfiguredWith = config
			return fakeBrokerClient, nil
		}
		brokerName = types.NamespacedName{Name: "broker-1", Namespace: "default"}

		expectedBroker = osbapiv1alpha1.Broker{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "broker-1",
				Namespace: "default",
			},
			Spec: osbapiv1alpha1.BrokerSpec{
				Name:     "broker-1",
				URL:      "broker-url",
				Username: "broker-username",
				Password: "broker-password",
			},
		}
		fakeKubeBrokerRepo.GetReturns(&expectedBroker, nil)

		fakeBrokerClient.GetCatalogReturns(&osbapi.CatalogResponse{
			Services: []osbapi.Service{
				{
					ID:          "id-service-1",
					Name:        "service-1",
					Description: "some fancy description",
					Plans:       []osbapi.Plan{{ID: "id-plan-1", Name: "plan-1"}},
				},
				{
					ID:          "id-service-2",
					Name:        "service-2",
					Description: "poorly written description",
					Plans:       []osbapi.Plan{{ID: "id-plan-2", Name: "plan-2"}, {ID: "id-plan-3", Name: "plan-3"}},
				},
			},
		}, nil)

		fakeKubeClient.StatusReturns(fakeKubeStatusWriter)
	})

	JustBeforeEach(func() {
		reconciler = NewBrokerReconciler(fakeKubeClient, createBrokerClient, fakeKubeBrokerRepo)

		_, err = reconciler.Reconcile(reconcile.Request{
			NamespacedName: brokerName,
		})
	})

	It("fetches the broker resource using the kube broker repo", func() {
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeKubeBrokerRepo.GetCallCount()).To(Equal(1))
		namespacedName := fakeKubeBrokerRepo.GetArgsForCall(0)
		Expect(namespacedName).To(Equal(types.NamespacedName{Name: "broker-1", Namespace: "default"}))
	})

	It("configures the broker client with correct options", func() {
		Expect(*brokerClientConfiguredWith).To(Equal(osbapi.ClientConfiguration{
			Name:                "broker-1",
			URL:                 "broker-url",
			APIVersion:          osbapi.LatestAPIVersion(),
			TimeoutSeconds:      60,
			EnableAlphaFeatures: false,
			AuthConfig: &osbapi.AuthConfig{
				BasicAuthConfig: &osbapi.BasicAuthConfig{
					Username: "broker-username",
					Password: "broker-password",
				},
			},
		}))
	})

	It("fetches the catalog using the broker client", func() {
		Expect(fakeBrokerClient.GetCatalogCallCount()).To(Equal(1))
	})

	It("updates the broker status to registered", func() {
		Expect(fakeKubeStatusWriter.UpdateCallCount()).To(Equal(1))
		_, obj := fakeKubeStatusWriter.UpdateArgsForCall(0)
		broker, ok := obj.(*osbapiv1alpha1.Broker)
		Expect(ok).To(BeTrue())
		Expect(broker.Status.State).To(Equal(osbapiv1alpha1.BrokerStateRegistered))
	})

	It("creates service resources using the kube client", func() {
		_, obj := fakeKubeClient.CreateArgsForCall(0)
		service, ok := obj.(*osbapiv1alpha1.BrokerService)
		Expect(ok).To(BeTrue())
		Expect(*service).To(Equal(osbapiv1alpha1.BrokerService{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "broker-1.id-service-1",
			},
			Spec: osbapiv1alpha1.BrokerServiceSpec{
				Name:        "service-1",
				Description: "some fancy description",
				//TODO: BrokerID    string `json:"brokerID"`
			},
		}))

		_, obj = fakeKubeClient.CreateArgsForCall(2)
		service, ok = obj.(*osbapiv1alpha1.BrokerService)
		Expect(ok).To(BeTrue())
		Expect(*service).To(Equal(osbapiv1alpha1.BrokerService{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "broker-1.id-service-2",
			},
			Spec: osbapiv1alpha1.BrokerServiceSpec{
				Name:        "service-2",
				Description: "poorly written description",
				//TODO: BrokerID    string `json:"brokerID"`
			},
		}))
	})

	It("creates plan resources using the kube client", func() {
		_, obj := fakeKubeClient.CreateArgsForCall(1)
		plan, ok := obj.(*osbapiv1alpha1.BrokerServicePlan)
		Expect(ok).To(BeTrue())
		Expect(*plan).To(Equal(osbapiv1alpha1.BrokerServicePlan{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "broker-1.id-service-1.id-plan-1",
			},
			Spec: osbapiv1alpha1.BrokerServicePlanSpec{
				Name: "plan-1",
				//TODO: ServiceID    string `json:"serviceID"`
			},
		}))

		_, obj = fakeKubeClient.CreateArgsForCall(3)
		plan, ok = obj.(*osbapiv1alpha1.BrokerServicePlan)
		Expect(ok).To(BeTrue())
		Expect(*plan).To(Equal(osbapiv1alpha1.BrokerServicePlan{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "broker-1.id-service-2.id-plan-2",
			},
			Spec: osbapiv1alpha1.BrokerServicePlanSpec{
				Name: "plan-2",
				//TODO: BrokerID    string `json:"brokerID"`
			},
		}))

		_, obj = fakeKubeClient.CreateArgsForCall(4)
		plan, ok = obj.(*osbapiv1alpha1.BrokerServicePlan)
		Expect(ok).To(BeTrue())
		Expect(*plan).To(Equal(osbapiv1alpha1.BrokerServicePlan{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "broker-1.id-service-2.id-plan-3",
			},
			Spec: osbapiv1alpha1.BrokerServicePlanSpec{
				Name: "plan-3",
				//TODO: BrokerID    string `json:"brokerID"`
			},
		}))
	})

	When("using different broker name and namespace", func() {
		BeforeEach(func() {
			brokerName = types.NamespacedName{Name: "other-broker", Namespace: "other"}
			expectedBroker.ObjectMeta.Name = brokerName.Name
			expectedBroker.ObjectMeta.Namespace = brokerName.Namespace
		})

		It("reflects that in the namespace and name of created service resources", func() {
			_, obj := fakeKubeClient.CreateArgsForCall(0)
			service, ok := obj.(*osbapiv1alpha1.BrokerService)
			Expect(ok).To(BeTrue())
			Expect(service.ObjectMeta).To(Equal(metav1.ObjectMeta{
				Namespace: "other",
				Name:      "other-broker.id-service-1",
			}))
		})

		It("reflects that in the namespace and name of created service plan resources", func() {
			_, obj := fakeKubeClient.CreateArgsForCall(1)
			plan, ok := obj.(*osbapiv1alpha1.BrokerServicePlan)
			Expect(ok).To(BeTrue())
			Expect(plan.ObjectMeta).To(Equal(metav1.ObjectMeta{
				Namespace: "other",
				Name:      "other-broker.id-service-1.id-plan-1",
			}))
		})
	})

	When("the broker state reports it is already registered", func() {
		BeforeEach(func() {
			expectedBroker.Status.State = osbapiv1alpha1.BrokerStateRegistered
		})

		It("doesn't call the broker", func() {
			Expect(fakeBrokerClient.GetCatalogCallCount()).To(Equal(0))
		})

		It("doesn't update the status", func() {
			Expect(fakeKubeStatusWriter.UpdateCallCount()).To(Equal(0))
		})

		It("still reconciles successfully ", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("updating the broker status errors", func() {
		BeforeEach(func() {
			fakeKubeStatusWriter.UpdateReturns(errors.New("error-updating-status"))
		})

		//TODO: test the state of service / plan creation here.
		It("returns the error", func() {
			Expect(err).To(MatchError("error-updating-status"))
		})
	})

	When("fetching the broker resource using the kube broker repo fails", func() {
		BeforeEach(func() {
			fakeKubeBrokerRepo.GetReturns(nil, errors.New("error-getting-broker"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-getting-broker"))
		})
	})

	When("configuring the broker client fails", func() {
		BeforeEach(func() {
			createBrokerClient = func(config *osbapi.ClientConfiguration) (osbapi.Client, error) {
				return nil, errors.New("error-configuring-broker-client")
			}
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-configuring-broker-client"))
		})
	})

	When("fetching the catalog using the broker client fails", func() {
		BeforeEach(func() {
			fakeBrokerClient.GetCatalogReturns(nil, errors.New("error-getting-catalog"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-getting-catalog"))
		})
	})

	When("creating service resource fails", func() {
		BeforeEach(func() {
			fakeKubeClient.CreateReturns(errors.New("error-creating-service"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-creating-service"))
		})
	})

	When("creating service plan resource fails", func() {
		BeforeEach(func() {
			fakeKubeClient.CreateReturnsOnCall(1, errors.New("error-creating-plan"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-creating-plan"))
		})
	})
})
