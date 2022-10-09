import {ColorModeScript} from "@chakra-ui/react"
import {createRoot} from 'react-dom/client';
import * as React from "react"
import {App} from "./App"
import {initializeFirebase} from "./analytics/firebase";

// Redirect to custom URL if coming from Firebase Hosting provided URLs
if (process.env.NODE_ENV === 'production' && window.location.hostname.indexOf((new URL(process.env.REACT_APP_DOMAIN!)).hostname) === -1) {
  window.location.replace(process.env.REACT_APP_DOMAIN!);
}

// Disable logs in production
if (process.env.NODE_ENV === 'production') {
  console.log = () => {
  }
  console.error = () => {
  }
  console.debug = () => {
  }
}

initializeFirebase()

const container = document.getElementById("root")
const root = createRoot(container!)
root.render(<React.StrictMode><ColorModeScript/><App/></React.StrictMode>)
