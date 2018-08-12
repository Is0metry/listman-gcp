# listman-gcp
A version of listman built to run in the GCP App Engine.

# How To Use
1. If you haven't already done so, create an account on [Google Cloud Platform](https://cloud.google.com), download the DevTools, and create a new project.
2. Run `go get https://github.com/is0metry/listman-gcp/`
3. Run `gcloud app deploy` to deploy app to the cloud. 
4. Profit.

# TODO
* Move datastore calls into Transactions.
* Multi-user support.
* Option to run as JSON backend.
