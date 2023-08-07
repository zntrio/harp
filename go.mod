module zntr.io/harp/v2

go 1.20

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.6.18
	github.com/opencontainers/runc => github.com/opencontainers/runc v1.1.5
	github.com/satori/go.uuid => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
)

// GHSA
require (
	github.com/docker/docker v23.0.3+incompatible // indirect
	github.com/opencontainers/image-spec v1.1.0-rc4 // indirect
	github.com/opencontainers/runc v1.1.5 // indirect
)

require (
	filippo.io/age v1.1.1
	filippo.io/edwards25519 v1.0.0
	github.com/MakeNowJust/heredoc/v2 v2.0.1
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/Masterminds/sprig/v3 v3.2.3
	github.com/alessio/shellescape v1.4.2
	github.com/awnumar/memguard v0.22.3
	github.com/basgys/goxml2json v1.1.0
	github.com/cloudflare/tableflip v1.2.3
	github.com/common-nighthawk/go-figure v0.0.0-20210622060536-734e95fb86be
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/dchest/uniuri v1.2.0
	github.com/essentialkaos/branca v1.3.4
	github.com/fatih/color v1.15.0
	github.com/fatih/structs v1.1.0
	github.com/fernet/fernet-go v0.0.0-20211208181803-9f70042a33ee
	github.com/go-akka/configuration v0.0.0-20200606091224-a002c0330665
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/go-zookeeper/zk v1.0.3
	github.com/gobwas/glob v0.2.3
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.3
	github.com/golang/snappy v0.0.4
	github.com/google/cel-go v0.17.1
	github.com/google/go-cmp v0.5.9
	github.com/google/go-github/v42 v42.0.0
	github.com/google/gofuzz v1.2.0
	github.com/google/gops v0.3.27
	github.com/gosimple/slug v1.13.1
	github.com/hashicorp/consul/api v1.23.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/hcl/v2 v2.17.0
	github.com/hashicorp/vault/api v1.9.2
	github.com/iancoleman/strcase v0.3.0
	github.com/imdario/mergo v0.3.16
	github.com/jmespath/go-jmespath v0.4.0
	github.com/klauspost/compress v1.16.7
	github.com/lytics/base62 v0.0.0-20180808010106-0ee4de5a5d6d
	github.com/magefile/mage v1.15.0
	github.com/mcuadros/go-defaults v1.2.0
	github.com/miscreant/miscreant.go v0.0.0-20200214223636-26d376326b75
	github.com/oklog/run v1.1.0
	github.com/open-policy-agent/opa v0.55.0
	github.com/ory/dockertest/v3 v3.10.0
	github.com/pelletier/go-toml v1.9.5
	github.com/pierrec/lz4 v2.6.1+incompatible
	github.com/pkg/errors v0.9.1
	github.com/psanford/memfs v0.0.0-20230130182539-4dbf7e3e865e
	github.com/sebdah/goldie v1.0.0
	github.com/sethvargo/go-diceware v0.3.0
	github.com/sethvargo/go-password v0.2.0
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.16.0
	github.com/stretchr/testify v1.8.4
	github.com/ulikunitz/xz v0.5.11
	github.com/xeipuuv/gojsonschema v1.2.0
	github.com/zclconf/go-cty v1.13.2
	gitlab.com/NebulousLabs/merkletree v0.0.0-20200118113624-07fbf710afc4
	go.etcd.io/etcd/client/v3 v3.5.9
	go.step.sm/crypto v0.34.0
	go.uber.org/zap v1.25.0
	golang.org/x/crypto v0.11.0
	golang.org/x/oauth2 v0.10.0
	golang.org/x/sync v0.3.0
	golang.org/x/sys v0.11.0
	golang.org/x/term v0.10.0
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
	gopkg.in/square/go-jose.v2 v2.6.0
	gopkg.in/yaml.v3 v3.0.1
	sigs.k8s.io/yaml v1.3.0
	zntr.io/paseto v1.2.0
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr/v4 v4.0.0-20230305170008-8188dc5388df // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200108200545-475eaeb16496 // indirect
	github.com/awnumar/memcall v0.1.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/cenkalti/backoff/v3 v3.0.0 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/containerd/continuity v0.3.0 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/docker/cli v23.0.3+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.6 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.6 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/huandu/xstrings v1.3.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/term v0.0.0-20221205130635-1aeaba878587 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.16.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/tchap/go-patricia/v2 v2.3.1 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/yashtewari/glob-intersection v0.2.0 // indirect
	go.etcd.io/etcd/api/v3 v3.5.9 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.9 // indirect
	go.opentelemetry.io/otel v1.16.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.opentelemetry.io/otel/sdk v1.16.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	go.uber.org/goleak v1.2.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29 // indirect
	golang.org/x/mod v0.9.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.7.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230706204954-ccb25ca9f130 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230706204954-ccb25ca9f130 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230720185612-659f7aaaa771 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
