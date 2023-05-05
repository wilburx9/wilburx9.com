let turnstileWidgetId

function renderCaptcha(siteKey) {
    let isDarkTheme = $('html').hasClass('dark');
    turnstile.ready(function () {
        turnstileWidgetId = turnstile.render('.form-captcha-container', {
            sitekey: siteKey,
            action: 'email-subscription',
            theme: isDarkTheme ? 'dark' : 'light',
            'response-field-name': 'captcha',
            callback: () => submitForm(),
        });
    });

    $('.subscription-content .form-cta-container').stop().fadeOut(300, 'linear')
    $('.subscription-content .form-captcha-container').stop().fadeIn(300, 'linear')

}

function submitForm() {

    let $progressDiv = $('.subscription-loader #loader')
    let animation = getProgressAnimation($progressDiv)

    $('.subscription-modal form').addClass('hide')
    $progressDiv.parent().removeClass('hide')

    // TODO: Replace this timeout with an API request to submit the form. If the request succeeds, call handleSubmitSuccess()
    setTimeout(function () {
        handleSubmitError($progressDiv, animation)
    }, 3000)
}

function handleSubmitSuccess(progressDiv, progressAnim) {
    let $successContainer = $('.subscription-success')
    let $animationDiv = $successContainer.find('#success-icon')
    let animation = getLottieAnimation($animationDiv, 'done')

    $('.subscription-success .success-cta').off('click').on('click', function () {
        animation.stop()
        animation.destroy()
        closeSubscriptionModal();
    });

    hideProgressUi(progressDiv, progressAnim, $successContainer)
}

function handleSubmitError(progressDiv, progressAnim) {
    $('.subscription-error .error-cta').off('click').on('click', function () {
        unHideForm()
    });

    let errorContainer = $('.subscription-error')
    hideProgressUi(progressDiv, progressAnim, errorContainer)
}

function hideProgressUi(progressDiv, progressAnim, showUI) {
    progressDiv.parent().addClass('hide') // Hide progress UI
    showUI.removeClass('hide') // Show error/success UI

    // Cleanup progress animation and its container
    progressAnim.stop()
    progressAnim.destroy()
    progressDiv.children().empty()
}

function setupSubscriptionForm(turnstileSiteKey) {
    handleOnFocus()
    $('.subscription-content form').on('submit', function (event) {
        event.preventDefault();
        if (validateForm()) renderCaptcha(turnstileSiteKey)
    });
}

function validateForm() {
    let $emailField = $('.subscription-modal input[type="email"]')
    let email = $emailField.val()

    if (email === '' || !/\S+@\S+/.test(email)) {
        handleFormError($emailField)
        return false
    } else {
        $emailField.removeClass('error');
        $emailField.closest('.subscription-modal').find('.form-error').stop().fadeOut(300, 'linear')
        $emailField.removeData('input-listener');
        $emailField.off('input');
        return true
    }
}

function handleFormError($emailField) {
    $emailField.addClass('error');

    let $errorField = $emailField.closest('.subscription-modal').find('.form-error');
    $errorField.text("Please enter a valid email address")
    $errorField.stop().fadeIn(300, 'linear')

    if ($emailField.data('input-listener')) return

    $emailField.data('input-listener', true);
    $emailField.on('input', function () {
        validateForm()
    });
}

function handleOnFocus() {
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

function setupSubscription(primaryTag) {
    $('#post-subscribe').click(function () {
        showSubscription(primaryTag)
    })

    // Listen for escape key
    $(document).keyup(function (event) {
        if (event.key === "Escape") closeSubscriptionModal()
    });

    // Listen for click events on the translucent background
    $(".subscription-modal").click(function (event) {
        if (event.target === this) closeSubscriptionModal()
    })


    // TODO: Remove after implementing subscription UI
    showSubscription(primaryTag)
}

function showSubscription(primaryTag) {
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

function closeSubscriptionModal() {
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
    unHideForm()
}

function unHideForm() {
    removeCaptcha()
    $('.subscription-content .form-cta-container').stop().fadeIn(0) // Un-hide the form CTA button
    $('.subscription-success').addClass('hide') // Hide the success UIs
    $('.subscription-error').addClass('hide') // Hide the error UIs
}

function removeCaptcha() {
    if (turnstileWidgetId) turnstile.remove(turnstileWidgetId)
    turnstileWidgetId = null
}
