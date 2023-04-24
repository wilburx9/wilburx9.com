// Script for processing html in default.hbs


!function () {
    handleDarkMode()
    setHrefs()
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
        document.getElementById("theme-switch").addEventListener("click", () => {
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
    }
}

function setHrefs() {
    document.addEventListener("DOMContentLoaded", () => {
        document.querySelectorAll('[id=about-me]').forEach(e => e.href = window.location.origin)
        let emailButton = document.querySelector('[id=email-me]')
        if (emailButton != null) emailButton.href = `mailto: me@${window.location.hostname}`
    });
}

// Get the details from a bookmark card and set the url
// and image to the post card url and featured image
function parseBookmark(content, linkId, imageId) {
    let c = document.createElement('div')
    c.innerHTML = content
    let bookmark = c.querySelector(':scope > figure.kg-bookmark-card')
    // A post which is just a reference to an external article
    // will contain nothing but the bookmark card.
    if (bookmark != null && c.children.length === 1) {
        let url = bookmark.querySelector('a.kg-bookmark-container').href
        let image = bookmark.querySelector('div.kg-bookmark-thumbnail').getElementsByTagName('img')[0].src
        let imageElement = document.getElementById(imageId);
        if (imageElement != null) imageElement.src = image
        document.getElementById(linkId).href = url
    }
}