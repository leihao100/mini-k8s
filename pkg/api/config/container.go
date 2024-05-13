package config

import (
	"github.com/docker/go-connections/nat"
)

type Container struct {
	Name         string              `json:"name,omitempty"`
	Args         []string            `json:"args,omitempty"`
	Cmd          []string            `json:"cmd,omitempty"`        // Command to run when starting the container
	Entrypoint   []string            `json:"entrypoint,omitempty"` // Entrypoint to run when starting the container
	Env          []string            `json:"env,omitempty"`        // List of environment variable to set in the container
	Image        string              `json:"image,omitempty"`      // Name of the image as it was passed by the operator (e.g. could be symbolic)
	Volumes      map[string]struct{} `json:"volumes,omitempty"`    // List of volumes (mounts) used for the container
	Labels       map[string]string   `json:"label,omitempty"`      // List of labels set to this container
	PortBindings nat.PortMap         `json:"portBindings,omitempty"`
	VolumesFrom  []string            `json:"volumesFrom,omitempty"`
	Binds        []string            `json:"binds,omitempty"`
	NetworkMode  string              `json:"networkMode,omitempty"`
	CPULimit     int64               `json:"CPULimit,omitempty"`
	MemLimit     int64               `json:"memLimit,omitempty"`
	Pause        string              `json:"pause,omitempty"`
}

//以下为官方文档中的config参数
//type Config struct {
//	Hostname        string              // Hostname
//	Domainname      string              // Domainname
//	User            string              // User that will run the command(s) inside the container, also support user:group
//	AttachStdin     bool                // Attach the standard input, makes possible user interaction
//	AttachStdout    bool                // Attach the standard output
//	AttachStderr    bool                // Attach the standard error
//	ExposedPorts    nat.PortSet         `json:",omitempty"` // List of exposed ports
//	Tty             bool                // Attach standard streams to a tty, including stdin if it is not closed.
//	OpenStdin       bool                // Open stdin
//	StdinOnce       bool                // If true, close stdin after the 1 attached client disconnects.
//	Env             []string            // List of environment variable to set in the container
//	Cmd             strslice.StrSlice   // Command to run when starting the container
//	Healthcheck     *HealthConfig       `json:",omitempty"` // Healthcheck describes how to check the container is healthy
//	ArgsEscaped     bool                `json:",omitempty"` // True if command is already escaped (Windows specific)
//	Image           string              // Name of the image as it was passed by the operator (e.g. could be symbolic)
//	Volumes         map[string]struct{} // List of volumes (mounts) used for the container
//	WorkingDir      string              // Current directory (PWD) in the command will be launched
//	Entrypoint      strslice.StrSlice   // Entrypoint to run when starting the container
//	NetworkDisabled bool                `json:",omitempty"` // Is network disabled
//	MacAddress      string              `json:",omitempty"` // Mac Address of the container
//	OnBuild         []string            // ONBUILD metadata that were defined on the image Dockerfile
//	Labels          map[string]string   // List of labels set to this container
//	StopSignal      string              `json:",omitempty"` // Signal to stop a container
//	StopTimeout     *int                `json:",omitempty"` // Timeout (in seconds) to stop a container
//	Shell           strslice.StrSlice   `json:",omitempty"` // Shell for shell-form of RUN, CMD, ENTRYPOINT
//}

//type HostConfig struct {
//	Binds           []string
//	ContainerIDFile string
//	LogConfig       LogConfig
//	NetworkMode     NetworkMode
//	PortBindings    nat.PortMap
//	RestartPolicy   RestartPolicy
//	AutoRemove      bool
//	VolumeDriver    string
//	VolumesFrom     []string
//	ConsoleSize     [2]uint
//	Annotations     map[string]string `json:",omitempty"`
//	CapAdd          strslice.StrSlice
//	CapDrop         strslice.StrSlice
//	CgroupnsMode    CgroupnsMode
//	DNS             []string `json:"Dns"`
//	DNSOptions      []string `json:"DnsOptions"`
//	DNSSearch       []string `json:"DnsSearch"`
//	ExtraHosts      []string
//	GroupAdd        []string
//	IpcMode         IpcMode
//	Cgroup          CgroupSpec
//	Links           []string
//	OomScoreAdj     int
//	PidMode         PidMode
//	Privileged      bool
//	PublishAllPorts bool
//	ReadonlyRootfs  bool
//	SecurityOpt     []string
//	StorageOpt      map[string]string `json:",omitempty"`
//	Tmpfs           map[string]string `json:",omitempty"`
//	UTSMode         UTSMode
//	UsernsMode      UsernsMode
//	ShmSize         int64
//	Sysctls         map[string]string `json:",omitempty"`
//	Runtime         string            `json:",omitempty"`
//	Isolation       Isolation
//	Resources
//	Mounts        []mount.Mount `json:",omitempty"`
//	MaskedPaths   []string
//	ReadonlyPaths []string
//	Init          *bool `json:",omitempty"`
//}
