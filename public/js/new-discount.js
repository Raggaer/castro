$('document').ready(function() {
   $('#uses-group').hide();
});
$('#uses').click(function() {
    $('#uses-group').fadeIn();
});
$('#unlimited').click(function() {
    $('#uses-group').fadeOut();
});