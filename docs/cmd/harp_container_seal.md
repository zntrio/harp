## harp container seal

Seal a secret container

```
harp container seal [flags]
```

### Options

```
      --dckd-master-key string      Master key used for deterministic container key derivation
      --dckd-target string          Target parameter for deterministic container key derivation
  -h, --help                        help for seal
      --identity stringArray        Identity allowed to unseal
      --identity-file stringArray   Files with identity allowed to unseal
      --in string                   Unsealed container input ('-' for stdin or filename)
      --json                        Display seal info as json
      --no-container-identity       Disable container identity
      --out string                  Sealed container output ('-' for stdout or filename)
      --pre-shared-key string       Use a pre-shared-key to seal the container to act as a second factor
      --seal-version uint           Select the sealing strategy version (1:modern, 2:fips-compliant) (default 1)
```

### Options inherited from parent commands

```
      --config string   config file
```

### SEE ALSO

* [harp container](harp_container.md)	 - Secret container commands

