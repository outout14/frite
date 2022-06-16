# FRITE 

F.R.I.T.E: Fast, Reliable lInk Transformation Engine 

## Goal 

FRITE aims to make it possible to create simple link shortcuts using a regular text file. 

This software is available as a simple binary and does not require any PHP or NodeJS engine.

### Use

#### Installation 

Get the latest binary or package for platform. 

#### APT Repository 

There is currently no support for GPG signing on the repository. However, the repository connexion is encrpyted using https (you might need to install ``apt-transport-https`` on your machine). 

```
deb [trusted=yes] https://repos.enpls.org/apt/ /
```

#### Links format 

By default, in packaged FRITE distributions, the systemd service will use ``/etc/frite/links.txt`` as the links file. 

```
demo    https://example.com
<short url>     <destination url>
``` 

You can separate links from shortcuts with tabs or spaces.

(Do not forget to put the protocol (eg : ``https://`` before the destination URL)

#### Configuration 

There is no config file for FRITE, everything has to be configured using CLI arguments.

Configuration of FRITE can be given using ``frite -h``

Output of ``frite -h``: 
```
Usage of ./frite:
  -debug
    	Enable debug logs
  -http-dir string
    	If proxied in subfolder (default "/")
  -http-host string
    	HTTP Listen IP (default "127.0.0.1")
  -http-port int
    	HTTP Listen port (default 8080)
  -links string
    	Path to the links file (default "links.txt")
  -test
    	Test links file for syntax
```

FRITE will automaticly add a trailing slash to your ``-http-dir`` path. 

### LICENSE 

This software is available under the [MIT license](/LICENSE)
