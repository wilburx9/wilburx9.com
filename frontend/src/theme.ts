import {extendTheme, theme as base, ThemeConfig, withDefaultVariant} from "@chakra-ui/react";
import {ThemeComponentProps} from "@chakra-ui/theme/dist/types/theme.types";

const config: ThemeConfig = {
  initialColorMode: 'dark',
  useSystemColorMode: true,
}

const inputStyles = {
  variants: {
    filled: (props: ThemeComponentProps) => ({
      field: {
        _focus: {
          borderWidth: '1px',
          boxShadow: 'none',
          borderColor: props.colorMode === 'light' ? 'black' : 'white'
        },
        _invalid: {
          borderColor: 'red.300', boxShadow: 'none', borderWidth: "1px"
        }
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
      Input: {...inputStyles},
    }
  },
  withDefaultVariant({
    variant: 'filled',
    components: ['Input', 'Textarea']
  }),
);

export default theme