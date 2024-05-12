package status

import "github.com/docker/docker/api/types"

type ContainerStatus struct {
	Name        string
	ContainerID string
	ImageID     string
	Image       string
	State       types.ContainerState
	//LastState    types.ContainerState
	Started bool
	//Ready        bool
	RestartCount int64
}

/*

containerID
string	Container's ID in the format '<type>://<container_id>'.
image
string	The image the container is running. More info: https://kubernetes.io/docs/concepts/containers/images.
imageID
string	ImageID of the container's image.
lastState
ContainerState	Details about the container's last termination condition.
name
string	This must be a DNS_LABEL. Each container in a pod must have a unique name. Cannot be updated.
ready
boolean	Specifies whether the container has passed its readiness probe.
restartCount
integer	The number of times the container has been restarted.
started
boolean	Specifies whether the container has passed its startup probe. Initialized as false, becomes true after startupProbe is considered successful. Resets to false when the container is restarted, or if kubelet loses state temporarily. Is always true when no startupProbe is defined.
state
ContainerState	Details about the container's current condition.
*/
