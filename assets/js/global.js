function logout() {
  clearCookies();
  document.location.href = "/";
}

function clearCookies() {  
  //Clear all cookies
  var cookies = document.cookie.split(";");

    for (var i = 0; i < cookies.length; i++) {
        var cookie = cookies[i];
        var eqPos = cookie.indexOf("=");
        var name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie;
        document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
    }
}

function showAlertDanger(message) {
  $(".alert .alert-dismissible").hide();
  $("alertDangerMessage").innerHTML = message;
  $("#alertDanger").slideToggle();
}

function showAlertWarning(message) {
  $(".alert .alert-dismissible").hide();
  $("alertWarningMessage").innerHTML = message;
  $("#alertWarning").slideToggle();
}

function showAlertSuccess(message) {
  $(".alert .alert-dismissible").hide();
  $("alertSuccessMessage").innerHTML = message;
  $("#alertSuccess").slideToggle();
}

function showAlertInfo(message) {
  $(".alert .alert-dismissible").hide();
  $("alertInfoMessage").innerHTML = message;
  $("#alertInfo").slideToggle();
}