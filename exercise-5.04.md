# Exercise 5.04

Comparison of two Kubernetes runtimes: [VMware Tanzu](https://tanzu.vmware.com/tanzu) and [RedHat OpenShift](https://www.openshift.com/).
For the sake of exercise, I assume that Tanzu is "better" and argue why it should be chosen.

## VMware Tanzu

* Works on top of VMware vSphere, making it possible to run Kubernetes within an existing on-premises virtualization infrastructure
* Basic edition included in vSphere, making initial investment small
* While being open-source aligned, provides preconfigured components that integrate into the VMware ecosystem, such as:
    * Logging (Fluent Bit)
    * Observability (Prometheus and Grafana or Wavefront)
    * Analytics
    * Service mesh (VMware Tanzu Service Mesh)
* Supports expanding to a hybrid cloud setup with Tanzu Kubernetes Grid

## RedHat OpenShift

* CoreOS on master node, RHEL on worker nodes
* Includes preconfigured components for logging, service mesh, observability, container registry and so on
* Supports multiple clouds

## Conclusion

Since Tanzu provides deeper integration with the VMware ecosystem, it is a better choice when VMware virtualization is already used.
