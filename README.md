# Kin

This is a small utility I have been needing for some time to quickly access various k8s clusters across Teleport.
This comes from a frustration that `tsh` doesn't provide shell completions.

The program simply pulls down the k8s clusters your user has access to and presents them in a nice selectable prompt.

### Releasing

The repo uses a release workflow via GitHub actions that invokes goreleaser.
To trigger a release, simply push a new semver tag to the origin.

```
git tag -a vX.X.X -m "Release note"
git push origin vX.X.X
```

Changelogs are generated automatically from git history.

### Installing

Installing the tool can be done directly from the GitHub releases for your desired version and architecture.
Here is an example of the process with a Darwin binary for my M1 Mac:

```
curl -OL https://github.com/VinceDeslo/kin/releases/download/vX.X.X/kin_Darwin_arm64.tar.gz
tar -xvzf kin_Darwin_arm64.tar.gz
chmod +x kin
sudo mv kin /usr/local/bin/
```

### Running

Now you can directly invoke it with your Teleport proxy and cluster values.

```
kin --proxy <my.teleport.proxy:port> --cluster <cluster-name>
```

You can also set it and forget it with a small shell alias. Here's an example for zsh.

```
echo 'alias kin="kin --proxy <my.teleport.proxy:port> --cluster <cluster-name>"' >> ~/.zshrc
source ~/.zshrc
```

### Uninstalling

Just remove the binary from your local bin.

```
rm /usr/local/bin/kin
```

> Make sure to clean up any alias values you set to invoke it.

### Notes

Yes I should probably write an install script, go package or homebrew tap of some sorts. 
Will I, maybe? Am I lazy though, definitely.
