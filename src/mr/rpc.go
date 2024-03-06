package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import "os"
import "strconv"

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}
type TaskArgs struct {
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.
type Coordinator struct {
	// Your definitions here.
	Nreduce int
	Mapchan chan *Job
	Redchan chan *Job
	Files []string
	State int //state of tasks
}

// Your code here -- RPC handlers for the worker to call.
//struct send to worker
type Jobtype string
type Job struct{
	Jobtype Jobtype
	Joblist []string
	Job_num int
	Nreduce int
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
