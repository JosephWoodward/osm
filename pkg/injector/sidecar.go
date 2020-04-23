package injector

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/open-service-mesh/osm/pkg/constants"
)

const (
	envoyBootstrapConfigFile  = "bootstrap.yaml"
	envoyProxyConfigPath      = "/etc/envoy"
	envoySidecarContainerName = "envoyproxy"
)

func getEnvoySidecarContainerSpec(data *EnvoySidecarData) corev1.Container {
	container := corev1.Container{
		Name:            data.Name,
		Image:           data.Image,
		ImagePullPolicy: corev1.PullAlways,
		SecurityContext: &corev1.SecurityContext{
			RunAsUser: func() *int64 {
				uid := constants.EnvoyUID
				return &uid
			}(),
		},
		Ports: []corev1.ContainerPort{
			{
				Name:          constants.EnvoyAdminPortName,
				ContainerPort: constants.EnvoyAdminPort,
			},
			{
				Name:          constants.EnvoyInboundListenerPortName,
				ContainerPort: constants.EnvoyInboundListenerPort,
			},
			{
				Name:          constants.EnvoyInboundPromethuesListenerPortName,
				ContainerPort: constants.EnvoyPrometheusInboundListenerPort,
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      envoyBootstrapConfigVolume,
				ReadOnly:  true,
				MountPath: envoyProxyConfigPath,
			},
		},
		Command: []string{
			"envoy",
		},
		Args: []string{
			"--log-level", "debug", // TODO: remove
			"--config-path", getEnvoyConfigPath(),
			"--service-node", data.Service,
			"--service-cluster", data.Service,
		},
	}

	return container
}
