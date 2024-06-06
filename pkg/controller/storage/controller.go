package storage

import (
	"MiniK8S/config"
	core "MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/meta"
	"MiniK8S/pkg/api/status"
	"MiniK8S/pkg/api/types"
	"MiniK8S/pkg/apiClient"
	"MiniK8S/pkg/controller/cache"
	"context"
	petname "github.com/dustinkirkland/golang-petname"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

type StorageController struct {
	storageClassClient            *apiClient.Client
	persistentVolumeClient        *apiClient.Client
	persistentVolumeClaimClient   *apiClient.Client
	storageClassInformer          *cache.Informer
	persistentVolumeInformer      *cache.Informer
	persistentVolumeClaimInformer *cache.Informer
	lock                          sync.Mutex
	mountedNFSServers             map[string]string
	//scQueue                       *cache.WorkQueue
	//pvQueue                       *cache.WorkQueue
	//pvcQueue                      *cache.WorkQueue
}

func NewController(si *cache.Informer, pi *cache.Informer, pci *cache.Informer, sc *apiClient.Client, pv *apiClient.Client, pvc *apiClient.Client) *StorageController {
	newSc := &StorageController{
		storageClassClient:            sc,
		persistentVolumeClient:        pv,
		persistentVolumeClaimClient:   pvc,
		storageClassInformer:          si,
		persistentVolumeInformer:      pi,
		persistentVolumeClaimInformer: pci,
		//scQueue:                       cache.NewWorkQueue(),
		//pvQueue:                       cache.NewWorkQueue(),
		//pvcQueue:                      cache.NewWorkQueue(),
		mountedNFSServers: make(map[string]string),
	}
	newSc.storageClassInformer.AddEventHandler(cache.EventHandlerFuncs{
		AddFunc:    newSc.AddStorageClass,
		UpdateFunc: newSc.UpdateStorageClass,
		DeleteFunc: newSc.DeleteStorageClass,
	})
	newSc.persistentVolumeInformer.AddEventHandler(cache.EventHandlerFuncs{
		AddFunc:    newSc.AddPersistentVolume,
		UpdateFunc: newSc.UpdatePersistentVolume,
		DeleteFunc: newSc.DeletePersistentVolume,
	})
	newSc.persistentVolumeClaimInformer.AddEventHandler(cache.EventHandlerFuncs{
		AddFunc:    newSc.AddPersistentVolumeClaim,
		UpdateFunc: newSc.UpdatePersistentVolumeClaim,
		DeleteFunc: newSc.DeletePersistentVolumeClaim,
	})
	return newSc
}

func (sc *StorageController) AddStorageClass(obj interface{}) {
	//scl := obj.(*core.StorageClass)
	//sc.scQueue.Add(scl)
	sc.handleStorageClass(obj.(*core.StorageClass))
}

func (sc *StorageController) UpdateStorageClass(oldObj, newObj interface{}) {
	//scl := newObj.(*core.StorageClass)
	//sc.scQueue.Add(scl)
}

func (sc *StorageController) DeleteStorageClass(obj interface{}) {}

func (sc *StorageController) AddPersistentVolume(obj interface{}) {
	//pv := obj.(*core.PersistentVolume)
	//sc.pvQueue.Add(pv)
	sc.handlePersistentVolume(obj.(*core.PersistentVolume))
}

func (sc *StorageController) UpdatePersistentVolume(oldObj, newObj interface{}) {
	//pv := newObj.(*core.PersistentVolume)
	//sc.pvQueue.Add(pv)
}

func (sc *StorageController) DeletePersistentVolume(obj interface{}) {}

func (sc *StorageController) AddPersistentVolumeClaim(obj interface{}) {
	//pvc := obj.(*core.PersistentVolumeClaim)
	//sc.pvcQueue.Add(pvc)
	sc.handlePersistentVolumeClaim(obj.(*core.PersistentVolumeClaim))
}

func (sc *StorageController) UpdatePersistentVolumeClaim(oldObj, newObj interface{}) {
	//pvc := newObj.(*core.PersistentVolumeClaim)
	//sc.pvcQueue.Add(pvc)
}

func (sc *StorageController) DeletePersistentVolumeClaim(obj interface{}) {}

func (sc *StorageController) handleStorageClass(scl *core.StorageClass) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	log.Println("[controller] Handling StorageClass, UUID:", scl.Metadata.Uid.String(), "Name:", scl.Metadata.Name)
	if scl.Spec.Provisioner != core.StorageClassNFSProvisioner {
		log.Println("[controller] StorageClass provisioner is not NFS")
		return
	}
	parentPath := "/mnt/minik8s-storage-mnt/"
	_, err := os.Stat(parentPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(parentPath, 0777)
			if err != nil {
				log.Println("[controller] Unable to create parent directory for storage class")
				return
			}
			log.Println("[controller] Created parent directory for storage class")
		} else {
			log.Println("[controller] Unable to get parent directory for storage class")
			return
		}
	}
	path := parentPath + scl.Metadata.Name + "_" + scl.Metadata.Uid.String()
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(path, 0777)
			if err != nil {
				log.Println("[controller] Unable to create directory for storage class")
				return
			}
			// Mount the server to the path
			remotePath := scl.Spec.Parameters.Server + ":" + scl.Spec.Parameters.Path
			log.Printf("[controller] Created directory for storage class %s, mounting %s to %s. \n", scl.Metadata.Name, remotePath, path)
			err = syscall.Mount(remotePath, path, "nfs", 0, "")
			if err != nil {
				log.Printf("[controller] Unable to mount NFS server to path, Error: %v\n", err)
				return
			}
		} else {
			log.Println("[controller] Unable to get directory for storage class")
			return
		}
	}
	if !info.IsDir() {
		log.Println("[controller] Path is not a directory")
		return
	}
	sc.mountedNFSServers[scl.Spec.Parameters.Server] = path
}

func (sc *StorageController) handlePersistentVolumeClaim(pvc *core.PersistentVolumeClaim) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	log.Println("[controller] Handling PersistentVolumeClaim, UUID:", pvc.Metadata.Uid.String(), "Name:", pvc.Metadata.Name)
	if pvc.Spec.StorageClassName == "" {
		log.Println("[controller] StorageClass name not specified")
		return
	}
	allPersistentVolumes, err := sc.persistentVolumeClient.GetAll()
	if err != nil {
		log.Println("[controller] Unable to get all persistent volumes")
		return
	}
	var selectedPv *core.PersistentVolume = nil
	for _, object := range allPersistentVolumes.GetItems() {
		pv := object.(*core.PersistentVolume)
		if pv.Status.Phase == status.PersistentVolumeAvailable {
			selectedPv = pv
			break
		}
	}
	if selectedPv != nil {
		selectedPv.Spec.ClaimRef = meta.OwnerReference{
			Name:       pvc.GetName(),
			UID:        pvc.GetUID(),
			APIGroup:   "",
			Kind:       string(types.PersistentVolumeClaimObjectType),
			Controller: false,
		}
		selectedPv.Status.Phase = status.PersistentVolumeBound
		pvc.Spec.VolumeName = selectedPv.GetName()
		pvc.Status.Phase = status.PersistentVolumeClaimBound
		code, err := sc.persistentVolumeClient.PutObject(selectedPv.GetName(), selectedPv)
		if err != nil {
			log.Printf("[controller] Unable to update persistent volume, error: %v, status code: %v\n", err, code)
			return
		}
		code, err = sc.persistentVolumeClaimClient.PutObject(pvc.GetName(), pvc)
		if err != nil {
			log.Printf("[controller] Unable to update persistent volume claim, error: %v, status code: %v\n", err, code)
			return
		}
		log.Printf("[controller] PersistentVolumeClaim %s bound to PersistentVolume %s\n", pvc.GetName(), selectedPv.GetName())
		return
	}
	log.Println("[controller] No currently available persistent volume in storage class, provisioning")
	// Provision a new persistent volume
	code, storageClassObj, err := sc.storageClassClient.GetObject(pvc.Spec.StorageClassName)
	if code != http.StatusOK {
		log.Printf("[controller] Unable to get storage class, status code: %v, error: %v\n", code, err)
		return
	}
	if err != nil {
		log.Printf("[controller] Unable to get storage class, error: %v\n", err)
		return
	}
	storageClass := storageClassObj.(*core.StorageClass)
	if storageClass.Spec.Provisioner != core.StorageClassNFSProvisioner {
		log.Println("[controller] StorageClass provisioner is not NFS, currently not supported")
		return
	}
	pvName := petname.Generate(2, "_")
	pv := &core.PersistentVolume{
		ApiVersion: config.Version(),
		Kind:       string(types.PersistentVolumeObjectType),
		Metadata: meta.ObjectMeta{
			Name:      "provisioned_volume_" + storageClass.GetName() + "_" + pvName,
			Namespace: "",
		},
		Spec: core.PersistentVolumeSpec{
			ClaimRef: meta.OwnerReference{
				Name:     pvc.GetName(),
				UID:      pvc.GetUID(),
				APIGroup: "",
				Kind:     string(types.PersistentVolumeClaimObjectType),
			},
			StorageClassName: storageClass.GetName(),
			Nfs: core.NFSVolumeSource{
				Server: storageClass.Spec.Parameters.Server,
				Path:   storageClass.Spec.Parameters.Path + pvName,
			},
		},
		Status: status.PersistentVolumeStatus{
			Phase: status.PersistentVolumeBound,
		},
	}
	pvc.Spec.VolumeName = pv.GetName()
	pvc.Status.Phase = status.PersistentVolumeClaimBound
	code, err = sc.persistentVolumeClient.PutObject(pv.GetName(), pv)
	if err != nil {
		log.Printf("[controller] Unable to put persistent volume, error: %v, status code: %v\n", err, code)
		return
	}
	code, err = sc.persistentVolumeClaimClient.PutObject(pvc.GetName(), pvc)
	if err != nil {
		log.Printf("[controller] Unable to update persistent volume claim, error: %v, status code: %v\n", err, code)
	}
}

func (sc *StorageController) handlePersistentVolume(pv *core.PersistentVolume) {
	sc.lock.Lock()
	defer sc.lock.Unlock()
	log.Println("[controller] Handling PersistentVolume, UUID:", pv.Metadata.Uid.String(), "Name:", pv.Metadata.Name)
	path, ok := sc.mountedNFSServers[pv.Spec.Nfs.Server]
	if !ok {
		log.Println("[controller] NFS server not mounted")
		return
	}
	volumePath := filepath.Join(path, pv.Spec.Nfs.Path)
	info, err := os.Stat(volumePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("[controller] Creating directory for persistent volume %s\n", pv.Metadata.Name)
			mkdirErr := os.MkdirAll(filepath.Join(path, pv.Spec.Nfs.Path), 0777)
			if mkdirErr != nil {
				log.Printf("[controller] Unable to create directory for persistent volume %s\n", pv.Metadata.Name)
				return
			}
		} else {
			log.Printf("[controller] Unable to get directory for persistent volume %s\n", pv.Metadata.Name)
			return
		}
	}
	if !info.IsDir() {
		log.Println("[controller] Path is not a directory")
		return
	}
}

func (sc *StorageController) init() {
	log.Println("[controller] Initializing StorageController")
	// Get all storage classes
	allStorageClasses, err := sc.storageClassClient.GetAll()
	if err != nil {
		log.Printf("[controller] Unable to get all storage classes, error: %v\n", err)
		return
	}
	for _, object := range allStorageClasses.GetItems() {
		sc.handleStorageClass(object.(*core.StorageClass))
	}
}

func (sc *StorageController) finalize() {
	for server, path := range sc.mountedNFSServers {
		log.Printf("[controller] Unmounting NFS server %s from path %s\n", server, path)
		err := syscall.Unmount(path, 0)
		if err != nil {
			log.Printf("[controller] Unable to unmount NFS server %s from path %s\n", server, path)
		}
	}
}

func (sc *StorageController) Run(ctx context.Context, cancel context.CancelFunc) {
	log.Println("[controller] Starting StorageController")
	sc.init()
	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				sc.finalize()
				return
			}
		}
	}()
}
