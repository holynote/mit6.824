package mr

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"sort"
	"strconv"
	"encoding/json"
)

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}


//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.
	task := gettask()
	// uncomment to send the Example RPC to the coordinator.
	// CallExample()
	DoMaptask(&task, mapf, reducef)
}
func gettask() Job{
	args := TaskArgs{}
	rely := Job{}
	if ok := call("Coordinator.Gettask", &args, &rely); ok{
		fmt.Printf("successfully get task-%d\n", rely.Job_num)	
	}else{
		fmt.Printf("fail to get task\n")
	}
	return rely
}

func DoMaptask(task *Job, mapf func(string, string) []KeyValue, reducef func(string, []string) string) {
	intermediate := []KeyValue{}
	for _, filename := range task.Joblist {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("cannot open %v", filename)
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("cannot read %v", filename)
		}
		file.Close()
		kva := mapf(filename, string(content))
		intermediate = append(intermediate, kva...)
	} //file of task changed to []KV

	sort.Sort(ByKey(intermediate))
	nreduce := task.Nreduce  //total nums of reduce task
	reduce_num := make([][]KeyValue, nreduce)
	for _, v := range intermediate{
		index := ihash(v.Key) % nreduce
		reduce_num[index] = append(reduce_num[index], v)
	} // reduce num correspond to []KY

	//create tmp file for each reduce num
	for i, _ := range reduce_num{
		oname := "mr-tmp-" + strconv.Itoa(task.Job_num) + strconv.Itoa(i)
		ofile, err := os.Create(oname)
		if err != nil{
			log.Fatal("create file failed", err)
		}
		enc := json.NewEncoder(ofile)
  		for _, kv := range reduce_num[i] {
    		err := enc.Encode(&kv)
			if err != nil{
				log.Fatal("encode failed", err)
			}
		}
		ofile.Close()
	}
}

//
// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
//
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
