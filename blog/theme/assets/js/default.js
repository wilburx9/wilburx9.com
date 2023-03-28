// Script for processing html in default.hbs


!function () {
    handleDarkMode()
}();

function handleDarkMode() {
    let modeQueryList = window.matchMedia('(prefers-color-scheme: dark)');
    const isDarkMode = () => localStorage.theme === 'dark' || (!('theme' in localStorage) && modeQueryList.matches);

    if (isDarkMode()) {
        document.documentElement.classList.add('dark')
    } else {
        document.documentElement.classList.remove('dark')
    }

    if (!('theme' in localStorage)) {
        modeQueryList.addEventListener('change', function () {
            if (!('theme' in localStorage)) changeMode(true, isDarkMode())
        });
    }

    document.addEventListener("DOMContentLoaded", () => {
        let wrapper = document.getElementById("theme-switch");
        if (isDarkMode()) {
            document.getElementById("theme-switch-dark").classList.toggle("hide")
        } else {
            document.getElementById("theme-switch-light").classList.toggle("hide")
        }
        wrapper.addEventListener("click", () => {
            changeMode(false, !isDarkMode())
        });
    });

    function changeMode(fromSystem, darkMode) {
        if (darkMode) {
            if (!fromSystem) localStorage.theme = 'dark'
            document.documentElement.classList.add('dark')
        } else {
            if (!fromSystem) localStorage.theme = 'light'
            document.documentElement.classList.remove('dark')
        }
        document.getElementById("theme-switch-dark").classList.toggle("hide")
        document.getElementById("theme-switch-light").classList.toggle("hide")
    }
}
