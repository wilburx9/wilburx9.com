import {extendTheme, theme as base, ThemeComponentProps, ThemeConfig, withDefaultVariant} from "@chakra-ui/react";

const config: ThemeConfig = {
  initialColorMode: 'dark',
  useSystemColorMode: true,
}

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
    fonts: {
      heading: `'Roboto Mono', monospace, ${base.fonts?.heading}`,
      body: `'Roboto Mono', monospace, ${base.fonts?.body}`,
    },
    components: {
      Textarea: {...inputStyles},
      Input: {...inputStyles}, // This doesn't work, and I don't know why. Apply theme on Textarea manually
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
