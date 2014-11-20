# pw
#### Command-line tool to store and retrieve passwords in a structured JSON file in users home

Learning Go programming this is one of my first tiny projects. pw stores and retrieves user passwords in ```~/.pw.json```. Run

```
git clone https://github.com/danbim/pw.git pw
cd pw
go build
go install
```

to install pw and make sure ```$GOPATH/bin``` is in your ```PATH``` environment variable. Then, run 

```
pw site password "site description"
```

to store the password ```password``` for site ```site``` with description ```site description``` (description is optional). Run

```
pw site
# password
```

to retrieve the password from the storage.

## Mac OS X comfort
If you're on Mac you can use the nifty ```pbcopy``` command that will copy its ```stdin``` to the Mac OS X clipboard so you can quickly paste it anywhere you need it (typically a Browser window). 

If you're additionally running the fish shell you might want to try this little function from ```~/.config/fish/functions/pwc.fish```:

```
function pwc
	set PASSWORD (pw $argv)
	echo -n $PASSWORD | pbcopy
end
```

that will do the trick for you so that running ```pwc site password``` and ```pwc site``` will automatically leave you with the password in the clipboard.

## License
```
-----------------------------------------------------------------------------
"THE BEER-WARE LICENSE" (Revision 42):
<daniel@bimschas.com> wrote this file.  As long as you retain this notice you
can do whatever you want with this stuff. If we meet some day, and you think
this stuff is worth it, you can buy me a beer in return.      Daniel Bimschas
-----------------------------------------------------------------------------
```

## Disclaimers
Seriously, this is just a tiny learning project for writing Go code. It just does the most basic things of JSON serialization and file disk access and is in no way a safe way to store passwords. I use it for storing non-critical site passwords as my disk is encrypted anyways and I don't want to sync my passwords with any Chrome cloud service.
