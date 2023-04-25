let loadMore

function setupInfiniteScroll(maxPages, pageTag) {
    if (loadMore) loadMore.cleanUp()
    loadMore = new LoadMore(maxPages, pageTag)
}

class LoadMore {
    animation
    loading = false
    offset = 300
    currentPage = 1
    rafId
    disposed = false

    constructor(maxPages, pageTag) {
        console.log(`${pageTag} :: Instantiated with ${maxPages} pages`)
        this.maxPages = maxPages
        this.pageTag = pageTag || ""
        this.baseUrl = `${document.location.origin}/blog/${pageTag}`
        this.$scrollingContent = $('.gh-content')

        this.cleanUp = this.cleanUp.bind(this)
        this.showLoader = this.showLoader.bind(this)
        this.hideLoader = this.hideLoader.bind(this)
        this.handleScroll = this.handleScroll.bind(this)
        this.handleResize = this.handleResize.bind(this)
        this.loadMorePosts = this.loadMorePosts.bind(this)
        this.setEventListeners = this.setEventListeners.bind(this)
        this.removeEventListeners = this.removeEventListeners.bind(this)

        this.setEventListeners()
        this.loadMorePosts()
    }

    setEventListeners() {
        console.log(`${this.pageTag} :: setEventListeners`)
        this.$scrollingContent.on('scroll', this.handleScroll)
        $(window).on('resize', this.handleResize)
    }

    removeEventListeners() {
        console.log(`${this.pageTag} :: removeEventListeners`)
        this.$scrollingContent.off('scroll', this.handleScroll)
        $(window).off('resize', this.handleResize)
    }

    loadMorePosts() {
        if (this.rafId) {
            cancelAnimationFrame(this.rafId)
        }

        if (this.loading) {
            console.log(`${this.pageTag} :: loadMorePosts :: loading Returning`)
            return
        }

        this.rafId = window.requestAnimationFrame(() => {
            const bottom = this.$scrollingContent[0].scrollHeight - this.$scrollingContent.innerHeight() - this.offset
            if (this.$scrollingContent.scrollTop() >= bottom) {
                this.fetchPosts()
            } else {
                console.log(`${this.pageTag} :: loadMorePosts :: Not at the bottom`)
            }
        })

    }

    fetchPosts() {
        this.loading = true
        this.showLoader()
        let nextPage = `page/${++this.currentPage}/`
        console.log(`${this.pageTag} :: fetchPosts :: ${nextPage}`)
        setTimeout(() => {
            $.ajax({
                url: this.baseUrl + nextPage,
                type: 'GET',
                success: data => {
                    console.log(`${this.pageTag} :: fetchPosts :: success`)
                    if (this.disposed) {
                        console.log(`${this.pageTag} :: fetchPosts :: not updating after success`)
                        return
                    }
                    const newPosts = $(data).find('.gh-card, .post');
                    newPosts.hide()
                    $('.gh-post-feed').append(newPosts)
                    this.hideLoader(newPosts)
                },
                error: (jqXHR) => {
                    if (this.disposed) {
                        console.log(`${this.pageTag} :: fetchPosts :: not updating after error`)
                        return
                    }
                    console.log(`${this.pageTag} :: fetchPosts :: error`)
                    if (jqXHR.status === 404) this.removeEventListeners()
                    this.hideLoader()
                },
                complete: () => {
                    console.log(`${this.pageTag} :: fetchPosts :: complete`)
                    if (this.currentPage >= this.maxPages) this.removeEventListeners()
                    this.loading = false
                }
            })
        }, 10000);
    }

    showLoader() {
        $('.gh-post-feed-footer-info').addClass('hide')
        let $loader = $('.gh-more-post-loader #more-loader')
        this.animation = bodymovin.loadAnimation({
            container: $loader[0],
            renderer: 'svg',
            loop: true,
            autoplay: true,
            path: '/assets/lottie/loading.json'
        })
        $loader.stop().fadeIn(300, 'linear')
    }

    hideLoader(newPosts) {
        let $loader = $('.gh-more-post-loader #more-loader')
        $loader.stop().fadeOut(300, 'linear', () => {
            if (newPosts) {
                newPosts.fadeIn(300, 'linear')
            }
            $('.gh-post-feed-footer-info').removeClass('hide')
            $loader.children().empty()
            if (this.animation) {
                this.animation.stop();
                this.animation.destroy();
            }
        })
    }

    handleScroll() {
        this.loadMorePosts()
    }

    handleResize() {
        this.loadMorePosts()
    }

    cleanUp() {
        this.disposed = true
        this.hideLoader()
        this.removeEventListeners()
    }

}