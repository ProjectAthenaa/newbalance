package module

import (
	"fmt"
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/base"
	"github.com/ProjectAthenaa/sonic-core/sonic/face"
	"github.com/ProjectAthenaa/sonic-core/sonic/frame"
	"github.com/prometheus/common/log"
)

var _ face.ICallback = (*Task)(nil)

type Task struct {
	*base.BTask
	productUrl           string
	productLineItemUUID  string
	originalShipmentUUID string
	shipmentUUID         string
	csrf_token           string
	PID                  string
	VariantId            string
}

func NewTask(data *module.Data) *Task {
	task := &Task{BTask: &base.BTask{Data: data}}
	task.Callback = task
	task.Init()
	return task
}

func (tk *Task) OnInit() {
	return
}
func (tk *Task) OnPreStart() error {
	return nil
}
func (tk *Task) OnStarting() {
	log.Info("hit on starting")
	tk.FastClient.CreateCookieJar()
	tk.Flow()
}
func (tk *Task) OnPause() error {
	return nil
}
func (tk *Task) OnStopping() {
	tk.FastClient.Destroy()
	//panic("")
	return
}

func (tk *Task) Flow() {
	pubsub, err := frame.SubscribeToChannel(tk.Data.Channels.MonitorChannel)
	if err != nil {
		log.Error(err)
		tk.SetStatus(module.STATUS_ERROR, "error listening to monitor")
		tk.Stop()
		return
	}
	fmt.Println(tk.Data.Channels.MonitorChannel)
	tk.SetStatus(module.STATUS_MONITORING)
	monitorData := <-pubsub.Chan(tk.Ctx)
	pubsub.Close()
	tk.VariantId = monitorData["variantid"].(string)
	tk.PID = monitorData["pid"].(string)
	tk.productUrl = monitorData["endpoint"].(string)

	tk.SetStatus(module.STATUS_PRODUCT_FOUND)

	funcarr := []func(){
		tk.GetProductPage,
		tk.ATC,
		tk.BeginCheckout,
		tk.Shipping,
		tk.Payment,
		tk.PlaceOrder,
	}

	for _, f := range funcarr {
		select {
		case <-tk.Ctx.Done():
			return
		default:
			f()
		}
	}
}
