package weblogic_buildpack

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/postal"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"os"
	"path/filepath"
)

func Build(logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		auxImageLayer, err := context.Layers.Get("auxiliary-image")
		if err != nil {
			return packit.BuildResult{}, err
		}
		auxImageLayer.Launch = true
		auxImageLayer.Build = true

		logger.Process("Prepare Application archive")
		archivePath, err := PrepareApplicationArchive(context.WorkingDir)
		logger.Action("Application archive placed under %s", archivePath)
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Action("Completed")
		logger.Break()

		logger.Process("Cleanup Working dir")
		err = CleanupDir(context.WorkingDir)
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Action("Completed")
		logger.Break()

		logger.Process("Install the WebLogic Deploy Tool")
		weblogicDeployToolLayer, err := InstallWeblogicDeployTool(context)
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Action("Completed")
		logger.Break()

		logger.Process("Provide the WebLogic domain models metadata")
		modelPath := filepath.Join(context.CNBPath, "app-model.yaml")
		modelsDir, err := ProvideModelMetadata(modelPath, context.WorkingDir)
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Action("Completed")
		logger.Break()

		logger.Process("Provide the application archive")
		err = ProvideApplicationArchive(archivePath, modelsDir)
		if err != nil {
			return packit.BuildResult{}, err
		}
		logger.Action("Completed")
		logger.Break()

		auxImageLayer.SharedEnv.Default("AUXILIARY_IMAGE_PATH", context.WorkingDir)

		return packit.BuildResult{
			Layers: []packit.Layer{auxImageLayer, weblogicDeployToolLayer},
			Build:  packit.BuildMetadata{},
			Launch: packit.LaunchMetadata{},
		}, nil
	}
}

func PrepareApplicationArchive(sourceDir string) (string, error) {
	tempDir, err := os.MkdirTemp("", "application")
	wlsDeployDir := filepath.Join(tempDir, "wlsdeploy")
	appDir := filepath.Join(wlsDeployDir, "applications", "app")
	archivePath := filepath.Join(tempDir, "app.zip")

	err = os.MkdirAll(appDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	err = fs.Copy(sourceDir, appDir)
	if err != nil {
		return "", err
	}

	err = RecursiveZip(wlsDeployDir, archivePath)
	if err != nil {
		return "", err
	}

	err = os.RemoveAll(wlsDeployDir)
	if err != nil {
		return "", err
	}

	return archivePath, nil
}

func InstallWeblogicDeployTool(context packit.BuildContext) (packit.Layer, error) {
	layer, err := context.Layers.Get("weblogic-deploy-ool")
	if err != nil {
		return packit.Layer{}, err
	}
	layer.Launch = false
	layer.Build = false

	_, err = layer.Reset()
	if err != nil {
		return packit.Layer{}, err
	}

	toolName := "weblogic-deploy"

	transport := cargo.NewTransport()
	service := postal.NewService(transport)

	buildpackDescPath := filepath.Join(context.CNBPath, "buildpack.toml")
	installationDir := filepath.Join(layer.Path, toolName)
	installationSymlink := filepath.Join(context.WorkingDir, toolName)

	deployTool, err := service.Resolve(buildpackDescPath, toolName, "2.1.0", context.Stack)
	if err != nil {
		return packit.Layer{}, err
	}

	err = service.Deliver(deployTool, context.CNBPath, layer.Path, context.Platform.Path)
	if err != nil {
		return packit.Layer{}, err
	}

	if err != nil {
		return packit.Layer{}, err
	}

	err = CopyContent(installationDir, installationSymlink)
	if err != nil {
		return packit.Layer{}, err
	}

	return layer, nil
}

func ProvideModelMetadata(modelPath string, destinationDir string) (string, error) {
	modelsDir := filepath.Join(destinationDir, "models")

	err := os.MkdirAll(modelsDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	err = CopyContent(modelPath, filepath.Join(modelsDir, filepath.Base(modelPath)))
	if err != nil {
		return "", err
	}

	return modelsDir, nil
}

func ProvideApplicationArchive(archivePath string, destinationDir string) error {
	return CopyContent(archivePath, filepath.Join(destinationDir, filepath.Base(archivePath)))
}
