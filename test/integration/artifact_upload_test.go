package integration_test

import (
	"context"
	"time"

	"github.com/benchkram/bob/bob"
	"github.com/benchkram/bob/bob/playbook"
	"github.com/benchkram/bob/pkg/boblog"
	restserverclient "github.com/benchkram/bobc/rest-server-client"
	"github.com/benchkram/bobc/restserver/generated"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Uploading artifacts", func() {
	var client *restserverclient.C

	When("bobfile `project` is a remote", func() {

		ctx := context.Background()

		It("should create rest-server client for this workflow", func() {
			var err error
			client, err = restserverclient.New(suite.restServer.Address(), suite.apiKey)
			Expect(err).NotTo(HaveOccurred())

			Eventually(client.Health, 60*time.Second, 1*time.Second).Should(BeTrue())
		})

		It("setup test with remote project", func() {
			bf.Project = remoteProject
			err := bf.BobfileSave(dir, "bob.yaml")
			Expect(err).NotTo(HaveOccurred())
		})

		var projectId uuid.UUID
		It("should create project", func() {
			p, err := client.ProjectCreate(generated.ProjectCreate{Name: projectName()})
			Expect(err).NotTo(HaveOccurred())

			projectId = uuid.MustParse(p.Id)
		})

		It("should assure project exists", func() {
			exists, err := client.ProjectExists(projectName())
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("should build without errors", func() {
			nixBuilder, err := NixBuilder()
			Expect(err).NotTo(HaveOccurred())
			b, err := bob.BobWithBaseStoreDir(
				storageDir,
				bob.WithDir(dir),
				bob.WithInsecure(true),
				bob.WithCachingEnabled(true),
				bob.WithNixBuilder(nixBuilder),
				bob.WithPushEnabled(true),
				bob.WithPullEnabled(true),
			)
			boblog.SetLogLevel(5)
			Expect(err).NotTo(HaveOccurred())

			err = b.CreateAuthContext("default", string(suite.apiKey))
			Expect(err).NotTo(HaveOccurred())

			capture()
			err = b.Build(ctx, "build")
			out := output()

			Expect(out).To(ContainSubstring("Artifact created."))
			Expect(err).NotTo(HaveOccurred())
		})

		var uploadedArtifactId string
		It("should have uploaded the artifacts on the project", func() {
			p, err := client.Project(projectId.String())
			Expect(err).NotTo(HaveOccurred())
			Expect(p.Hashes).NotTo(BeNil())
			Expect(*p.Hashes).To(HaveLen(1))

			artifact := (*p.Hashes)[0]
			uploadedArtifactId = artifact.Id
		})

		It("cleans the local cache and removes the artifact", func() {
			nixBuilder, err := NixBuilder()
			Expect(err).NotTo(HaveOccurred())
			cleaner, err := bob.BobWithBaseStoreDir(
				storageDir,
				bob.WithDir(dir),
				bob.WithNixBuilder(nixBuilder),
			)
			Expect(err).NotTo(HaveOccurred())
			err = cleaner.Clean()
			Expect(err).NotTo(HaveOccurred())

			exists, err := artifactExists(uploadedArtifactId)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("can run a new build and artifact will be downloaded from remote", func() {
			nixBuilder, err := NixBuilder()
			Expect(err).NotTo(HaveOccurred())
			newBuilder, err := bob.BobWithBaseStoreDir(
				storageDir,
				bob.WithDir(dir),
				bob.WithInsecure(true),
				bob.WithCachingEnabled(true),
				bob.WithNixBuilder(nixBuilder),
				bob.WithPushEnabled(true),
				bob.WithPullEnabled(true),
			)
			boblog.SetLogLevel(5)
			Expect(err).NotTo(HaveOccurred())
			capture()
			err = newBuilder.Build(ctx, "build")
			Expect(err).NotTo(HaveOccurred())

			exists, err := artifactExists(uploadedArtifactId)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			out := output()

			Expect(out).To(ContainSubstring("pulling artifact %s", uploadedArtifactId))
			// build command is marked as cached (since it had downloaded a cached artifact)
			state := playbook.StateNoRebuildRequired
			Expect(out).To(ContainSubstring(state.Short()))
		})
	})

	When("bobfile `project` is a not a remote", func() {
		It("replacing project with a regular name", func() {
			bf.Project = "bob-playground"
			err := bf.BobfileSave(dir, "bob.yaml")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should build without errors and not sync from local to remote", func() {
			nixBuilder, err := NixBuilder()
			Expect(err).NotTo(HaveOccurred())
			b, err := bob.BobWithBaseStoreDir(
				storageDir,
				bob.WithDir(dir),
				bob.WithInsecure(true),
				bob.WithCachingEnabled(true),
				bob.WithNixBuilder(nixBuilder),
				bob.WithPushEnabled(true),
				bob.WithPullEnabled(true),
			)
			Expect(err).NotTo(HaveOccurred())
			boblog.SetLogLevel(5)

			capture()
			err = b.Build(context.Background(), "build")
			Expect(err).NotTo(HaveOccurred())
			Expect(output()).To(Not(ContainSubstring("failed to sync from remote to local")))
		})
	})
})
