# Pivotal Cloud Foundry tile for SnappyData

## Building the tile

- Install the [tile generator](https://docs.pivotal.io/tiledev/2-1/tile-generator.html) tool on your local machine.

- You also need to install [latest bosh CLI](https://bosh.io/docs/#downloads).

- Move to [tiles/snappydata](snappydata/) directory of the cloned repository where the tile.yml is present and build the tile.
    - `$ tile build`

- Optionally, one could specify version for the tile.
    - `$ tile build 1.0`

- Upload the generated .pivotal file from product/ folder to Pivotal Ops Manager.

## Configuring the tile

Once the .pivotal file imported into Pivotal Ops Manager, add the tile to the dashboard by clicking on the '+' sign.
Click on the tile to configure it.

Here, at a minimum, you need to specify 1. the credentials for the Kubernetes/PKS cluster you would launch the SnappyData chart on and also 2. select the appropriate network.
The credentials to connect to Kubernetes/PKS cluster include CA Cert, Cluster Token and Cluster url. (Note that the CA cert needs to be base64 decoded after fetching it from your kubeconfig file.)

Save these configurations and return to the installation dashboard and hit 'Apply changes'.


## Creating and consuming a service

Once installed, users can create the service instance of SnappyData which essentially installs the Helm chart on the
Kubernetes/PKS cluster provided during tile configuration.

- Install the [CF CLI](https://docs.cloudfoundry.org/cf-cli/install-go-cli.html) and log in to your Cloud Foundry's API server.
    - `$ cf login -a https://<api-server-name> --skip-ssl-validation`
- You can view that the SnappyData service broker is now visible.
    - `$ cf service-brokers`
- Create a service instance using the broker. It'll launch the SnappyData cluster on your configured Kubernetes/PKS cluster.
- Currently, we have three plans for the SnappyData service: 1. small (default), 2. medium and 3. large. All of these start the cluster with one locator, one lead and two servers but with different memory.
    - `$ cf create-service snappydata small snappydata_small`
- We'll add to or change these plans in future.
- Now that the service is available, you can bind it to any of your running apps.
    - `$ cf bs <app-name> snappydata_small`
    - `$ cf restage <app-name>`
- Now, your app has access to information about the service via environment variable VCAP_SERVICES. Typically, all the Kubernetes services created by the chart are included in this environment variable.
    - `$ cf env <app-name>`

