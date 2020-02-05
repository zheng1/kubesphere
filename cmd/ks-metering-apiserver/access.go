package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"k8s.io/klog"
	"kubesphere.io/kubesphere/pkg/models/metrics"
)

type AccessServer struct {
	server *http.Server
}

// Access Server parameters
type AccessSvrParameters struct {
	port int // access server port
}

type MeterData struct {
	ResourceId string    `json:"resource_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	MeterId    string    `json:"meter_id"`
	MeterValue float64   `json:"meter_value"`
}

type GetMeteringResponseBody struct {
	Total    int         `json:"total"`
	MeterSet []MeterData `json:"meter_set"`
}

func startAccessServer(accessSvrParameters AccessSvrParameters) *AccessServer {
	accessServer := &AccessServer{
		server: &http.Server{
			Addr: fmt.Sprintf(":%v", accessSvrParameters.port),
		},
	}

	// define http server and server handler
	mux := http.NewServeMux()
	mux.HandleFunc("/metering_callback", accessServer.serve)
	accessServer.server.Handler = mux

	// start access server in new routine
	go func() {
		if err := accessServer.server.ListenAndServe(); err != nil {
			klog.Errorf("Failed to listen and serve access server: %v", err)
		}
	}()
	return accessServer
}

func getMeter(resourceIds, meters []string, startTime, endTime time.Time) GetMeteringResponseBody {
	query := url.Values{}
	query.Set("start", strconv.Itoa(int(startTime.Unix())))
	query.Set("end", strconv.Itoa(int(endTime.Unix())))
	query.Set("step", "60s")

	respBody := GetMeteringResponseBody{}
	for _, resourceId := range resourceIds {
		resource := GetResourceFromString(resourceId)

		params := metrics.RequestParams{
			MetricsFilter: strings.Join(meters, "|"),
			NamespaceName: resource.Namespace,
			WorkloadName:  resource.Name,
			WorkloadKind:  "deployment",
			QueryType:     metrics.RangeQuery,
			QueryParams:   query,
		}
		klog.Infof("params: [%+v]", params)
		resp := metrics.GetPodMetrics(params)
		klog.Infof("resp: [%+v]", resp)
		for _, result := range resp.Results {
			meterId := result.MetricName
			if len(result.Data.Result) == 0 {
				continue
			}
			for _, v := range result.Data.Result[0].Values {
				t := v[0].(float64)
				value := v[1].(string)
				v, _ := strconv.ParseFloat(value, 64)
				st := time.Unix(int64(t), 0)
				md := MeterData{
					ResourceId: resourceId,
					MeterId:    meterId,
					StartTime:  st,
					EndTime:    st,
					MeterValue: v,
				}
				respBody.MeterSet = append(respBody.MeterSet, md)
			}
		}
	}
	respBody.Total = len(respBody.MeterSet)
	return respBody
}

var loc = time.FixedZone("CST", 8*3600)

// Serve method for webhook server
func (asvr *AccessServer) serve(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	resourceIds := query.Get("resource_ids")
	startTimeStr := query.Get("start_time")
	endTimeStr := query.Get("end_time")
	meters := query.Get("meters")

	resourceIdSet := strings.Split(resourceIds, "|")
	metersSet := strings.Split(meters, "|")
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startTimeStr, loc)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endTimeStr, loc)
	klog.Infof("Get meter with resource ids [%+v] meters [%+v] start_time [%v] end_time [%v]",
		resourceIdSet, metersSet, startTime, endTime)
	respMeter := getMeter(resourceIdSet, metersSet, startTime, endTime)

	resp, err := json.Marshal(respMeter)
	if err != nil {
		klog.Errorf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	klog.Infof("Ready to write response: %s ...", string(resp))
	if _, err := w.Write(resp); err != nil {
		klog.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
