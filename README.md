# alias-resolve-ldap
alias-resolve-ldap was written to be used as a hook for the [chasquid](https://github.com/albertito/chasquid) SMTP server  
to perform alias address lookups in LDAP  
It returns a single attribute value for a given mail address  
The example config includes a filter for the qmail LDAP schema but can be adapted for others

## Configuration
Configuration is read from a YAML file called 'config.yaml'  
An example can be found in config.example.yaml

## Usage
```
# Return 'cn' for a given alias address
$ ./alias-resolve-ldap alias@example.com
ada
```
