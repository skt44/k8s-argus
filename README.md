<p align="center"><a href=""><img src="./logo.png"></a></p>

---

<p align="center">
  <a href="https://godoc.org/github.com/logicmonitor/k8s-argus"><img alt="GoDoc" src="http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square"></a>
  <a href="https://travis-ci.org/logicmonitor/k8s-argus"><img alt="Travis" src="https://img.shields.io/travis/logicmonitor/k8s-argus.svg?style=flat-square"></a>
  <a href="https://codecov.io/gh/logicmonitor/k8s-argus"><img alt="Codecov" src="https://img.shields.io/codecov/c/github/logicmonitor/k8s-argus.svg?style=flat-square"></a>
  <a href="https://goreportcard.com/report/github.com/logicmonitor/k8s-argus"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/logicmonitor/k8s-argus?style=flat-square"></a>
  <a href="https://github.com/logicmonitor/k8s-argus/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/logicmonitor/argus.svg?style=flat-square"></a>
  <a href="https://github.com/logicmonitor/k8s-argus/releases/latest"><img alt="GitHub (pre-)release" src="https://img.shields.io/github/release/logicmonitor/argus/all.svg?style=flat-square"></a>
</p>

**Argus** is a tool for monitoring Kubernetes with [LogicMonitor](https://www.logicmonitor.com). Some of the key features of Argus are:
-   **Automated Device Management:** Discovers etcd members, cluster nodes, services, and pods, and automates the management of their lifecylce as [Devices](https://www.logicmonitor.com/support/devices/) in LogicMonitor.
-   **Automated Device Property Management:** Discovers Kubernetes resource metadata and maps them to [Device Properties](https://www.logicmonitor.com/support/devices/adding-managing-devices/device-properties/).
-   **Automated Device Tree Management:** Leverages Kubernetes events to provide real-time accuracy of a cluster's resources in LogicMonitor, and organizes them by resource type and namespace in the [Device Tree](https://www.logicmonitor.com/support/devices/devices-page-overview/navigating-devices/).

See the [documentation](https://logicmonitor.github.io/k8s-argus) to discover more about Argus.
> **Note:** Argus is a community driven project. LogicMonitor support will not assist in any issues related to Argus.

Developing Argus
----------------
To build Argus, run the following in the root directory:
```
$ make
```
> **Note:** The Dockerfile uses multi-stage builds. Docker 17.05.0 or greater is required.

To build the documentation, run the following in the `docs` directory:
```
$ hugo
```


### License
[![license](https://img.shields.io/github/license/logicmonitor/k8s-argus.svg?style=flat-square)](https://github.com/logicmonitor/k8s-argus/blob/master/LICENSE)
