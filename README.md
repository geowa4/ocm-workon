# OCM Workon

A utility to help people working on multiple OpenShift clusters--potentially at the same time.

## Building

Build with make.

```shell
make clean build
```

This will create a file `./bin/ocm-workon`.
Copy this onto your `$PATH` to run it.

```shell
ocm workon
```

## Configuration

Most arguments can be read from environment variables or via a configuration file located in `${$XDG_CONFIG_DIR:-$HOME/.config}/ocm/workon.yaml`.
While not required, [direnv](https://direnv.net/) is highly recommended.
Other useful additions are [asdf](https://asdf-vm.com/) and [asdf-direnv](https://github.com/asdf-community/asdf-direnv).

This is a good starting point if you have zsh installed at /bin/zsh.

```yaml
cluster_base_directory: /home/geowa4/Clusters
cluster_use_direnv: true
```

If you only ever work on production clusters, set `cluster_production: true` as well.

## Usage

Before running any other commands it is recommended you initialize the base directory where cluster configurations will be stored.
In this directory will be a .zshrc file that you can use to run certain commands when you start working on a cluster.

```shell
ocm workon init
```

Supply one of the ID, external ID, or name of the cluster and a zsh shell will 

```shell
ocm workon cluster geowa4-test
```

This will record that you have worked on this cluster in a SQLite database in your cluster base directory.
To query that, run the list command.

```shell
ocm workon recent
```

## TODO

- [ ] command to note that the cluster I'm working on should generate a compliance alert
- [ ] search for recent clusters worked that I expected to have generated a compliance alert
- [ ] use a logger, especially to report errors and warnings
