package integration_test

import (
	"time"

	"github.com/benchkram/bobc/pkg/rnd"
	restserverclient "github.com/benchkram/bobc/rest-server-client"
	"github.com/benchkram/bobc/restserver/generated"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("should test the project lifecycle", func() {

	var client *restserverclient.C

	Context("check if tests are running", func() {

		It("should create rest-server client for this workflow", func() {
			var err error
			client, err = restserverclient.New(suite.restServer.Address(), suite.apiKey)
			Expect(err).NotTo(HaveOccurred())

			Eventually(client.Health, 60*time.Second, 1*time.Second).Should(BeTrue())
		})

		It("should wait for the server to be reachable", func() {
			Eventually(client.Health, 60*time.Second, 1*time.Second).Should(BeTrue())
		})

		var project *generated.ExtendedProject
		projectName := rnd.RandStringBytesMaskImprSrc(10)
		It("should create project", func() {
			p, err := client.ProjectCreate(generated.ProjectCreate{Name: projectName})
			Expect(err).NotTo(HaveOccurred())

			project = p
		})

		It("should assure project exists", func() {
			exists, err := client.ProjectExists(projectName)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("should get project", func() {
			p, err := client.Project(project.Id)
			Expect(err).NotTo(HaveOccurred())
			Expect(p.Name).To(Equal(projectName))
		})

		It("should delete project", func() {
			err := client.ProjectDelete(project.Id)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should assure project was deleted", func() {
			exists, err := client.ProjectExists(projectName)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

	})

})
