package actors_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestActors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Actors Suite")
}
