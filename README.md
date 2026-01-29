# testsuite-diff

This is a small program to answer the question "What testsuites are running in this product scenario that aren't in this other product scenario", takes a source and target yaml (from openqa's job templates) and figures out the diff between product scenarios

## building

There is a makefile for convenience with a single target however if preferred:

```
go build -o testsuite-diff ./cmd/testsuite-diff
```

Then copy the binary either somewhere in $PATH or to where you want to use it

## Usage

Assuming you are in a local copy of https://github.com/os-autoinst/opensuse-jobgroups and 

```
testsuite-diff --source job_groups/opensuse_tumbleweed.yaml \
 --source-product opensuse-Tumbleweed-agama-installer-x86_64 \
 --target job_groups/opensuse_leap_16.0.yaml \
 --target-product opensuse-offline-installer-x86_64
```

Will generate the following output

```
Diff: What is in 'job_groups/opensuse_tumbleweed.yaml' (opensuse-Tumbleweed-agama-installer-x86_64) but missing in 'job_groups/opensuse_leap_16.0.yaml' (opensuse-offline-installer-x86_64)?
--------------------------------------------------------------------------------
--- In 'job_groups/opensuse_tumbleweed.yaml' but MISSING in 'job_groups/opensuse_leap_16.0.yaml' ---
- autofs_client
- autofs_server
- base_os_gnome
- base_os_kde
- base_os_xfce
- boot_to_snapshot
<...>
- system_check
- systemd-networkd
- textmode_extra_tests_development
- textmode_extra_tests_kernel
- textmode_extra_tests_networking
- textmode_extra_tests_utils
- toolchain_zypper
- toolkits
- toolkits-kde
--- In 'job_groups/opensuse_tumbleweed.yaml' but MISSING in 'job_groups/opensuse_leap_16.0.yaml' ---
+ gnome-agama
+ kde-agama

```