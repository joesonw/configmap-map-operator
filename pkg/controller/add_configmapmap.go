package controller

import (
	"github.com/dstreamcloud/configmap-map-operator/pkg/controller/configmapmap"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, configmapmap.Add)
}
