import {extendTheme, theme as base, ThemeComponentProps, ThemeConfig, withDefaultVariant} from "@chakra-ui/react";
import {mode, StyleFunctionProps} from "@chakra-ui/theme-tools";
import {Dict} from "@chakra-ui/utils";

const config: ThemeConfig = {
  initialColorMode: 'dark',
  useSystemColorMode: true,
}

// TODO: Remove input styles after removing contact form
const inputStyles = {
  variants: {
    filled: (props: ThemeComponentProps) => ({
      field: {
        ...getInputFilledStyle(props.colorMode === 'light')
      }
    })
  }
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
      Textarea: {...inputStyles},
      Input: {...inputStyles}, // This doesn't work, and I don't know why. Apply theme on Textarea manually
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
  withDefaultVariant({
    variant: 'filled',
    components: ['Input', 'Textarea']
  }),
);

export function getInputFilledStyle(isLightMode: boolean) {
  return {
    bg: isLightMode ? 'gray.200' : 'whiteAlpha.50',
    _focus: {
      borderWidth: '1px',
      boxShadow: 'none',
      borderColor: isLightMode ? 'black' : 'white'
    },
    _invalid: {
      borderColor: 'red.300', boxShadow: 'none', borderWidth: "1px"
    }
  }
}

export default theme
