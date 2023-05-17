class Subscription {
    turnstileWidgetId

    constructor(turnstileSiteKey, baseApiUrl, primaryTag) {
        this.apiUrl = baseApiUrl
        this.setupSubscriptionForm(turnstileSiteKey)
        this.setupSubscription(primaryTag)
    }

    renderCaptcha(siteKey) {
        turnstile.ready(() => {
            this.turnstileWidgetId = turnstile.render('.form-captcha-container', {
                sitekey: siteKey,
                action: 'email-subscription',
                theme: isDarkTheme() ? 'dark' : 'light',
                'response-field-name': 'captcha',
                callback: () => this.submitForm(),
            });
        });

        $('.subscription-content .form-cta-container').stop().fadeOut(300, 'linear')
        $('.subscription-content .form-captcha-container').stop().fadeIn(300, 'linear')

    }

    submitForm() {

        let $progressDiv = $('.subscription-loader #loader')
        let animation = getProgressAnimation($progressDiv)

        const $form = $('.subscription-modal form');
        $form.addClass('hide')
        $progressDiv.parent().removeClass('hide')

        // Convert the form data to JSON
        let formData = new FormData($form[0])
        let jsonData = Object.fromEntries(
            Array.from(formData.keys()).map(key => [
                key, key === "tags" ? formData.getAll(key) : formData.get(key)
            ])
        )

        $.ajax({
            url: `${this.apiUrl}/subscribe`,
            type: 'POST',
            data: JSON.stringify(jsonData),
            contentType: 'application/json',
            complete: (xhr) => {
                if (xhr.status >= 200 && xhr.status <= 299) {
                    this.handleSubmitSuccess($progressDiv, animation)
                } else {
                    this.handleSubmitError($progressDiv, animation)
                }
            }
        })

    }

    handleSubmitSuccess(progressDiv, progressAnim) {
        let $successContainer = $('.subscription-success')
        let $animationDiv = $successContainer.find('#success-icon')
        let animation = getLottieAnimation($animationDiv, 'done', 1)

        $('.subscription-success .success-cta').off('click').on('click', () => {
            animation.stop()
            animation.destroy()
            this.closeSubscriptionModal();
        });

        this.hideProgressUi(progressDiv, progressAnim, $successContainer)
    }

    handleSubmitError(progressDiv, progressAnim) {
        $('.subscription-error .error-cta').off('click').on('click', () => {
            this.unHideForm()
        });

        let errorContainer = $('.subscription-error')
        this.hideProgressUi(progressDiv, progressAnim, errorContainer)
    }

    hideProgressUi(progressDiv, progressAnim, showUI) {
        progressDiv.parent().addClass('hide') // Hide progress UI
        showUI.removeClass('hide') // Show error/success UI

        // Cleanup progress animation and its container
        progressAnim.stop()
        progressAnim.destroy()
        progressDiv.children().empty()
    }

    setupSubscriptionForm(turnstileSiteKey) {
        this.handleOnFocus()
        $('.subscription-content form').on('submit', event => {
            event.preventDefault();
            if (this.validateForm()) this.renderCaptcha(turnstileSiteKey)
        });
    }

    validateForm() {
        let $emailField = $('.subscription-modal input[type="email"]')
        let email = $emailField.val()

        if (email === '' || !/\S+@\S+/.test(email)) {
            this.handleFormError($emailField)
            return false
        } else {
            $emailField.removeClass('error');
            $emailField.closest('.subscription-modal').find('.form-error').stop().fadeOut(300, 'linear')
            $emailField.removeData('input-listener');
            $emailField.off('input');
            return true
        }
    }

    handleFormError($emailField) {
        $emailField.addClass('error');

        let $errorField = $emailField.closest('.subscription-modal').find('.form-error');
        $errorField.text("Please enter a valid email address")
        $errorField.stop().fadeIn(300, 'linear')

        if ($emailField.data('input-listener')) return

        $emailField.data('input-listener', true);
        $emailField.on('input', () => {
            this.validateForm()
        });
    }

    handleOnFocus() {
        // Add 'active' class to the email field when it looses focus, and it has texts.
        $('.subscription-modal input[type="email"]').on('blur', function () {
            let label = $(this).next('label');
            if ($(this).val().trim() === '') {
                label.removeClass('active');
            } else {
                label.addClass('active');
            }
        });
    }

    setupSubscription(primaryTag) {
        $('#post-subscribe').click(() => {
            this.showSubscriptionModal(primaryTag)
        })

        // Listen for escape key
        $(document).keyup(event => {
            if (event.key === "Escape") this.closeSubscriptionModal()
        });

        // Listen for click events on the translucent background
        $(".subscription-modal").click(event => {
            if (event.target === event.currentTarget) this.closeSubscriptionModal()
        })
    }

    showSubscriptionModal(primaryTag) {
        // Fade in the translucent background
        $('.subscription-modal').fadeIn()
        // Slide the modal content from the bottom
        $(".subscription-content").css({
            "transform": "translateY(0%) translateX(-50%)",
            "opacity": 1
        });

        // Select the chip of primary tag of the post that is currently being read
        let checkbox = $(`.subscription-content #${primaryTag}`);
        if (checkbox.length > 0) {
            checkbox.prop('checked', true);
        }
    }

    closeSubscriptionModal() {
        // Don't close the dialog if the progress UI is currently displayed
        let progressLoader = '.subscription-loader';
        if (!$(progressLoader).hasClass("hide")) return

        // Slide down the modal content
        $(".subscription-content").css({
            "transform": "translateY(50%) translateX(-50%)",
            "opacity": 0.0
        });

        $('.subscription-modal').fadeOut() // Fade out the translucent background
        let $form = $('.subscription-modal form');
        $form.removeClass('hide') // Un-hide the form
        $form[0].reset() // Reset all inputs in the form to their default values

        $('.subscription-modal input[type="email"]').next('label').removeClass('active') // Reset the email active state
        $(progressLoader).addClass('hide') // Hide the progress UI
        this.unHideForm()
    }

    unHideForm() {
        // Remove the captcha from the document
        if (this.turnstileWidgetId) turnstile.remove(this.turnstileWidgetId)
        this.turnstileWidgetId = null

        $('.subscription-content .form-cta-container').stop().fadeIn(0) // Un-hide the form CTA button
        $('.subscription-success').addClass('hide') // Hide the success UIs
        $('.subscription-error').addClass('hide') // Hide the error UIs
        $('.subscription-modal form').removeClass('hide') // Show the form
    }
}