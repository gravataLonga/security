# Build    
```
env GOOS=linux GOARCH=amd64 go build -o ccsecurity main.go  
```  

# Upload
```
sftp <user>@<server>  
> put ccsecurity  
> exit  
```

# Put binary globally  
```
mv ccsecurity /usr/local/bin
```

# Create a new DIGEST file
```
ccsecurity digest -c -o file.chk "./\<path\>/\*\*/\*.php"
```

# Check Digest
```
ccsecurity digest -o file.chk "./\<path\>/\*\*/\*.php"
```
