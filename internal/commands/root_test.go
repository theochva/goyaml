package commands

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRootCommand(t *testing.T) {
	RegisterFailHandler(Fail)
}

var _ = Describe("Root command scenarios", func() {
	When("No params specified", func() {
		It("prints out help", func() {

			out, err := runCommand("", "--help")
			Expect(err).ToNot(HaveOccurred())
			Expect(out).To(Equal(getHelpTextForCommand("")))
		})
	})
})
