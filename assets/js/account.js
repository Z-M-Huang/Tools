function onPasswordChange(id) {
  id = "#" + id;
  var ele = $(id);
  var passwordValid = isPasswordValid();
  if (!passwordValid) {
    if (!ele.hasClass("is-invalid")) {
      ele.removeClass("is-valid");
      ele.addClass("is-invalid");
      $("#passwordInvalidDiv").show();
    }
  } else if (!ele.hasClass("is-valid")) {
    ele.removeClass("is-invalid");
    ele.addClass("is-valid");
    $("#passwordInvalidDiv").hide();
  }
}

function isPasswordValid() {
  var ele = $("#password");
  if (ele.val().length < 12) {
    return false;
  }
  return true;
}

function onConfirmePasswordChange(id) {
  id = "#" + id;
  var ele = $(id);
  var passwordInput = $("#password");
  var passwordValid = isPasswordValid();

  if (ele.val() != "" && ele.val() == passwordInput.val()) {
    if (!ele.hasClass("is-valid")) {
      ele.removeClass("is-invalid");
      ele.addClass("is-valid");
      $("#confirmPasswordInvalidDiv").hide();
    }
  } else {
    ele.removeClass("is-valid");
    ele.addClass("is-invalid");
    $("#confirmPasswordInvalidDiv").show();
  }
}

function onUsernameChange(id) {
  id = "#" + id;
  var ele = $(id);
  if (ele.val().length < 6) {
    if (!ele.hasClass("is-invalid")) {
      ele.removeClass("is-valid");
      ele.addClass("is-invalid");
      $("#usernameInvalidDiv").show();
    }
  } else if (!ele.hasClass("is-valid")) {
    ele.removeClass("is-invalid");
    ele.addClass("is-valid");
    $("#usernameInvalidDiv").hide();
  }
}

function validateEmail(email) {
  var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return re.test(String(email).toLowerCase());
}

function onEmailChange(id) {
  id = "#" + id;
  var ele = $(id);
  if (!validateEmail(ele.val())) {
    if (!ele.hasClass("is-invalid")) {
      ele.removeClass("is-valid");
      ele.addClass("is-invalid");
      $("#emailInvalidDiv").show();
    }
  } else if (!ele.hasClass("is-valid")) {
    ele.removeClass("is-invalid");
    ele.addClass("is-valid");
    $("#emailInvalidDiv").hide();
  }
}