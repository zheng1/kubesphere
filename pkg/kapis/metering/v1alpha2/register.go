/*

 Copyright 2019 The KubeSphere Authors.

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
package v1alpha2

import (
	"net/http"

	"kubesphere.io/kubesphere/pkg/models/monitoring"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"kubesphere.io/kubesphere/pkg/apiserver/metering"
	"kubesphere.io/kubesphere/pkg/apiserver/runtime"
	"kubesphere.io/kubesphere/pkg/constants"
)

const (
	GroupName = "metering.kubesphere.io"
	RespOK    = "ok"
)

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha2"}

func AddToContainer(c *restful.Container) error {
	ws := runtime.NewWebService(GroupVersion)

	ws.Route(ws.GET("/cluster").To(metering.MonitorCluster).
		Doc("Get cluster-level metric data.").
		Param(ws.QueryParameter("metrics_filter", "The metric name filter consists of a regexp pattern. It specifies which metric data to return. For example, the following filter matches both cluster CPU usage and disk usage: `cluster_cpu_usage|cluster_disk_size_usage`. View available metrics at [kubesphere.io](https://docs.kubesphere.io/advanced-v2.0/zh-CN/api-reference/metering-metrics/).").DataType("string").Required(false)).
		Param(ws.QueryParameter("start", "Start time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1559347200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("end", "End time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1561939200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("step", "Time interval. Retrieve metric data at a fixed interval within the time range of start and end. It requires both **start** and **end** are provided. The format is [0-9]+[smhdwy]. Defaults to 10m (i.e. 10 min).").DataType("string").DefaultValue("10m").Required(false)).
		Param(ws.QueryParameter("time", "A timestamp in Unix time format. Retrieve metric data at a single point in time. Defaults to now. Time and the combination of start, end, step are mutually exclusive.").DataType("string").Required(false)).
		Param(ws.QueryParameter("type", "Additional operations. Currently available types is statistics. It retrieves the total number of workspaces, devops projects, namespaces, accounts in the cluster at the moment.").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.ClusterMetricsTag}).
		Writes(monitoring.Metrics{}).
		Returns(http.StatusOK, RespOK, monitoring.Metrics{})).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/workspaces").To(metering.MonitorWorkspace).
		Doc("Get workspace-level metric data of all workspaces.").
		Param(ws.QueryParameter("metrics_filter", "The metric name filter consists of a regexp pattern. It specifies which metric data to return. For example, the following filter matches both workspace CPU usage and memory usage: `workspace_cpu_usage|workspace_memory_usage`. View available metrics at [kubesphere.io](https://docs.kubesphere.io/advanced-v2.0/zh-CN/api-reference/metering-metrics/).").DataType("string").Required(false)).
		Param(ws.QueryParameter("resources_filter", "The workspace filter consists of a regexp pattern. It specifies which workspace data to return.").DataType("string").Required(false)).
		Param(ws.QueryParameter("start", "Start time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1559347200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("end", "End time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1561939200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("step", "Time interval. Retrieve metric data at a fixed interval within the time range of start and end. It requires both **start** and **end** are provided. The format is [0-9]+[smhdwy]. Defaults to 10m (i.e. 10 min).").DataType("string").DefaultValue("10m").Required(false)).
		Param(ws.QueryParameter("time", "A timestamp in Unix time format. Retrieve metric data at a single point in time. Defaults to now. Time and the combination of start, end, step are mutually exclusive.").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_metric", "Sort workspaces by the specified metric. Not applicable if **start** and **end** are provided.").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_type", "Sort order. One of asc, desc.").DefaultValue("desc.").DataType("string").Required(false)).
		Param(ws.QueryParameter("page", "The page number. This field paginates result data of each metric, then returns a specific page. For example, setting **page** to 2 returns the second page. It only applies to sorted metric data.").DataType("integer").Required(false)).
		Param(ws.QueryParameter("limit", "Page size, the maximum number of results in a single page. Defaults to 5.").DataType("integer").Required(false).DefaultValue("5")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.WorkspaceMetricsTag}).
		Writes(monitoring.Metrics{}).
		Returns(http.StatusOK, RespOK, monitoring.Metrics{})).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/workspaces/{workspace}").To(metering.MonitorWorkspace).
		Doc("Get workspace-level metric data of a specific workspace.").
		Param(ws.PathParameter("workspace", "Workspace name.").DataType("string").Required(true)).
		Param(ws.QueryParameter("metrics_filter", "The metric name filter consists of a regexp pattern. It specifies which metric data to return. For example, the following filter matches both workspace CPU usage and memory usage: `workspace_cpu_usage|workspace_memory_usage`. View available metrics at [kubesphere.io](https://docs.kubesphere.io/advanced-v2.0/zh-CN/api-reference/metering-metrics/).").DataType("string").Required(false)).
		Param(ws.QueryParameter("start", "Start time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1559347200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("end", "End time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1561939200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("step", "Time interval. Retrieve metric data at a fixed interval within the time range of start and end. It requires both **start** and **end** are provided. The format is [0-9]+[smhdwy]. Defaults to 10m (i.e. 10 min).").DataType("string").DefaultValue("10m").Required(false)).
		Param(ws.QueryParameter("time", "A timestamp in Unix time format. Retrieve metric data at a single point in time. Defaults to now. Time and the combination of start, end, step are mutually exclusive.").DataType("string").Required(false)).
		Param(ws.QueryParameter("type", "Additional operations. Currently available types is statistics. It retrieves the total number of namespaces, devops projects, members and roles in this workspace at the moment.").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.WorkspaceMetricsTag}).
		Writes(monitoring.Metrics{}).
		Returns(http.StatusOK, RespOK, monitoring.Metrics{})).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/workspaces/{workspace}/namespaces").To(metering.MonitorNamespace).
		Doc("Get namespace-level metric data of a specific workspace.").
		Param(ws.PathParameter("workspace", "Workspace name.").DataType("string").Required(true)).
		Param(ws.QueryParameter("metrics_filter", "The metric name filter consists of a regexp pattern. It specifies which metric data to return. For example, the following filter matches both namespace CPU usage and memory usage: `namespace_cpu_usage|namespace_memory_usage`. View available metrics at [kubesphere.io](https://docs.kubesphere.io/advanced-v2.0/zh-CN/api-reference/metering-metrics/).").DataType("string").Required(false)).
		Param(ws.QueryParameter("resources_filter", "The namespace filter consists of a regexp pattern. It specifies which namespace data to return. For example, the following filter matches both namespace test and kube-system: `test|kube-system`.").DataType("string").Required(false)).
		Param(ws.QueryParameter("start", "Start time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1559347200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("end", "End time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1561939200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("step", "Time interval. Retrieve metric data at a fixed interval within the time range of start and end. It requires both **start** and **end** are provided. The format is [0-9]+[smhdwy]. Defaults to 10m (i.e. 10 min).").DataType("string").DefaultValue("10m").Required(false)).
		Param(ws.QueryParameter("time", "A timestamp in Unix time format. Retrieve metric data at a single point in time. Defaults to now. Time and the combination of start, end, step are mutually exclusive.").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_metric", "Sort namespaces by the specified metric. Not applicable if **start** and **end** are provided.").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_type", "Sort order. One of asc, desc.").DefaultValue("desc.").DataType("string").Required(false)).
		Param(ws.QueryParameter("page", "The page number. This field paginates result data of each metric, then returns a specific page. For example, setting **page** to 2 returns the second page. It only applies to sorted metric data.").DataType("integer").Required(false)).
		Param(ws.QueryParameter("limit", "Page size, the maximum number of results in a single page. Defaults to 5.").DataType("integer").Required(false).DefaultValue("5")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceMetricsTag}).
		Writes(monitoring.Metrics{}).
		Returns(http.StatusOK, RespOK, monitoring.Metrics{})).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces").To(metering.MonitorNamespace).
		Doc("Get namespace-level metric data of all namespaces.").
		Param(ws.QueryParameter("metrics_filter", "The metric name filter consists of a regexp pattern. It specifies which metric data to return. For example, the following filter matches both namespace CPU usage and memory usage: `namespace_cpu_usage|namespace_memory_usage`. View available metrics at [kubesphere.io](https://docs.kubesphere.io/advanced-v2.0/zh-CN/api-reference/metering-metrics/).").DataType("string").Required(false)).
		Param(ws.QueryParameter("resources_filter", "The namespace filter consists of a regexp pattern. It specifies which namespace data to return. For example, the following filter matches both namespace test and kube-system: `test|kube-system`.").DataType("string").Required(false)).
		Param(ws.QueryParameter("start", "Start time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1559347200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("end", "End time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1561939200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("step", "Time interval. Retrieve metric data at a fixed interval within the time range of start and end. It requires both **start** and **end** are provided. The format is [0-9]+[smhdwy]. Defaults to 10m (i.e. 10 min).").DataType("string").DefaultValue("10m").Required(false)).
		Param(ws.QueryParameter("time", "A timestamp in Unix time format. Retrieve metric data at a single point in time. Defaults to now. Time and the combination of start, end, step are mutually exclusive.").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_metric", "Sort namespaces by the specified metric. Not applicable if **start** and **end** are provided.").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_type", "Sort order. One of asc, desc.").DefaultValue("desc.").DataType("string").Required(false)).
		Param(ws.QueryParameter("page", "The page number. This field paginates result data of each metric, then returns a specific page. For example, setting **page** to 2 returns the second page. It only applies to sorted metric data.").DataType("integer").Required(false)).
		Param(ws.QueryParameter("limit", "Page size, the maximum number of results in a single page. Defaults to 5.").DataType("integer").Required(false).DefaultValue("5")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceMetricsTag}).
		Writes(monitoring.Metrics{}).
		Returns(http.StatusOK, RespOK, monitoring.Metrics{})).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{namespace}").To(metering.MonitorNamespace).
		Doc("Get namespace-level metric data of the specific namespace.").
		Param(ws.PathParameter("namespace", "The name of the namespace.").DataType("string").Required(true)).
		Param(ws.QueryParameter("metrics_filter", "The metric name filter consists of a regexp pattern. It specifies which metric data to return. For example, the following filter matches both namespace CPU usage and memory usage: `namespace_cpu_usage|namespace_memory_usage`. View available metrics at [kubesphere.io](https://docs.kubesphere.io/advanced-v2.0/zh-CN/api-reference/metering-metrics/).").DataType("string").Required(false)).
		Param(ws.QueryParameter("start", "Start time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1559347200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("end", "End time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1561939200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("step", "Time interval. Retrieve metric data at a fixed interval within the time range of start and end. It requires both **start** and **end** are provided. The format is [0-9]+[smhdwy]. Defaults to 10m (i.e. 10 min).").DataType("string").DefaultValue("10m").Required(false)).
		Param(ws.QueryParameter("time", "A timestamp in Unix time format. Retrieve metric data at a single point in time. Defaults to now. Time and the combination of start, end, step are mutually exclusive.").DataType("string").Required(false)).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NamespaceMetricsTag}).
		Writes(monitoring.Metrics{}).
		Returns(http.StatusOK, RespOK, monitoring.Metrics{})).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{namespace}/workloads").To(metering.MonitorWorkload).
		Doc("Get workload-level metric data of a specific namespace's workloads.").
		Param(ws.PathParameter("namespace", "The name of the namespace.").DataType("string").Required(true)).
		Param(ws.QueryParameter("metrics_filter", "The metric name filter consists of a regexp pattern. It specifies which metric data to return. For example, the following filter matches both workload CPU usage and memory usage: `workload_cpu_usage|workload_memory_usage`. View available metrics at [kubesphere.io](https://docs.kubesphere.io/advanced-v2.0/zh-CN/api-reference/metering-metrics/).").DataType("string").Required(false)).
		Param(ws.QueryParameter("resources_filter", "The workload filter consists of a regexp pattern. It specifies which workload data to return. For example, the following filter matches any workload whose name begins with prometheus: `prometheus.*`.").DataType("string").Required(false)).
		Param(ws.QueryParameter("start", "Start time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1559347200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("end", "End time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1561939200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("step", "Time interval. Retrieve metric data at a fixed interval within the time range of start and end. It requires both **start** and **end** are provided. The format is [0-9]+[smhdwy]. Defaults to 10m (i.e. 10 min).").DataType("string").DefaultValue("10m").Required(false)).
		Param(ws.QueryParameter("time", "A timestamp in Unix time format. Retrieve metric data at a single point in time. Defaults to now. Time and the combination of start, end, step are mutually exclusive.").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_metric", "Sort workloads by the specified metric. Not applicable if **start** and **end** are provided.").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_type", "Sort order. One of asc, desc.").DefaultValue("desc.").DataType("string").Required(false)).
		Param(ws.QueryParameter("page", "The page number. This field paginates result data of each metric, then returns a specific page. For example, setting **page** to 2 returns the second page. It only applies to sorted metric data.").DataType("integer").Required(false)).
		Param(ws.QueryParameter("limit", "Page size, the maximum number of results in a single page. Defaults to 5.").DataType("integer").Required(false).DefaultValue("5")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.WorkloadMetricsTag}).
		Writes(monitoring.Metrics{}).
		Returns(http.StatusOK, RespOK, monitoring.Metrics{})).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/namespaces/{namespace}/workloads/{kind}").To(metering.MonitorWorkload).
		Doc("Get workload-level metric data of all workloads which belongs to a specific kind.").
		Param(ws.PathParameter("namespace", "The name of the namespace.").DataType("string").Required(true)).
		Param(ws.PathParameter("kind", "Workload kind. One of deployment, daemonset, statefulset.").DataType("string").Required(true)).
		Param(ws.QueryParameter("metrics_filter", "The metric name filter consists of a regexp pattern. It specifies which metric data to return. For example, the following filter matches both workload CPU usage and memory usage: `workload_cpu_usage|workload_memory_usage`. View available metrics at [kubesphere.io](https://docs.kubesphere.io/advanced-v2.0/zh-CN/api-reference/metering-metrics/).").DataType("string").Required(false)).
		Param(ws.QueryParameter("resources_filter", "The workload filter consists of a regexp pattern. It specifies which workload data to return. For example, the following filter matches any workload whose name begins with prometheus: `prometheus.*`.").DataType("string").Required(false)).
		Param(ws.QueryParameter("start", "Start time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1559347200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("end", "End time of query. Use **start** and **end** to retrieve metric data over a time span. It is a string with Unix time format, eg. 1561939200. ").DataType("string").Required(false)).
		Param(ws.QueryParameter("step", "Time interval. Retrieve metric data at a fixed interval within the time range of start and end. It requires both **start** and **end** are provided. The format is [0-9]+[smhdwy]. Defaults to 10m (i.e. 10 min).").DataType("string").DefaultValue("10m").Required(false)).
		Param(ws.QueryParameter("time", "A timestamp in Unix time format. Retrieve metric data at a single point in time. Defaults to now. Time and the combination of start, end, step are mutually exclusive.").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_metric", "Sort workloads by the specified metric. Not applicable if **start** and **end** are provided.").DataType("string").Required(false)).
		Param(ws.QueryParameter("sort_type", "Sort order. One of asc, desc.").DefaultValue("desc.").DataType("string").Required(false)).
		Param(ws.QueryParameter("page", "The page number. This field paginates result data of each metric, then returns a specific page. For example, setting **page** to 2 returns the second page. It only applies to sorted metric data.").DataType("integer").Required(false)).
		Param(ws.QueryParameter("limit", "Page size, the maximum number of results in a single page. Defaults to 5.").DataType("integer").Required(false).DefaultValue("5")).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.WorkloadMetricsTag}).
		Writes(monitoring.Metrics{}).
		Returns(http.StatusOK, RespOK, monitoring.Metrics{})).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	c.Add(ws)
	return nil
}
