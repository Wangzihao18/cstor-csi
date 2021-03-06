package iscsi

import (
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	apis "github.com/openebs/cstor-csi/pkg/apis/openebs.io/core/v1alpha1"
	"k8s.io/kubernetes/pkg/util/mount"
)

func getISCSIInfo(vol *apis.CSIVolume) (*iscsiDisk, error) {
	portal := portalMounter(vol.Spec.ISCSI.TargetPortal)
	var portals []string
	portals = append(portals, portal)

	chapDiscovery := false
	chapSession := false

	return &iscsiDisk{
		VolName:       vol.Spec.Volume.Name,
		Portals:       portals,
		Iqn:           vol.Spec.ISCSI.Iqn,
		lun:           vol.Spec.ISCSI.Lun,
		Iface:         vol.Spec.ISCSI.IscsiInterface,
		chapDiscovery: chapDiscovery,
		chapSession:   chapSession,
	}, nil
}

func getISCSIInfoFromPV(req *csi.NodePublishVolumeRequest) (*iscsiDisk, error) {
	volName := req.GetVolumeId()
	tp := req.GetVolumeContext()["targetPortal"]
	iqn := req.GetVolumeContext()["iqn"]
	lun := req.GetVolumeContext()["lun"]
	if tp == "" || iqn == "" || lun == "" {
		return nil, fmt.Errorf("iSCSI target information is missing")
	}

	//portalList := req.GetVolumeContext()["portals"]
	secretParams := req.GetVolumeContext()["secret"]
	secret := parseSecret(secretParams)

	portal := portalMounter(tp)
	var portals []string
	portals = append(portals, portal)

	iface := req.GetVolumeContext()["iscsiInterface"]
	initiatorName := req.GetVolumeContext()["initiatorName"]
	chapDiscovery := false
	if req.GetVolumeContext()["discoveryCHAPAuth"] == "true" {
		chapDiscovery = true
	}

	chapSession := false
	if req.GetVolumeContext()["sessionCHAPAuth"] == "true" {
		chapSession = true
	}

	return &iscsiDisk{
		VolName:       volName,
		Portals:       portals,
		Iqn:           iqn,
		lun:           lun,
		Iface:         iface,
		chapDiscovery: chapDiscovery,
		chapSession:   chapSession,
		secret:        secret,
		InitiatorName: initiatorName}, nil
}

func getISCSIDiskUnmounter(req *csi.NodeUnpublishVolumeRequest) *iscsiDiskUnmounter {
	return &iscsiDiskUnmounter{
		iscsiDisk: &iscsiDisk{
			VolName: req.GetVolumeId(),
		},
		mounter: mount.New(""),
		exec:    mount.NewOsExec(),
	}
}
