package handlers

import (
	"MiniK8S/pkg/api/address"
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/status"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/etcd"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io"
	"net/http"
)

func HandleGetApiObjects(c *gin.Context, ty types.ApiObjectType) {
	var etcdPath string
	switch ty {
	case types.PodObjectType:
		etcdPath = "/api/pods/"
	case types.NodeObjectType:
		etcdPath = "/api/nodes/"
	case types.ServiceObjectType:
		etcdPath = "/api/services/"
	case types.DeploymentObjectType:
		etcdPath = "/api/deployments/"
	case types.HorizontalPodAutoscalerObjectType:
		etcdPath = "/api/hpas/"
	}
	buf, err := etcd.GetAllWithPrefix(etcdPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
	} else {
		apiObjectList := config.NewApiObjectList(ty)
		err := apiObjectList.AppendItems(buf)
		if err != nil {
			fmt.Printf("[apiserver] Error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, apiObjectList)
	}
}

func HandleGetApiObject(c *gin.Context, ty types.ApiObjectType) {
	UID := c.Param("name")
	var etcdPath string
	switch ty {
	case types.PodObjectType:
		etcdPath = "/api/pods/"
	case types.NodeObjectType:
		etcdPath = "/api/nodes/"
	case types.ServiceObjectType:
		etcdPath = "/api/services/"
	case types.DeploymentObjectType:
		etcdPath = "/api/deployments/"
	case types.HorizontalPodAutoscalerObjectType:
		etcdPath = "/api/hpas/"
	}
	etcdPath += UID
	buf, err := etcd.Get(etcdPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
	} else if buf == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": "ERR", "error": fmt.Sprintf("No %v with UID: %v", ty, UID)})
	} else {
		apiObject := config.NewApiObject(ty)
		err := apiObject.JsonUnmarshal([]byte(buf))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		} else {
			c.JSON(http.StatusOK, apiObject)
		}
	}
}

func HandleCreateApiObject(c *gin.Context, ty types.ApiObjectType) {

	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
	}
	newApiObject := config.NewApiObject(ty)
	err = newApiObject.JsonUnmarshal(buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
	UID := uuid.New()
	newApiObject.SetUID(UID)
	if ty == types.PodObjectType {
		pod := newApiObject.(*config.Pod)
		pod.Status = status.PodStatus{
			ContainerStatuses: nil,
			HostIP:            "",
			Phase:             "Pending",
			PodIP:             "",
		}
	}
	if ty == types.NodeObjectType {
		node := newApiObject.(*config.Node)
		node.Status = status.NodeStatus{
			Addresses: address.NodeAddress{
				Type:    "",
				Address: "",
			},
			DaemonEndpoints: 0,
			Phase:           "Pending",
		}
	}
	etcd.VersionLock.Lock()
	defer etcd.VersionLock.Unlock()

	version := etcd.RVM.GetResourceVersion()
	newApiObject.SetResourceVersion(etcd.RVM.GetNextResourceVersion())
	buf, err = newApiObject.JsonMarshal()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
	var etcdPath string
	switch ty {
	case types.PodObjectType:
		etcdPath = "/api/pods/"
	case types.NodeObjectType:
		etcdPath = "/api/nodes/"
	case types.ServiceObjectType:
		etcdPath = "/api/services/"
	case types.DeploymentObjectType:
		etcdPath = "/api/deployments/"
	case types.HorizontalPodAutoscalerObjectType:
		etcdPath = "/api/hpas/"
	}
	etcdPath += UID.String()
	newVersion, err := etcd.Put(etcdPath, string(buf))
	fmt.Printf("[apiServer] generate new %v: expected version:%v, actual version:%v\n", ty, version+1, newVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "uid": UID})
	}
}

func HandleModifyApiObject(c *gin.Context, ty types.ApiObjectType) {
	UID := c.Param("name")
	var etcdPath string
	switch ty {
	case types.PodObjectType:
		etcdPath = "/api/pods/"
	case types.NodeObjectType:
		etcdPath = "/api/nodes/"
	case types.ServiceObjectType:
		etcdPath = "/api/services/"
	case types.DeploymentObjectType:
		etcdPath = "/api/deployments/"
	case types.HorizontalPodAutoscalerObjectType:
		etcdPath = "/api/hpas/"
	}
	etcdPath += UID
	exist, version, err := etcd.ExistWithVersion(etcdPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
	if !exist {
		c.JSON(http.StatusNotFound, gin.H{"status": "ERR", "error": fmt.Sprintf("No %v with UID: %v", ty, UID)})
		return
	}
	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
	}
	newApiObject := config.NewApiObject(ty)
	err = newApiObject.JsonUnmarshal(buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
	clientVersion := newApiObject.GetResourceVersion()
	if version != clientVersion {
		c.JSON(http.StatusConflict, gin.H{"status": "FAILED", "error": fmt.Sprintf("client version %v unmatch host version %v", clientVersion, version)})
		return
	}
	etcd.VersionLock.Lock()
	defer etcd.VersionLock.Unlock()
	newApiObject.SetResourceVersion(etcd.RVM.GetNextResourceVersion())
	buf, err = newApiObject.JsonMarshal()
	if err != err {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
	}
	newVersion, success, err := etcd.PutWithVersion(etcdPath, string(buf), version)
	if !success {
		c.JSON(http.StatusConflict, gin.H{"status": "FAILED", "error": fmt.Sprintf("client version %v unmatch host version %v", clientVersion, version)})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "uid": UID, "resourceVersion": newVersion})
	}
}

func HandleDeleteApiObjects(c *gin.Context, ty types.ApiObjectType) {
	var etcdPath string
	switch ty {
	case types.PodObjectType:
		etcdPath = "/api/pods/"
	case types.NodeObjectType:
		etcdPath = "/api/nodes/"
	case types.ServiceObjectType:
		etcdPath = "/api/services/"
	case types.DeploymentObjectType:
		etcdPath = "/api/deployments/"
	case types.HorizontalPodAutoscalerObjectType:
		etcdPath = "/api/hpas/"
	}
	err := etcd.DeleteAllWithPrefix(etcdPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	}
}

func HandleDeleteApiObject(c *gin.Context, ty types.ApiObjectType) {
	UID := c.Param("name")
	var etcdPath string
	switch ty {
	case types.PodObjectType:
		etcdPath = "/api/pods/"
	case types.NodeObjectType:
		etcdPath = "/api/nodes/"
	case types.ServiceObjectType:
		etcdPath = "/api/services/"
	case types.DeploymentObjectType:
		etcdPath = "/api/deployments/"
	case types.HorizontalPodAutoscalerObjectType:
		etcdPath = "/api/hpas/"
	}
	etcdPath += UID
	exist, err := etcd.Exist(etcdPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
	if !exist {
		c.JSON(http.StatusNotFound, gin.H{"status": "ERR", "error": fmt.Sprintf("No %v with UID: %v", ty, UID)})
	}
	err = etcd.Delete(etcdPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "uid": UID})
	}
}

func HandleWatchApiObjects(c *gin.Context, ty types.ApiObjectType) {
	var etcdPath string
	switch ty {
	case types.PodObjectType:
		etcdPath = "/api/pods/"
	case types.NodeObjectType:
		etcdPath = "/api/nodes/"
	case types.ServiceObjectType:
		etcdPath = "/api/services/"
	case types.DeploymentObjectType:
		etcdPath = "/api/deployments/"
	case types.HorizontalPodAutoscalerObjectType:
		etcdPath = "/api/hpas/"
	}
	fmt.Printf("[apiServer]start watch,type: %v etcdPath: %v\n", ty, etcdPath)
	cancel, ch := etcd.WatchAllWithPrefix(etcdPath)
	flusher, _ := c.Writer.(http.Flusher)
	for {
		select {
		case ev := <-ch:
			event, err := json.Marshal(ev)
			if err != nil {
				defer cancel()
				fmt.Printf("[apiServer] Error: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
				return
			}
			_, err = fmt.Fprintf(c.Writer, "%v\n", string(event))
			if err != nil {
				defer cancel()
				fmt.Printf("[apiServer] Error: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
				return
			}
			switch ev.Type {
			case clientv3.EventTypeDelete:
				fmt.Printf("[apiServer] apiObjest delete, type: %v\n", ty)
			case clientv3.EventTypePut:
				fmt.Printf("[apiServer] apiObject put, type: %v\n", ty)
			}
			flusher.Flush()
		case <-c.Request.Context().Done():
			defer cancel()
			fmt.Printf("[apiServer] watch connection closed\n")
			c.JSON(http.StatusOK, gin.H{"status": "OK"})
			return
		}
	}

}

func HandleWatchApiObject(c *gin.Context, ty types.ApiObjectType) {
	UID := c.Param("name")
	var etcdPath string
	switch ty {
	case types.PodObjectType:
		etcdPath = "/api/pods/"
	case types.NodeObjectType:
		etcdPath = "/api/nodes/"
	case types.ServiceObjectType:
		etcdPath = "/api/services/"
	case types.DeploymentObjectType:
		etcdPath = "/api/deployments/"
	case types.HorizontalPodAutoscalerObjectType:
		etcdPath = "/api/hpas/"
	}
	etcdPath += UID
	exixt, err := etcd.Exist(etcdPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
	if !exixt {
		c.JSON(http.StatusNotFound, gin.H{"status": "ERR", "error": fmt.Sprintf("No %v with UID: %v", ty, UID)})
		return
	}
	fmt.Printf("[apiServer]start watch,type: %v etcdPath: %v\n", ty, etcdPath)
	cancel, ch := etcd.WatchAllWithPrefix(etcdPath)
	flusher, _ := c.Writer.(http.Flusher)
	for {
		select {
		case ev := <-ch:
			event, err := json.Marshal(ev)
			if err != nil {
				defer cancel()
				fmt.Printf("[apiServer] Error: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
				return
			}
			_, err = fmt.Fprintf(c.Writer, "%v\n", string(event))
			if err != nil {
				defer cancel()
				fmt.Printf("[apiServer] Error: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
				return
			}
			switch ev.Type {
			case clientv3.EventTypeDelete:
				defer cancel()
				flusher.Flush()
				fmt.Printf("[apiServer] apiObjest delete, type: %v UID: %v\n", ty, UID)
				c.JSON(http.StatusOK, gin.H{"status": "OK"})
				return
			case clientv3.EventTypePut:
				fmt.Printf("[apiServer] apiObject put, type: %v UID: %v\n", ty, UID)
				flusher.Flush()
				// c.JSON(http.StatusOK, gin.H{"status": "OK"})
			}
		case <-c.Request.Context().Done():
			defer cancel()
			fmt.Printf("[apiServer] watch connection closed\n")
			c.JSON(http.StatusOK, gin.H{"status": "OK"})
			return
		}
	}
}

func HandleGetApiObjectStatus(c *gin.Context, ty types.ApiObjectType) {
	UID := c.Param("name")
	var etcdPath string
	switch ty {
	case types.PodObjectType:
		etcdPath = "/api/pods/"
	case types.NodeObjectType:
		etcdPath = "/api/nodes/"
	case types.ServiceObjectType:
		etcdPath = "/api/services/"
	case types.DeploymentObjectType:
		etcdPath = "/api/deployments/"
	case types.HorizontalPodAutoscalerObjectType:
		etcdPath = "/api/hpas/"
	}
	etcdPath += UID
	buf, err := etcd.Get(etcdPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
	} else if buf == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": "ERR", "error": fmt.Sprintf("No %v with UID: %v", ty, UID)})
	} else {
		apiObeject := config.NewApiObject(ty)
		err := apiObeject.JsonUnmarshal([]byte(buf))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, apiObeject.GetStatus())
	}
}

func HandleModifyApiObjectStatus(c *gin.Context, ty types.ApiObjectType) {
	UID := c.Param("name")
	var etcdPath string
	switch ty {
	case types.PodObjectType:
		etcdPath = "/api/pods/"
	case types.NodeObjectType:
		etcdPath = "/api/nodes/"
	case types.ServiceObjectType:
		etcdPath = "/api/services/"
	case types.DeploymentObjectType:
		etcdPath = "/api/deployments/"
	case types.HorizontalPodAutoscalerObjectType:
		etcdPath = "/api/hpas/"
	}
	etcdPath += UID
	exixt, version, err := etcd.ExistWithVersion(etcdPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
	if !exixt {
		c.JSON(http.StatusNotFound, gin.H{"status": "ERR", "error": fmt.Sprintf("No %v with UID: %v", ty, UID)})
		return
	}
	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
	apiObejectStatus := config.NewApiObjectStatus(ty)
	err = apiObejectStatus.JsonUnmarshal(buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
	apiObjectJson, err := etcd.Get(etcdPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
	apiObject := config.NewApiObject(ty)
	err = apiObject.JsonUnmarshal([]byte(apiObjectJson))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
	success := apiObject.SetStatus(apiObejectStatus)
	if !success {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": fmt.Sprintf("set status error, type: %v UID: %v", ty, UID)})
		return
	}
	clientVersion := apiObject.GetResourceVersion()
	if version != clientVersion {
		c.JSON(http.StatusConflict, gin.H{"status": "FAILED", "error": fmt.Sprintf("client version %v unmatch host version %v", clientVersion, version)})
		return
	}
	etcd.VersionLock.Lock()
	defer etcd.VersionLock.Unlock()

	apiObject.SetResourceVersion(etcd.RVM.GetNextResourceVersion())
	buf, err = apiObject.JsonMarshal()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
		return
	}
	newVersion, success, err := etcd.PutWithVersion(etcdPath, string(buf), version)
	if !success {
		c.JSON(http.StatusConflict, gin.H{"status": "FAILED", "error": fmt.Sprintf("client version %v unmatch host version %v", clientVersion, version)})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERR", "error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "uid": UID, "resourceVersion": newVersion})
	}
}
