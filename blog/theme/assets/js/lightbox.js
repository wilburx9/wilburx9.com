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
    })
}();
