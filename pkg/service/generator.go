package service

import (
	"fmt"
	"reflect"

	"github.com/emqx/emqx-operator/apis/apps/v1beta1"
	"github.com/emqx/emqx-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewSecretForCR(emqx v1beta1.EmqxEnterprise, labels map[string]string, ownerRefs []metav1.OwnerReference) *corev1.Secret {
	stringData := map[string]string{"emqx.lic": emqx.GetLicense()}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Labels:          labels,
			Name:            util.GetLicense(&emqx)["name"],
			Namespace:       emqx.GetNamespace(),
			OwnerReferences: ownerRefs,
		},
		Type:       corev1.SecretTypeOpaque,
		StringData: stringData,
	}
}

func NewRBAC(emqx v1beta1.Emqx, labels map[string]string, ownerRefs []metav1.OwnerReference) (*corev1.ServiceAccount, *rbacv1.Role, *rbacv1.RoleBinding) {
	meta := metav1.ObjectMeta{
		Name:            emqx.GetServiceAccountName(),
		Namespace:       emqx.GetNamespace(),
		Labels:          labels,
		OwnerReferences: ownerRefs,
	}

	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: meta,
	}

	role := &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "Role",
		},
		ObjectMeta: meta,
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:     []string{"get", "watch", "list"},
				APIGroups: []string{""},
				Resources: []string{"endpoints"},
			},
		},
	}

	roleBinding := &rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "RoleBinding",
		},
		ObjectMeta: meta,
		Subjects: []rbacv1.Subject{
			{
				Kind:      sa.Kind,
				Name:      sa.Name,
				Namespace: sa.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     role.Kind,
			Name:     role.Name,
		},
	}

	return sa, role, roleBinding
}

func NewHeadLessSvcForCR(emqx v1beta1.Emqx, labels map[string]string, ownerRefs []metav1.OwnerReference) *corev1.Service {
	ports, _, _ := generatePorts(emqx)

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels:          labels,
			Name:            emqx.GetHeadlessServiceName(),
			Namespace:       emqx.GetNamespace(),
			OwnerReferences: ownerRefs,
		},
		Spec: corev1.ServiceSpec{
			Ports:     ports,
			Selector:  labels,
			ClusterIP: corev1.ClusterIPNone,
		},
	}
}

func NewListenerSvcForCR(emqx v1beta1.Emqx, labels map[string]string, ownerRefs []metav1.OwnerReference) *corev1.Service {
	listener := emqx.GetListener()
	ports, _, _ := generatePorts(emqx)

	listenerSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels:          labels,
			Name:            emqx.GetName(),
			Namespace:       emqx.GetNamespace(),
			OwnerReferences: ownerRefs,
		},
		Spec: corev1.ServiceSpec{
			Type:                     listener.Type,
			LoadBalancerIP:           listener.LoadBalancerIP,
			LoadBalancerSourceRanges: listener.LoadBalancerSourceRanges,
			ExternalIPs:              listener.ExternalIPs,
			Ports:                    ports,
			Selector:                 labels,
		},
	}
	return listenerSvc
}

func NewConfigMapForAcl(emqx v1beta1.Emqx, labels map[string]string, ownerRefs []metav1.OwnerReference) *corev1.ConfigMap {
	acl := util.GetACL(emqx)
	cmForAcl := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Labels:          labels,
			Name:            acl["name"],
			Namespace:       emqx.GetNamespace(),
			OwnerReferences: ownerRefs,
		},
		Data: map[string]string{"acl.conf": acl["conf"]},
	}

	return cmForAcl
}

func NewConfigMapForLoadedModules(emqx v1beta1.Emqx, labels map[string]string, ownerRefs []metav1.OwnerReference) *corev1.ConfigMap {
	modules := util.GetLoadedModules(emqx)
	cmForPM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Labels:          labels,
			Name:            modules["name"],
			Namespace:       emqx.GetNamespace(),
			OwnerReferences: ownerRefs,
		},
		Data: map[string]string{"loaded_modules": modules["conf"]},
	}

	return cmForPM
}

func NewConfigMapForLoadedPlugins(emqx v1beta1.Emqx, labels map[string]string, ownerRefs []metav1.OwnerReference) *corev1.ConfigMap {
	plugins := util.GetLoadedPlugins(emqx)
	cmForPG := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Labels:          labels,
			Name:            plugins["name"],
			Namespace:       emqx.GetNamespace(),
			OwnerReferences: ownerRefs,
		},
		Data: map[string]string{"loaded_plugins": plugins["conf"]},
	}

	return cmForPG
}

func NewEmqxStatefulSet(emqx v1beta1.Emqx, labels map[string]string, ownerRefs []metav1.OwnerReference) *appsv1.StatefulSet {
	_, ports, env := generatePorts(emqx)

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            emqx.GetName(),
			Namespace:       emqx.GetNamespace(),
			Labels:          labels,
			OwnerReferences: ownerRefs,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: emqx.GetHeadlessServiceName(),
			Replicas:    emqx.GetReplicas(),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			PodManagementPolicy: appsv1.ParallelPodManagement,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: emqx.GetAnnotations(),
				},
				Spec: corev1.PodSpec{
					Affinity:           emqx.GetAffinity(),
					ServiceAccountName: emqx.GetServiceAccountName(),
					SecurityContext:    generateSecurityContext(),
					Tolerations:        emqx.GetToleRations(),
					NodeSelector:       emqx.GetNodeSelector(),
					Containers: []corev1.Container{
						{
							Name:            emqx.GetName(),
							Image:           emqx.GetImage(),
							ImagePullPolicy: emqx.GetImagePullPolicy(),
							SecurityContext: generateContainerSecurityContext(),
							Resources:       emqx.GetResource(),
							Env:             mergeEnv(env, emqx.GetEnv()),
							Ports:           ports,
							VolumeMounts:    generateEmqxVolumeMounts(emqx),
							ReadinessProbe: &corev1.Probe{
								InitialDelaySeconds: 5,
								PeriodSeconds:       5,
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/status",
										Port: intstr.IntOrString{
											IntVal: emqx.GetListener().Ports.API,
										},
									},
								},
							},
						},
					},
					Volumes: generateEmqxVolumes(emqx),
				},
			},
		},
	}

	storageSpec := emqx.GetStorage()
	if !reflect.ValueOf(storageSpec).IsNil() {
		sts.Spec.VolumeClaimTemplates = append(
			sts.Spec.VolumeClaimTemplates,
			generateVolumeClaimTemplate(emqx, util.GetDataVolume(emqx)["name"]),
			generateVolumeClaimTemplate(emqx, util.GetLogVolume(emqx)["name"]),
		)
	}

	return sts
}

func generateSecurityContext() *corev1.PodSecurityContext {
	emqxUserGroup := int64(1000)
	emqxUser := int64(1000)
	runAsNonRoot := true
	fsGroupChangeAlways := corev1.FSGroupChangeAlways

	return &corev1.PodSecurityContext{
		FSGroup:             &emqxUserGroup,
		FSGroupChangePolicy: &fsGroupChangeAlways,
		RunAsNonRoot:        &runAsNonRoot,
		RunAsUser:           &emqxUser,
		SupplementalGroups:  []int64{emqxUserGroup},
	}
}
func generateContainerSecurityContext() *corev1.SecurityContext {
	emqxUser := int64(1000)
	runAsNonRoot := true

	return &corev1.SecurityContext{
		RunAsNonRoot: &runAsNonRoot,
		RunAsUser:    &emqxUser,
	}
}

func generateEmqxVolumeMounts(emqx v1beta1.Emqx) []corev1.VolumeMount {
	volumeMounts := []corev1.VolumeMount{}

	data := util.GetDataVolume(emqx)
	volumeMounts = append(volumeMounts,
		corev1.VolumeMount{
			Name:      data["name"],
			MountPath: data["mountPath"],
		},
	)

	log := util.GetLogVolume(emqx)
	volumeMounts = append(volumeMounts,
		corev1.VolumeMount{
			Name:      log["name"],
			MountPath: log["mountPath"],
		},
	)

	extraVolumeMounts := emqx.GetExtraVolumeMounts()
	volumeMounts = append(volumeMounts, extraVolumeMounts...)

	acl := util.GetACL(emqx)
	volumeMounts = append(volumeMounts,
		corev1.VolumeMount{
			Name:      acl["name"],
			MountPath: acl["mountPath"],
			SubPath:   acl["subPath"],
		},
	)
	modules := util.GetLoadedModules(emqx)
	volumeMounts = append(volumeMounts,
		corev1.VolumeMount{
			Name:      modules["name"],
			MountPath: modules["mountPath"],
			SubPath:   modules["subPath"],
		},
	)
	plugins := util.GetLoadedPlugins(emqx)
	volumeMounts = append(volumeMounts,
		corev1.VolumeMount{
			Name:      plugins["name"],
			MountPath: plugins["mountPath"],
			SubPath:   plugins["subPath"],
		},
	)
	emqxEnterprise, ok := emqx.(*v1beta1.EmqxEnterprise)
	if ok && emqxEnterprise.GetLicense() != "" {
		license := util.GetLicense(emqx)
		volumeMounts = append(volumeMounts,
			corev1.VolumeMount{
				Name:      license["name"],
				MountPath: license["mountPath"],
				SubPath:   license["subPath"],
				ReadOnly:  true,
			},
		)
	}
	return volumeMounts
}

func generateEmqxVolumes(emqx v1beta1.Emqx) []corev1.Volume {
	volumes := []corev1.Volume{}
	storageSpec := emqx.GetStorage()
	if reflect.ValueOf(storageSpec).IsNil() {
		volumes = append(volumes,
			corev1.Volume{
				Name: util.GetDataVolume(emqx)["name"],
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		)
		volumes = append(volumes,
			corev1.Volume{
				Name: util.GetLogVolume(emqx)["name"],
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			},
		)
	}

	extraVolumes := emqx.GetExtraVolumes()
	volumes = append(volumes, extraVolumes...)

	acl := util.GetACL(emqx)
	volumes = append(volumes,
		corev1.Volume{
			Name: acl["name"],
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: acl["name"],
					},
					Items: []corev1.KeyToPath{
						{
							Key:  acl["subPath"],
							Path: acl["subPath"],
						},
					},
				},
			},
		},
	)

	plugins := util.GetLoadedPlugins(emqx)
	volumes = append(volumes,
		corev1.Volume{
			Name: plugins["name"],
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: plugins["name"],
					},
					Items: []corev1.KeyToPath{
						{
							Key:  plugins["subPath"],
							Path: plugins["subPath"],
						},
					},
				},
			},
		},
	)

	modules := util.GetLoadedModules(emqx)
	volumes = append(volumes,
		corev1.Volume{
			Name: modules["name"],
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: modules["name"],
					},
					Items: []corev1.KeyToPath{
						{
							Key:  modules["subPath"],
							Path: modules["subPath"],
						},
					},
				},
			},
		},
	)
	emqxEnterprise, ok := emqx.(*v1beta1.EmqxEnterprise)
	if ok && emqxEnterprise.GetLicense() != "" {
		volumes = append(volumes,
			corev1.Volume{
				Name: util.GetLicense(emqx)["name"],
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: util.GetLicense(emqx)["name"],
						Items: []corev1.KeyToPath{
							{
								Key:  "emqx.lic",
								Path: "emqx.lic",
							},
						},
					},
				},
			},
		)
	}
	return volumes
}

func generateVolumeClaimTemplate(emqx v1beta1.Emqx, Name string) corev1.PersistentVolumeClaim {
	template := emqx.GetStorage().VolumeClaimTemplate
	pvc := corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: template.APIVersion,
			Kind:       template.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Name,
			Namespace: emqx.GetNamespace(),
		},
		Spec:   template.Spec,
		Status: template.Status,
	}
	if pvc.Spec.AccessModes == nil {
		pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
	}
	if pvc.Spec.VolumeMode == nil {
		fileSystem := corev1.PersistentVolumeFilesystem
		pvc.Spec.VolumeMode = &fileSystem
	}
	return pvc
}

func generatePorts(emqx v1beta1.Emqx) ([]corev1.ServicePort, []corev1.ContainerPort, []corev1.EnvVar) {
	var servicePorts []corev1.ServicePort
	var containerPorts []corev1.ContainerPort
	var env []corev1.EnvVar
	listener := emqx.GetListener()
	if !reflect.ValueOf(listener.Ports.MQTT).IsZero() {
		env = append(env, corev1.EnvVar{
			Name:  "EMQX_LISTENER__TCP__EXTERNAL",
			Value: fmt.Sprint(listener.Ports.MQTT),
		})
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          "mqtt",
			ContainerPort: listener.Ports.MQTT,
		})
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:     "mqtt",
			Port:     listener.Ports.MQTT,
			Protocol: "TCP",
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: listener.Ports.MQTT,
			},
		})
	}
	if !reflect.ValueOf(listener.Ports.MQTTS).IsZero() {
		env = append(env, corev1.EnvVar{
			Name:  "EMQX_LISTENER__SSL__EXTERNAL",
			Value: fmt.Sprint(listener.Ports.MQTTS),
		})
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          "mqtts",
			ContainerPort: listener.Ports.MQTTS,
		})
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:     "mqtts",
			Port:     listener.Ports.MQTTS,
			Protocol: "TCP",
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: listener.Ports.MQTTS,
			},
		})
	}
	if !reflect.ValueOf(listener.Ports.WS).IsZero() {
		env = append(env, corev1.EnvVar{
			Name:  "EMQX_LISTENER__WS__EXTERNAL",
			Value: fmt.Sprint(listener.Ports.WS),
		})
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          "ws",
			ContainerPort: listener.Ports.WS,
		})
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:     "ws",
			Port:     listener.Ports.WS,
			Protocol: "TCP",
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: listener.Ports.WS,
			},
		})
	}
	if !reflect.ValueOf(listener.Ports.WSS).IsZero() {
		env = append(env, corev1.EnvVar{
			Name:  "EMQX_LISTENER__WSS__EXTERNAL",
			Value: fmt.Sprint(listener.Ports.WSS),
		})
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          "wss",
			ContainerPort: listener.Ports.WSS,
		})
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:     "wss",
			Port:     listener.Ports.WSS,
			Protocol: "TCP",
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: listener.Ports.WSS,
			},
		})
	}
	if !reflect.ValueOf(listener.Ports.Dashboard).IsZero() {
		env = append(env, corev1.EnvVar{
			Name:  "EMQX_DASHBOARD__LISTENER__HTTP",
			Value: fmt.Sprint(listener.Ports.Dashboard),
		})
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          "dashboard",
			ContainerPort: listener.Ports.Dashboard,
		})
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:     "dashboard",
			Port:     listener.Ports.Dashboard,
			Protocol: "TCP",
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: listener.Ports.Dashboard,
			},
		})
	}
	if !reflect.ValueOf(listener.Ports.API).IsZero() {
		env = append(env, corev1.EnvVar{
			Name:  "EMQX_MANAGEMENT__LISTENER__HTTP",
			Value: fmt.Sprint(listener.Ports.API),
		})
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          "api",
			ContainerPort: listener.Ports.API,
		})
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:     "api",
			Port:     listener.Ports.API,
			Protocol: "TCP",
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: listener.Ports.API,
			},
		})
	}
	return servicePorts, containerPorts, env
}

func contains(Env []corev1.EnvVar, Name string) int {
	for index, value := range Env {
		if value.Name == Name {
			return index
		}
	}
	return -1
}

func mergeEnv(env1, env2 []corev1.EnvVar) []corev1.EnvVar {
	for _, value := range env2 {
		r := contains(env1, value.Name)
		if r == -1 {
			env1 = append(env1, value)
		}
	}
	return env1
}
