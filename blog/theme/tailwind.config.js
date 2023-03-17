const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
    content: ['./*.hbs', './**/*.hbs'],
    theme: {
        extend: {
            colors: {
                'cardSet': {
                    DEFAULT: '#F6F6F6',
                    dark: '#0A1520',
                },
                'boxSet': {
                    DEFAULT: '#EEEEEE',
                    dark: '#152234',
                },
                'greySet': {
                    DEFAULT: '#3D3D3D',
                    dark: '#E2E2E2',
                },
                'greyLightSet': {
                    DEFAULT: '#595959',
                    dark: '#ADADAD',
                },
                'blackSet': {
                    DEFAULT: '#000000',
                    dark: '#FFFFFF',
                },
                'bgSet': {
                    DEFAULT: '#FFFFFF',
                    dark: '#0A1018',
                },
                'orangeSet': {
                    DEFAULT: '#DF7902',
                    dark: '#F59300',
                }
            }
        },
        fontFamily: {
            'sans': ['Poppins', ...defaultTheme.fontFamily.sans],
            'mono': ['JetBrains Mono', ...defaultTheme.fontFamily.mono],
        },
        fontSize: {
            'body1': ['16px', {
                'lineHeight': '20px',
                'fontWeight': '400'
            }],
            'body1Bold': ['16px', {
                'lineHeight': '20px',
                'fontWeight': '700'
            }],
            'body2': ['14px', {
                'lineHeight': '20px',
                'fontWeight': '400'
            }],
            'body3': ['12px', {
                'lineHeight': '16px',
                'fontWeight': '400'
            }],
        }
    },
    plugins: [
        require('@tailwindcss/typography'),
        require('@tailwindcss/line-clamp')
    ],
}