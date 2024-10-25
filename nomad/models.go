package nomad 

type ResourceUsage struct {
	MemoryStats struct {
		RSS            int      `json:"RSS"`
		Cache          int      `json:"Cache"`
		Swap           int      `json:"Swap"`
		MappedFile     int      `json:"MappedFile"`
		Usage          int      `json:"Usage"`
		MaxUsage       int      `json:"MaxUsage"`
		KernelUsage    int      `json:"KernelUsage"`
		KernelMaxUsage int      `json:"KernelMaxUsage"`
		Measured       []string `json:"Measured"`
	} `json:"MemoryStats"`
	CpuStats struct {
		SystemMode       float64  `json:"SystemMode"`
		UserMode         float64  `json:"UserMode"`
		TotalTicks       float64  `json:"TotalTicks"`
		ThrottledPeriods int      `json:"ThrottledPeriods"`
		ThrottledTime    int      `json:"ThrottledTime"`
		Percent          float64  `json:"Percent"`
		Measured         []string `json:"Measured"`
	} `json:"CpuStats"`
	DeviceStats []interface{} `json:"DeviceStats"`
}

type JobDefinition struct {
	Stop        bool     `json:"Stop"`
	Region      string   `json:"Region"`
	Namespace   string   `json:"Namespace"`
	ID          string   `json:"ID"`
	ParentID    string   `json:"ParentID"`
	Name        string   `json:"Name"`
	Type        string   `json:"Type"`
	Priority    int      `json:"Priority"`
	AllAtOnce   bool     `json:"AllAtOnce"`
	Datacenters []string `json:"Datacenters"`
	NodePool    string   `json:"NodePool"`
	Constraints []struct {
		LTarget string `json:"LTarget"`
		RTarget string `json:"RTarget"`
		Operand string `json:"Operand"`
	} `json:"Constraints"`
	Affinities interface{} `json:"Affinities"`
	Spreads    interface{} `json:"Spreads"`
	TaskGroups []struct {
		Name   string `json:"Name"`
		Count  int    `json:"Count"`
		Update struct {
			Stagger          int    `json:"Stagger"`
			MaxParallel      int    `json:"MaxParallel"`
			HealthCheck      string `json:"HealthCheck"`
			MinHealthyTime   int    `json:"MinHealthyTime"`
			HealthyDeadline  int    `json:"HealthyDeadline"`
			ProgressDeadline int    `json:"ProgressDeadline"`
			AutoRevert       bool   `json:"AutoRevert"`
			AutoPromote      bool   `json:"AutoPromote"`
			Canary           int    `json:"Canary"`
		} `json:"Update"`
		Migrate struct {
			MaxParallel     int    `json:"MaxParallel"`
			HealthCheck     string `json:"HealthCheck"`
			MinHealthyTime  int    `json:"MinHealthyTime"`
			HealthyDeadline int    `json:"HealthyDeadline"`
		} `json:"Migrate"`
		Constraints []struct {
			LTarget string `json:"LTarget"`
			RTarget string `json:"RTarget"`
			Operand string `json:"Operand"`
		} `json:"Constraints"`
		Scaling       interface{} `json:"Scaling"`
		RestartPolicy struct {
			Attempts        int    `json:"Attempts"`
			Interval        int    `json:"Interval"`
			Delay           int    `json:"Delay"`
			Mode            string `json:"Mode"`
			RenderTemplates bool   `json:"RenderTemplates"`
		} `json:"RestartPolicy"`
		Tasks []struct {
			Name   string `json:"Name"`
			Driver string `json:"Driver"`
			User   string `json:"User"`
			Config struct {
				Mount      interface{} `json:"mount"`
				ForcePull  bool        `json:"force_pull"`
				Privileged bool        `json:"privileged"`
				Image      string      `json:"image"`
			} `json:"Config"`
			Env         map[string]string `json:"Env"`
			Services    interface{}       `json:"Services"`
			Vault       interface{}       `json:"Vault"`
			Consul      interface{}       `json:"Consul"`
			Templates   interface{}       `json:"Templates"`
			Constraints interface{}       `json:"Constraints"`
			Affinities  interface{}       `json:"Affinities"`
			Resources   struct {
				CPU         int         `json:"CPU"`
				Cores       int         `json:"Cores"`
				MemoryMB    int         `json:"MemoryMB"`
				MemoryMaxMB int         `json:"MemoryMaxMB"`
				DiskMB      int         `json:"DiskMB"`
				IOPS        int         `json:"IOPS"`
				Networks    interface{} `json:"Networks"`
				Devices     interface{} `json:"Devices"`
				NUMA        interface{} `json:"NUMA"`
			} `json:"Resources"`
			RestartPolicy struct {
				Attempts        int    `json:"Attempts"`
				Interval        int    `json:"Interval"`
				Delay           int    `json:"Delay"`
				Mode            string `json:"Mode"`
				RenderTemplates bool   `json:"RenderTemplates"`
			} `json:"RestartPolicy"`
			DispatchPayload interface{} `json:"DispatchPayload"`
			Lifecycle       interface{} `json:"Lifecycle"`
			Meta            interface{} `json:"Meta"`
			KillTimeout     int         `json:"KillTimeout"`
			LogConfig       struct {
				MaxFiles      int  `json:"MaxFiles"`
				MaxFileSizeMB int  `json:"MaxFileSizeMB"`
				Disabled      bool `json:"Disabled"`
			} `json:"LogConfig"`
			Artifacts       interface{} `json:"Artifacts"`
			Leader          bool        `json:"Leader"`
			ShutdownDelay   int         `json:"ShutdownDelay"`
			VolumeMounts    interface{} `json:"VolumeMounts"`
			ScalingPolicies interface{} `json:"ScalingPolicies"`
			KillSignal      string      `json:"KillSignal"`
			Kind            string      `json:"Kind"`
			CSIPluginConfig interface{} `json:"CSIPluginConfig"`
			Identity        struct {
				Name         string   `json:"Name"`
				Audience     []string `json:"Audience"`
				ChangeMode   string   `json:"ChangeMode"`
				ChangeSignal string   `json:"ChangeSignal"`
				Env          bool     `json:"Env"`
				File         bool     `json:"File"`
				ServiceName  string   `json:"ServiceName"`
				TTL          int      `json:"TTL"`
			} `json:"Identity"`
			Identities interface{} `json:"Identities"`
			Actions    interface{} `json:"Actions"`
		} `json:"Tasks"`
		EphemeralDisk struct {
			Sticky  bool `json:"Sticky"`
			SizeMB  int  `json:"SizeMB"`
			Migrate bool `json:"Migrate"`
		} `json:"EphemeralDisk"`
		Meta             interface{} `json:"Meta"`
		ReschedulePolicy struct {
			Attempts      int    `json:"Attempts"`
			Interval      int    `json:"Interval"`
			Delay         int    `json:"Delay"`
			DelayFunction string `json:"DelayFunction"`
			MaxDelay      int    `json:"MaxDelay"`
			Unlimited     bool   `json:"Unlimited"`
		} `json:"ReschedulePolicy"`
		Affinities interface{} `json:"Affinities"`
		Spreads    interface{} `json:"Spreads"`
		Networks   []struct {
			Mode          string      `json:"Mode"`
			Device        string      `json:"Device"`
			CIDR          string      `json:"CIDR"`
			IP            string      `json:"IP"`
			Hostname      string      `json:"Hostname"`
			MBits         int         `json:"MBits"`
			DNS           interface{} `json:"DNS"`
			ReservedPorts interface{} `json:"ReservedPorts"`
			DynamicPorts  []struct {
				Label       string `json:"Label"`
				Value       int    `json:"Value"`
				To          int    `json:"To"`
				HostNetwork string `json:"HostNetwork"`
			} `json:"DynamicPorts"`
		} `json:"Networks"`
		Consul struct {
			Namespace string `json:"Namespace"`
			Cluster   string `json:"Cluster"`
			Partition string `json:"Partition"`
		} `json:"Consul"`
		Services []struct {
			Name              string      `json:"Name"`
			TaskName          string      `json:"TaskName"`
			PortLabel         string      `json:"PortLabel"`
			AddressMode       string      `json:"AddressMode"`
			Address           string      `json:"Address"`
			EnableTagOverride bool        `json:"EnableTagOverride"`
			Tags              []string    `json:"Tags"`
			CanaryTags        interface{} `json:"CanaryTags"`
			Checks            []struct {
				Name          string      `json:"Name"`
				Type          string      `json:"Type"`
				Command       string      `json:"Command"`
				Args          interface{} `json:"Args"`
				Path          string      `json:"Path"`
				Protocol      string      `json:"Protocol"`
				PortLabel     string      `json:"PortLabel"`
				Expose        bool        `json:"Expose"`
				AddressMode   string      `json:"AddressMode"`
				Interval      int         `json:"Interval"`
				Timeout       int         `json:"Timeout"`
				InitialStatus string      `json:"InitialStatus"`
				TLSServerName string      `json:"TLSServerName"`
				TLSSkipVerify bool        `json:"TLSSkipVerify"`
				Method        string      `json:"Method"`
				Header        interface{} `json:"Header"`
				CheckRestart  struct {
					Limit          int  `json:"Limit"`
					Grace          int  `json:"Grace"`
					IgnoreWarnings bool `json:"IgnoreWarnings"`
				} `json:"CheckRestart"`
				GRPCService            string `json:"GRPCService"`
				GRPCUseTLS             bool   `json:"GRPCUseTLS"`
				TaskName               string `json:"TaskName"`
				SuccessBeforePassing   int    `json:"SuccessBeforePassing"`
				FailuresBeforeCritical int    `json:"FailuresBeforeCritical"`
				FailuresBeforeWarning  int    `json:"FailuresBeforeWarning"`
				Body                   string `json:"Body"`
				OnUpdate               string `json:"OnUpdate"`
			} `json:"Checks"`
			Connect struct {
				Native         bool `json:"Native"`
				SidecarService struct {
					Tags  interface{} `json:"Tags"`
					Port  string      `json:"Port"`
					Proxy struct {
						LocalServiceAddress string      `json:"LocalServiceAddress"`
						LocalServicePort    int         `json:"LocalServicePort"`
						Upstreams           interface{} `json:"Upstreams"`
						Expose              interface{} `json:"Expose"`
						Config              interface{} `json:"Config"`
					} `json:"Proxy"`
					DisableDefaultTCPCheck bool        `json:"DisableDefaultTCPCheck"`
					Meta                   interface{} `json:"Meta"`
				} `json:"SidecarService"`
				SidecarTask interface{} `json:"SidecarTask"`
				Gateway     interface{} `json:"Gateway"`
			} `json:"Connect"`
			Meta            interface{} `json:"Meta"`
			CanaryMeta      interface{} `json:"CanaryMeta"`
			TaggedAddresses interface{} `json:"TaggedAddresses"`
			Namespace       string      `json:"Namespace"`
			OnUpdate        string      `json:"OnUpdate"`
			Provider        string      `json:"Provider"`
			Cluster         string      `json:"Cluster"`
			Identity        interface{} `json:"Identity"`
		} `json:"Services"`
		Volumes                   interface{} `json:"Volumes"`
		ShutdownDelay             interface{} `json:"ShutdownDelay"`
		StopAfterClientDisconnect interface{} `json:"StopAfterClientDisconnect"`
		MaxClientDisconnect       interface{} `json:"MaxClientDisconnect"`
		PreventRescheduleOnLost   bool        `json:"PreventRescheduleOnLost"`
	} `json:"TaskGroups"`
	Update struct {
		Stagger          int    `json:"Stagger"`
		MaxParallel      int    `json:"MaxParallel"`
		HealthCheck      string `json:"HealthCheck"`
		MinHealthyTime   int    `json:"MinHealthyTime"`
		HealthyDeadline  int    `json:"HealthyDeadline"`
		ProgressDeadline int    `json:"ProgressDeadline"`
		AutoRevert       bool   `json:"AutoRevert"`
		AutoPromote      bool   `json:"AutoPromote"`
		Canary           int    `json:"Canary"`
	} `json:"Update"`
	Multiregion              interface{} `json:"Multiregion"`
	Periodic                 interface{} `json:"Periodic"`
	ParameterizedJob         interface{} `json:"ParameterizedJob"`
	Dispatched               bool        `json:"Dispatched"`
	DispatchIdempotencyToken string      `json:"DispatchIdempotencyToken"`
	Payload                  interface{} `json:"Payload"`
	Meta                     interface{} `json:"Meta"`
	ConsulToken              string      `json:"ConsulToken"`
	ConsulNamespace          string      `json:"ConsulNamespace"`
	VaultToken               string      `json:"VaultToken"`
	VaultNamespace           string      `json:"VaultNamespace"`
	NomadTokenID             string      `json:"NomadTokenID"`
	Status                   string      `json:"Status"`
	StatusDescription        string      `json:"StatusDescription"`
	Stable                   bool        `json:"Stable"`
	Version                  int         `json:"Version"`
	SubmitTime               int         `json:"SubmitTime"`
	CreateIndex              int         `json:"CreateIndex"`
	ModifyIndex              int         `json:"ModifyIndex"`
	JobModifyIndex           int         `json:"JobModifyIndex"`
}

type NomadAllocations struct {
	NomadAllocations []NomadAlloc
}

// serialize me the json
type NomadAlloc struct {
	ID        string        `json:"ID"`
	Namespace string        `json:"Namespace"`
	EvalID    string        `json:"EvalID"`
	Name      string        `json:"Name"`
	NodeID    string        `json:"NodeID"`
	NodeName  string        `json:"NodeName"`
	JobID     string        `json:"JobID"`
	Job       JobDefinition `json:"Job"`
	TaskGroup string        `json:"TaskGroup"`
	Resources struct {
		CPU         int `json:"CPU"`
		Cores       int `json:"Cores"`
		MemoryMB    int `json:"MemoryMB"`
		MemoryMaxMB int `json:"MemoryMaxMB"`
		DiskMB      int `json:"DiskMB"`
		IOPS        int `json:"IOPS"`
		Networks    []struct {
			Mode          string      `json:"Mode"`
			Device        string      `json:"Device"`
			CIDR          string      `json:"CIDR"`
			IP            string      `json:"IP"`
			Hostname      string      `json:"Hostname"`
			MBits         int         `json:"MBits"`
			DNS           interface{} `json:"DNS"`
			ReservedPorts []struct {
				Label       string `json:"Label"`
				Value       int    `json:"Value"`
				To          int    `json:"To"`
				HostNetwork string `json:"HostNetwork"`
			} `json:"ReservedPorts"`
			DynamicPorts interface{} `json:"DynamicPorts"`
		} `json:"Networks"`
		Devices interface{} `json:"Devices"`
		NUMA    interface{} `json:"NUMA"`
	} `json:"Resources"`
	SharedResources struct {
		CPU         int `json:"CPU"`
		Cores       int `json:"Cores"`
		MemoryMB    int `json:"MemoryMB"`
		MemoryMaxMB int `json:"MemoryMaxMB"`
		DiskMB      int `json:"DiskMB"`
		IOPS        int `json:"IOPS"`
		Networks    []struct {
			Mode          string      `json:"Mode"`
			Device        string      `json:"Device"`
			CIDR          string      `json:"CIDR"`
			IP            string      `json:"IP"`
			Hostname      string      `json:"Hostname"`
			MBits         int         `json:"MBits"`
			DNS           interface{} `json:"DNS"`
			ReservedPorts []struct {
				Label       string `json:"Label"`
				Value       int    `json:"Value"`
				To          int    `json:"To"`
				HostNetwork string `json:"HostNetwork"`
			} `json:"ReservedPorts"`
			DynamicPorts interface{} `json:"DynamicPorts"`
		} `json:"Networks"`
		Devices interface{} `json:"Devices"`
		NUMA    interface{} `json:"NUMA"`
	} `json:"SharedResources"`
}
