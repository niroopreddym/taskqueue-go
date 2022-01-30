package enums

//Status is the enum
type Status int

const (
	//Untouched ...
	Untouched = iota + 1
	//Completed task execution completed
	Completed
	//Failed task execution failed
	Failed
	//Timeout is context timeout
	Timeout
)

func (v Status) String() string {
	dictMap := map[Status]string{
		1: "Untouched",
		2: "Completed",
		3: "Failed",
		4: "Timeout",
	}

	return dictMap[v]
}
