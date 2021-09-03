package module

import (
	"github.com/ProjectAthenaa/sonic-core/protos/module"
	"github.com/ProjectAthenaa/sonic-core/sonic/base"
	"github.com/ProjectAthenaa/sonic-core/sonic/face"
	"github.com/ProjectAthenaa/sonic-core/sonic/frame"
)

var _ face.ICallback = (*Task)(nil)

type Task struct {
	*base.BTask
	Monitor *frame.PubSub
	productLineItemUUID  string
	originalShipmentUUID string
	shipmentUUID         string
	csrf_token           string
	PID					 string
	VariantId    		 string
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
	tk.SetStatus(module.STATUS_STARTING, "starting task")
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
	if err != nil{
		tk.Stop()
		return
	}
	defer pubsub.Close()

	monitorData := <- pubsub.Chan(tk.Ctx)
	tk.VariantId = monitorData["variantid"].(string)
	tk.PID = monitorData["pid"].(string)

	funcarr := []func(){
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
