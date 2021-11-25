package constants

const (
	//  The default value for EMQ X Cluster Name
	EMQX_NAME = "emqx"

	// The constant value for configmap
	EMQX_LIC_NAME    = "emqx-lic"
	EMQX_LIC_DIR     = "/opt/emqx/etc/emqx.lic"
	EMQX_LIC_SUBPATH = "emqx.lic"

	EMQX_DATA_DIR = "/opt/emqx/data"
	EMQX_LOG_DIR  = "/opt/emqx/log"

	// The default value for service ports
	EMQX_LISTENERS__TCP__EXTERNAL_NAME = "mqtt"
	EMQX_LISTENERS__TCP__EXTERNAL_PORT = 1883

	EMQX_LISTENERS__SSL__EXTERNAL_NAME = "mqtts"
	EMQX_LISTENERS__SSL__EXTERNAL_PORT = 8883

	EMQX_LISTENERS__WS__EXTERNAL_NAME = "ws"
	EMQX_LISTENERS__WS__EXTERNAL_PORT = 8083

	EMQX_LISTENERS__WSS__EXTERNAL_NAME = "wss"
	EMQX_LISTENERS__WSS__EXTERNAL_PORT = 8084

	EMQX_DASHBOARD__LISTENER__HTTP_NAME = "dashboard"
	EMQX_DASHBOARD__LISTENER__HTTP_PORT = 18083

	EMQX_MANAGEMENT__LISTENER__HTTP_NAME = "api"
	EMQX_MANAGEMENT__LISTENER__HTTP_PORT = 8081

	// The constant key-value for labels
	OPERATOR_NAME        = "emqx-operator"
	LABEL_MANAGED_BY_KEY = "apps.emqx.io/managed-by"
	LABEL_NAME_KEY       = "emqx-operator/v1alpha2"
)
