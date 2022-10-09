## harp share put

Put secret in Vault Cubbyhole and return a wrapped token

```
harp share put [flags]
```

### Options

```
  -h, --help               help for put
      --in string          Input path ('-' for stdin or filename) (default "-")
      --json               Display result as json
      --namespace string   Vault namespace
      --prefix string      Vault backend prefix (default "cubbyhole")
      --ttl duration       Token expiration (default 30s)
```

### Options inherited from parent commands

```
      --config string   config file
```

### SEE ALSO

* [harp share](harp_share.md)	 - Share secret using Vault Cubbyhole

