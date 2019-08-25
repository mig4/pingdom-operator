/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package check

import (
	"github.com/go-logr/logr"
	"github.com/russellcardullo/go-pingdom/pingdom"
	observabilityv1alpha1 "gitlab.com/mig4/pingdom-operator/api/v1alpha1"
	"gitlab.com/mig4/pingdom-operator/controllers/resources"
)

/*
Config is a structure holding data needed to create a new ResourceReconciller
for Check objects.
*/
type Config struct {
	Logger   logr.Logger
	PdClient *pingdom.Client
	Check    *observabilityv1alpha1.Check
}

type checkReconciler struct {
	log      logr.Logger
	pdClient *pingdom.Client
	check    *observabilityv1alpha1.Check

	didWork bool
}

/*
New returns a new ResourceReconciller for a Pingdom Check external resource.
*/
func New(config *Config) resources.ResourceReconciler {
	return &checkReconciler{
		log: config.Logger.WithName("resource-reconciler").WithValues(
			"name", config.Check.GetName(),
		),
		pdClient: config.PdClient,
		check:    config.Check,
		didWork:  false,
	}
}
