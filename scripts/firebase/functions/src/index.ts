import * as functions from "firebase-functions";
import admin from "firebase-admin";
import * as https from "https";

admin.initializeApp();

export const cacheUpdateViaHttps = functions
    .https
    .onRequest(async (req, res) => {
      if (req.method == "POST") {
        const auth = req.get("Authorization");
        // Just confirm that the request has an authorization,
        // the backend will validate the key
        if (auth && auth?.length > 0) {
          const result = await makeRequest( auth);
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
    const req = https.request(options, (res) => {
      if (res.statusCode != 200) {
        const message = `Request failed with status: ${res.statusCode}`;
        return reject(new Error(message));
      }
      res.on("end", () => {
        return resolve(`Requested succeeded with status: ${res.statusCode}`);
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
