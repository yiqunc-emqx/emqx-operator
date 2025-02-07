package v1beta1

//+kubebuilder:object:generate=true
type Plugin struct {
	Name   string `json:"name,omitempty"`
	Enable bool   `json:"enable,omitempty"`
}

func generatePlugins(plugins []Plugin) []Plugin {
	if plugins == nil {
		return defaultLoadedPlugins()
	}

	contains := func(plugins []Plugin) int {
		for index, value := range plugins {
			if value.Name == "emqx_management" {
				return index
			}
		}
		return -1
	}

	if contains(plugins) == -1 {
		plugins = append(plugins, Plugin{Name: "emqx_management", Enable: true})
	}

	return plugins
}

func defaultLoadedPlugins() []Plugin {
	return []Plugin{
		{
			Name:   "emqx_management",
			Enable: true,
		},
		{
			Name:   "emqx_recon",
			Enable: true,
		},
		{
			Name:   "emqx_retainer",
			Enable: true,
		},
		{
			Name:   "emqx_dashboard",
			Enable: true,
		},
		{
			Name:   "emqx_telemetry",
			Enable: true,
		},
		{
			Name:   "emqx_rule_engine",
			Enable: true,
		},
	}
}
