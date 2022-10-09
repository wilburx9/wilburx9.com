# wilburx9.dev
[![Test Workflow](https://github.com/wilburt/wilburx9.dev/actions/workflows/test.yaml/badge.svg)](https://github.com/wilburt/wilburx9.dev/actions/workflows/test.yaml)
[![Coverage Status](https://coveralls.io/repos/github/wilburt/wilburx9.dev/badge.svg?branch=develop)](https://coveralls.io/github/wilburt/wilburx9.dev?branch=develop)


My website.



## Running Locally

1. Fill out the backend [config file](./backend/configs/config.yaml) or add keys (prepend them with `WILBURX9`) in the linked file to environment variable.
2. Next, fill out the frontend [config file](./frontend/.env.development) or add the keys to environment variable.
3. Ensure [Go is installed](https://go.dev/doc/install) and execute `go run backend/cmd/main.go` from the root project directory.
4. Then, send a POST request to `http://localhost:$PORT/api/protected/cache` to populate the local db.
5. Finally, run `cd frontend && npm install && npm run build`.

## Live Deployment
1. Add the backend config variables to [Cloud Run](https://cloud.google.com/run/docs/configuring/environment-variables#setting).
2. Add the following [GitHub Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets#creating-encrypted-secrets-for-a-repository):
   
   | Secret      | Description                                                                                                                                                         |
   |-----------------| ------------- |
   | FIREBASE_SA_KEY | For deploying the frontend to Firebase Hosting. See [how to generate](https://github.com/FirebaseExtended/action-hosting-deploy/blob/main/docs/service-account.md). |
   | FIREBASE_TOKEN | For deploying the Cloud Functions. See [how to generate](https://firebase.google.com/docs/cli#cli-ci-systems).                                                      |
   | GCP_PROJECT | The project id on Google Cloud Platform                                                                                                                             |
   | GCP_SA_KEY | The service account for deploying the backend to Cloud Run. See [how to setup](https://github.com/google-github-actions/deploy-cloudrun#setup).                     |
   | REACT_APP_DOMAIN | The domain issued by Firebase Hosting.                                                                                                                              |
   | REACT_APP_FIREBASE_API_KEY | For configuring Firebase Analytics. See [documentation](https://firebase.google.com/docs/analytics/get-started?platform=web#add-sdk).                               |
   | REACT_APP_FIREBASE_APP_ID | Same as above                                                                                                                                                       |                                                                                           |
3. Any push to the `live` branch triggers a deployment to Cloud Run, Firebase Hosting and Cloud Functions for Firebase. See [GitHub Workflow](./.github/workflows/build_and_deploy.yaml).


## Contribution
Pull-requests are not accepted. Feature requests may or may not be worked on.
