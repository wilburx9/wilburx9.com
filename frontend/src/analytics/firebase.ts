import {initializeApp} from "firebase/app";
import {getAnalytics, logEvent} from "firebase/analytics";

const firebaseConfig = {
  apiKey: process.env.REACT_APP_FIREBASE_API_KEY,
  projectId: process.env.REACT_APP_GCP_PROJECT_ID,
  appId: process.env.REACT_APP_FIREBASE_APP_ID,
};

export function initializeFirebase() {
  let app = initializeApp(firebaseConfig);
  getAnalytics(app)
}

export function logAnalyticsEvent(name: string, params?: Map<string, any>) {
  logEvent(getAnalytics(), name, params)
}

export function getAnalyticsParams(): Map<string, any> {
  return new Map<string, any>()
}