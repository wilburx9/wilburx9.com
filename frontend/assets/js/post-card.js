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
