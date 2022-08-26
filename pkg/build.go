package docatl

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	util "github.com/docat-org/docatl/internal"
	"github.com/mholt/archiver/v3"
	"gopkg.in/yaml.v2"
)

const metadataFileName = ".docatl.meta.yaml"

type BuildMetadata struct {
	Host    string `yaml:"host,omitempty"`
	Project string `yaml:"project,omitempty"`
	Version string `yaml:"version,omitempty"`
}

func Build(docsPath string, meta BuildMetadata) (string, error) {
	if !util.IsDirectory(docsPath) {
		return "", fmt.Errorf("the given documentation path must be a directory")
	}

	docsPath = util.ResolvePath(docsPath)

	// NOTE(TF): the `archiver` package does not have an option to not create a top-level directory
	filesInDocsPath, err := os.ReadDir(docsPath)
	if err != nil {
		return "", fmt.Errorf("cannot list the contents within the given documentation directory: %w", err)
	}
	filesToArchive := make([]string, 0)
	for _, f := range filesInDocsPath {
		filesToArchive = append(filesToArchive, filepath.Join(docsPath, f.Name()))
	}

	if meta.Project != "" && meta.Version != "" {
		metadataFile, err := generateMetadataFile(meta)
		if err != nil {
			log.Fatal(err)
		}
		filesToArchive = append(filesToArchive, metadataFile)
	}

	outputPath := generateArtifactFileName(docsPath, meta)

	z := archiver.Zip{OverwriteExisting: true,  FileMethod: archiver.BZIP2}
	err = z.Archive(filesToArchive, outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to archive docs: %w", err)
	}

	return outputPath, nil
}

func generateArtifactFileName(docsPath string, meta BuildMetadata) string {
	if meta.Project == "" && meta.Version == "" {
		return fmt.Sprintf("%s.zip", filepath.Base(docsPath))
	}

	if meta.Project != "" && meta.Version == "" {
		return fmt.Sprintf("docs_%s.zip", meta.Project)
	}

	return fmt.Sprintf("docs_%s_%s.zip", meta.Project, meta.Version)
}

func generateMetadataFile(meta BuildMetadata) (string, error) {
	tmpDir, err := os.MkdirTemp("", "docatl-*")
	if err != nil {
		return "", fmt.Errorf("unable to create temp directory for metadatafile: %w", err)
	}
	metadataFile := filepath.Join(tmpDir, metadataFileName)

	doc, err := yaml.Marshal(&meta)
	if err != nil {
		return "", fmt.Errorf("unable to generate metadata file for data: %v: %w", meta, err)
	}

	err = os.WriteFile(metadataFile, doc, 0755)
	if err != nil {
		return "", fmt.Errorf("unabel to write metadata to file %s: %w", metadataFile, err)
	}

	return metadataFile, nil
}

func ExtractMetadata(docsPath string) (BuildMetadata, error) {
	var meta BuildMetadata
	err := archiver.Walk(docsPath, func(f archiver.File) error {
		if zfh, ok := f.Header.(zip.FileHeader); ok {
			if zfh.Name == metadataFileName {
				contents := make([]byte, f.Size())
				_, err := io.ReadFull(f, contents)
				if err != nil {
					return fmt.Errorf("unable to extract metadata: %w", err)
				}
				err = yaml.Unmarshal(contents, &meta)
				if err != nil {
					return fmt.Errorf("unable to read metadata file contents as YAML: %v: %w", contents, err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return BuildMetadata{}, err
	}

	return meta, nil
}
