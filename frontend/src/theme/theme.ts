import {extendTheme, theme as base, ThemeComponentProps, ThemeConfig} from "@chakra-ui/react";
import {mode, StyleFunctionProps} from "@chakra-ui/theme-tools";
import {Dict} from "@chakra-ui/utils";

const config: ThemeConfig = {
  initialColorMode: 'dark',
  useSystemColorMode: true,
}

const theme = extendTheme({
    config,
    styles: {
      global: (props: Dict | StyleFunctionProps) => ({
        body: {
          bg: mode("white", "#1A1B22")(props),
        }
      })
    },
    colors: {
      secondary: {
        50: '#fde8ff',
        100: '#ebbef5',
        200: '#db95eb',
        300: '#cc6ae1',
        400: '#bd41d8',
        500: '#a327be',
        600: '#7f1e95',
        700: '#5b146b',
        800: '#370b42',
        900: '#16011a',
      }
    },
    fonts: {
      heading: `'Inter', sans-serif, ${base.fonts?.heading}`,
      body: `'Inter', sans-serif, ${base.fonts?.body}`,
    },
    components: {
      Heading: {
        baseStyle: (props: ThemeComponentProps) => {
          return ({
            color: mode(
              "#1A1B22",
              "white"
            )(props),
          });
        },
      }
    }
  },
);

export default theme
