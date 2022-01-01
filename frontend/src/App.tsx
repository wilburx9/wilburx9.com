import * as React from "react"
import {ChakraProvider} from "@chakra-ui/react"
import theme from "./theme/theme";
import "./theme/style.css"
import DataProvider from "./DataProvider";
import {ContentComponent} from "./ContentComponent";

export const App = () => (
  <ChakraProvider theme={theme}>
    <DataProvider>
      <ContentComponent/>
    </DataProvider>
  </ChakraProvider>
)

