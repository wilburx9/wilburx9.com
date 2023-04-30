const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
    darkMode: 'class',
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
                },
                'borderSet': {
                    DEFAULT: '#A3A3A3',
                    dark: '#38567D',
                },
                'greenSet': {
                    DEFAULT: '#0F8437',
                    dark: '#15C35B',
                },
                'redSet': {
                    DEFAULT: '#AA3B3B',
                    dark: '#DC4A4A',
                }
            },
        },
        fontFamily: {
            'sans': ['Poppins', ...defaultTheme.fontFamily.sans],
            'mono': ['JetBrains Mono', ...defaultTheme.fontFamily.mono],
        },
        fontSize: {
            'text1': ['20px', {
                'lineHeight': '32px',
            }],
            'text2': ['16px', {
                'lineHeight': '28px',
            }],
            'text3': ['14px', {
                'lineHeight': '22px',
            }],
            'body1': ['16px', {
                'lineHeight': '20px',
            }],
            'body2': ['14px', {
                'lineHeight': '20px',
            }],
            'body3': ['12px', {
                'lineHeight': '16px',
            }],
            'headline': ['36px', {
                'lineHeight': '54px',
                'fontWeight': '700'
            }],
            'headline1': ['24px', {
                'lineHeight': '28px',
            }],
            'headline2': ['20px', {
                'lineHeight': '24px',
            }],
            'button': ['16px', {
                'lineHeight': '20px',
                'fontWeight': '500'
            }],
            'caption': ['12px', {
                'lineHeight': '18px',
                'fontWeight': '400',
            }],
            'exif': ['14px', {
                'lineHeight': '20px',
                'fontWeight': '300',
            }],
            'code': ['20px', {
                'lineHeight': '32px',
                'fontWeight': '400',
            }],
            'title': ['58px', {
                'lineHeight': '68px',
                'fontWeight': '700',
            }],
            'largeTitle': ['62px', {
                'lineHeight': '72px',
                'fontWeight': '700',
            }],
        }
    },
    plugins: [
        require('@tailwindcss/line-clamp')
    ],
}
