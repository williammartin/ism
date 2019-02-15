package kube_test

import (
	"context"

	"k8s.io/client-go/kubernetes/scheme"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/pivotal-cf/ism/kube"
	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

var _ = Describe("Plan", func() {
	var (
		kubeClient client.Client
		plan       *Plan
	)

	BeforeEach(func() {
		var err error
		kubeClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())

		plan = &Plan{
			KubeClient: kubeClient,
		}
	})

	Describe("FindByService", func() {
		var (
			plans []*osbapi.Plan
			err   error
		)

		BeforeEach(func() {
			planResource := &v1alpha1.BrokerServicePlan{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "plan-1",
					Namespace: "default",
				},
				Spec: v1alpha1.BrokerServicePlanSpec{
					Name:      "plan-1",
					ServiceID: "service-1",
				},
			}
			Expect(kubeClient.Create(context.TODO(), planResource)).To(Succeed())

			planResource2 := &v1alpha1.BrokerServicePlan{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "plan-2",
					Namespace: "default",
				},
				Spec: v1alpha1.BrokerServicePlanSpec{
					Name:      "plan-2",
					ServiceID: "service-2",
				},
			}
			Expect(kubeClient.Create(context.TODO(), planResource2)).To(Succeed())
		})

		JustBeforeEach(func() {
			plans, err = plan.FindByService("service-1")
		})

		AfterEach(func() {
			deletePlans(kubeClient, "plan-1")
		})

		It("returns plans by service id", func() {
			Expect(err).NotTo(HaveOccurred())

			Expect(plans).To(Equal([]*osbapi.Plan{{
				Name:      "plan-1",
				ServiceID: "service-1",
			}}))
		})
	})
})

func deletePlans(kubeClient client.Client, planNames ...string) {
	for _, p := range planNames {
		pToDelete := &v1alpha1.BrokerServicePlan{
			ObjectMeta: metav1.ObjectMeta{
				Name:      p,
				Namespace: "default",
			},
		}
		Expect(kubeClient.Delete(context.TODO(), pToDelete)).To(Succeed())
	}
}
