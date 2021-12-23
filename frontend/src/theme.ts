import {extendTheme, theme as base, ThemeConfig} from "@chakra-ui/react";

const config: ThemeConfig = {
  initialColorMode: 'dark',
  useSystemColorMode: true,
}

const theme = extendTheme({
  config,
  fonts: {
    heading: `'Roboto Mono', monospace, ${base.fonts?.heading}`,
    body: `'Roboto Mono', monospace, ${base.fonts?.body}`
  },
});

export default theme