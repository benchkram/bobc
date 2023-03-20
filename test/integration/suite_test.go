package integration_test

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/benchkram/bob/bob/bobfile"
	"github.com/benchkram/bob/bob/global"
	nixbuilder "github.com/benchkram/bob/bob/nix-builder"
	"github.com/benchkram/bob/pkg/nix"
	"github.com/benchkram/bobc/restserver"
	"github.com/benchkram/bobc/test/environment"
	"github.com/benchkram/errz"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	suite *Suite

	dir           string
	bf            *bobfile.Bobfile
	artifactDir   string
	storageDir    string
	cleanup       func() error
	remoteProject string

	stdout *os.File
	stderr *os.File
	pr     *os.File
	pw     *os.File
)

// bobfileRemote as present in bob.yaml
var bobfileRemote = "localhost:8100"

var _ = BeforeSuite(func() {
	var err error
	suite, err = setup()
	Expect(err).NotTo(HaveOccurred())

	dir, storageDir, cleanup, err = suiteTestDirs("artifacts-upload")
	Expect(err).NotTo(HaveOccurred())
	artifactDir = filepath.Join(storageDir, global.BobCacheArtifactsDir)

	bf, err = bobfile.BobfileRead(".")
	Expect(err).NotTo(HaveOccurred())

	err = os.Chdir(dir)
	Expect(err).NotTo(HaveOccurred())

	remoteProject = strings.ReplaceAll(bf.Project, bobfileRemote, suite.restServer.Address())
})

var _ = AfterSuite(func() {
	if suite != nil {
		err := suite.Stop()
		Expect(err).NotTo(HaveOccurred())
	}

	err := cleanup()
	Expect(err).NotTo(HaveOccurred())

	for _, file := range tmpFiles {
		err := os.Remove(file)
		Expect(err).NotTo(HaveOccurred())
	}
})

func TestUpload(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artifact Upload Suite")
}

type Suite struct {
	restServer   *restserver.S
	postgres     *environment.Postgres
	minioStorage *environment.MinIO
	apiKey       []byte
}

func (s *Suite) Stop() error {
	if suite.restServer != nil {
		err := suite.restServer.Stop()
		Expect(err).NotTo(HaveOccurred())
	}

	if suite.postgres != nil {
		err := suite.postgres.Stop(true)
		Expect(err).NotTo(HaveOccurred())
	}

	if suite.minioStorage != nil {
		err := suite.minioStorage.Stop(true)
		Expect(err).NotTo(HaveOccurred())
	}

	return nil
}

func setup() (*Suite, error) {
	// setup database
	pg, err := environment.NewPostgresStarted()
	if err != nil {
		return nil, err
	}

	conf := pg.Config()

	minio, err := environment.NewMinioStarted()
	if err != nil {
		return nil, err
	}

	app, err := environment.InitApp(&conf, minio.Config())
	if err != nil {
		return nil, err
	}

	apiKey := []byte("debug-signing-key")

	s, err := environment.InitRestServer(app, apiKey)
	if err != nil {
		return nil, err
	}

	err = s.Start()
	Expect(err).NotTo(HaveOccurred())

	return &Suite{
		restServer:   s,
		postgres:     pg,
		minioStorage: minio,
		apiKey:       apiKey,
	}, nil
}

func projectName() string {
	return path.Base(bf.Project)
}

// TestDirs creates a general test dir and an "out-of-tree" storage dir used in tests.
// Call cleanup() to delete all dirs at the end of the test.
func suiteTestDirs(testName string) (testDir, storageDir string, cleanup func() error, _ error) {
	plain := func() error { return nil }

	testDir, err := ioutil.TempDir("", "bob-test-"+testName+"-*")
	if err != nil {
		return testDir, storageDir, plain, err
	}

	storageDir, err = ioutil.TempDir("", "bob-test-"+testName+"-storage-*")
	if err != nil {
		return testDir, storageDir, plain, err
	}

	err = os.Chmod(testDir, os.ModePerm)
	errz.Fatal(err)

	err = os.Chmod(storageDir, os.ModePerm)
	errz.Fatal(err)

	cleanup = func() (err error) {
		err = os.RemoveAll(testDir)
		if err != nil {
			return err
		}
		err = os.RemoveAll(storageDir)
		if err != nil {
			return err
		}
		return nil
	}
	return testDir, storageDir, cleanup, nil
}

// artifactExists checks if aN artifact exists in the local artifact store
func artifactExists(id string) (exist bool, _ error) {
	fs, err := os.ReadDir(artifactDir)
	if err != nil {
		return false, err
	}

	for _, f := range fs {
		if f.Name() == id {
			exist = true
			break
		}
	}

	return exist, nil
}

func capture() {
	stdout = os.Stdout
	stderr = os.Stderr

	var err error
	pr, pw, err = os.Pipe()
	Expect(err).NotTo(HaveOccurred())

	os.Stdout = pw
	os.Stderr = pw
}

func output() string {
	pw.Close()

	b, err := io.ReadAll(pr)
	Expect(err).NotTo(HaveOccurred())

	pr.Close()

	os.Stdout = stdout
	os.Stderr = stderr

	return string(b)
}

// tmpFiles tracks temporarily created files in these tests
// to be cleaned up at the end.
var tmpFiles []string

func NixBuilder() (*nixbuilder.NB, error) {
	file, err := ioutil.TempFile("", ".nix_cache*")
	if err != nil {
		return nil, err
	}
	name := file.Name()
	file.Close()

	tmpFiles = append(tmpFiles, name)
	cache, err := nix.NewCacheStore(nix.WithPath(name))
	if err != nil {
		return nil, err
	}

	return nixbuilder.New(nixbuilder.WithCache(cache)), nil
}
