// Script for processing html in post.hbs

!function () {
    const images = document.querySelectorAll(".kg-image-card img");
    for (let i = 0; i < images.length; i++) {
        let image = images[i]
        let lightBoxId = `lightbox__photo__${i}`
        let imgWrapper = document.createElement("div")
        imgWrapper.style.backgroundImage = `url("${image.currentSrc || image.src}")`
        imgWrapper.classList.add("group")
        let container = image.parentElement
        let imgW = Number(image.getAttribute("width"))
        let imgH = Number(image.getAttribute("height"))
        let maxH = `${imgH}px`;
        let imgAspectRatio = `${imgW}/${imgH}`
        let containerW = imgW;

        if ((imgW / imgH) < 0.80) containerW = imgH * 0.80

        let maxSide = Math.max(containerW, imgH);
        let minSide = Math.min(containerW, imgH);
        container.style.aspectRatio = `${maxSide}/${minSide}`
        container.style.maxHeight = maxH
        image.style.aspectRatio = imgAspectRatio
        image.style.maxWidth = `${imgW}px`
        image.style.maxHeight = maxH
        imgWrapper.style.maxHeight = maxH


        container.insertBefore(imgWrapper, image.parentElement.firstChild)
        imgWrapper.append(image)
        imgWrapper.insertAdjacentHTML("afterbegin", `<span onclick='showLightBox("${lightBoxId}");' class='photo-zoom-handle'><svg width="18" height="18" viewBox="0 0 20 20" fill="none"><path stroke-linecap="round" stroke-linejoin="round" d="m19 19-4.35-4.35M9 6v6M6 9h6m5 0A8 8 0 1 1 1 9a8 8 0 0 1 16 0Z"/></svg></span>`)
        container.insertAdjacentHTML("beforebegin", `<div class='photo-lightbox' id='${lightBoxId}'><div style="aspect-ratio: ${imgAspectRatio};"><span onclick='closeLightBox("${lightBoxId}");' class='photo-zoom-handle'><svg width="18" height="18" viewBox="0 0 20 20" fill="none"><path stroke-linecap="round" stroke-linejoin="round" d="m19 19-4.35-4.35M6 9h6m5 0A8 8 0 1 1 1 9a8 8 0 0 1 16 0Z"/></svg></span><img src="${image.src}" alt="${image.alt}" style="aspect-ratio: ${imgAspectRatio}"/></div></div>`)
    }
}();

!function () {
    let codes = document.querySelectorAll('code[class*="language-"]')
    for (let i = 0; i < codes.length; i++) {
        let code = codes[i]
        let pre = code.parentElement
        if (pre.tagName.toLowerCase() !== 'pre') continue
        pre.classList.add("group")

        let copied = `<span class="hide">Copied<svg class="fill-greenSet dark:fill-greenSet-dark" width="20" height="20" stroke="none" fill="none"><circle cx="10" cy="10" r="10"/><g clip-path="url(#a)"><path fill="#fff" d="M8.438 12.188 6.25 10l-.73.73 2.918 2.916 6.25-6.25-.73-.73-5.52 5.521Z"/></g><defs><clipPath id="a"><path fill="#fff" d="M3.75 3.75h12.5v12.5H3.75z"/></clipPath></defs></svg></span>`
        let copy = `<span onclick='copyCode(this.parentElement)'>Copy<svg width="22" height="22" fill="none"><path stroke-linecap="round" stroke-linejoin="round" d="M9.5 1.003c-.675.009-1.08.048-1.408.215a2 2 0 0 0-.874.874c-.167.328-.206.733-.215 1.408M18.5 1.003c.675.009 1.08.048 1.408.215a2 2 0 0 1 .874.874c.167.328.206.733.215 1.408m0 9c-.009.675-.048 1.08-.215 1.408a2 2 0 0 1-.874.874c-.328.167-.733.206-1.408.215M21 7v2m-8-8h2M4.2 21h7.6c1.12 0 1.68 0 2.108-.218a2 2 0 0 0 .874-.874C15 19.48 15 18.92 15 17.8v-7.6c0-1.12 0-1.68-.218-2.108a2 2 0 0 0-.874-.874C13.48 7 12.92 7 11.8 7H4.2c-1.12 0-1.68 0-2.108.218a2 2 0 0 0-.874.874C1 8.52 1 9.08 1 10.2v7.6c0 1.12 0 1.68.218 2.108a2 2 0 0 0 .874.874C2.52 21 3.08 21 4.2 21Z"/></svg></span>`
        pre.insertAdjacentHTML("afterbegin", `<div' class="code-copy">${copied}${copy}</div>`)
    }
}();

!function () {
    document.getElementById("back_icon").parentElement.href = `${window.location.origin}/blog`
    document.getElementById("post-link-copy").addEventListener("click", (event) => {
        let e = event.currentTarget;
        copy(e, window.location.href, () => {
            e.classList.toggle("copied")
        })
    })
}();


function closeLightBox(id) {
    document.getElementById(id).style.display = "none"
    document.body.style.overflowY = 'auto'
}

function showLightBox(id) {
    document.getElementById(id).style.display = 'flex'
    document.body.style.overflowY = 'hidden'
}

function copyCode(e) {
    let code = e.parentElement.getElementsByTagName('code')[0]
    let text = code.innerText || code.textContent
    copy(e, text)
}

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
