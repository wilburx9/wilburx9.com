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

function isDarkTheme() {
    return $('html').hasClass('dark');
}

function setHrefs() {
    document.addEventListener("DOMContentLoaded", () => {
        document.querySelectorAll('[id=about-me]').forEach(e => e.href = window.location.origin)
        let emailButton = document.querySelector('[id=email-me]')
        if (emailButton != null) emailButton.href = `mailto: me@${window.location.hostname}`
    });
}

// Get the details from a bookmark card and set the url
// and image to the post-card url and featured image
function postProcessBookmarkCard(content, tagString, linkId, imageId) {
    let tags = tagString.split(',').map(tag => tag.trim());
    if (!tags.includes("#external")) return

    let c = $('<div></div>').html(content)
    let bookmark = $(c).find('figure.kg-bookmark-card')

    let url = bookmark.find('a.kg-bookmark-container').attr('href');
    let image = bookmark.find('div.kg-bookmark-thumbnail img:first').attr('src');
    let readingTime = $(c).find('.external-reading-time').text();


    // The anchor tag wrapper of the post-card
    let anchor = $('#' + linkId);
    // This image will not be found if the post has a feature image
    let img = anchor.find('#' + imageId);
    if (img.length) img.attr('src', image);

    anchor.find('#external-tag').css('display', 'flex');
    anchor.find('#reading-time').text(readingTime + " read");
    anchor.attr('href', url);
    anchor.attr('target', '_blank');

}

function getProgressAnimation(container) {
    let name = isDarkTheme() ? 'loading_dark' : 'loading_light'
    return getLottieAnimation(container, name, true)
}

function getLottieAnimation(container, name, loop) {
    return bodymovin.loadAnimation({
        container: container[0],
        renderer: 'svg',
        loop: loop,
        autoplay: true,
        path: `/assets/lottie/${name}.json`
    })
}