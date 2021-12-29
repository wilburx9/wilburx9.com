import * as functions from "firebase-functions";
import admin from "firebase-admin";
import axios from "axios";

admin.initializeApp();

const http = axios.create({
  baseURL: `${functions.config().cacheUpdate.domain}/api/protected`,
  headers: {
    "Content-type": "application/json",
    "Authorization": functions.config().cacheUpdate.apiKey,
  },
});

export const cacheUpdate = functions
    .pubsub.schedule("every 5 minutes")
    .timeZone("Africa/Lagos")
    .onRun(() => {
      http
          .get("/cache")
          .then((response) => {
            console.log(`Update success:: ${response.data}`);
          })
          .catch((e) => {
            console.log(`Update failure:: ${e}`);
          });
    });

