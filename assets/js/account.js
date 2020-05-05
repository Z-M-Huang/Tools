function onNewPasswordKeyPress() {
  var ele = $(this);
  var passwordValid = isPasswordValid()
  if (!passwordValid && !ele.hasClass("is-invalid")) {
    ele.addClass("is-invalid");
    $("#newPasswordInvalidDiv").show();
  } else if (passwordValid) {
    if (ele.hasClass("is-invalid")) {
      ele.removeClass("is-invalid")
    }
    ele.addClass("is-valid")
    $("#newPasswordInvalidDiv").hide();
  }
}

function isPasswordValid() {
  var ele = $("#newPassword")
  if (ele.val().length < 12) {
    return false
  }
  return true
}

function onConfirmPasswordKeyPress() {
  var ele = $(this);
  var passwordInput = $("#newPassword")
  var passwordValid = isPasswordValid()
  if (passwordValid && ele.val() === passwordInput.val() && !ele.hasClass("is-invalid")) {
    ele.addClass("is-invalid");
    $("#confirmPasswordInvalidDiv").show();
  } else if (passwordValid && ele.val() === passwordInput.val()) {
    if (ele.hasClass("is-invalid")) {
      ele.removeClass("is-invalid")
    }
    ele.addClass("is-valid")
    $("#confirmPasswordInvalidDiv").hide();
  }
}