# go-exploitdb
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/mozqnet/go-exploitdb/blob/master/LICENSE)

[
![](https://images.microbadger.com/badges/version/princechrismc/go-exploitdb.svg)
![](https://img.shields.io/docker/cloud/automated/princechrismc/go-exploitdb.svg)
![](https://img.shields.io/docker/cloud/build/princechrismc/go-exploitdb.svg?logo=docker)
![](https://img.shields.io/docker/pulls/princechrismc/go-exploitdb.svg)
![](https://img.shields.io/docker/stars/princechrismc/go-exploitdb.svg)
![](https://images.microbadger.com/badges/image/princechrismc/go-exploitdb.svg)
](https://hub.docker.com/r/princechrismc/go-exploitdb)

This is a tool for searching Exploits from some Exploit Databases.
Exploits are inserted at sqlite database(go-exploitdb) can be searched by command line interface.
In server mode, a simple Web API can be used.

As the following vulnerabilities database

1. [ExploitDB(OffensiveSecurity)](https://www.exploit-db.com/) by CVE number or Exploit Database ID.
2. [GitHub Repositories](https://github.com/search?o=desc&q=CVE&s=&type=Repositories)
3. [Awesome Cve Poc](https://github.com/qazbnm456/awesome-cve-poc#toc473)

### Docker Deployment
There's a Docker image available `docker pull princechrismc/go-exploitdb`. When using the container, it takes the same arguments as the [normal command line](#Usage).

### Installation for local deployment
###### Requirements
go-exploitdb requires the following packages.
- git
- SQLite3, MySQL, PostgreSQL, Redis
- lastest version of go
    - https://golang.org/doc/install

###### Install go-exploitdb
```bash
$ mkdir -p $GOPATH/src/github.com/mozqnet
$ cd $GOPATH/src/github.com/mozqnet
$ git clone https://github.com/mozqnet/go-exploitdb.git
$ cd go-exploitdb
$ make install
```

----

### Usage: Fetch and Insert Exploit
```bash
$ Fetch the data of exploit

Usage:
  go-exploitdb fetch [command]

Available Commands:
  awesomepoc  Fetch the data of Awesome Poc
  exploitdb   Fetch the data of offensive security exploit db
  githubrepos Fetch the data of github repos

Flags:
  -h, --help   help for fetch

Global Flags:
      --config string       config file (default is $HOME/.go-exploitdb.yaml)
      --dbpath string       /path/to/sqlite3 or SQL connection string
      --dbtype string       Database type to store data in (sqlite3, mysql, postgres, or redis supported)
      --debug               debug mode (default: false)
      --debug-sql           SQL debug mode
      --deep                deep mode extract cve-id from github sources
      --http-proxy string   http://proxy-url:port (default: empty)
      --log-dir string      /path/to/log
      --log-json            output log as JSON
      --quiet               quiet mode (no output)

Use "go-exploitdb fetch [command] --help" for more information about a command.
```

###### Fetch and Insert Offensive Security ExploitDB
```bash
$ go-exploitdb fetch exploitdb
```

###### Deep Fetch and Insert Exploit
- This is very time consuming
- We will further increase the mapping rate between exploit and cveID.
- The number of exploits that can be detected remains unchanged
```bash
$ go-exploitdb fetch exploitdb -deep
```

### Usage: Search Exploits
```bash
$ go-exploitdb search -h

Search the data of exploit

Usage:
  go-exploitdb search [flags]

Flags:
  -h, --help            help for search
      --param string   All Exploits: None  |  by CVE: [CVE-xxxx]  | by ID: [xxxx]  (default: None)
      --type string    All Exploits by CVE: CVE  |  by ID: ID (default: CVE)

Global Flags:
      --config string       config file (default is $HOME/.go-exploitdb.yaml)
      --dbpath string       /path/to/sqlite3 or SQL connection string
      --dbtype string       Database type to store data in (sqlite3, mysql, postgres or redis supported)
      --debug               debug mode (default: false)
      --debug-sql           SQL debug mode
      --deep                deep mode extract cve-id from github sources
      --http-proxy string   http://proxy-url:port (default: empty)
      --log-dir string      /path/to/log
      --log-json            output log as JSON
      --quiet               quiet mode (no output)
```

###### Search Exploits by CVE(ex. CVE-2009-4091)
```bash
$ go-exploitdb search --type CVE --param CVE-2009-4091

Results:
---------------------------------------

[*]CVE-ExploitID Reference:
  CVE: CVE-2009-4091
  Exploit Type: OffensiveSecurity
  Exploit Unique ID: 10180
  URL: https://www.exploit-db.com/exploits/10180
  Description: Simplog 0.9.3.2 - Multiple Vulnerabilities

[*]Exploit Detail Info:
  [*]OffensiveSecurity:
  - Document:
    Path: https://github.com/offensive-security/exploitdb/exploits/php/webapps/10180.txt
    File Type: webapps
---------------------------------------
```

###### Search Exploits by ExploitDB-ID(ex. ExploitDB-ID: 10180)
```bash
$ go-exploitdb search --type ID --param 10180

Results:
---------------------------------------

[*]CVE-ExploitID Reference:
  CVE: CVE-2009-4091
  Exploit Type: OffensiveSecurity
  Exploit Unique ID: 10180
  URL: https://www.exploit-db.com/exploits/10180
  Description: Simplog 0.9.3.2 - Multiple Vulnerabilities

[*]Exploit Detail Info:
  [*]OffensiveSecurity:
  - Document:
    Path: https://github.com/offensive-security/exploitdb/exploits/php/webapps/10180.txt
    File Type: webapps
---------------------------------------

[*]CVE-ExploitID Reference:
  CVE: CVE-2009-4092
  Exploit Type: OffensiveSecurity
  Exploit Unique ID: 10180
  URL: https://www.exploit-db.com/exploits/10180
  Description: Simplog 0.9.3.2 - Multiple Vulnerabilities

[*]Exploit Detail Info:
  [*]OffensiveSecurity:
  - Document:
    Path: https://github.com/offensive-security/exploitdb/exploits/php/webapps/10180.txt
    File Type: webapps
---------------------------------------

[*]CVE-ExploitID Reference:
  CVE: CVE-2009-4093
  Exploit Type: OffensiveSecurity
  Exploit Unique ID: 10180
  URL: https://www.exploit-db.com/exploits/10180
  Description: Simplog 0.9.3.2 - Multiple Vulnerabilities

[*]Exploit Detail Info:
  [*]OffensiveSecurity:
  - Document:
    Path: https://github.com/offensive-security/exploitdb/exploits/php/webapps/10180.txt
    File Type: webapps
---------------------------------------
```

### Usage: Start go-exploitdb as server mode
```bash
$ go-exploitdb server -h

Start go-exploitdb HTTP server

Usage:
  go-exploitdb server [flags]

Flags:
      --bind string   HTTP server bind to IP address (default: loop back interface
  -h, --help          help for server
      --port string   HTTP server port number (default: 1326

Global Flags:
      --config string       config file (default is $HOME/.go-exploitdb.yaml)
      --dbpath string       /path/to/sqlite3 or SQL connection string
      --dbtype string       Database type to store data in (sqlite3, mysql, postgres or redis supported)
      --debug               debug mode (default: false)
      --debug-sql           SQL debug mode
      --deep                deep mode extract cve-id from github sources
      --http-proxy string   http://proxy-url:port (default: empty)
      --log-dir string      /path/to/log
      --log-json            output log as JSON
      --quiet               quiet mode (no output)
```

###### Starting Server
```bash
$ go-exploitdb server
INFO[09-30|15:05:57] Starting HTTP Server...
INFO[09-30|15:05:57] Listening...                             URL=127.0.0.1:1326
```

###### Search Exploits Get by cURL for CVE(ex. CVE-2006-2896)
```
$ curl http://127.0.0.1:1326/cves/CVE-2006-2896 | jq

  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:100   666  100   666    0     0  39340      0 --:--:-- --:--:-- --:--:-- 41625
[
  {
    "ID": 325173,
    "exploit_type": "OffensiveSecurity",
    "exploit_unique_id": "1875",
    "url": "https://www.exploit-db.com/exploits/1875",
    "description": "FunkBoard CF0.71 - 'profile.php' Remote User Pass Change",
    "cve_id": "CVE-2006-2896",
    "offensive_security": {
      "ID": 325173,
      "ExploitID": 325173,
      "exploit_unique_id": "1875",
      "document": {
        "OffensiveSecurityID": 325173,
        "exploit_unique_id": "1875",
        "document_url": "https://github.com/offensive-security/exploitdb/exploits/php/webapps/1875.html",
        "description": "FunkBoard CF0.71 - 'profile.php' Remote User Pass Change",
        "date": "0001-01-01T00:00:00Z",
        "author": "ajann",
        "type": "webapps",
        "palatform": "php",
        "port": ""
      },
      "shell_code": null,
    }
  }
]
```

###### Search Exploits by Unique ID(ex. Exploit Unique ID: 10180)
```
$ curl http://127.0.0.1:1326/id/10180 | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:100  1936  100  1936    0     0  52643      0 --:--:-- --:--:-- --:--:-- 53777
[
  {
    "ID": 334917,
    "exploit_type": "OffensiveSecurity",
    "exploit_unique_id": "10180",
    "url": "https://www.exploit-db.com/exploits/10180",
    "description": "Simplog 0.9.3.2 - Multiple Vulnerabilities",
    "cve_id": "CVE-2009-4091",
    "offensive_security": {
      "ID": 334917,
      "ExploitID": 334917,
      "exploit_unique_id": "10180",
      "document": {
        "OffensiveSecurityID": 334917,
        "exploit_unique_id": "10180",
        "document_url": "https://github.com/offensive-security/exploitdb/exploits/php/webapps/10180.txt",
        "description": "Simplog 0.9.3.2 - Multiple Vulnerabilities",
        "date": "0001-01-01T00:00:00Z",
        "author": "Amol Naik",
        "type": "webapps",
        "palatform": "php",
        "port": ""
      },
      "shell_code": null,
    }
  },
  {
    "ID": 334918,
    "exploit_type": "OffensiveSecurity",
    "exploit_unique_id": "10180",
    "url": "https://www.exploit-db.com/exploits/10180",
    "description": "Simplog 0.9.3.2 - Multiple Vulnerabilities",
    "cve_id": "CVE-2009-4092",
    "offensive_security": {
      "ID": 334917,
      "ExploitID": 334917,
      "exploit_unique_id": "10180",
      "document": {
        "OffensiveSecurityID": 334917,
        "exploit_unique_id": "10180",
        "document_url": "https://github.com/offensive-security/exploitdb/exploits/php/webapps/10180.txt",
        "description": "Simplog 0.9.3.2 - Multiple Vulnerabilities",
        "date": "0001-01-01T00:00:00Z",
        "author": "Amol Naik",
        "type": "webapps",
        "palatform": "php",
        "port": ""
      },
      "shell_code": null,
    }
  },
  {
    "ID": 334919,
    "exploit_type": "OffensiveSecurity",
    "exploit_unique_id": "10180",
    "url": "https://www.exploit-db.com/exploits/10180",
    "description": "Simplog 0.9.3.2 - Multiple Vulnerabilities",
    "cve_id": "CVE-2009-4093",
    "offensive_security": {
      "ID": 334917,
      "ExploitID": 334917,
      "exploit_unique_id": "10180",
      "document": {
        "OffensiveSecurityID": 334917,
        "exploit_unique_id": "10180",
        "document_url": "https://github.com/offensive-security/exploitdb/exploits/php/webapps/10180.txt",
        "description": "Simplog 0.9.3.2 - Multiple Vulnerabilities",
        "date": "0001-01-01T00:00:00Z",
        "author": "Amol Naik",
        "type": "webapps",
        "palatform": "php",
        "port": ""
      },
      "shell_code": null,
    }
  }
]
```
