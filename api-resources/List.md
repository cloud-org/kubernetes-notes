<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [List](#list)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

### List

I think this is because List is not actually a "resource". Each resource can have an associated list type, PodList or CronJobList, but those are not actually resources. And each of those lists is represented in yaml by kind: List.

- https://github.com/kubernetes/kubectl/issues/837
