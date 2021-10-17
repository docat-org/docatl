package docatl

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	util "github.com/docat-org/docatl/internal"
	"github.com/mholt/archiver/v3"
)

type BuildMetadata struct {
	Project string
	Version string
}

func Build(docsPath string, meta BuildMetadata) (string, error) {
	if !util.IsDirectory(docsPath) {
		return "", fmt.Errorf("the given documentation path must be a directory")
	}

	docsPath = util.ResolvePath(docsPath)

	// NOTE(TF): the `archiver` package does not have an option to not create a top-level directory
	filesInDocsPath, err := ioutil.ReadDir(docsPath)
	if err != nil {
		return "", fmt.Errorf("cannot list the contents within the given documentation directory: %w", err)
	}
	filesToArchive := make([]string, 0)
	for _, f := range filesInDocsPath {
		filesToArchive = append(filesToArchive, filepath.Join(docsPath, f.Name()))
	}

	outputPath := generateArtifactFileName(docsPath, meta)

	z := archiver.Zip{OverwriteExisting: true}
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
