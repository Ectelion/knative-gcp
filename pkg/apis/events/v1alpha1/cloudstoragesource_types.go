/*
Copyright 2019 Google LLC.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	duckv1alpha1 "github.com/google/knative-gcp/pkg/apis/duck/v1alpha1"
	kngcpduck "github.com/google/knative-gcp/pkg/duck/v1alpha1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/webhook/resourcesemantics"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudStorageSource is a specification for a Google Cloud CloudStorageSource Source resource.
type CloudStorageSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CloudStorageSourceSpec   `json:"spec"`
	Status CloudStorageSourceStatus `json:"status"`
}

// Verify that CloudStorageSource matches various duck types.
var (
	_ apis.Convertible             = (*CloudStorageSource)(nil)
	_ apis.Defaultable             = (*CloudStorageSource)(nil)
	_ apis.Validatable             = (*CloudStorageSource)(nil)
	_ runtime.Object               = (*CloudStorageSource)(nil)
	_ kmeta.OwnerRefable           = (*CloudStorageSource)(nil)
	_ resourcesemantics.GenericCRD = (*CloudStorageSource)(nil)
	_ kngcpduck.Identifiable       = (*CloudStorageSource)(nil)
	_ kngcpduck.PubSubable         = (*CloudStorageSource)(nil)
)

// CloudStorageSourceSpec is the spec for a CloudStorageSource resource.
type CloudStorageSourceSpec struct {
	// This brings in the PubSub based Source Specs. Includes:
	// Sink, CloudEventOverrides, Secret and Project
	duckv1alpha1.PubSubSpec `json:",inline"`

	// Bucket to subscribe to.
	Bucket string `json:"bucket"`

	// EventTypes to subscribe to. If unspecified, then subscribe to all events.
	// +optional
	EventTypes []string `json:"eventTypes,omitempty"`

	// ObjectNamePrefix limits the notifications to objects with this prefix
	// +optional
	ObjectNamePrefix string `json:"objectNamePrefix,omitempty"`

	// PayloadFormat specifies the contents of the message payload.
	// See https://cloud.google.com/storage/docs/pubsub-notifications#payload.
	// +optional
	PayloadFormat string `json:"payloadFormat,omitempty"`
}

const (
	// CloudStorageSourceConditionReady has status True when the CloudStorageSource is ready to send events.
	CloudStorageSourceConditionReady = apis.ConditionReady

	// NotificationReady has status True when GCS has been configured properly to
	// send Notification events.
	NotificationReady apis.ConditionType = "NotificationReady"
)

var storageCondSet = apis.NewLivingConditionSet(
	duckv1alpha1.PullSubscriptionReady,
	duckv1alpha1.TopicReady,
	NotificationReady)

// CloudStorageSourceStatus is the status for a GCS resource.
type CloudStorageSourceStatus struct {
	// This brings in the Status for our GCP PubSub event sources.
	// duck/v1beta1 Status, SinkURI, ProjectID, TopicID and SubscriptionID
	duckv1alpha1.PubSubStatus `json:",inline"`

	// NotificationID is the ID that GCS identifies this notification as.
	// +optional
	NotificationID string `json:"notificationId,omitempty"`
}

func (storage *CloudStorageSource) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("CloudStorageSource")
}

// Methods for identifiable interface.
// IdentitySpec returns the IdentitySpec portion of the Spec.
func (s *CloudStorageSource) IdentitySpec() *duckv1alpha1.IdentitySpec {
	return &s.Spec.IdentitySpec
}

// IdentityStatus returns the IdentityStatus portion of the Status.
func (s *CloudStorageSource) IdentityStatus() *duckv1alpha1.IdentityStatus {
	return &s.Status.IdentityStatus
}

// ConditionSet returns the apis.ConditionSet of the embedding object.
func (s *CloudStorageSource) ConditionSet() *apis.ConditionSet {
	return &storageCondSet
}

// Methods for pubsubable interface.

// PubSubSpec returns the PubSubSpec portion of the Spec.
func (s *CloudStorageSource) PubSubSpec() *duckv1alpha1.PubSubSpec {
	return &s.Spec.PubSubSpec
}

// PubSubStatus returns the PubSubStatus portion of the Status.
func (s *CloudStorageSource) PubSubStatus() *duckv1alpha1.PubSubStatus {
	return &s.Status.PubSubStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudStorageSourceList is a list of CloudStorageSource resources.
type CloudStorageSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CloudStorageSource `json:"items"`
}
