package integration_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/benchkram/bobc/pkg/rnd"
	restserverclient "github.com/benchkram/bobc/rest-server-client"
	"github.com/benchkram/bobc/restserver/generated"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("should test the artifact lifecycle", func() {

	var client *restserverclient.C

	Context("check if tests are running", func() {

		It("should create rest-server client for this workflow", func() {
			var err error
			client, err = restserverclient.New(suite.restServer.Address(), suite.apiKey)
			Expect(err).NotTo(HaveOccurred())

			Eventually(client.Health, 60*time.Second, 1*time.Second).Should(BeTrue())
		})

		var project *generated.ExtendedProject
		projectName := rnd.RandStringBytesMaskImprSrc(10)
		It("should create project", func() {
			p, err := client.ProjectCreate(generated.ProjectCreate{Name: projectName})
			Expect(err).NotTo(HaveOccurred())

			project = p
		})

		sha1Hash := rnd.RandSHA1(8)
		It("should add a artifact to project", func() {
			// create a random file
			filepath := filepath.Join(os.TempDir(), "file.test")
			bigBuff := make([]byte, 750)
			err := ioutil.WriteFile(filepath, bigBuff, 0666)
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(filepath)

			err = client.ArtifactCreate(project.Name, sha1Hash, filepath)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should assure artifact exists", func() {
			exists, err := client.ArtifactExists(project.Name, sha1Hash)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("should get valid artifact object", func() {
			artifact, err := client.Artifact(project.Name, sha1Hash)
			Expect(err).NotTo(HaveOccurred())
			Expect(artifact.Location).NotTo(BeNil())
			Expect(*artifact.Location).NotTo(BeEmpty())
		})

		It("should get valid artifact object", func() {
			err := client.ArtifactDelete(project.Name, sha1Hash)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should assure artifact does not exist", func() {
			exists, err := client.ArtifactExists(project.Name, sha1Hash)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
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
