import {extendTheme, theme as base} from "@chakra-ui/react";

const theme = extendTheme({
  fonts: {
    heading: `'Roboto Mono', monospace, ${base.fonts?.heading}`,
    body: `'Roboto Mono', monospace, ${base.fonts?.body}`
  }
});

export default theme