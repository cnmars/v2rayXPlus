package core

import (
	"gitlab.com/xiayesuifeng/v2rayxplus/conf"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func CopyFile(src, target string) error {
	srcConfig, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcConfig.Close()

	targetConfig, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer targetConfig.Close()

	_, err = io.Copy(targetConfig, srcConfig)
	if err != nil {
		return err
	}

	return nil
}

func GetConfigName() (name, path string) {
	for i := 0; ; i++ {
		tmp := ""
		if i > 0 {
			tmp = strconv.FormatInt(int64(i), 10)
		}
		path := conf.ConfigPath + "/v2ray/config" + tmp + ".json"
		_, err := os.Stat(path)
		if err != nil && os.IsNotExist(err) {
			return "config" + tmp, path
		}
	}
}

func StartV2rayXPlusSerive(config string) bool {
	bytes, err := exec.Command("systemd-escape", "-p", conf.V2rayConfigPath+"/"+config+".json").CombinedOutput()
	if err != nil {
		return false
	}
	return StartService("v2rayxplus@" + string(bytes))
}

func RestartV2rayXPlusSerive(config string) bool {
	bytes, err := exec.Command("systemd-escape", "-p", conf.V2rayConfigPath+"/"+config+".json").CombinedOutput()
	if err != nil {
		return false
	}
	if !StopService("v2rayxplus@" + string(bytes)) {
		return false
	}

	return StartService("v2rayxplus@" + string(bytes))
}

func StatusV2rayXPlusSerive(config string) (exited, enable bool) {
	bytes, err := exec.Command("systemd-escape", "-p", conf.V2rayConfigPath+"/"+config+".json").CombinedOutput()
	if err != nil {
		return false, false
	}

	return StatusService("v2rayxplus@" + string(bytes))
}

func StopV2rayXPlusSerive(config string) bool {
	bytes, err := exec.Command("systemd-escape", "-p", conf.V2rayConfigPath+"/"+config+".json").CombinedOutput()
	if err != nil {
		return false
	}
	return StopService("v2rayxplus@" + string(bytes))
}

func StartService(service string) bool {
	return exec.Command("systemctl", "start", service).Run() == nil
}

func RestartService(service string) bool {
	return exec.Command("systemctl", "restart", service).Run() == nil
}

func StopService(service string) bool {
	return exec.Command("systemctl", "stop", service).Run() == nil
}

func StatusService(service string) (exited, enable bool) {
	bytes, err := exec.Command("systemctl", "status", service).CombinedOutput()
	if err != nil {
		return exited, enable
	}

	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Loaded:") {
			tmp := strings.Split(line, ";")
			if len(tmp) > 1 {
				if strings.Contains(tmp[1], "enable") {
					enable = true
				}
			}
		} else if strings.HasPrefix(line, "Active:") {
			exited = strings.Contains(line, "exited")
		}
	}

	return exited, enable
}

func InitIpTables() {
	if exited, _ := StatusService("iptables"); !exited {
		if _, err := os.Stat("/etc/iptable/iptables.rules"); os.IsNotExist(err) {
			if _, err := os.Stat("/etc/iptable/empty.rules"); err == nil {
				CopyFile("/etc/iptable/empty.rules", "/etc/iptable/iptables.rules")
				StartService("iptables")
			}
		}
	}

	bytes, err := exec.Command("sysctl", "net.ipv4.ip_forward").CombinedOutput()
	result := strings.TrimSpace(string(bytes))
	if err == nil {
		if strings.HasSuffix(result, "0") {
			log.Println(string(bytes))
			exec.Command("sysctl", "net.ipv4.ip_forward=1").Run()
		}
	}
}
