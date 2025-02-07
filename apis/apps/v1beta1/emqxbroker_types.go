/*
Copyright 2021.

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

package v1beta1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EmqxBrokerSpec defines the desired state of EmqxBroker
type EmqxBrokerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The fields of Broker.
	//The replicas of emqx broker
	//+kubebuilder:validation:Minimum=3
	Replicas *int32 `json:"replicas,omitempty"`

	//+kubebuilder:validation:Required
	Image string `json:"image,omitempty"`

	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// The service account name which is being bind with the service
	// account of the crd instance.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	Storage *Storage `json:"storage,omitempty"`

	// TODO: waiting to be deleted, should use meta.labels
	Labels Labels `json:"labels,omitempty"`
	// TODO: waiting to be deleted, should use meta.annotations
	Annotations map[string]string `json:"annotations,omitempty"`

	Listener Listener `json:"listener,omitempty"`

	Affinity        *corev1.Affinity    `json:"affinity,omitempty"`
	ToleRations     []corev1.Toleration `json:"toleRations,omitempty"`
	NodeSelector    map[string]string   `json:"nodeSelector,omitempty"`
	ImagePullPolicy corev1.PullPolicy   `json:"imagePullPolicy,omitempty"`

	ExtraVolumes      []corev1.Volume      `json:"extraVolumes,omitempty"`
	ExtraVolumeMounts []corev1.VolumeMount `json:"extraVolumeMounts,omitempty"`

	Env []corev1.EnvVar `json:"env,omitempty"`

	ACL []ACL `json:"acl,omitempty"`

	Plugins []Plugin `json:"plugins,omitempty"`

	Modules []EmqxBrokerModules `json:"modules,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=emqx
//+kubebuilder:subresource:status
//+kubebuilder:storageversion
//+kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas
// EmqxBroker is the Schema for the emqxbrokers API
type EmqxBroker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmqxBrokerSpec `json:"spec,omitempty"`
	Status `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// EmqxBrokerList contains a list of EmqxBroker
type EmqxBrokerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EmqxBroker `json:"items"`
}

func (emqx *EmqxBroker) String() string {
	return fmt.Sprintf("EmqxBroker instance [%s],Image [%s]",
		emqx.ObjectMeta.Name,
		emqx.Spec.Image,
	)
}

func init() {
	SchemeBuilder.Register(&EmqxBroker{}, &EmqxBrokerList{})
}

func (emqx *EmqxBroker) GetAPIVersion() string        { return emqx.APIVersion }
func (emqx *EmqxBroker) SetAPIVersion(version string) { emqx.APIVersion = version }
func (emqx *EmqxBroker) GetKind() string              { return emqx.Kind }
func (emqx *EmqxBroker) SetKind(kind string)          { emqx.Kind = kind }

func (emqx *EmqxBroker) GetReplicas() *int32 {
	return emqx.Spec.Replicas
}
func (emqx *EmqxBroker) SetReplicas(replicas *int32) { emqx.Spec.Replicas = replicas }

func (emqx *EmqxBroker) GetImage() string      { return emqx.Spec.Image }
func (emqx *EmqxBroker) SetImage(image string) { emqx.Spec.Image = image }

func (emqx *EmqxBroker) GetServiceAccountName() string {
	return emqx.Spec.ServiceAccountName
}
func (emqx *EmqxBroker) SetServiceAccountName(serviceAccountName string) {
	emqx.Spec.ServiceAccountName = serviceAccountName
}

func (emqx *EmqxBroker) GetResource() corev1.ResourceRequirements { return emqx.Spec.Resources }
func (emqx *EmqxBroker) SetResource(resource corev1.ResourceRequirements) {
	emqx.Spec.Resources = resource
}

func (emqx *EmqxBroker) GetStorage() *Storage        { return emqx.Spec.Storage }
func (emqx *EmqxBroker) SetStorage(storage *Storage) { emqx.Spec.Storage = storage }

func (emqx *EmqxBroker) GetAnnotations() map[string]string { return emqx.Spec.Annotations }
func (emqx *EmqxBroker) SetAnnotations(annotations map[string]string) {
	emqx.Spec.Annotations = annotations
}

func (emqx *EmqxBroker) GetListener() Listener { return emqx.Spec.Listener }
func (emqx *EmqxBroker) SetListener(listener Listener) {
	emqx.Spec.Listener = listener
}

func (emqx *EmqxBroker) GetAffinity() *corev1.Affinity         { return emqx.Spec.Affinity }
func (emqx *EmqxBroker) SetAffinity(affinity *corev1.Affinity) { emqx.Spec.Affinity = affinity }

func (emqx *EmqxBroker) GetToleRations() []corev1.Toleration { return emqx.Spec.ToleRations }
func (emqx *EmqxBroker) SetToleRations(tolerations []corev1.Toleration) {
	emqx.Spec.ToleRations = tolerations
}

func (emqx *EmqxBroker) GetNodeSelector() map[string]string { return emqx.Spec.NodeSelector }
func (emqx *EmqxBroker) SetNodeSelector(nodeSelector map[string]string) {
	emqx.Spec.NodeSelector = nodeSelector
}

func (emqx *EmqxBroker) GetImagePullPolicy() corev1.PullPolicy { return emqx.Spec.ImagePullPolicy }
func (emqx *EmqxBroker) SetImagePullPolicy(pullPolicy corev1.PullPolicy) {
	emqx.Spec.ImagePullPolicy = pullPolicy
}

func (emqx *EmqxBroker) GetExtraVolumes() []corev1.Volume { return emqx.Spec.ExtraVolumes }
func (emqx *EmqxBroker) GetExtraVolumeMounts() []corev1.VolumeMount {
	return emqx.Spec.ExtraVolumeMounts
}

func (emqx *EmqxBroker) GetACL() []ACL { return emqx.Spec.ACL }
func (emqx *EmqxBroker) SetACL(acl []ACL) {
	emqx.Spec.ACL = acl
}

func (emqx *EmqxBroker) GetEnv() []corev1.EnvVar { return emqx.Spec.Env }
func (emqx *EmqxBroker) SetEnv(env []corev1.EnvVar) {
	emqx.Spec.Env = env
}

func (emqx *EmqxBroker) GetPlugins() []Plugin { return emqx.Spec.Plugins }
func (emqx *EmqxBroker) SetPlugins(plugins []Plugin) {
	emqx.Spec.Plugins = plugins
}

func (emqx *EmqxBroker) GetModules() []EmqxBrokerModules { return emqx.Spec.Modules }
func (emqx *EmqxBroker) SetModules(modules []EmqxBrokerModules) {
	emqx.Spec.Modules = modules
}

func (emqx *EmqxBroker) GetHeadlessServiceName() string {
	return fmt.Sprintf("%s-%s", emqx.Name, "headless")
}
