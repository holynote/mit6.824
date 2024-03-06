package mr

import "log"
import "net"
import "os"
import "net/rpc"
import "net/http"



//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}


//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}
func (c *Coordinator) Gettask(args TaskArgs, rely *Job) {
	*rely = *<-c.Mapchan
}
//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.


	return ret
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{
		Files: files,
		Nreduce: nReduce,
		Mapchan: make(chan *Job, len(files)),
		Redchan: make(chan *Job, nReduce),
		State: 0,
	}
	// Your code here.
	c.Make_maptask()
	c.server()
	return &c
}
func (c *Coordinator) Make_maptask(){
	for i, f := range c.Files{
		//one task deal with one file
		task := Job{
			Jobtype: "map",
			Joblist: []string{f},
			Job_num: i,
			Nreduce: c.Nreduce,
		}
		c.Mapchan <- &task
	}
}