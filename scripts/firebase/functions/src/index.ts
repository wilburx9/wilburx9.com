import * as functions from "firebase-functions";
import admin from "firebase-admin";
import * as https from "https";

require("firebase-functions/lib/logger/compat");

admin.initializeApp();

export const cacheUpdateViaHttps = functions
    .https
    .onRequest(async (req, res) => {
      if (req.method == "POST") {
        console.log("Is POST");
        const auth = req.get("Authorization");
        // Just confirm that the request has an authorization,
        // the backend will validate the key
        if (auth && auth?.length > 0) {
          console.log("Has Authorization");
          const result = await makeRequest( auth);
          res.status(200).send(result);
          return;
        }
        console.log("No Authorization");
        res.status(403).send("Forbidden!");
        return;
      }

      console.log("Not POST");
      res.status(404).send("Yo! Have you lost your way?");
      return;
    });

export const cacheUpdateViaPubSub = functions
    .pubsub.schedule("every saturday 03:00")
    .timeZone("Africa/Lagos")
    .onRun(() => {
      return makeRequest(functions.config().cacheupdate.key);
    });

const makeRequest = (auth: string): Promise<string> => {
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
    console.log(`Options :: ${JSON.stringify(options)}`);
    const req = https.request(options, (res) => {
      if (res.statusCode != 200) {
        console.log(`Is not 200:: ${res.statusCode}`);
        const message = `Request failed with status: ${res.statusCode}`;
        return reject(new Error(message));
      }
      console.log("Is 200");
      const body: Uint8Array[] = [];
      res.on("data", (chunk) => body.push(chunk));
      res.on("end", () => {
        const resStr = Buffer.concat(body).toString();
        return resolve(resStr);
      });
    });

    req.on("error", (err) => {
      console.log(`Network error:: ${err}`);
      reject(err);
    });

    req.on("timeout", () => {
      console.log("Timeout error");
      req.destroy();
      reject(new Error("Request time out"));
    });

    req.write("");
    req.end();
  }));
};
