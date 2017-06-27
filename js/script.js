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
    $('#table0').show();


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
                    $('#table0').show(); break;
                case "franklunch":
                    $('#table1').show(); break;
                case "frankdinner":
                    $('#table2').show(); break;
                case "frarybreak":
                    $('#table3').show(); break;
                case "frarylunch":
                    $('#table4').show(); break;
                case "frarydinner":
                    $('#table5').show(); break;
                case "oldenborgbreak":
                    $('#table6').show(); break;                                                                         
                case "oldenborglunch":
                    $('#table7').show(); break;
                case "oldenborgdinner":
                    $('#table8').show(); break;
                case "collinsbreak":
                    $('#table9').show(); break;             
                case "collinslunch":
                    $('#table10').show(); break;
                case "collinsdinner":
                    $('#table11').show(); break;             
                case "scrippsbreak":
                    $('#table12').show(); break;             
                case "scrippslunch":
                    $('#table13').show(); break;                    
                case "scrippsdinner":
                    $('#table14').show(); break; 
                case "pitzerbreak":
                    $('#table15').show(); break;            
                case "pitzerlunch":
                    $('#table16').show(); break;                    
                case "pitzerdinner":
                    $('#table17').show(); break;
                case "muddbreak":
                    $('#table18').show(); break;             
                case "muddlunch":
                    $('#table19').show(); break;                    
                case "mudddinner":
                    $('#table20').show(); break;                                                 
            }
        }
        currentTab = str;
    }
});
