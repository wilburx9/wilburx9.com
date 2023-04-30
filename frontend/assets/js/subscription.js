$(function () {
    $('.subscription-modal input[type="email"]').on('blur', function () {
        let label = $(this).next('label');
        if ($(this).val().trim() === '') {
            label.removeClass('active');
        } else {
            label.addClass('active');
        }
    });
})

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

    $(document).keyup(function(event) {
        if (event.key === "Escape") closeModal()
    });

    $(".subscription-modal").click(function(event) {
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

function validateForm() {
    
}