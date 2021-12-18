import * as React from "react"
import {ChakraProvider, Box, Flex} from "@chakra-ui/react"
import {ColorModeSwitcher} from "./components/ColorModeSwitcher"
import {ArticlesComponent} from "./components/ArticlesComponent";
import theme from "./theme";
import "./style.css"
import DataProvider from "./DataProvider";
import {ReposComponent} from "./components/ReposComponent";

export const App = () => (
  <ChakraProvider theme={theme}>
    <DataProvider>
      <Box
        textAlign="center"
        fontSize="xl">
        <Flex flexDir="column">
          <ColorModeSwitcher ml='auto'/>
          <Box
            flex='1'
            maxW="container.xl"
            px={[5, null, 10]}
            alignSelf="center">
            <ArticlesComponent/>
            <ReposComponent/>
          </Box>
        </Flex>
      </Box>
    </DataProvider>
  </ChakraProvider>
)
