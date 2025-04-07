package main

import (
	"bufio"
	"encoding/json"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var (
	goos       = runtime.GOOS
	shell      = [2]string{}
	commandMap = new(sync.Map)
)

func initShell() {
	var ok bool
	shells := []string{"sh", "bash", "cmd", "powershell"}
	for _, item := range shells {
		if ok = checkShell(item); ok {
			shell[0] = item
			if item == "cmd" || item == "powershell" {
				shell[1] = "/c"
			} else {
				shell[1] = "-c"
			}
			break
		}
	}
}

func checkShell(shell string) bool {
	_, err := exec.LookPath(shell)
	if err == nil {
		return true
	}
	return false
}

func asyncLog(reader io.ReadCloser, f func(msg string)) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		msg := scanner.Text()
		log.Printf("async log: %s\n", msg)
		if f != nil && len(msg) > 0 {
			f(msg)
		}
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF && !strings.Contains(err.Error(), "closed") {
			log.Printf("Error reading command output: %s\n", err.Error())
			if f != nil {
				f(err.Error())
			}
		}
	}
	return nil
}

func cancel(id string) error {
	od, ok := commandMap.Load(id)
	if !ok {
		return nil
	}
	if err := od.(*exec.Cmd).Process.Kill(); err != nil {
		log.Printf("Error killing command: %s\n", err.Error())
		return err
	}
	return nil
}
func execute(id, src string, f func(msg string)) error {
	cmd := exec.Command(shell[0], shell[1], src)
	commandMap.Store(id, cmd)
	defer commandMap.Delete(id)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %s......", err.Error())
		return err
	}

	go asyncLog(stdout, f)
	go asyncLog(stderr, f)

	if err := cmd.Wait(); err != nil {
		log.Printf("Error waiting for command execution: %s......", err.Error())
		return err
	}
	return nil
}
func connExecute(w ws, src string) {
	id := gjson.Get(src, "id").String()
	data := gjson.Get(src, "data").String()
	if id == "" || data == "" {
		_ = w.SendMsg(SERVER_EXEC_COMMAND, "参数校验失败")
		return
	}
	if err := execute(id, data, func(msg string) {
		body := map[string]interface{}{
			"id":  id,
			"msg": msg,
		}
		js, _ := json.Marshal(body)
		err := w.SendMsg(SERVER_EXEC_COMMAND, js)
		if err != nil {
			log.Printf("Error sending command output: %s\n", err.Error())
		}
	}); err != nil {
		_ = w.SendMsg(SERVER_EXEC_COMMAND, "执行失败:"+err.Error())
	} else {
		_ = w.SendMsg(SERVER_EXEC_COMMAND, "执行成功")
	}
}

func connExecuteCancel(w ws, id string) {
	if err := cancel(id); err != nil {
		_ = w.SendMsg(SERVER_CANCEL_EXEC_COMMAND, "取消执行脚本失败:"+err.Error())
	} else {
		_ = w.SendMsg(SERVER_CANCEL_EXEC_COMMAND, "取消执行脚本成功")
	}
}
