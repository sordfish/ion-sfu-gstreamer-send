## Usage

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

  helm repo add sordfish https://sordfish.github.io/helm-charts

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages.  You can then run `helm search repo
sordfish` to see the charts.

To install the ion-sfu-gstreamer-send chart:

    helm install my-ion-sfu-gstreamer-send sordfish/ion-sfu-gstreamer-send

To uninstall the chart:

    helm delete my-ion-sfu-gstreamer-send