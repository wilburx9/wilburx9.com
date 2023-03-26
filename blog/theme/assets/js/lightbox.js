!function () {
    const images = document.querySelectorAll(".kg-image-card img");
    for (let i = 0; i < images.length; i++) {
        let image = images[i]
        let lightBoxId = `lightbox__photo__${i}`
        let imgWrapper = document.createElement("div")
        imgWrapper.style.backgroundImage = `url("${image.currentSrc || image.src}")`
        let container = image.parentElement
        let imgW = Number(image.getAttribute("width"))
        let imgH = Number(image.getAttribute("height"))
        let imgAspectRatio = `${imgW}/${imgH}`
        let containerW = imgW;

        if ((imgW / imgH) < 0.80) containerW = imgH * 0.80

        let maxSide = Math.max(containerW, imgH);
        let minSide = Math.min(containerW, imgH);
        container.style.aspectRatio = `${maxSide}/${minSide}`
        container.style.maxHeight = `${imgH}px`
        image.style.aspectRatio = imgAspectRatio
        image.style.maxWidth = `${imgW}px`
        image.style.maxHeight = `${imgH}px`


        container.insertBefore(imgWrapper, image.parentElement.firstChild)
        imgWrapper.append(image)
        imgWrapper.insertAdjacentHTML("afterbegin", `<span onclick='showLightBox("${lightBoxId}");' class='photo-zoom-handle'><svg width="20" height="20" fill="none"><path stroke-linecap="round" stroke-linejoin="round" d="m19 19-4.35-4.35M9 6v6M6 9h6m5 0A8 8 0 1 1 1 9a8 8 0 0 1 16 0Z"/></svg></span>`)
        container.insertAdjacentHTML("beforebegin", `<div class='photo-lightbox' id='${lightBoxId}'><div style="aspect-ratio: ${imgAspectRatio};"><span onclick='closeLightBox("${lightBoxId}");' class='photo-zoom-handle'><svg width="20" height="20" fill="none"><path stroke-linecap="round" stroke-linejoin="round" d="m19 19-4.35-4.35M6 9h6m5 0A8 8 0 1 1 1 9a8 8 0 0 1 16 0Z"/></svg></span><img src="${image.src}" alt="${image.alt}" style="aspect-ratio: ${imgAspectRatio}"/></div></div>`)
    }
}();

function closeLightBox(id) {
    document.getElementById(id).style.display = "none"
}

function showLightBox(id) {
    console.log(`Clicked :: ${id}`)
    document.getElementById(id).style.display = 'flex'
}
