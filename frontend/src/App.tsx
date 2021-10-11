import * as React from "react"
import {
  ChakraProvider,
  Box,
  Flex,
} from "@chakra-ui/react"
import {ColorModeSwitcher} from "./ColorModeSwitcher"
import {Article} from "./Article";
import theme from "./theme";
import "./style.css"

export const App = () => (
  <ChakraProvider theme={theme}>
    <Box textAlign="center" fontSize="xl">
      <Flex flexDir="column">
        <ColorModeSwitcher justifySelf="flex-end" alignSelf="end"/>
        <Box flex="1">
          <Article/>
        </Box>
      </Flex>
    </Box>
  </ChakraProvider>
)
