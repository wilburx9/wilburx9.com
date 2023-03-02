const defaultTheme = require("tailwindcss/defaultTheme")

module.exports = {
    darkMode: "class",
    content: [
        "./src/**/*.{js,jsx,ts,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                "bgSet": {
                    "light": "#FFFFFF",
                    "dark": "#0A1018"
                },
                "blackSet": {
                    "light": "#000000",
                    "dark": "#FFFFFF"
                },
                "fontGrey": {
                    "light": "#3D3D3D",
                    "dark": "#E2E2E2"
                },
                "orangeSet": {
                    "light": "#DF7902",
                    "dark": "#F59300"
                },
                "pizazz": "#F59300",
                "electricViolet": "#C41CFF",
            },
        },
        fontFamily: {
            "sans": ["Poppins", ...defaultTheme.fontFamily.sans]
        },
        fontSize: {
            'largeTitle': "62px",
            'headline4': "20px",
            "button": "16px"
        },
        lineHeight: {
            '72': '72px',
            '24': '24px',
            '20': '20px',
        }

    },
    plugins: [],
}
