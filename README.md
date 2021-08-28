# gofort

Yet another [`fortune`](https://wiki.archlinux.org/title/Fortune) clone, written in Go. 

Besides that, what's different here?

1. Uses simple text files. No effort has been done to re-implement parsing the strfiles as in original fortune and some other clones.
2. Uses a streaming approach combined with reservoir sampling instead of reading whole files into memory first.
3. Self-contained binary, only using the standard library. Fortune files are bundled via `go:embed`.

> **tl;dr:** simple, efficient and self-contained.

## Bootstrapping

Just some ideas where you could use this.

### initramfs

Hook it into [intramfs](https://wiki.archlinux.org/title/mkinitcpio) when you're feeling extra lucky.
Take a look [at this repository](https://github.com/kdevo/mkinitcpio-asciilogo).

### Bash

```sh
echo "~/path/to/gofort" >> ~/.bashrc
```

### Zsh

```sh
echo "~/path/to/gofort" >> ~/.zshrc
```

### motd

The motd setup is distro-specific:

- On Debian and Ubuntu, the [`update-motd` mechanism can be used](https://wiki.ubuntu.com/UpdateMotd) with a script that launches motd. 
- On Arch, it's a bit more complicated: 
    - One way to do it is described [in this repo](https://github.com/lfelipe1501/Arch-MOTD/wiki/Installation-Guide) (replace the `update_motd.sh` with your own script). 
    - If you don't want to use PAM, you can of course also use a cronjob or systemd timer to regularly populate `/etc/motd` with a fortune by executing gofort.