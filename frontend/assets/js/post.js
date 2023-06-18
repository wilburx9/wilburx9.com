// Script for processing html in post.hbs

// Redirect if the post is just a reference to an external blog post
// Otherwise add "group" class to it.
!function () {
    let container = $('.gh-post-content');
    let bookmark = container.find(':first');
    // A post which is just a reference to an external article
    // will contain nothing but the bookmark card and the reading time.
    if (bookmark.length > 0 && container.children().length === 2) {
        let url = bookmark.find('a.kg-bookmark-container').attr('href');
        if (url) window.location.replace(url);
        return
    }

    $('figure.kg-bookmark-card').each(function () {
        $(this).find('.kg-bookmark-thumbnail').append('<div><svg width="18" height="18" viewBox="0 0 16 16" fill="none"><path stroke-width="2" class="stroke-orangeSet dark:stroke-orangeSet-dark" stroke-linecap="round" stroke-linejoin="round" d="M15.5 5.5v-5m0 0h-5m5 0L8.833 7.167m-2.5-5H4.5c-1.4 0-2.1 0-2.635.272A2.5 2.5 0 0 0 .772 3.532C.5 4.066.5 4.767.5 6.167V11.5c0 1.4 0 2.1.272 2.635a2.5 2.5 0 0 0 1.093 1.092C2.4 15.5 3.1 15.5 4.5 15.5h5.333c1.4 0 2.1 0 2.635-.273a2.5 2.5 0 0 0 1.093-1.092c.272-.535.272-1.235.272-2.635V9.667"/></svg><span>Open</span></div>');
        $(this).wrap('<div class="group"></div>');
    });
}()


class ImageProcessor {
    isPhotography

    constructor(primaryTag) {
        this.isPhotography = primaryTag === 'photography'
        this.postProcess()
    }

    // Loop through every figure image tag and apply modifications to them
    postProcess() {
        $(".kg-image-card img").each((i, image) => {
            let $wrapper = $("<div></div>");
            let $figure = $(image).parent(); // The <figure> container of the image
            let width = Number($(image).attr("width")) || image.naturalWidth;
            let height = Number($(image).attr("height")) || image.naturalHeight;

            this.resizeAndWrap(image, $figure, $wrapper, width, height)
            if (!this.isPhotography) return // Don't add lightbox and exif data on images for non-photography posts

            this.addLightBox(image, $figure, $wrapper, width, height, `lightbox__photo__${i}`)
        });
    }

    // Wrap the image in a blurred background and add zoom-in handle
    resizeAndWrap(image, $figure, $wrapper, width, height) {

        // Ensure the container height is not larger than the image
        $wrapper.css({
            "background-image": `url("${image.currentSrc || image.src}")`,
            "aspect-ratio": Math.max((width / height), this.getMinAspectRatio()).toString(),
            "max-height": `${height}px`
        }).addClass("group"); // Add group for Tailwind group hover

        $(image).css({
            "aspect-ratio": `${width}/${height}`,
            "max-width": `${width}px`,
            "max-height": `${height}px`
        });

        $figure.prepend($wrapper);
        $wrapper.append(image);
        return $figure
    }

    // Add a lightbox and zoom-out handle to the image
    addLightBox(image, $figure, $wrapper, width, height, lightBoxId) {
        const lightBox = `<div class='photo-lightbox' id='${lightBoxId}'>
                <div class="photo-lightbox-content">
                    <div class="group">
                        <span class='photo-zoom-out-handle'>
                            <svg width="18" height="18" viewBox="0 0 20 20" fill="none">
                                <path stroke-linecap="round" stroke-linejoin="round"
                                      d="m19 19-4.35-4.35M6 9h6m5 0A8 8 0 1 1 1 9a8 8 0 0 1 16 0Z"/>
                            </svg>
                        </span>
                            <img src="${image.src}" alt="${image.alt}" style="aspect-ratio: ${width}/${height}"/>
                    </div>
                </div>
            </div>`;
        const zoomInIcon = `<span class='photo-zoom-in-handle'>
                <svg width="18" height="18" viewBox="0 0 20 20" fill="none">
                    <path stroke-linecap="round" stroke-linejoin="round"
                          d="m19 19-4.35-4.35M9 6v6M6 9h6m5 0A8 8 0 1 1 1 9a8 8 0 0 1 16 0Z"/>
                </svg>
            </span>`;
        // Add the zoom-in icon as the last child of the image wrapper.
        $wrapper.prepend(zoomInIcon);
        // Add the lightbox above the image figure. After this, the image figure and lightbox share the same parent.
        $figure.before(lightBox);

        // Show the lightbox when the zoom-in icon on the image is clicked.
        $figure.find('.photo-zoom-in-handle').click(() => {
            this.showLightBox(lightBoxId, $figure, width, height)
        })

        // Close the lightbox when the zoom-out icon on the image is clicked.
        $figure.parent().find(`#${lightBoxId} .photo-zoom-out-handle`).click(() => {
            this.closeLightBox(lightBoxId, $figure)
        })

        // Listen for click events on the content background.
        $(`#${lightBoxId} .photo-lightbox-content`).click(event => {
            // Close the lightbox only if it was the content background that is clicked.
            if (event.target === event.currentTarget) this.closeLightBox(lightBoxId, $figure)
        })
    }

    closeLightBox(id, figure) {
        $(document).off(`keyup.${id}`);
        $(`#${id}`).fadeOut()
        figure.find('img.kg-image').fadeIn()
        $(`#${id} .photo-lightbox-content img`).removeClass('scale-full')
    }

    showLightBox(id, figure, imgWidth, imgHeight) {
        // Listen for escape key
        $(document).on(`keyup.${id}`, event => {
            if (event.key === "Escape") this.closeLightBox(id, figure);
        });
        figure.find('img.kg-image').fadeOut()
        $(`#${id}`).fadeIn()
        $(`#${id} .photo-lightbox-content img`).addClass('scale-full')
        $(`#${id} .photo-lightbox-content > div`).attr('style', this.getZoomImgWrapperStyle(imgWidth, imgHeight))
    }

    getZoomImgWrapperStyle(imgW, imgH) {
        let minW = Math.min(imgW, $(window).width())
        let minH = Math.min(imgH, $(window).height())
        let style = `aspect-ratio: ${imgW / imgH}; `
        if (minH > minW) {
            style += `height: auto; max-height: 100%; width: ${minW}px;`
        } else {
            style += `width: auto; max-width: 100%; height: ${minH}px;`
        }
        return style
    }

    getMinAspectRatio() {
        // 768 is tailwinds md breakpoint: https://tailwindcss.com/docs/responsive-design
        if ($(window).width() > 768) return 1.5
        return 0.6
    }

}


// Add Copy button to code blocks
!function () {
    let codes = document.querySelectorAll('code[class*="language-"]')
    for (let i = 0; i < codes.length; i++) {
        let code = codes[i]
        let pre = code.parentElement
        if (pre.tagName.toLowerCase() !== 'pre') continue

        let copied = `<span class="hide" id="copied"><span>Copied</span><svg class="fill-greenSet dark:fill-greenSet-dark" width="20" height="20" stroke="none" fill="none"><circle cx="10" cy="10" r="10"/><g clip-path="url(#a)"><path fill="#fff" d="M8.438 12.188 6.25 10l-.73.73 2.918 2.916 6.25-6.25-.73-.73-5.52 5.521Z"/></g><defs><clipPath id="a"><path fill="#fff" d="M3.75 3.75h12.5v12.5H3.75z"/></clipPath></defs></svg></span>`
        let copy = `<span id="copy"><span>Copy</span><svg width="22" height="22" fill="none"><path stroke-linecap="round" stroke-linejoin="round" d="M9.5 1.003c-.675.009-1.08.048-1.408.215a2 2 0 0 0-.874.874c-.167.328-.206.733-.215 1.408M18.5 1.003c.675.009 1.08.048 1.408.215a2 2 0 0 1 .874.874c.167.328.206.733.215 1.408m0 9c-.009.675-.048 1.08-.215 1.408a2 2 0 0 1-.874.874c-.328.167-.733.206-1.408.215M21 7v2m-8-8h2M4.2 21h7.6c1.12 0 1.68 0 2.108-.218a2 2 0 0 0 .874-.874C15 19.48 15 18.92 15 17.8v-7.6c0-1.12 0-1.68-.218-2.108a2 2 0 0 0-.874-.874C13.48 7 12.92 7 11.8 7H4.2c-1.12 0-1.68 0-2.108.218a2 2 0 0 0-.874.874C1 8.52 1 9.08 1 10.2v7.6c0 1.12 0 1.68.218 2.108a2 2 0 0 0 .874.874C2.52 21 3.08 21 4.2 21Z"/></svg></span>`
        pre.insertAdjacentHTML("afterbegin", `<div class="code-copy"><div>${copied}${copy}</div></div>`)
    }

    function copyCode(e) {
        let code = e.parentElement.parentElement.getElementsByTagName('code')[0]
        let text = code.innerText || code.textContent
        copy(e, text)
    }

    // Set click listeners on all the copy buttons
    let copyButtons = document.querySelectorAll('#copy');
    copyButtons.forEach((button) => {
        button.addEventListener('click', function () {
            copyCode(this.parentElement);
        });
    });
}();

// Set click listeners for the back and share buttons.
!function () {
    document.getElementById("back_icon").parentElement.href = `${window.location.origin}/blog`
    document.getElementById("post-link-copy").addEventListener("click", (event) => {
        let e = event.currentTarget;
        copy(e, window.location.href, () => {
            e.classList.toggle("copied")
        })
    })
}();

// Copy text from the element
function copy(element, text, toggle) {
    if (element.children[0].className !== 'hide') return
    navigator.clipboard.writeText(text).then(function () {
        element.children[0].classList.toggle("hide")
        element.children[1].classList.toggle("hide")
        if (typeof toggle === 'function') toggle()
        setTimeout(() => {
            element.children[0].classList.toggle("hide")
            element.children[1].classList.toggle("hide")
            if (typeof toggle === 'function') toggle()
        }, 2000);
    });
}