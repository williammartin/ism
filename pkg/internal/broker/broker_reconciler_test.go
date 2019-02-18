package broker_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	osbapiv1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/broker"
	"github.com/pivotal-cf/ism/pkg/internal/broker/brokerfakes"
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
		expectedBroker             osbapiv1alpha1.Broker
		kubeGetStub                = func(_ context.Context, name types.NamespacedName, result runtime.Object) error {
			t, ok := result.(*osbapiv1alpha1.Broker)
			Expect(ok).To(BeTrue())
			*t = expectedBroker
			return nil
		}
	)

	BeforeEach(func() {
		fakeKubeClient = &brokerfakes.FakeKubeClient{}
		fakeKubeStatusWriter = &brokerfakes.FakeKubeStatusWriter{}
		fakeBrokerClient = &brokerfakes.FakeBrokerClient{}
		fakeKubeClient.GetCalls(kubeGetStub)
		createBrokerClient = func(config *osbapi.ClientConfiguration) (osbapi.Client, error) {
			brokerClientConfiguredWith = config
			return fakeBrokerClient, nil
		}

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

		fakeKubeClient.StatusReturns(fakeKubeStatusWriter)
	})

	JustBeforeEach(func() {
		reconciler = NewBrokerReconciler(fakeKubeClient, createBrokerClient)

		_, err = reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{Name: "broker-1", Namespace: "default"},
		})
	})

	It("fetches the broker resource using the kube client", func() {
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeKubeClient.GetCallCount()).To(Equal(1))
		_, namespacedName, _ := fakeKubeClient.GetArgsForCall(0)
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

	When("fetching the broker resource using the kube client fails", func() {
		BeforeEach(func() {
			fakeKubeClient.GetReturns(errors.New("error-getting-broker"))
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
})
