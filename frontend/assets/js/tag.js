$(function () {
    let request = null

    function onRequestComplete(feed, footer, loader, anim) {
        feed.stop().fadeIn(300)
        footer.stop().fadeIn(300)
        loader.stop().fadeOut(300, function () {
            loader.children().empty()
        })
        anim.stop();
        anim.destroy();
        request = null
    }

    $('.header-tag').on('click', function (event) {
        event.preventDefault() // prevent the link from navigating to a new page

        let slug = $(this).data('slug')

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

        let postFeed = $('.gh-post-feed')
        let postFeedFooter = $('.gh-post-feed-footer')
        let postLoader = $('.gh-post-loader #loader')

        postFeed.fadeOut(300)
        postFeedFooter.fadeOut(300)
        postLoader.fadeIn(300)

        let animation = bodymovin.loadAnimation({
            container: postLoader[0],
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
                $('.gh-post-feed').html($(data).find('.gh-post-feed').html())
                onRequestComplete(postFeed, postFeedFooter, postLoader, animation)

            },
            error: function () {
                onRequestComplete(postFeed, postFeedFooter, postLoader, animation)
            }
        })


    })
})