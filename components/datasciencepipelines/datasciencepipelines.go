// Package datasciencepipelines provides utility functions to config Data Science Pipelines: Pipeline solution for end to end MLOps workflows that support the Kubeflow Pipelines SDK and Tekton
package datasciencepipelines

import (
	dsci "github.com/opendatahub-io/opendatahub-operator/v2/apis/dscinitialization/v1"
	"github.com/opendatahub-io/opendatahub-operator/v2/components"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/deploy"
	operatorv1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ComponentName = "data-science-pipelines-operator"
	Path          = deploy.DefaultManifestPath + "/" + ComponentName + "/base"
)

var imageParamMap = map[string]string{
	"IMAGES_APISERVER":         "RELATED_IMAGE_ODH_ML_PIPELINES_API_SERVER_IMAGE",
	"IMAGES_ARTIFACT":          "RELATED_IMAGE_ODH_ML_PIPELINES_ARTIFACT_MANAGER_IMAGE",
	"IMAGES_PERSISTENTAGENT":   "RELATED_IMAGE_ODH_ML_PIPELINES_PERSISTENCEAGENT_IMAGE",
	"IMAGES_SCHEDULEDWORKFLOW": "RELATED_IMAGE_ODH_ML_PIPELINES_SCHEDULEDWORKFLOW_IMAGE",
	"IMAGES_CACHE":             "RELATED_IMAGE_ODH_ML_PIPELINES_CACHE_IMAGE",
	"IMAGES_DSPO":              "RELATED_IMAGE_ODH_DATA_SCIENCE_PIPELINES_OPERATOR_CONTROLLER_IMAGE",
}

type DataSciencePipelines struct {
	components.Component `json:""`
}

func (d *DataSciencePipelines) SetImageParamsMap(imageMap map[string]string) map[string]string {
	imageParamMap = imageMap
	return imageParamMap
}

func (d *DataSciencePipelines) GetComponentName() string {
	return ComponentName
}

// Verifies that Dashboard implements ComponentInterface
var _ components.ComponentInterface = (*DataSciencePipelines)(nil)

func (d *DataSciencePipelines) ReconcileComponent(cli client.Client, owner metav1.Object, dscispec *dsci.DSCInitializationSpec) error {
	enabled := d.GetManagementState() == operatorv1.Managed

	if enabled {
		// check if the dependent operator installed is done in dashboard

		// Update image parameters only when we do not have customized manifests set
		if dscispec.DevFlags.ManifestsUri == "" {
			if err := deploy.ApplyImageParams(Path, imageParamMap); err != nil {
				return err
			}
		}
	}

	err := deploy.DeployManifestsFromPath(cli, owner, Path, dscispec.ApplicationsNamespace, d.GetComponentName(), enabled)
	return err
}

func (d *DataSciencePipelines) DeepCopyInto(target *DataSciencePipelines) {
	*target = *d
	target.Component = d.Component
}
