sc create ArangoDBAgent1 type=own start=auto error=normal binPath="\"C:\Program Files\Arangodb3e 3.1.1\usr\bin\arangod.exe\" --start-service true --configuration c:/Users/jenkins/max/agent1.conf"
sc description ArangoDBAgent1 ArangoDBAgent1
sc create ArangoDBAgent2 type=own start=auto error=normal binPath="\"C:\Program Files\Arangodb3e 3.1.1\usr\bin\arangod.exe\" --start-service true --configuration c:/Users/jenkins/max/agent2.conf"
sc description ArangoDBAgent2 ArangoDBAgent2
sc create ArangoDBAgent3 type=own start=auto error=normal binPath="\"C:\Program Files\Arangodb3e 3.1.1\usr\bin\arangod.exe\" --start-service true --configuration c:/Users/jenkins/max/agent3.conf"
sc description ArangoDBAgent3 ArangoDBAgent3
sc create ArangoDBDBServer1 type=own start=auto error=normal binPath="\"C:\Program Files\Arangodb3e 3.1.1\usr\bin\arangod.exe\" --start-service true --configuration c:/Users/jenkins/max/dbserver1.conf" depend=ArangoDBAgent1/ArangoDBAgent2/ArangoDBAgent3
sc description ArangoDBDBServer1 ArangoDBDBServer1
sc create ArangoDBDBServer2 type=own start=auto error=normal binPath="\"C:\Program Files\Arangodb3e 3.1.1\usr\bin\arangod.exe\" --start-service true --configuration c:/Users/jenkins/max/dbserver2.conf" depend=ArangoDBAgent1/ArangoDBAgent2/ArangoDBAgent3
sc description ArangoDBDBServer2 ArangoDBDBServer2
sc create ArangoDBCoordinator1 type=own start=auto error=normal binPath="\"C:\Program Files\Arangodb3e 3.1.1\usr\bin\arangod.exe\" --start-service true --configuration c:/Users/jenkins/max/coordinator1.conf" depend=ArangoDBAgent1/ArangoDBAgent2/ArangoDBAgent3
sc description ArangoDBCoordinator1 ArangoDBCoordinator1
