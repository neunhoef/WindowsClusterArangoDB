package main

import (
	"fmt"
	"os"
	"strconv"
)

var ipAddresses = []string{
	"10.192.43.36",
	"10.192.43.39",
	"10.192.43.41",
	"10.192.43.38",
	"10.192.43.42",
}

var dbserverNames = []string{
	"dbserver4",
	"dbserver3",
	"dbserver1",
	"dbserver5",
	"dbserver2",
}

var coordinatorNames = []string{
	"coordinator4",
	"coordinator3",
	"coordinator1",
	"coordinator5",
	"coordinator2",
}

var agentNames = []string{
	"agent0",
	"agent1",
	"agent2",
}

var installationDir = "E:/arangodb/"
var configDir = "D:/arangodb/configuration/"
var executable = `D:/arangodb/ArangoDB3e-3.1.2-1_win64/usr/bin/arangod.exe`

var agentSkeleton = `# ArangoDB configuration file

[database]
directory = %s

[server]
endpoint = tcp://0.0.0.0:%d
authentication = false
statistics = false
threads = 5

[javascript]
startup-directory = @ROOTDIR@/usr/share/arangodb3/js
app-path = %s

[foxx]
queues = false

[log]
level = info
file = %s

[agency]
activate = true
size = 3
supervision = true
my-address = %s
election-timeout-min = 1.0
election-timeout-max = 5.0
`

var serverSkeleton = `# ArangoDB configuration file

[database]
directory = %s

[server]
endpoint = tcp://0.0.0.0:%d
authentication = false
statistics = true
threads = 5

[javascript]
startup-directory = @ROOTDIR@/usr/share/arangodb3/js
app-path = %s

[foxx]
queues = true

[log]
level = info
file = %s

[cluster]
my-address = %s
my-local-info = %s
my-role = %s
agency-endpoint = %s
agency-endpoint = %s
agency-endpoint = %s
`

func agentConfigs() {
	for i := 0; i < 3; i++ {
		name := agentNames[i]
		out, _ := os.Create(name + ".conf")
		fmt.Fprintf(out, agentSkeleton, installationDir+name, 4001,
			installationDir+name+"-apps",
			installationDir+name+".log",
			"tcp://"+ipAddresses[i]+":4001")
		for j := 0; j < i; j++ {
			fmt.Fprintf(out, "endpoint = %s\n", "tcp://"+ipAddresses[j]+":4001")
		}
		out.Close()
	}
}

func dbserverConfigs() {
	for i := 0; i < len(ipAddresses); i++ {
		name := dbserverNames[i]
		out, _ := os.Create(name + ".conf")
		fmt.Fprintf(out, serverSkeleton, installationDir+name, 8629,
			installationDir+name+"-apps",
			installationDir+name+".log",
			"tcp://"+ipAddresses[i]+":8629",
			ipAddresses[i]+":8629", "PRIMARY",
			"tcp://"+ipAddresses[0]+":4001",
			"tcp://"+ipAddresses[1]+":4001",
			"tcp://"+ipAddresses[2]+":4001")
		out.Close()
	}
}

func coordinatorConfigs() {
	for i := 0; i < len(ipAddresses); i++ {
		name := coordinatorNames[i]
		out, _ := os.Create(name + ".conf")
		fmt.Fprintf(out, serverSkeleton, installationDir+name, 8530,
			installationDir+name+"-apps",
			installationDir+name+".log",
			"tcp://"+ipAddresses[i]+":8530",
			ipAddresses[i]+":8530", "COORDINATOR",
			"tcp://"+ipAddresses[0]+":4001",
			"tcp://"+ipAddresses[1]+":4001",
			"tcp://"+ipAddresses[2]+":4001")
		out.Close()
	}
}

func makeBatFiles(typ string) {
	var nr = len(ipAddresses)
	if typ == "agent" {
		nr = 3
	}
	for i := 0; i < nr; i++ {
		var name string
		if typ == "agent" {
			name = agentNames[i]
		} else if typ == "dbserver" {
			name = dbserverNames[i]
		} else if typ == "coordinator" {
			name = coordinatorNames[i]
		}
		out, _ := os.Create(name + ".bat")
		fmt.Fprintf(out, "\"%s\" --configuration %s\n", executable,
			configDir+name+".conf")
		out.Chmod(0755)
		out.Close()
	}
}

func serviceCreateBat() {
	for i := 0; i < len(ipAddresses); i++ {
		out, _ := os.Create("createServices" + strconv.Itoa(i) + ".bat")
		if i < 3 {
			fmt.Fprintf(out, `sc create ArangoDB%s type=own start=auto error=normal binPath="\"%s\" --start-service true --configuration %s%s.conf"`,
				agentNames[i], executable, configDir, agentNames[i])
			fmt.Fprintf(out, "\r\nsc description ArangoDB%s ArangoDB%s\r\n",
				agentNames[i], agentNames[i])
		}
		fmt.Fprintf(out, `sc create ArangoDB%s type=own start=auto error=normal binPath="\"%s\" --start-service true --configuration %s%s.conf"`,
			coordinatorNames[i], executable, configDir, coordinatorNames[i])
		fmt.Fprintf(out, "\r\nsc description ArangoDB%s ArangoDB%s\r\n",
			coordinatorNames[i], coordinatorNames[i])
		fmt.Fprintf(out, `sc create ArangoDB%s type=own start=auto error=normal binPath="\"%s\" --start-service true --configuration %s%s.conf"`,
			dbserverNames[i], executable, installationDir, dbserverNames[i])
		fmt.Fprintf(out, "\r\nsc description ArangoDB%s ArangoDB%s\r\n",
			dbserverNames[i], dbserverNames[i])
		out.Chmod(0755)
		out.Close()
	}
}

func serviceDeleteBat() {
	for i := 0; i < len(ipAddresses); i++ {
		out, _ := os.Create("deleteServices" + strconv.Itoa(i) + ".bat")
		if i < 3 {
			fmt.Fprintf(out, "sc delete ArangoDB%s\r\n", agentNames[i])
		}
		fmt.Fprintf(out, "sc delete ArangoDB%s\r\n", coordinatorNames[i])
		fmt.Fprintf(out, "sc delete ArangoDB%s\r\n", dbserverNames[i])
		out.Chmod(0755)
		out.Close()
	}
}

func main() {
	agentConfigs()
	makeBatFiles("agent")
	dbserverConfigs()
	makeBatFiles("dbserver")
	coordinatorConfigs()
	makeBatFiles("coordinator")
	serviceCreateBat()
	serviceDeleteBat()
}
