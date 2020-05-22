# Add tolration for all the pods in the namesapce.

This is a mutating web hook that will add ml-pod:Equals:true toleration to all the pods for hiwch the webhook is configured to look into.
When we deploy this webhook, we will configure it via namespace selector.

# Build
USe the attached Dockerfile to build and push the imahe to public repo (quay.io e.g.)
```bash
docker build -t pod-toleration-mutation-webhook:1.0.0 .
docker tag pod-toleration-mutation-webhook:1.0.0 quay.io/ml-aml-workshop/pod-toleration-mutation-webhook:1.0.0
docker push quay.io/ml-aml-workshop/pod-toleration-mutation-webhook:1.0.0

```
# Configure WebHook in minikube
## Create a self signed cert first and deploy as a secret
The name of the service is webhook-server.ml-workshop.svc. so create self signed cert accordingly.
```bash
#create CA
openssl req -nodes -new -x509 -keyout controller_ca.key -out controller_ca.crt -subj "/CN=POd Toleration Mutating Admission Controller Webhook CA"

# create Priv Key
openssl genrsa -out tls.key 2048

# sign the cert with teh CA- see the subject is same as service name
openssl req -new -key tls.key -subj "/CN=webhook-server.ml-workshop.svc" \
    | openssl x509 -req -CA controller_ca.crt -CAkey controller_ca.key -CAcreateserial -out tls.crt

# create k8s secret
kubectl -n ml-workshop create secret tls webhook-server-tls \
    --cert "tls.crt" \
    --key "tls.key"


#### now create a base 64 encoded form of the CA cert to create the trust
openssl x509 -inform PEM -in controller_ca.crt > controller_ca.crt.pem
openssl base64 -in controller_ca.crt.pem -out controller_ca-base64.crt.pem
#optional use following to make it one liner so it can be pasted o the nutating-webhook-config.yaml
cat controller_ca-base64.crt.pem | tr -d '\n' > onelinecert.pem

```
## Use the .openshift folder to create, configure and deploy webhook.
Make sure to replace the caBundle field in the webhook config to use the oneline cert

## Create a new pod in the ns ml-workshop and see the toleration is added automatically.
Make sure to add label in the ml-workshop ns as per the controller. applyMLToleration=true

Use the file targetpod.yaml to create the pod. the pod will have a tolration ml-pod. Note that if there is no taint 
then the pod will not be scheduled.

## Enjoy!
