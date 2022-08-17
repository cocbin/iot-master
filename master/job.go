package master

import (
	"fmt"
	"iot-master/model"
	"iot-master/pkg/cron"
	"iot-master/pkg/events"
)

//Job 任务
type Job struct {
	model.Job

	job *cron.Job

	events.EventEmitter
}

//Start 启动任务
func (j *Job) Start() error {
	if j.job != nil {
		j.job.Cancel()
		//return errors.New("任务已经启动")
	}

	hours := j.Clock / 60
	minutes := j.Clock % 60

	var err error
	j.job, err = cron.ClockWithWeekdays(hours, minutes, j.Weekdays, func() {
		j.Execute()
	})
	return err
}

//Execute 执行任务
func (j *Job) Execute() {
	//for _, i:= range j.Invokes {
	//	j.events.Publish("invoke", i)
	//}
	//TODO 避免拥堵计时器(需要确认)
	go j.Emit("invoke")
}

//Stop 取消任务
func (j *Job) Stop() {
	if j.job != nil {
		j.job.Cancel()
	}
}

//String 任务描述
func (j *Job) String() string {
	hours := j.Clock / 60
	minutes := j.Clock % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
