# logrotate-anonymize-ip

## Usage
```
Anon-IP

Usage: anonip [OPTION]

Options: 
  -d    Prints debug output
  -h    Show help message
```

Any number of files can be attached to the script:

`anonip /var/log/nginx/access.log /var/log/nginx/error.log`

or

`anonip /var/log/nginx/*.log`

## Example

```
/var/log/nginx/access.log
/var/log/nginx/error.log {
	daily
    missingok
    rotate 7
    compress
    delaycompress
    notifempty
    create 0640 www-data adm
    sharedscripts
    # anonip block
    prerotate
        /usr/local/bin/anonip $1
    endscript
    # anonip block
    postrotate
        invoke-rc.d nginx rotate >/dev/null 2>&1
    endscript
}
```