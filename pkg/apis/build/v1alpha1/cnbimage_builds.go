package v1alpha1

import (
	"github.com/knative/pkg/kmeta"
	"github.com/pborman/uuid"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (im *CNBImage) BuildNeeded(lastBuild *CNBBuild, builder *CNBBuilder) bool {
	if lastBuild == nil {
		return true
	}

	if im.configMatches(lastBuild) && builtWithBuilderBuildpacks(builder, lastBuild) {
		return false
	}

	return true
}

func builtWithBuilderBuildpacks(builder *CNBBuilder, build *CNBBuild) bool {
	for _, bp := range build.Status.BuildMetadata {
		if !builder.Spec.BuilderMetadata.Include(bp) {
			return false
		}
	}

	return true
}

func (im *CNBImage) configMatches(build *CNBBuild) bool {
	return im.Spec.Image == build.Spec.Image &&
		im.Spec.GitURL == build.Spec.GitURL &&
		im.Spec.GitRevision == build.Spec.GitRevision
}

func (im *CNBImage) CreateBuild(builder *CNBBuilder) *CNBBuild {
	return &CNBBuild{
		ObjectMeta: v1.ObjectMeta{
			Name: im.Name + "-build-" + uuid.New(),
			OwnerReferences: []v1.OwnerReference{
				*kmeta.NewControllerRef(im),
			},
		},
		Spec: CNBBuildSpec{
			Image:          im.Spec.Image,
			Builder:        builder.Spec.Image,
			ServiceAccount: im.Spec.ServiceAccount,
			GitURL:         im.Spec.GitURL,
			GitRevision:    im.Spec.GitRevision,
		},
	}
}
