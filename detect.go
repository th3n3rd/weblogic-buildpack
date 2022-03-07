package weblogic_buildpack

import "github.com/paketo-buildpacks/packit/v2"

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		_, err := IsJvmApplicationPackage(context.WorkingDir)
		if err != nil {
			return packit.DetectResult{}, err
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: "weblogic-model-in-image"},
				},
				Requires: []packit.BuildPlanRequirement{
					{Name: "weblogic-model-in-image"},
					{Name: "jvm-application-package"},
				},
			},
		}, nil
	}
}
