!function () {
    const images = document.querySelectorAll(".kg-image-card img");
    images.forEach(function (image) {
        let imgWrapper = document.createElement("div")
        imgWrapper.style.backgroundImage = "url(" + image.currentSrc + ")"
        let container = image.parentElement
        let imgW = Number(image.getAttribute("width"))
        let imgH = Number(image.getAttribute("height"))
        let containerW = imgW;

        if ((imgW / imgH) < 0.80) containerW = imgH * 0.80

        let maxSide = Math.max(containerW, imgH);
        let minSide = Math.min(containerW, imgH);
        container.style.aspectRatio = maxSide + "/" + minSide
        image.style.aspectRatio = imgW + "/" + imgH


        container.insertBefore(imgWrapper, image.parentElement.firstChild)
        imgWrapper.append(image)
        imgWrapper.insertAdjacentHTML("afterbegin", "<span id='zoom-handle'><svg width=\"20\" height=\"20\" fill=\"none\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"m19 19-4.35-4.35M9 6v6M6 9h6m5 0A8 8 0 1 1 1 9a8 8 0 0 1 16 0Z\"/></svg></span>")
    })
}();
