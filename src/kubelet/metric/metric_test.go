package metric

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/definition"
	"github.com/stretchr/testify/assert"
)

var responseOKDataSample = `{ "pods": [ { "podRef": { "name": "newrelic-infra-monitoring-pjp0v", "namespace": "kube-system", "uid": "b5a9c98f-d34f-11e7-95fe-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:37Z", "containers": [ { "name": "kube-state-metrics", "startTime": "2017-11-30T09:12:51Z", "cpu": { "time": "2017-11-30T14:48:10Z", "usageNanoCores": 184087, "usageCoreNanoSeconds": 4284675040 }, "memory": { "time": "2017-11-30T14:48:10Z", "usageBytes": 22552576, "workingSetBytes": 15196160, "rssBytes": 7352320, "pageFaults": 4683, "majorPageFaults": 152 }, "rootfs": { "time": "2017-11-30T14:48:10Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 35000320, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 24 }, "logs": { "time": "2017-11-30T14:48:10Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 20480, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null }, { "name": "newrelic-infra", "startTime": "2017-11-30T09:12:44Z", "cpu": { "time": "2017-11-30T14:48:12Z", "usageNanoCores": 13046199, "usageCoreNanoSeconds": 303855795298 }, "memory": { "time": "2017-11-30T14:48:12Z", "usageBytes": 243638272, "workingSetBytes": 38313984, "rssBytes": 15785984, "pageFaults": 10304448, "majorPageFaults": 217 }, "rootfs": { "time": "2017-11-30T14:48:12Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 1305837568, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 52 }, "logs": { "time": "2017-11-30T14:48:12Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 657747968, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null } ], "network": { "time": "2017-11-30T14:48:12Z", "rxBytes": 15741653, "rxErrors": 0, "txBytes": 19551073, "txErrors": 0 }, "volume": [ { "time": "2017-11-30T09:13:29Z", "availableBytes": 1048637440, "capacityBytes": 1048649728, "usedBytes": 12288, "inodesFree": 256009, "inodes": 256018, "inodesUsed": 9, "name": "default-token-7cg8m" } ] }, { "podRef": { "name": "kube-dns-910330662-pflkj", "namespace": "kube-system", "uid": "a6f2130b-a21e-11e7-8db6-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:36Z", "containers": [ { "name": "dnsmasq", "startTime": "2017-11-30T09:12:43Z", "cpu": { "time": "2017-11-30T14:48:07Z", "usageNanoCores": 208374, "usageCoreNanoSeconds": 3653471654 }, "memory": { "time": "2017-11-30T14:48:07Z", "usageBytes": 19812352, "workingSetBytes": 12828672, "rssBytes": 5201920, "pageFaults": 3376, "majorPageFaults": 139 }, "rootfs": { "time": "2017-11-30T14:48:07Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 42041344, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 20 }, "logs": { "time": "2017-11-30T14:48:07Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 20480, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null } ], "network": { "time": "2017-11-30T14:48:07Z", "rxBytes": 14447980, "rxErrors": 0, "txBytes": 15557657, "txErrors": 0 }, "volume": [ { "time": "2017-11-30T09:13:29Z", "availableBytes": 1048637440, "capacityBytes": 1048649728, "usedBytes": 12288, "inodesFree": 256009, "inodes": 256018, "inodesUsed": 9, "name": "default-token-7cg8m" } ] } ] }`
var responseContainerWithTheSameName = `{ "pods": [ { "podRef": { "name": "newrelic-infra-monitoring-pjp0v", "namespace": "kube-system", "uid": "b5a9c98f-d34f-11e7-95fe-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:37Z", "containers": [ { "name": "kube-state-metrics", "startTime": "2017-11-30T09:12:51Z", "cpu": { "time": "2017-11-30T14:48:10Z", "usageNanoCores": 184087, "usageCoreNanoSeconds": 4284675040 }, "memory": { "time": "2017-11-30T14:48:10Z", "usageBytes": 22552576, "workingSetBytes": 15196160, "rssBytes": 7352320, "pageFaults": 4683, "majorPageFaults": 152 }, "rootfs": { "time": "2017-11-30T14:48:10Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 35000320, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 24 }, "logs": { "time": "2017-11-30T14:48:10Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 20480, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null }, { "name": "newrelic-infra", "startTime": "2017-11-30T09:12:44Z", "cpu": { "time": "2017-11-30T14:48:12Z", "usageNanoCores": 13046199, "usageCoreNanoSeconds": 303855795298 }, "memory": { "time": "2017-11-30T14:48:12Z", "usageBytes": 243638272, "workingSetBytes": 38313984, "rssBytes": 15785984, "pageFaults": 10304448, "majorPageFaults": 217 }, "rootfs": { "time": "2017-11-30T14:48:12Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 1305837568, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 52 }, "logs": { "time": "2017-11-30T14:48:12Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 657747968, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null } ], "network": { "time": "2017-11-30T14:48:12Z", "rxBytes": 15741653, "rxErrors": 0, "txBytes": 19551073, "txErrors": 0 }, "volume": [ { "time": "2017-11-30T09:13:29Z", "availableBytes": 1048637440, "capacityBytes": 1048649728, "usedBytes": 12288, "inodesFree": 256009, "inodes": 256018, "inodesUsed": 9, "name": "default-token-7cg8m" } ] }, { "podRef": { "name": "kube-dns-910330662-pflkj", "namespace": "kube-system", "uid": "a6f2130b-a21e-11e7-8db6-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:36Z", "containers": [ { "name": "kube-state-metrics", "startTime": "2017-11-30T09:12:51Z", "cpu": { "time": "2017-11-30T14:48:10Z", "usageNanoCores": 184087, "usageCoreNanoSeconds": 4284675040 }, "memory": { "time": "2017-11-30T14:48:10Z", "usageBytes": 22552576, "workingSetBytes": 15196160, "rssBytes": 7352320, "pageFaults": 4683, "majorPageFaults": 152 }, "rootfs": { "time": "2017-11-30T14:48:10Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 35000320, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 24 }, "logs": { "time": "2017-11-30T14:48:10Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 20480, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null }, { "name": "dnsmasq", "startTime": "2017-11-30T09:12:43Z", "cpu": { "time": "2017-11-30T14:48:07Z", "usageNanoCores": 208374, "usageCoreNanoSeconds": 3653471654 }, "memory": { "time": "2017-11-30T14:48:07Z", "usageBytes": 19812352, "workingSetBytes": 12828672, "rssBytes": 5201920, "pageFaults": 3376, "majorPageFaults": 139 }, "rootfs": { "time": "2017-11-30T14:48:07Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 42041344, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 20 }, "logs": { "time": "2017-11-30T14:48:07Z", "availableBytes": 6911750144, "capacityBytes": 17293533184, "usedBytes": 20480, "inodesFree": 9574871, "inodes": 9732096, "inodesUsed": 157225 }, "userDefinedMetrics": null } ], "network": { "time": "2017-11-30T14:48:07Z", "rxBytes": 14447980, "rxErrors": 0, "txBytes": 15557657, "txErrors": 0 }, "volume": [ { "time": "2017-11-30T09:13:29Z", "availableBytes": 1048637440, "capacityBytes": 1048649728, "usedBytes": 12288, "inodesFree": 256009, "inodes": 256018, "inodesUsed": 9, "name": "default-token-7cg8m" } ] } ] }`
var responseMissingContainerName = `{ "pods": [ { "podRef": { "name": "newrelic-infra-monitoring-pjp0v", "namespace": "kube-system", "uid": "b5a9c98f-d34f-11e7-95fe-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:37Z", "containers": [ { "startTime": "2017-11-30T09:12:51Z", "cpu": { "time": "2017-11-30T14:48:10Z", "usageNanoCores": 184087, "usageCoreNanoSeconds": 4284675040 }, "memory": { "time": "2017-11-30T14:48:10Z", "usageBytes": 22552576, "workingSetBytes": 15196160, "rssBytes": 7352320, "pageFaults": 4683, "majorPageFaults": 152 } } ], "network": { "time": "2017-11-30T14:48:12Z", "rxBytes": 15741653, "txBytes": 52463212, "rxErrors": 0,  "txErrors": 0 } } ] }`
var responseMissingPodName = `{ "pods": [ { "podRef": { "namespace": "kube-system", "uid": "b5a9c98f-d34f-11e7-95fe-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:37Z", "containers": [ { "name": "kube-state-metrics", "startTime": "2017-11-30T09:12:51Z", "cpu": { "time": "2017-11-30T14:48:10Z", "usageNanoCores": 184087, "usageCoreNanoSeconds": 4284675040 }, "memory": { "time": "2017-11-30T14:48:10Z", "usageBytes": 22552576, "workingSetBytes": 15196160, "rssBytes": 7352320, "pageFaults": 4683, "majorPageFaults": 152 } } ], "network": { "time": "2017-11-30T14:48:12Z", "rxBytes": 15741653, "txBytes": 52463212, "rxErrors": 0,  "txErrors": 0 } } ] }`
var responseMissingRxBytesForPod = `{ "pods": [ { "podRef": { "name": "newrelic-infra-monitoring-pjp0v", "namespace": "kube-system", "uid": "b5a9c98f-d34f-11e7-95fe-62d16fb0cc7f" }, "startTime": "2017-11-30T09:12:37Z", "containers": [ { "name": "kube-state-metrics", "startTime": "2017-11-30T09:12:51Z", "cpu": { "time": "2017-11-30T14:48:10Z", "usageNanoCores": 184087, "usageCoreNanoSeconds": 4284675040 }, "memory": { "time": "2017-11-30T14:48:10Z", "usageBytes": 22552576, "workingSetBytes": 15196160, "rssBytes": 7352320, "pageFaults": 4683, "majorPageFaults": 152 } } ], "network": { "time": "2017-11-30T14:48:12Z", "txBytes": 52463212, "rxErrors": 0,  "txErrors": 0 } } ] }`

func toSummary(response string) (*Summary, error) {
	var summary = new(Summary)
	err := json.Unmarshal([]byte(response), summary)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling the response body. Got error: %v", err.Error())
	}
	return summary, nil
}

func TestGroupStatsSummary_CorrectValue(t *testing.T) {
	expectedRawData := definition.RawGroups{
		"pod": {
			"kube-system_newrelic-infra-monitoring-pjp0v": definition.RawMetrics{
				"podName":   "newrelic-infra-monitoring-pjp0v",
				"namespace": "kube-system",
				"rxBytes":   15741653,
				"errors":    0,
				"txBytes":   19551073,
			},
			"kube-system_kube-dns-910330662-pflkj": definition.RawMetrics{
				"podName":   "kube-dns-910330662-pflkj",
				"namespace": "kube-system",
				"rxBytes":   14447980,
				"errors":    0,
				"txBytes":   15557657,
			},
		},
		"container": {
			"kube-system_newrelic-infra-monitoring-pjp0v_kube-state-metrics": definition.RawMetrics{
				"containerName":  "kube-state-metrics",
				"usageBytes":     22552576,
				"usageNanoCores": 184087,
				"podName":        "newrelic-infra-monitoring-pjp0v",
				"namespace":      "kube-system",
			},
			"kube-system_newrelic-infra-monitoring-pjp0v_newrelic-infra": definition.RawMetrics{
				"containerName":  "newrelic-infra",
				"usageBytes":     243638272,
				"usageNanoCores": 13046199,
				"podName":        "newrelic-infra-monitoring-pjp0v",
				"namespace":      "kube-system",
			},
			"kube-system_kube-dns-910330662-pflkj_dnsmasq": definition.RawMetrics{
				"containerName":  "dnsmasq",
				"usageBytes":     19812352,
				"usageNanoCores": 208374,
				"podName":        "kube-dns-910330662-pflkj",
				"namespace":      "kube-system",
			},
		},
	}
	summary, err := toSummary(responseOKDataSample)
	assert.NoError(t, err)

	rawData, errs := GroupStatsSummary(summary)
	assert.Empty(t, errs)
	assert.Equal(t, expectedRawData, rawData)
}

func TestGroupStatsSummary_ContainerWithTheSameName(t *testing.T) {
	expectedRawData := definition.RawGroups{
		"pod": {
			"kube-system_newrelic-infra-monitoring-pjp0v": definition.RawMetrics{
				"podName":   "newrelic-infra-monitoring-pjp0v",
				"namespace": "kube-system",
				"rxBytes":   15741653,
				"errors":    0,
				"txBytes":   19551073,
			},
			"kube-system_kube-dns-910330662-pflkj": definition.RawMetrics{
				"podName":   "kube-dns-910330662-pflkj",
				"namespace": "kube-system",
				"rxBytes":   14447980,
				"errors":    0,
				"txBytes":   15557657,
			},
		},
		"container": {
			"kube-system_newrelic-infra-monitoring-pjp0v_kube-state-metrics": definition.RawMetrics{
				"containerName":  "kube-state-metrics",
				"usageBytes":     22552576,
				"usageNanoCores": 184087,
				"podName":        "newrelic-infra-monitoring-pjp0v",
				"namespace":      "kube-system",
			},
			"kube-system_kube-dns-910330662-pflkj_kube-state-metrics": definition.RawMetrics{
				"containerName":  "kube-state-metrics",
				"usageBytes":     22552576,
				"usageNanoCores": 184087,
				"podName":        "kube-dns-910330662-pflkj",
				"namespace":      "kube-system",
			},
			"kube-system_newrelic-infra-monitoring-pjp0v_newrelic-infra": definition.RawMetrics{
				"containerName":  "newrelic-infra",
				"usageBytes":     243638272,
				"usageNanoCores": 13046199,
				"podName":        "newrelic-infra-monitoring-pjp0v",
				"namespace":      "kube-system",
			},
			"kube-system_kube-dns-910330662-pflkj_dnsmasq": definition.RawMetrics{
				"containerName":  "dnsmasq",
				"usageBytes":     19812352,
				"usageNanoCores": 208374,
				"podName":        "kube-dns-910330662-pflkj",
				"namespace":      "kube-system",
			},
		},
	}
	summary, err := toSummary(responseContainerWithTheSameName)
	assert.NoError(t, err)

	rawData, errs := GroupStatsSummary(summary)
	assert.Empty(t, errs)
	assert.Equal(t, expectedRawData, rawData)
}

func TestGroupStatsSummary_IncompleteStatsSummaryMessage_NoContainerName(t *testing.T) {
	expectedRawData := definition.RawGroups{
		"pod": {
			"kube-system_newrelic-infra-monitoring-pjp0v": definition.RawMetrics{
				"podName":   "newrelic-infra-monitoring-pjp0v",
				"namespace": "kube-system",
				"rxBytes":   15741653,
				"txBytes":   52463212,
				"errors":    0,
			},
		},
		"container": map[string]definition.RawMetrics{},
	}

	summary, err := toSummary(responseMissingContainerName)
	assert.NoError(t, err)

	rawData, errs := GroupStatsSummary(summary)
	assert.Len(t, errs, 1, "Not expected length of errors")
	assert.Equal(t, expectedRawData, rawData)
}

func TestGroupStatsSummary_IncompleteStatsSummaryMessage_NoPodName(t *testing.T) {
	expectedRawData := definition.RawGroups{
		"pod":       map[string]definition.RawMetrics{},
		"container": map[string]definition.RawMetrics{},
	}

	summary, err := toSummary(responseMissingPodName)
	assert.NoError(t, err)

	rawData, errs := GroupStatsSummary(summary)
	assert.Len(t, errs, 1, "Not expected length of errors")
	assert.EqualError(t, errs[0], "empty pod identifier, fetching pod data skipped")
	assert.Len(t, rawData, 2, "Not expected length of rawData for pods and containers")
	assert.Equal(t, expectedRawData, rawData)
	assert.Empty(t, rawData["pod"])
	assert.Empty(t, rawData["container"])
}

func TestGroupStatsSummary_IncompleteStatsSummaryMessage_NoRxBytesForPod_ReportedAsZero(t *testing.T) {
	expectedRawData := definition.RawGroups{
		"pod": {
			"kube-system_newrelic-infra-monitoring-pjp0v": definition.RawMetrics{
				"podName":   "newrelic-infra-monitoring-pjp0v",
				"namespace": "kube-system",
				"rxBytes":   0,
				"errors":    0,
				"txBytes":   52463212,
			},
		},
		"container": {
			"kube-system_newrelic-infra-monitoring-pjp0v_kube-state-metrics": definition.RawMetrics{
				"containerName":  "kube-state-metrics",
				"usageBytes":     22552576,
				"usageNanoCores": 184087,
				"podName":        "newrelic-infra-monitoring-pjp0v",
				"namespace":      "kube-system",
			},
		},
	}

	summary, err := toSummary(responseMissingRxBytesForPod)
	assert.NoError(t, err)

	rawData, errs := GroupStatsSummary(summary)

	assert.Empty(t, errs)
	assert.Equal(t, expectedRawData, rawData)
}

func TestGroupStatsSummary_EmptyStatsSummaryMessage(t *testing.T) {
	var summary = new(Summary)

	rawData, errs := GroupStatsSummary(summary)
	assert.Empty(t, errs)
	assert.Len(t, rawData, 2, "Not expected length of rawData for pods and containers")
	assert.Empty(t, rawData["pod"])
	assert.Empty(t, rawData["container"])
}
