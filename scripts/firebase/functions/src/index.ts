import * as functions from "firebase-functions";
import admin from "firebase-admin";
import * as https from "https";

require("firebase-functions/lib/logger/compat");

admin.initializeApp();

export const cacheUpdateViaHttps = functions
    .https
    .onRequest(async (req, res) => {
      if (req.method == "POST") {
        const auth = req.get("Authorization");
        // Just confirm that the request has an authorization,
        // the backend will validate the key
        if (auth && auth?.length > 0) {
          const result = await makeRequest("HTTP", auth);
          res.status(200).send(result);
          return;
        }
        res.status(403).send("Forbidden!");
        return;
      }

      res.status(404).send("Yo! Have you lost your way?");
      return;
    });

export const cacheUpdateViaPubSub = functions
    .pubsub.schedule("every 5 minutes")
    .timeZone("Africa/Lagos")
    .onRun(() => {
      return makeRequest("PubSub", functions.config().cacheupdate.key);
    });

const makeRequest = (trigger: string, auth: string): Promise<string> => {
  const options = {
    hostname: functions.config().cacheupdate.domain,
    path: "/api/protected/cache",
    method: "POST",
    headers: {
      "Content-type": "application/json",
      "Authorization": auth,
    },
  };

  return new Promise<string>(((resolve, reject) => {
    const req = https.request(options, (res) => {
      if (res.statusCode != 200) {
        const message = `Request failed with status: ${res.statusCode}`;
        return reject(new Error(message));
      }
      const body: Uint8Array[] = [];
      res.on("data", (chunk) => body.push(chunk));
      res.on("end", () => {
        const resStr = Buffer.concat(body).toString();
        return resolve(resStr);
      });
    });

    req.on("error", (err) => {
      reject(err);
    });

    req.on("timeout", () => {
      req.destroy();
      reject(new Error("Request time out"));
    });

    req.write("");
    req.end();
  }));
};
