import * as React from "react"
import {ChakraProvider, Box, Flex} from "@chakra-ui/react"
import {ColorModeSwitcher} from "./components/ColorModeSwitcher"
import {ArticlesComponent} from "./components/ArticlesComponent";
import theme from "./theme";
import "./style.css"
import DataProvider from "./DataProvider";
import {ReposComponents} from "./components/ReposComponents";

export const App = () => (
  <ChakraProvider theme={theme}>
    <DataProvider>
      <Box
        textAlign="center"
        fontSize="xl">
        <Flex flexDir="column">
          <ColorModeSwitcher
            justifySelf="flex-end"
            alignSelf="end"/>
          <Box
            flex='1'
            maxW="container.xl"
            alignSelf="center">
            <ArticlesComponent/>
            <ReposComponents/>
          </Box>
        </Flex>
      </Box>
    </DataProvider>
  </ChakraProvider>
)
