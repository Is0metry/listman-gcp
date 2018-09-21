# listman-gcp
A simple recursive list app built to run on the GCP app engine. 

# How To Install
1. If you haven't already done so, create an account on [Google Cloud Platform](https://cloud.google.com), download the DevTools, and create a new project.
2. Run `go get https://github.com/is0metry/listman-gcp/`
3. Run `gcloud app deploy` to deploy app to the cloud.
## OR
  run command `devapp_server.py app.yaml` in source directory (note: if running on the dev server, imports in the source files must remove file path from $GOPATH, i.e. `github.com/Is0metry/listman-gcp/lists` would become simply `lists`. I don't understand why either).
4. Profit.

# Use
Your front end should send JSON requests (see `handlers/jsontype.go` for `OperationRequest` format.
# TODO
* Move datastore calls into Transactions.
* Multi-user support.
* Enable use of template-based webapp.
