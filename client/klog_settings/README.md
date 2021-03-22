<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [klog-settings](#klog-settings)
  - [usage](#usage)
  - [acknowledgement](#acknowledgement)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## klog-settings 

```sh
$ ./klog_settings -h 
Usage of ./klog_settings:
  -add_dir_header
        If true, adds the file directory to the header of the log messages
  -alsologtostderr
        log to standard error as well as files
  -log_backtrace_at value
        when logging hits line file:N, emit a stack trace
  -log_dir string
        If non-empty, write log files in this directory
  -log_file string
        If non-empty, use this log file
  -log_file_max_size uint
        Defines the maximum size a log file can grow to. Unit is megabytes. If the value is 0, the maximum file size is unlimited. (default 1800)
  -logtostderr
        log to standard error instead of files (default true)
  -skip_headers
        If true, avoid header prefixes in the log messages
  -skip_log_headers
        If true, avoid headers when opening log files
  -stderrthreshold value
        logs at or above this threshold go to stderr (default 2)
  -v value
        number for the log level verbosity
  -vmodule value
        comma-separated list of pattern=N settings for file-filtered logging

```

### usage

需要设置 `-logtostderr=false` 才会 flush 成功

### acknowledgement

- [k8s-klog](https://bamboox.online/2019/10/01/k8s-klog/)
