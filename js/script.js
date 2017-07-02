$(document).ready(function() {

    /*------------------------------------*/
    /*           Intializations           */
    /*------------------------------------*/

    // Loading all of the reviews
    $('.menu_table').accordion({
        active: false,
        header: '.food_cell',
        icons: false,
        collapsible: true,
        heightStyle: 'panel'
    });
    $('.menu_table').hide();
    $('#table00').show();


    $('.rating_stars').each(function(i, star) {
        var rateNum = $(star).siblings('.rating_number').text();
        rateNum = parseFloat(rateNum.substring(1, rateNum.length - 3));
        $(star).rateYo({
            rating: rateNum,
            readOnly: true
        });
    });

    $('span').each(function(i, nameSpan) {
        var name = $(nameSpan).text();
        if(name.length >= 35) {
            $(nameSpan).addClass('too_long');
        }
    });

    $('.food_cell').each(function(i, foodCell) {
        var foodID = $(foodCell).attr('id');
        $.ajax({
            url: '/getReviews/' + foodID,
            method: 'GET',
            contentType: 'application/json',
            success: function(response) {
                if(response.length != 0) {
                    $.each(response, function(i, review) {
                        $(foodCell).next().append(
                            "<div class='review_bubble'><div class='review_stars'></div>"
                            + "<h3>Created at: " + review.created_at + "</h3>"
                            + "<h2>" + review.review_text + "</h2>"
                            + "</div>");
                        var star = $('.review_stars');
                        $(star).rateYo({
                            rating: review.rating,
                            readOnly: true,
                            starWidth: "20px"
                        })
                    });

                } else {
                    $(foodCell).next().append("<h1>No Reviews Yet.</h1>")
                }
            }
        })
    })

    /*------------------------------------*/
    /*          Handling Clicks           */
    /*------------------------------------*/

    // School tab click
    $('.schooltab').click(function() {
        $('.schooltab').removeClass('active');
        $(this).addClass('active');
        handleTableSelect();
    })

    // Meal tab click
    $('.mealtab').click(function() {
        $('.mealtab').removeClass('active');
        $(this).addClass('active');
        handleTableSelect();
    })

    /*------------------------------------*/
    /*              Handlers              */
    /*------------------------------------*/
    // Handling selecting a different meal or school
    var currentTab = "frankbreak";
    var handleTableSelect = function() {
        var str = "";
        $('.tabs .active').each(function(i, id) {
            str += $(id).attr('id');
        });
        console.log(str);
        if(str != currentTab) {
            $('.menu_table').hide();
            switch(str) {
                case "frankbreak":
                    $('#table00').show(); break;
                case "franklunch":
                    $('#table01').show(); break;
                case "frankdinner":
                    $('#table02').show(); break;
                case "frarybreak":
                    $('#table10').show(); break;
                case "frarylunch":
                    $('#table11').show(); break;
                case "frarydinner":
                    $('#table12').show(); break;
                case "oldenborgbreak":
                    $('#table20').show(); break;
                case "oldenborglunch":
                    $('#table21').show(); break;
                case "oldenborgdinner":
                    $('#table22').show(); break;
                case "collinsbreak":
                    $('#table30').show(); break;
                case "collinslunch":
                    $('#table31').show(); break;
                case "collinsdinner":
                    $('#table32').show(); break;
                case "scrippsbreak":
                    $('#table40').show(); break;
                case "scrippslunch":
                    $('#table41').show(); break;
                case "scrippsdinner":
                    $('#table42').show(); break;
                case "pitzerbreak":
                    $('#table50').show(); break;
                case "pitzerlunch":
                    $('#table51').show(); break;
                case "pitzerdinner":
                    $('#table52').show(); break;
                case "muddbreak":
                    $('#table60').show(); break;
                case "muddlunch":
                    $('#table61').show(); break;
                case "mudddinner":
                    $('#table62').show(); break;
            }
        }
        currentTab = str;
    }
});
