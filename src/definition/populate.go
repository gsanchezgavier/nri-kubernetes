package definition

import (
	"fmt"

	"github.com/newrelic/infra-integrations-sdk/data/attribute"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
)

// GuessFunc guesses from data.
type GuessFunc func(clusterName, groupLabel, entityID string, groups RawGroups) (string, error)

// PopulateFunc populates raw metric groups using your specs
type PopulateFunc func(RawGroups, SpecGroups) (bool, []error)

// MetricSetManipulator manipulates the MetricSet for a given entity and clusterName
type MetricSetManipulator func(ms *metric.Set, entityMeta *integration.EntityMetadata, clusterName string) error

func populateCluster(i *integration.Integration, clusterName string, k8sVersion fmt.Stringer) error {
	e, err := i.Entity(clusterName, "k8s:cluster")
	if err != nil {
		return err
	}
	ms := e.NewMetricSet("K8sClusterSample")

	e.Inventory.SetItem("cluster", "name", clusterName)
	err = ms.SetMetric("clusterName", clusterName, metric.ATTRIBUTE)
	if err != nil {
		return err
	}

	k8sVersionStr := k8sVersion.String()
	e.Inventory.SetItem("cluster", "k8sVersion", k8sVersionStr)
	return ms.SetMetric("clusterK8sVersion", k8sVersionStr, metric.ATTRIBUTE)
}

// IntegrationPopulator populates an integration with the given metrics and definition.
func IntegrationPopulator(
	i *integration.Integration,
	clusterName string,
	k8sVersion fmt.Stringer,
	msTypeGuesser GuessFunc,
) PopulateFunc {
	return func(groups RawGroups, specs SpecGroups) (bool, []error) {
		var populated bool
		var errs []error
		var msEntityType string
		for groupLabel, entities := range groups {
			for entityID := range entities {

				// Only populate specified groups.
				if _, ok := specs[groupLabel]; !ok {
					continue
				}

				msEntityID := entityID
				if generator := specs[groupLabel].IDGenerator; generator != nil {
					generatedEntityID, err := generator(groupLabel, entityID, groups)
					if err != nil {
						errs = append(errs, fmt.Errorf("error generating entity ID for %s: %s", entityID, err))
						continue
					}
					msEntityID = generatedEntityID
				}

				if generatorType := specs[groupLabel].TypeGenerator; generatorType != nil {
					generatedEntityType, err := generatorType(groupLabel, entityID, groups, clusterName)
					if err != nil {
						errs = append(errs, fmt.Errorf("error generating entity type for %s: %s", entityID, err))
						continue
					}
					msEntityType = generatedEntityType
				}

				e, err := i.Entity(msEntityID, msEntityType)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				// Add entity attributes, which will propagate to all metric.Sets.
				// This was previously (on sdk v2) done by msManipulators.
				e.AddAttributes(
					attribute.Attr("clusterName", clusterName),
					attribute.Attr("displayName", e.Metadata.Name),
				)

				msType, err := msTypeGuesser(clusterName, groupLabel, entityID, groups)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				ms := e.NewMetricSet(msType)

				wasPopulated, populateErrs := metricSetPopulateFunc(ms, groupLabel, entityID)(groups, specs)
				if len(populateErrs) != 0 {
					for _, err := range populateErrs {
						errs = append(errs, fmt.Errorf("error populating metric for entity ID %s: %s", entityID, err))
					}
				}

				if wasPopulated {
					populated = true
				}
			}
		}
		if populated {
			err := populateCluster(i, clusterName, k8sVersion)
			if err != nil {
				errs = append(errs, err)
			}
		}
		return populated, errs
	}
}

func metricSetPopulateFunc(ms *metric.Set, groupLabel, entityID string) PopulateFunc {
	return func(groups RawGroups, specs SpecGroups) (populated bool, errs []error) {
		for _, ex := range specs[groupLabel].Specs {
			val, err := ex.ValueFunc(groupLabel, entityID, groups)
			if err != nil {
				if !ex.Optional {
					errs = append(errs, fmt.Errorf("cannot fetch value for metric %q: %w", ex.Name, err))
				}
				continue
			}

			if multiple, ok := val.(FetchedValues); ok {
				for k, v := range multiple {
					err := ms.SetMetric(k, v, ex.Type)
					if err != nil {
						if !ex.Optional {
							errs = append(errs, fmt.Errorf("cannot set metric %s with value %v in metric set, %s", k, v, err))
						}
						continue
					}

					populated = true
				}
			} else {
				err := ms.SetMetric(ex.Name, val, ex.Type)
				if err != nil {
					if !ex.Optional {
						errs = append(errs, fmt.Errorf("cannot set metric %s with value %v in metric set, %s", ex.Name, val, err))
					}
					continue
				}

				populated = true
			}
		}

		return
	}
}
