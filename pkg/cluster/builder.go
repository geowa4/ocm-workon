package cluster

import (
	_ "embed"
	"fmt"
	"github.com/geowa4/ocm-workon/pkg/config"
	"github.com/geowa4/ocm-workon/pkg/utils"
	"os"
	"text/template"
	"time"
)

const ProductionEnvironment = "prod"
const StagingEnvironment = "stage"
const ProductionOcmUrl = "https://api.openshift.com"
const StagingOcmUrl = "https://api.stage.openshift.com"
const ProductionHCPNamespacePrefix = "ocm-production-"
const StagingHCPNamespacePrefix = "ocm-staging-"

//go:embed envrc
var envrcContent string

//go:embed dotenv
var dotenvContent string

type WorkConfig struct {
	ClusterBase         string
	Environment         string
	OcmUrl              string
	HCPNamespacePrefix  string
	UseDirenv           bool
	OcmConfigFile       string
	BackplaneConfigFile string
	ClusterData         *NormalizedClusterData
}

func (w *WorkConfig) Validate() error {
	if w.Environment != StagingEnvironment && w.Environment != ProductionEnvironment {
		return fmt.Errorf("environment must be one of %s or %s", ProductionEnvironment, StagingEnvironment)
	}
	fileInfo, err := os.Lstat(w.ClusterBase)
	if err != nil || !fileInfo.IsDir() {
		return fmt.Errorf("the directory '%s' does not exist or is not a directory: %q", w.ClusterBase, err)
	}
	return nil
}

func (w *WorkConfig) setEnvDependentFields() {
	if w.Environment == ProductionEnvironment {
		w.OcmUrl = ProductionOcmUrl
		w.HCPNamespacePrefix = ProductionHCPNamespacePrefix
	} else {
		w.OcmUrl = StagingOcmUrl
		w.HCPNamespacePrefix = StagingHCPNamespacePrefix
	}
	w.OcmConfigFile = config.GetOcmConfigFile(w.Environment)
	w.BackplaneConfigFile = config.GetBackplaneConfigFile(w.Environment)
}

func (w *WorkConfig) getClusterDirName() string {
	return w.ClusterBase + utils.PathSep + w.Environment + utils.PathSep + w.ClusterData.Name
}

func (w *WorkConfig) Build() (string, error) {
	if err := w.Validate(); err != nil {
		return "", err
	}

	w.setEnvDependentFields()

	clusterDir := w.getClusterDirName()
	if err := os.MkdirAll(clusterDir, 0744); err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("could not create cluster directory in %s: %q", w.ClusterBase, err)
	}

	if err := makeNotesFile(clusterDir); err != nil {
		return "", err
	}

	if err := makeKubeconfig(clusterDir); err != nil {
		return "", err
	}

	if err := w.makeDotenv(clusterDir); err != nil {
		return "", err
	}

	if w.UseDirenv {
		if err := w.makeEnvrc(clusterDir); err != nil {
			return "", err
		}
	}

	return clusterDir, nil
}

func (w *WorkConfig) makeDotenv(clusterDir string) error {
	dotenvFilePath := clusterDir + utils.PathSep + ".env.cluster"
	if _, err := os.Lstat(dotenvFilePath); os.IsNotExist(err) {
		dotenvFile, err := os.OpenFile(dotenvFilePath, os.O_CREATE, 0744)
		if err != nil {
			return err
		}
		_ = dotenvFile.Close()
	} else if err != nil {
		return fmt.Errorf("error creating dotenv file in %s: %q", clusterDir, err)
	}
	dotenvFile, err := os.OpenFile(dotenvFilePath, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening dotenv file in %s: %q", clusterDir, err)
	}
	defer utils.CloseFileAndIgnoreErrors(dotenvFile)
	dotenvTemplate, err := template.New("dotenv").Parse(dotenvContent)
	if err != nil {
		return fmt.Errorf("error generating content for dotenv file in %s: %q", clusterDir, err)
	}
	if err = dotenvTemplate.Execute(dotenvFile, w); err != nil {
		return fmt.Errorf("error writing to dotenv file in %s: %q", clusterDir, err)
	}
	return nil
}

func (w *WorkConfig) makeEnvrc(clusterDir string) error {
	envrcFilePath := clusterDir + utils.PathSep + ".envrc"
	if fileInfo, err := os.Lstat(envrcFilePath); os.IsNotExist(err) || fileInfo.Size() == 0 {
		envrcFile, err := os.OpenFile(envrcFilePath, os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		_ = envrcFile.Close()
	} else if err != nil {
		return fmt.Errorf("error creating envrc file in %s: %q", clusterDir, err)
	}
	envrcFile, err := os.OpenFile(envrcFilePath, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening envrc file in %s: %q", clusterDir, err)
	}
	defer utils.CloseFileAndIgnoreErrors(envrcFile)

	envrcTemplate, err := template.New("envrc").Parse(envrcContent)
	if err != nil {
		return fmt.Errorf("error generating content for envrc file in %s: %q", clusterDir, err)
	}
	if err = envrcTemplate.Execute(envrcFile, w); err != nil {
		return fmt.Errorf("error writing to envrc file in %s: %q", clusterDir, err)
	}
	return nil
}

func makeKubeconfig(clusterDir string) error {
	kubeFilePath := clusterDir + utils.PathSep + "kubeconfig"
	if _, err := os.Lstat(kubeFilePath); os.IsNotExist(err) {
		kubeFile, err := os.OpenFile(kubeFilePath, os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		_ = kubeFile.Close()
	}
	return nil
}

func makeNotesFile(clusterDir string) error {
	desiredFile := "notes"
	notesFilePath := clusterDir + utils.PathSep + desiredFile
	notesIsNew := false
	if _, err := os.Lstat(notesFilePath); os.IsNotExist(err) {
		notesFile, err := os.OpenFile(notesFilePath, os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		_ = notesFile.Close()
		notesIsNew = true
	} else if err != nil {
		return fmt.Errorf("error creating %s file in %s: %q", desiredFile, clusterDir, err)
	}

	notesFile, err := os.OpenFile(notesFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening %s file in %s: %q", desiredFile, clusterDir, err)
	}
	defer utils.CloseFileAndIgnoreErrors(notesFile)

	notesPrefix := ""
	if !notesIsNew {
		notesPrefix = "\n\n"
	}
	_, err = notesFile.WriteString(fmt.Sprintf("%s---\n%s\n\n", notesPrefix, time.Now().Format("2006-01-02 15:04:05 Monday")))
	if err != nil {
		return fmt.Errorf("error writing to %s file in %s: %q", desiredFile, clusterDir, err)
	}
	return nil
}
