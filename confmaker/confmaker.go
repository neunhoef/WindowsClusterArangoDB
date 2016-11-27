package main

import (
  "fmt"
	"os"
	"strconv"
)

var ipAddresses = []string{
  "192.168.173.95",
  "192.168.173.224",
  "192.168.173.221",
}

var installationDir = "c:/Users/jenkins/ArangoDB/"
var executable = `C:/Program Files/ArangoDB3e 3.1.2/usr/bin/arangod.exe`

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
		name := "agent" + strconv.Itoa(i)
		out, _ := os.Create(name + ".conf")
		fmt.Fprintf(out, agentSkeleton, installationDir + name, 4001,
		            installationDir + name + "-apps",
								installationDir + name + ".log",
								"tcp://" + ipAddresses[i] + ":4001")
		for j := 0; j < i; j++ {
			fmt.Fprintf(out, "endpoint = %s\n", "tcp://" + ipAddresses[j] + ":4001")
		}
		out.Close()
	}
}

func dbserverConfigs() {
	for i := 0; i < len(ipAddresses); i++ {
		name := "dbserver" + strconv.Itoa(i)
		out, _ := os.Create(name + ".conf")
		fmt.Fprintf(out, serverSkeleton, installationDir + name, 8629,
		            installationDir + name + "-apps",
								installationDir + name + ".log",
								"tcp://" + ipAddresses[i] + ":8629",
								ipAddresses[i] + ":8629", "PRIMARY",
								"tcp://" + ipAddresses[0] + ":4001",
								"tcp://" + ipAddresses[1] + ":4001",
								"tcp://" + ipAddresses[2] + ":4001")
		out.Close()
	}
}

func coordinatorConfigs() {
	for i := 0; i < len(ipAddresses); i++ {
		name := "coordinator" + strconv.Itoa(i)
		out, _ := os.Create(name + ".conf")
		fmt.Fprintf(out, serverSkeleton, installationDir + name, 8530,
		            installationDir + name + "-apps",
								installationDir + name + ".log",
								"tcp://" + ipAddresses[i] + ":8530",
								ipAddresses[i] + ":8629", "COORDINATOR",
								"tcp://" + ipAddresses[0] + ":4001",
								"tcp://" + ipAddresses[1] + ":4001",
								"tcp://" + ipAddresses[2] + ":4001")
		out.Close()
	}
}

func makeBatFiles(typ string) {
	var nr = len(ipAddresses)
	if typ == "agent" {
		nr = 3
	}
	for i := 0; i < nr; i++ {
		name := typ + strconv.Itoa(i)
		out, _ := os.Create(name + ".bat")
		fmt.Fprintf(out, "\"%s\" --configuration %s\n", executable,
		            installationDir + name + ".conf")
		out.Chmod(0755)
		out.Close()
	}
}

func serviceCreateBat() {
	for i := 0; i < len(ipAddresses); i++ {
		out, _ := os.Create("createServices" + strconv.Itoa(i) + ".bat")
		if i < 3 {
			fmt.Fprintf(out, `sc create ArangoDBAgent type=own start=auto error=normal binPath="\"%s\" --start-service true --configuration %sagent%d.conf"`,
			            executable, installationDir, i)
			fmt.Fprintf(out, "\r\nsc description ArangoDBAgent ArangoDBAgent%d\r\n",
			            i)
		}
		fmt.Fprintf(out, `sc create ArangoDBCoordinator type=own start=auto error=normal binPath="\"%s\" --start-service true --configuration %scoordinator%d.conf"`,
								executable, installationDir, i)
		fmt.Fprintf(out, "\r\nsc description ArangoDBCoordinator ArangoDBCoordinator%d\r\n",
								i)
		fmt.Fprintf(out, `sc create ArangoDBDBserver type=own start=auto error=normal binPath="\"%s\" --start-service true --configuration %sdbserver%d.conf"`,
								executable, installationDir, i)
		fmt.Fprintf(out, "\r\nsc description ArangoDBDBserver ArangoDBDBserver%d\r\n",
								i)
		out.Chmod(0755)
		out.Close()
  }
}

func serviceDeleteBat() {
	for i := 0; i < len(ipAddresses); i++ {
		out, _ := os.Create("deleteServices" + strconv.Itoa(i) + ".bat")
		if i < 3 {
			fmt.Fprintf(out, "sc delete ArangoDBAgent%d\r\n", i)
		}
	  fmt.Fprintf(out, "sc delete ArangoDBCoordinator%d\r\n", i)
	  fmt.Fprintf(out, "sc delete ArangoDBDBserver%d\r\n", i)
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
