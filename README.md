[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/zntrio/harp)](https://goreportcard.com/report/github.com/zntrio/harp)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![GitHub release](https://img.shields.io/github/release/zntrio/harp.svg)](https://github.com/zntrio/harp/releases/)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://zntr.io/harp/graphs/commit-activity)

- [Harp](#harp)
  - [TL;DR.](#tldr)
  - [Visual overview](#visual-overview)
  - [Why harp?](#why-harp)
  - [Use cases](#use-cases)
  - [How does it work?](#how-does-it-work)
    - [Like a Data pipeline but for secret](#like-a-data-pipeline-but-for-secret)
    - [Immutable transformation](#immutable-transformation)
  - [What can I do?](#what-can-i-do)
  - [FAQ](#faq)
  - [License](#license)
- [Build instructions](#build-instructions)
  - [Clone repository](#clone-repository)
  - [Setup dev environment](#setup-dev-environment)
    - [With nix flake](#with-nix-flake)
    - [Non-nix managed environment](#non-nix-managed-environment)
      - [Check your go version](#check-your-go-version)
      - [Install mage](#install-mage)
        - [From source](#from-source)
      - [Bootstrap tools](#bootstrap-tools)
  - [Mage targets](#mage-targets)
- [Plugins](#plugins)
- [Community](#community)

# Harp

Harp is for Harpocrates (Ancient Greek: Ἁρποκράτης) the god of silence, secrets
and confidentiality in the Hellenistic religion. - [Wikipedia](https://en.wikipedia.org/wiki/Harpocrates)

> This tool was initially developed while I was at Elastic, to be able to continue
> to maintain Harp without the upstream dependency, I decided to do a hard-fork
> of the Elastic repository.
>
> I'm going to introduce breaking changes from the Elastic original version.

## TL;DR.

Harp is an innovative toolset that emphasizes `secret management through contracts`. Its primary objective revolves around mitigating value-centric management by offering a structured approach to handling secret data in a reproducible manner. Harp aims to enhance security and efficiency in managing sensitive information by providing a technical stack that describes `contract-managed values` within `pipelines`.

One of Harp's standout features is its ability to establish consistent associations between secrets and `predictable identifiers`. This ensures referencable secrets can be accessed within the system, contributing to a more organized and controlled secret management environment. Including metadata associated with each secret provides comprehensive insights into the nature and context of the managed data, empowering developers with a clear understanding of their data.

Furthermore, Harp leverages a concept known as `Bundles` stored in `immutable containers`, which serve as pivotal elements in facilitating communication between different components within the system. These Bundles enable seamless interaction among various modules, promoting cohesion and integrity in secret management operations.

In addition to its core functionalities, Harp offers a `template engine` that empowers users to generate diverse confidence values such as passwords, passphrases, encryption keys, and more. This feature enhances Harp's flexibility and versatility by enabling users to create tailored configurations based on specific requirements and security considerations.

Harp provides a `robust SDK` that allows developers to integrate its functionalities into their applications seamlessly. This fosters seamless integration and interoperability with existing systems and promotes collaboration and innovation within the software development ecosystem. This aspect of Harp opens up exciting possibilities for developers, inspiring them to explore and create.

In conclusion, Harp represents a comprehensive solution for enhancing secret management practices through contract-based mechanisms. By offering a range of features such as predictable identifiers, metadata associations, bundle storage in immutable containers, template engine capabilities, and an SDK for integration, Harp stands out as a valuable toolset for safeguarding sensitive data and promoting efficient workflows in information security.

## Visual overview

![Visual overview](docs/harp/img/HARP_FLOW.png)

## Why harp?

* Secret management is in essence a collection of processes that must be
  auditable, executable and reproducible for infosec and operation requirements;
* Secret provisioning must be designed with secret rotation as a day one task,
  due to the fact that secret data must be rotated periodically to keep its
  secret property;
* `Developers` should negotiate secret value for the secret consumer they are
  currently developing, by the contract based on a path (reference to the secret)
  and a value specification (for code contract) without the knowledge of the
  final deployed value;
* `Secret Operators` use different set of tools to achieve secret
  management operation which increases the error/secret exposure probability due to
  tool count involved in the process (incompatibility, changes, etc.);
* Without a defined secret naming convention, the secret storage becomes difficult to
  handle in time (naming is hard) and secret naming could not be helped to
  get a consistent, reliable and flexible secret tree;
* Secret storage backend can use various implementations in different environments
  and should be provisioned consistently;
* When you use `Terraform` for secret management, you have the cleartext value
  stored in the state. To protect the state you have to deploy a complex infrastructure.
  To simplify this we use harp for secret provisioning and use the secret reference
  in the Terraform topology.

## Use cases

* You want to have a `single secret value` and you are asking yourself
  `how to generate a strong password` - Harp has a template engine with secret
  value generation functions to allow you to generate such values.
* You have `thousands secrets` to handle to deploy your platform/customers
  `on multiple cloud providers` with `different secret storages` - Harp will help you
  to define consistent secret provisioning bundles and pipelines.
* You need a `ephemeral secret storage` to `bootstrap` your long term cloud
  secret storage - Harp will help you to create
  secret containers that can be consumed on deployment.
* You want to `migrate massively` your secrets from one secret storage to
  another - Harp provides you a secret container to store these secrets while
  they are going to be distributed in other secret storage implementations.
* You have to `alter/modifiy` a secret (rotation/deprecation/renewal) - Harp
  provides you a `GitOps-able` secret `storage agnostic operation set`, so that you
  can define a specification to describe how your secret operation is going to
  be applied offline on the secret container.

## How does it work?

![Secret management Pipeline](docs/harp/img/SM-HARP-PIPELINE.png)

### Like a Data pipeline but for secret

`harp` allows you to handle secrets using deterministic pipelines expressed
using an atomic series of CLI operations applied to a commonly shared container
immutable and standalone file system used to store secret collection (Bundle)
generated from a template engine via user specification, or external secret
value coming from files or external secret storage.

![Pipelines](docs/harp/img/SM-HARP.png)

These pipelines use the immutable container file system as a data exchange
protocol and could be extended for new input, intermediary operation or output
via plugins created with the `harp` SDK.

### Immutable transformation

Each applied transformation creates a container with transformed data inside.
This will enforce container reproducibility by eliminating cumulative
side effects applied to the same container.

The container handles for you the confidentiality and integrity protection applied
to the secret collection stored inside and manipulated by copy during the
pipeline execution.

## What can I do?

> New to harp, let's start with [onboarding tutorial](docs/onboarding/README.md) !
> TL;DR - [Features overview](FEATURES.md)

Harp provides :

* A methodology to design your secret management;
  * Secret naming convention (CSO);
  * A defined common language and complete processes to achieve secret management
    operations;
* A SDK to create your own tools to orchestrate your secret management pipelines;
  * A container manipulation library exposed as `zntr.io/harp/v2/pkg/container`;
  * A secret bundle specification to store and manipulate secrets exposed as `zntr.io/harp/v2/pkg/bundle`;
  * An `on-steroid` template engine exposed as `zntr.io/harp/v2/pkg/template`
  * A path name validation library exposed as `zntr.io/harp/v2/pkg/cso`
* A CLI for secret management implementation
  * CI/CD integration;
  * Based on human-readable definitions (YAML);
  * In order to create auditable and reproducible pipelines.
  * An extensible tool which can be enhanced via [plugins](https://github.com/zntrio/harp-plugins).

And allows :

* Bundle level operations
  * Create a bundle from scratch / template / JSON (more via plugins);
  * Generate a complete bundle using a YAML Descriptor (`BundleTemplate`) to describe secret and their usages;
  * Read value stored in the K/V virtual file system;
  * Update the K/V virtual file system;
  * Reproducible patch applied on immutable container (copy-on-write);
  * Import / Export to Vault.
* Immutable container level operations
  * Seal / Unseal a container for integrity and confidentiality property conservation
    to enforce at-rest encryption (aes256-gcm96 or chacha20-poly1305);
  * Multiple identities sealing algorithm;

## FAQ

* Is it used internally at zntrio? - Yes. It is used to generate bootstrap
  secrets used to bootstrap the new region infrastructure components.
  #ChickenEggProblem

* Harp is only supporting `Vault`? - No, it has been published with only vault
  support built-in, but it supports many other secret storage implementations via
  plugins.

* What's the difference with `Vault`? - HashiCorp Vault is an encrypted highly
  available K/V store with advanced authorization engine, it doesn't handle
  secret provisioning for you. You can't ask Vault to generate secrets for your
  application and store them using a defined logic. Harp is filling this
  requirement.

## License

`harp` artifacts and source code is released under [Apache 2.0 Software License](LICENSE).

# Build instructions

Download a [release](https://github.com/zntrio/harp/releases) or build from source.

## Clone repository

```sh
$ git clone git@github.com:zntrio/harp.git
$ export HARP_REPOSITORY=$(pwd)/harp
```

## Setup dev environment

### With nix flake

Install `nix` on your system, if not already installed.

```sh
$ sudo install -d -m755 -o $(id -u) -g $(id -g) /nix
$ curl -L https://nixos.org/nix/install | sh
```

> More information? - <https://nixos.wiki/wiki/Nix_Installation_Guide>

```sh
$ cd $HARP_REPOSITORY
$ nix develop
```

### Non-nix managed environment

#### Check your go version

> Only last 2 minor versions of a major are supported.

`Harp` is compiled with :

```sh
$ go version
go version go1.21 linux/amd64
```

> Simple go version manager - <https://github.com/stefanmaric/g>

#### Install mage

[Mage](https://magefile.org/) is an alternative to Make where language used is Go.
You can install it using 2 different methods.

##### From source

```sh
# Install mage
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
```

#### Bootstrap tools

```sh
# Go to tools submodule
cd $HARP_REPOSITORY/tools
# Resolve dependencies
go mod tidy
go mod vendor
# Pull tools sources, compile them and install executable in tools/bin
mage
```

## Mage targets

```sh
❯ mage -l
Targets:
  api:generate     protobuf objects from proto definitions.
  build*           harp executable.
  code:format      source code and process imports.
  code:generate    SDK code (mocks, tests, etc.)
  code:licenser    apply copyright banner to source code.
  code:lint        code using golangci-lint.
  compile          harp code to create an executable.
  docker:harp      build harp docker image
  docker:tools     prepares docker images with go toolchain and project tools.
  homebrew         generates homebrew formula from compiled artifacts.
  release          harp version and cross-compile code to produce all artifacts.
  releaser:harp    releases harp artifacts using docker pipeline.
  test:cli         Test harp application.
  test:unit        Test harp application.

* default target
```

# Plugins

You can find more Harp feature extensions - <https://github.com/zntrio/harp-plugins>

# Community

Here is the list of external projects used as inspiration :

* [Kubernetes](https://github.com/kubernetes/)
* [Helm](https://github.com/helm/)
* [Open Policy Agent ConfTest](https://github.com/open-policy-agent/conftest)
* [SaltPack](https://github.com/keybase/saltpack)
* [Hashicorp Vault](https://github.com/hashicorp/vault)
* [AWS SDK Go](https://github.com/aws/aws-sdk-go)
