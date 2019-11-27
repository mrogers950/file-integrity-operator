// +build !ignore_autogenerated

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/mrogers950/file-integrity-operator/pkg/apis/fileintegrity/v1alpha1.FileIntegrity":       schema_pkg_apis_fileintegrity_v1alpha1_FileIntegrity(ref),
		"github.com/mrogers950/file-integrity-operator/pkg/apis/fileintegrity/v1alpha1.FileIntegrityConfig": schema_pkg_apis_fileintegrity_v1alpha1_FileIntegrityConfig(ref),
		"github.com/mrogers950/file-integrity-operator/pkg/apis/fileintegrity/v1alpha1.FileIntegritySpec":   schema_pkg_apis_fileintegrity_v1alpha1_FileIntegritySpec(ref),
		"github.com/mrogers950/file-integrity-operator/pkg/apis/fileintegrity/v1alpha1.FileIntegrityStatus": schema_pkg_apis_fileintegrity_v1alpha1_FileIntegrityStatus(ref),
	}
}

func schema_pkg_apis_fileintegrity_v1alpha1_FileIntegrity(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "FileIntegrity is the Schema for the fileintegrities API",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/mrogers950/file-integrity-operator/pkg/apis/fileintegrity/v1alpha1.FileIntegritySpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/mrogers950/file-integrity-operator/pkg/apis/fileintegrity/v1alpha1.FileIntegrityStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/mrogers950/file-integrity-operator/pkg/apis/fileintegrity/v1alpha1.FileIntegritySpec", "github.com/mrogers950/file-integrity-operator/pkg/apis/fileintegrity/v1alpha1.FileIntegrityStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_fileintegrity_v1alpha1_FileIntegrityConfig(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "FileIntegrityConfig defines the name, namespace, and data key for an AIDE config to use for integrity checking.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"name": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"namespace": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"key": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_fileintegrity_v1alpha1_FileIntegritySpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "FileIntegritySpec defines the desired state of FileIntegrity",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"config": {
						SchemaProps: spec.SchemaProps{
							Description: "INSERT ADDITIONAL SPEC FIELDS - desired state of cluster Important: Run \"operator-sdk generate k8s\" to regenerate code after modifying this file Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html",
							Ref:         ref("github.com/mrogers950/file-integrity-operator/pkg/apis/fileintegrity/v1alpha1.FileIntegrityConfig"),
						},
					},
				},
				Required: []string{"config"},
			},
		},
		Dependencies: []string{
			"github.com/mrogers950/file-integrity-operator/pkg/apis/fileintegrity/v1alpha1.FileIntegrityConfig"},
	}
}

func schema_pkg_apis_fileintegrity_v1alpha1_FileIntegrityStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "FileIntegrityStatus defines the observed state of FileIntegrity",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"phase": {
						SchemaProps: spec.SchemaProps{
							Description: "INSERT ADDITIONAL STATUS FIELD - define observed state of cluster Important: Run \"operator-sdk generate k8s\" to regenerate code after modifying this file Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
	}
}
