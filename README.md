# SSH WEB PROXY

The intention of this tool is to simplify viewing of a static web content on a remote server when the remote server only allows ssh connection.

Example is using of a development VM in the cloud. Tools like code coverage checkers, test runners or security scnaners might produce resutls in a html/css format. Instead of downloading them manually you can use SSH WEB PROXY and view the files in your local browser.

When you run the tool, it will create a local HTTP server. When you issue a request to the HTTP server it wil use the SFTP client to retrieve the content from your specified host and serve it over HTTP.

*NOTE: Tested on Linux and MacOS only but compiled for Windows as well*

## Usage
You can download the source code and build/run the solution yourself or you can can download one of the [pre-built releases for your platform](https://github.com/rvidis/ssh-web-proxy/releases)

Once you have the binary you can run it from your terminal to see all the options:
```
./swp -h
```

To run the server replace with your own details:
```
./swp -i ~/.ssh/my_id_rsa -hd /my/site -hu myuser
```

Once HTTP server is running you can request your content, for example from your browser:
```
http://localhost:8080/index.html
```

## Options
```
-ha string
        Host Address (default "0.0.0.0")
-hd string
    Host base Directory - directory containing the content on the host
-hp int
    Host Port - port to use with remote host (default 22)
-hu string
    Host User
-i string
    Identity file
-sp int
    Local HTTP server port (default 8080)
```

- **ha - Host Address** - specify host address if it is different from the default 0.0.0.0
- **hd - Host Base Directory** - this is the directory that contains the content you wish to serve over HTTP. The full path gets worked out by combining this path to gether with the path from the URL. For example, if you set `-hd /my/site` and navigate to `http://localhost:22/content/index.html` the path will resolve to `/my/site/content/index.html`
- **hp - Host Port** - specify host port if is is different from the default port 22
- **hu - Host User** - Host user to be used during the SSH authentication
- **i - Identity File** - Path to the identity file that contains SSH key, i.e. `-i ~/.ssh/id_rsa`
- **sp - Local HTTP Server Port** - Local HTTP server will always open on a 0.0.0.0 address (localhost) and you can specify port if you wish for it to be different from the default port 8080