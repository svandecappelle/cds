+++
title = "CDS - The origins"
weight = 1

+++


Something is missing in that landscape.

You will find below the six requirements of the OVH Continuous Delivery team on the choice of the continuous 
integration / delivery tool.
As you have understood, no existing tool on the market met all needs and CDS was created. 

CDS is used in OVH Production since November 2015 and was opensourced since November 2016.

## 1 - Self-service

At OVH, things move fast: New services and applications are started every week. This generates an ever-growing whishlist of build, test and deployment environments.

Centralizing the handling of these requirements in a Continuous Delivery team is the best way to create an organizational SPOF.

The CDS team believes that this activity must be distributed and that users need to be autonomous to create, configure, fine-tune and troubleshoot their environments. To implement this strong *self-service culture* files an issue everytime a user needs the help of a CDS administrator to achieve a simple day-to-day task.

## 2 - Scalability

CDS is built to scale. And this capability is assessed everyday in a large-scale production environment.
This ability to scale has been made possible thanks to a couple of design principles:

### 2.a – Stateless API servers @ Scale

CDS's API servers are stateless. They do not store anything on the fileSystem. Therefore, they can be deployed as much times as reauired to support the load.

Instances can be spawned and decommissioned dynamically to handle usage surges when required while keeping the cost at its lowest when the platform is underused.

### 2.b – Workers @ Scale – Hatcheries

« 100,000 jobs are launched every week »

{{< figure src="/images/concepts_why_cds_hachery.png" title="a CDS Hatchery" >}}

« Can build and deploy 150 microservices in 8 min »

While respecting the principle that a worker is completely destroyed at the end of his work, achieving the goal of
to deploy 150 microservices in 8 min implies the need for automation on the creation and removal of Workers.

Hatcheries support this role and are dedicated to the technologies used. To date, CDS has:
 
- A **Docker Swarm hatchery**: allows to start Workers CDS automatically on a cluster swarm (or a host where there is a docker).
- An **Openstack hatchery**: allows to start virtual machines. For example, these VMs are mainly used to make of the docker build in a completely isolated way.
- A **VSphere hatchery**: same as the hatchery openstack, this hatchery allows to start virtual machines using VSphere.
- A **Mesos / Marathon hatchery**: allows to start Docker containers with Marathon.
- A **Local hatchery**: allows to launch Workers on a host.
- A **Kubernetes hatchery**: allows to launch Workers CDS in Kubernetes Pods.
 
{{< figure src="/images/concepts_why_cds_hacheries.png" title="CDS Hatcheries" >}}

## 3 - High Availability

« The possible loss of a machine on which the CDS API is deployed is a non-event ».

CDS is a very active OpenSource project, we deploy new versions of CDS in production several times a day
without impacting the work of the users. It is not conceivable to ask our developers to stop
their work so that we can update CDS. This implies some basic principles:

- There is no instance "master": we can update any instance CDS at any time.
- No data is stored on the FileSystem of the CDS API. User data is stored in a PostgreSQL database.
- To limit frequent SQL queries on the database, some data is stored in a Redis.

## 4 - Self-Hosting

You can install CDS on your own infrastructure. You will need: a PostgreSQL database and a Redis. 
You only have to manage backups for the database.
Consult the documentation https://ovh.github.io/cds/hosting/ for more information.

## 5 - Reuse Pipelines

"You have 150 applications built and deployed in the same way, you will have only two pipelines to maintain:
the build pipeline and the deployment pipeline."

This feature is essential, allowing teams to quickly deploy new applications. A system of variables related to applications,
environments, the context of pipelines is available to users to allow to have a minimum of pipeline to maintain.

## 6 - REST API

"Everything must be scriptable, automatable"

An API, a "Command Line", an SDK .. CDS gives the power to its users to create / manage / deploy quickly new workflows on CDS.


# Continuous Delivery Principles

- The process for releasing/deploying software MUST be repeatable and reliable.
- Automate everything!
- If somethings difficult or painful, do it more often
- Keep everything in source control
- Dev Done means “released”
- Build quality in!
- Everybody has responsibility for the release process
- Improve continuously
- Build binaries only once
- Use precisely the same mechanism to deploy to every environment
- Smoke test your deployment
- If anything fails, stop the line

ref. http://devopsnet.com/2011/08/04/continuous-delivery/

Some of these points need to be managed by the Continuous Delivery tool. To these, we add
the principles of pipelines / stages / jobs described above.

CDS natively manages jobs, stages, pipelines, workflows, but also user access rights to environments,
applications, pipelines. Each user, if he wishes, can implement the principles mentioned above using CDS.




=====================


CDS's design empowers users and forsters their freedom. Users are not forced to follow a mapped out path.

CDS ensures its users, who are essentially developers or sysadmin, a great freedom.
Without asking for CDS administrators, any user is able to create new projects, 
applications, pipelines, workflows ... to manage ACLs on all elements.

{{< figure src="/images/concepts_why_cds_languages.png" title="Some Languages used at OVH - non-exhaustive list" >}}

{{< figure src="/images/concepts_why_cds_platforms.png" title="Some Platforms used at OVH - non-exhaustive list" >}}

A CDS Job is executed by a CDS Worker. A worker CDS is a simple binary, you do not need libraries
or particular JVMs on your machine to run it. CDS workers are compatible with Linux, Darwin, OSX,
FreeBSD as well as Windows, in architectures 386, amd64 or arm. The only prerequisite for a CDS Worker to work
is that it can access the CDS API. Any CDS user can therefore launch a CDS Worker where he wishes for him
to execute the jobs he wants - and this - without intervention from the CDS administrators. The user
can thus use its own resources if he wishes.

{{< figure src="/images/concepts_why_cds_workers.png" title="a CDS Worker" >}}


Python 2.7, Python 3.4, Golang 1.9.3 ... Python 2.7, but with this or that library. CDS workers are binaries
executed in a docker image, a virtual machine. Its Build environments are called "Worker Model". If a
user wants to Rust, Scala on a specific version, it is autonomous to create his own
Worker Model and launch it on shared infrastructure. Workers CDS are isolated from each other,
isolation is essentially done by Docker or virtualization, see Scalability below.

Security is therefore an essential element, it is not acceptable to be able to navigate on any system of
file to retrieve files from other teams. By default, when a CDS Worker completes his job, all
is deleted: the temporary files, the container or the virtual machine.