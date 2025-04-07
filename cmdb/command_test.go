package main

import "testing"

func Test_execute(t *testing.T) {
	shell = [2]string{"/bin/bash", "-c"}
	err := execute("test", `
#!/bin/bash
 
#print time
for((i=0;i<10;i++))
do
    sleep 1
    echo $(date +"%Y-%m-%d %H:%M:%S")
done
`, nil)
	t.Log(err)
}
