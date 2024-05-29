package status

import (
	"encoding/json"

	"github.com/docker/docker/api/types"
)

type ContainerStatus struct {
	Name         string               `json:"name,omitempty"`
	ContainerID  string               `json:"containerID,omitempty"`
	ImageID      string               `json:"imageID,omitempty"`
	Image        string               `json:"image,omitempty"`
	State        types.ContainerState `json:"state,omitempty"`   //LastState    types.ContainerState
	Started      bool                 `json:"started,omitempty"` //Ready bool
	RestartCount int64                `json:"restartCount,omitempty"`
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
func (c *ContainerStatus) JsonMarshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *ContainerStatus) JsonUnmarshal(data []byte) error {
	return json.Unmarshal(data, &c)
}
