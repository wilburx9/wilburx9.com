import * as React from "react"
import {ChakraProvider, Box, Flex} from "@chakra-ui/react"
import {ColorModeSwitcher} from "./ColorModeSwitcher"
import {ArticleComponent} from "./ArticleComponent";
import theme from "./theme";
import "./style.css"
import DataProvider from "./DataProvider";

export const App = () => (
  <ChakraProvider theme={theme}>
    <DataProvider>
      <Box textAlign="center" fontSize="xl">
        <Flex flexDir="column">
          <ColorModeSwitcher justifySelf="flex-end" alignSelf="end"/>
          <Box  maxW="container.xl" alignSelf="center">
            <ArticleComponent />
          </Box>
        </Flex>
      </Box>
    </DataProvider>
  </ChakraProvider>
)
