// Package node provides the logic for mapping a Kubernetes Node to a
// LogicMonitor device.
package node

import (
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	lm "github.com/logicmonitor/lm-sdk-go"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	resource = "nodes"
)

// Watcher represents a watcher type that watches nodes.
type Watcher struct {
	types.DeviceManager
	DeviceGroups map[string]int32
	LMClient     *lm.DefaultApi
}

// Resource is a function that implements the Watcher interface.
func (w *Watcher) Resource() string {
	return resource
}

// ObjType is a function that implements the Watcher interface.
func (w *Watcher) ObjType() runtime.Object {
	return &v1.Node{}
}

// AddFunc is a function that implements the Watcher interface.
func (w *Watcher) AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		node := obj.(*v1.Node)

		log.Debugf("received ADD event: %s", node.Name)

		// Require an IP address.
		if getInternalAddress(node.Status.Addresses) == nil {
			return
		}
		w.add(node)
	}
}

// UpdateFunc is a function that implements the Watcher interface.
func (w *Watcher) UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		old := oldObj.(*v1.Node)
		new := newObj.(*v1.Node)

		log.Debugf("received UPDATE event: %s", old.Name)

		// If the old node does not have an IP, then there is no way we could
		// have added it to LogicMonitor. Therefore, it must be a new device.
		oldInternalAddress := getInternalAddress(old.Status.Addresses)
		newInternalAddress := getInternalAddress(new.Status.Addresses)
		if oldInternalAddress == nil && newInternalAddress != nil {
			w.add(new)
			return
		}
		// Covers the case when the old node is in the process of terminating
		// and the new node is coming up to replace it.
		if oldInternalAddress.Address != newInternalAddress.Address {
			w.update(old, new)
		}
	}
}

// DeleteFunc is a function that implements the Watcher interface.
// nolint: dupl
func (w *Watcher) DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		node := obj.(*v1.Node)

		log.Debugf("received DELETE event: %s", node.Name)

		// Delete the node.
		internalAddress := getInternalAddress(node.Status.Addresses).Address
		if w.Config().DeleteDevices {
			if err := w.DeleteByName(internalAddress); err != nil {
				log.Errorf("Failed to delete node: %v", err)
				return
			}
			log.Infof("Deleted node %s", internalAddress)
			return
		}

		// Move the node.
		w.move(node)
	}
}

// nolint: dupl
func (w *Watcher) add(node *v1.Node) {
	if _, err := w.Add(w.args(node, constants.NodeCategory)...); err != nil {
		log.Errorf("Failed to add node %q: %v", node.Name, err)
	} else {
		log.Infof("Added node %q", node.Name)
	}

	w.createRoleDeviceGroup(node.Labels)
}

func (w *Watcher) update(old, new *v1.Node) {
	if _, err := w.UpdateAndReplaceByName(old.Name, w.args(new, constants.NodeCategory)...); err != nil {
		log.Errorf("Failed to update node %q: %v", new.Name, err)
	} else {
		log.Infof("Updated node %q", old.Name)
	}

	// determine if we need to add a new node role device group
	oldLabel, _ := utilities.GetLabelByPrefix(constants.LabelNodeRole, old.Labels)
	newLabel, _ := utilities.GetLabelByPrefix(constants.LabelNodeRole, new.Labels)
	if oldLabel != newLabel {
		w.createRoleDeviceGroup(new.Labels)
	}
}

// nolint: dupl
func (w *Watcher) move(node *v1.Node) {
	if _, err := w.UpdateAndReplaceFieldByName(node.Name, constants.CustomPropertiesFieldName, w.args(node, constants.NodeDeletedCategory)...); err != nil {
		log.Errorf("Failed to move node %q: %v", node.Name, err)
		return
	}
	log.Infof("Moved node %q", node.Name)
}

func (w *Watcher) args(node *v1.Node, category string) []types.DeviceOption {
	categories := utilities.BuildSystemCategoriesFromLabels(category, node.Labels)

	return []types.DeviceOption{
		w.Name(getInternalAddress(node.Status.Addresses).Address),
		w.ResourceLabels(node.Labels),
		w.DisplayName(node.Name),
		w.SystemCategories(categories),
		w.Auto("name", node.Name),
		w.Auto("selflink", node.SelfLink),
		w.Auto("uid", string(node.UID)),
	}
}

// getInternalAddress finds the node's internal address.
func getInternalAddress(addresses []v1.NodeAddress) *v1.NodeAddress {
	for _, address := range addresses {
		if address.Type == v1.NodeInternalIP {
			return &address
		}
	}

	return nil
}

func (w *Watcher) createRoleDeviceGroup(labels map[string]string) {
	label, _ := utilities.GetLabelByPrefix(constants.LabelNodeRole, labels)
	if label == "" {
		return
	}
	role := strings.Replace(label, constants.LabelNodeRole, "", -1)

	if devicegroup.Exists(w.DeviceGroups[constants.ClusterDeviceGroupPrefix+w.Config().ClusterName], role, w.LMClient) {
		log.Infof("Device group for node role %q already exists", role)
		return
	}

	opts := &devicegroup.Options{
		ParentID:              w.DeviceGroups[constants.NodeDeviceGroupName],
		Name:                  role,
		DisableAlerting:       true,
		AppliesTo:             devicegroup.NewAppliesToBuilder().HasCategory(label + "=").And().Auto("clustername").Equals(w.Config().ClusterName),
		Client:                w.LMClient,
		DeleteDevices:         w.Config().DeleteDevices,
		AppliesToDeletedGroup: devicegroup.NewAppliesToBuilder().HasCategory(label + "=").And().Auto("clustername").Equals(w.Config().ClusterName),
	}

	log.Debugf("%v", opts)

	_, err := devicegroup.Create(opts)
	if err != nil {
		log.Errorf("Failed to add device group for node role to %q: %v", role, err)
		return
	}

	log.Printf("Added device group for node role %q", role)
}
