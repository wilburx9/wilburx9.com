$(function () {
    let request = null
    let currentSlug = ""

    function onRequestComplete(feed, footer, loader, anim) {
        feed.stop().fadeIn(300, 'linear')
        footer.stop().fadeIn(300, 'linear')
        loader.stop().fadeOut(300, 'linear', function () {
            loader.children().empty()
        })
        anim.stop();
        anim.destroy();
        request = null
    }

    $('.header-tag').on('click', function (event) {
        event.preventDefault() // prevent the link from navigating to a new page

        let slug = $(this).data('slug')
        currentSlug = slug

        if (request) request.abort()

        let indicator = $(this).children('.gh-tag-nav-indicator')

        // Check if the tag is already selected
        if (indicator.attr('id')) {
            indicator.removeAttr('id')
            slug = "" // Empty the slug so the home page will be loaded
        } else {
            indicator.attr('id', 'active')
        }
        $('.header-tag').not(this).children('.gh-tag-nav-indicator').removeAttr('id')

        let $postFeed = $('.gh-post-feed')
        let $postFeedFooter = $('.gh-post-feed-footer')
        let $postLoader = $('.gh-post-loader #loader')

        $postFeed.stop().fadeOut(300, 'linear')
        $postFeedFooter.stop().fadeOut(300, 'linear')
        $postLoader.stop().fadeIn(300, 'linear')

        let animation = bodymovin.loadAnimation({
            container: $postLoader[0],
            renderer: 'svg',
            loop: true,
            autoplay: true,
            path: '/assets/lottie/loading.json'
        })

        request = $.ajax({
            url: `/blog/${slug}`,
            type: 'GET',
            dataType: 'html',
            success: function (data) {
                if (currentSlug !==  slug) return
                let html = $(data);
                $('.gh-post-feed').html(html.find('.gh-post-feed').html())
                onRequestComplete($postFeed, $postFeedFooter, $postLoader, animation)
                setupInfiniteScroll(html.find('.total-page-count').text(), `${slug}/`)

            },
            error: function () {
                if (currentSlug !==  slug) return
                onRequestComplete($postFeed, $postFeedFooter, $postLoader, animation)
            }
        })


    })
})