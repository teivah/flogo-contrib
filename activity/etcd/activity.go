package etcd

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"strings"
	"github.com/coreos/etcd/client"
	"time"
	"fmt"
	"context"
)

// ActivityLog is the default logger for the Log Activity
var activityLog = logger.GetLogger("activity-tibco-etcd")

const (
	serverDelimiter = ";"

	methodCreate = "Create"
	methodGet    = "Get"
	methodUpdate = "Update"
	methodDelete = "Delete"

	ivKey     = "key"
	ivValue   = "value"
	ivMethod  = "method"
	ivServers = "servers"

	ovOutput = "output"
)

func init() {
	activityLog.SetLogLevel(logger.InfoLevel)
}

// Integration with Couchbase
// inputs: {data, method, expiry, server, username, password, bucket, bucketPassword}
// outputs: {output, status}
type EtcdActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &EtcdActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *EtcdActivity) Metadata() *activity.Metadata {
	return a.metadata
}

type A struct {
}

// Eval implements api.Activity.Eval - Couchbase integration
func (a *EtcdActivity) Eval(ctx activity.Context) (done bool, err error) {
	key, _ := ctx.GetInput(ivKey).(string)
	value, _ := ctx.GetInput(ivValue).(string)
	method, _ := ctx.GetInput(ivMethod).(string)
	serverList, _ := ctx.GetInput(ivServers).(string)

	servers := strings.Split(serverList, serverDelimiter)

	cfg := client.Config{
		Endpoints:               servers,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)

	if err != nil {
		activityLog.Errorf("etcd connection error: %v", err)
		return false, err
	}

	kapi := client.NewKeysAPI(c)

	switch method {
	case methodGet:
		resp, err := kapi.Get(context.Background(), key, nil)
		if err != nil {
			activityLog.Errorf("Get error: %v", err)
			return false, err
		} else {
			ctx.SetOutput(ovOutput, resp.Node.Value)
			return true, nil
		}
	case methodCreate:
		resp, err := kapi.Create(context.Background(), key, value)
		if err != nil {
			activityLog.Errorf("Create error: %v", err)
			return false, err
		} else {
			ctx.SetOutput(ovOutput, resp.Node.Value)
			return true, nil
		}
	case methodUpdate:
		resp, err := kapi.Update(context.Background(), key, value)
		if err != nil {
			activityLog.Errorf("Update error: %v", err)
			return false, err
		} else {
			ctx.SetOutput(ovOutput, resp.Node.Value)
			return true, nil
		}
	case methodDelete:
		resp, err := kapi.Delete(context.Background(), key, nil)
		if err != nil {
			activityLog.Errorf("Delete error: %v", err)
			return false, err
		} else {
			ctx.SetOutput(ovOutput, resp.Node.Value)
			return true, nil
		}
	default:
		activityLog.Errorf("Method %v not recognized.", method)
		return false, fmt.Errorf("method %v not recognized", method)
	}

	return true, nil
}
