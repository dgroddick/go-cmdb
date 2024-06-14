# Go CMDB

Go utility for working with ServiceNow CMDB. Provides an interface to perform simple CRUD actions. I'd previously written this using Python, I wanted to try rewriting it in Go. This is still a work in progress.

* Importing Linux hosts into ServiceNow from Ansible. 
* Updating ServiceNow Configuration Items when host information (in Ansible) changes.
* Deleting a host from ServiceNow CMDB.

### Installation

You can either clone the repo, and build or run the project

```
$ git clone https://github.com/dgroddick/go-cmdb.git && cd go-cmdb
$ go run main.go

or

$ go build
$ ./gocmdb

```

### Instructions

Assuming you've installed the package, run:

```
$ ./gocmdb <arg>
```

The ServiceNow instance connection details are taken from Linux environment variables. You need to have a running ServiceNow instance with access to the CMDB.
You can then add your authentication details to the .env file and run

```
$ source .env
```

Should be able to do most actions with cli arguments.
Currently there's the following:

```
-output <hostname>
-list-all
-add
-update
-delete <hostname>
-delete-all
```

output and output-all dump out the host contents from SN. Output requires the hostname as an addition argument, output-all just dumps everything.
Add and update require an 'out' folder to be present and contain Ansible host files in JSON format. These can be generated from Ansible using:

```
ansible -i inventory all -m setup --tree out/
```

These arguments will read all Ansible hostfiles in the 'out' folder and either add hosts to the CMDB if they don't already exist, or update hosts if they do. Update checks for a 'sys_id' first from ServiceNow.

Delete and delete-all are obvious. Delete needs the hostname as an extra argument. Delete-all doesn't, it just deletes everything in the CMDB.

### TODO

There's still a lot to do. In particular, I'm currently generating the ansible out/ folder manually and then ensuring the snow-cmdb utility can find the folder.
This is a bit tedious. I'll do something better.