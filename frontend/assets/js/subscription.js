$(function () {
    handleOnFocus()
    handleFormSubmission()
})

function handleFormSubmission() {
    $('.subscription-content form').on('submit', function (event) {
        event.preventDefault();
        if (!validateForm()) return

    });
}

function validateForm() {
    let $emailField = $('.subscription-modal input[type="email"]')
    let email = $emailField.val()

    if (email === '' || !/\S+@\S+/.test(email)) {
        handleFormError($emailField)
    } else {
        $emailField.removeClass('error');
        $emailField.closest('.subscription-modal').find('.subscription-error').stop().fadeOut(300, 'linear')
        $emailField.removeData('input-listener');
        $emailField.off('input');
    }
}

function handleFormError($emailField) {
    $emailField.addClass('error');

    let $errorField = $emailField.closest('.subscription-modal').find('.subscription-error');
    $errorField.text("Please enter a valid email address")
    $errorField.stop().fadeIn(300, 'linear')

    if ($emailField.data('input-listener')) return

    $emailField.data('input-listener', true);
    $emailField.on('input', function () {
        validateForm()
    });
}

function handleOnFocus() {
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

    function closeModal() {
        $(".subscription-content").css({
            "transform": "translateY(50%) translateX(-50%)",
            "opacity": 0.0
        });
        $('.subscription-modal').fadeOut()
    }

    $('#post-subscribe').click(function () {
        $('.subscription-modal').fadeIn()
        $(".subscription-content").css({
            "transform": "translateY(0%) translateX(-50%)",
            "opacity": 1
        });
    })

    $(document).keyup(function (event) {
        if (event.key === "Escape") closeModal()
    });

    $(".subscription-modal").click(function (event) {
        if (event.target === this) closeModal()
    })

    let checkbox = $(`.subscription-content #${primaryTag}`);
    if (checkbox.length > 0) {
        checkbox.prop('checked', true);
    }

    // TODO: Remove after implement subscription UI
    $('.subscription-modal').fadeIn()
    $(".subscription-content").css({
        "transform": "translateY(0%) translateX(-50%)",
        "opacity": 1
    });
}
