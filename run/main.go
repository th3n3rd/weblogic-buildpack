package main

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/th3n3rd/weblogic-buildpack"
	"os"
)

func main() {
	emitter := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))
	packit.Run(
		weblogic_buildpack.Detect(),
		weblogic_buildpack.Build(emitter),
	)
}
