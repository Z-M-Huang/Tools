var appCache = window.applicationCache;

if (appCache != undefined && appCache != null) {
  appCache.update(); // Attempt to update the user's cache.

  if (appCache.status == window.applicationCache.UPDATEREADY) {
    appCache.swapCache();  // The fetch was successful, swap in the new cache.
  }
}

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
    document.cookie = name + "=; Max-age=0; path=/; domain=" + location.host;
  }
}

function bindForm(id, url, callback) {
  id = "#" + id;
  $(id).on("submit", (e) => {
    e.preventDefault();
    var data = parseFormToJSON(id);
    console.log(data);
    $.ajax({
      type: "POST",
      url: url,
      data: JSON.stringify(data),
      dataType: "json",
      contentType: "application/json",
      beforeSend: (xhr) => {
        var sessionToken = getCookieValue("session_token");
        if (
          sessionToken != "" &&
          sessionToken != null &&
          sessionToken != undefined
        ) {
          xhr.setRequestHeader("Authorization", "Bearer " + sessionToken);
        }
      },
      success: (data) => {
        if (data.Alert.Message != "") {
          showAlertCondition(data.Alert);
        } else if (callback != null && callback != undefined) {
          callback(data.Data);
        }
      },
      error: (xhr, status, error) => {
        console.log(xhr.status + ":" + xhr.statusCode + ":" + xhr.statusText);
        showAlertDanger(
          "Failed to receive success response, please try again later."
        );
      },
    });
  });
}

function getCookieValue(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length == 2) return parts.pop().split(";").shift();
}

/*******************************************************
 *                    Alert Section
 *******************************************************/
function showAlertDanger(message) {
  $(".alert").hide();
  $("#alertDangerMessage").text(message);
  $("#alertDanger").slideToggle();
}

function showAlertWarning(message) {
  $(".alert").hide();
  $("#alertWarningMessage").text(message);
  $("#alertWarning").slideToggle();
}

function showAlertSuccess(message) {
  $(".alert").hide();
  $("#alertSuccessMessage").text(message);
  $("#alertSuccess").slideToggle();
}

function showAlertInfo(message) {
  $(".alert").hide();
  $("#alertInfoMessage").text(message);
  $("#alertInfo").slideToggle();
}

function showAlertCondition(alert) {
  if (alert != "" && alert != undefined && alert != null) {
    if (alert.IsDanger) {
      showAlertDanger(alert.Message);
    } else if (alert.IsWarning) {
      showAlertWarning(alert.Message);
    } else if (alert.IsSuccess) {
      showAlertSuccess(alert.Message);
    } else if (alert.IsInfo) {
      showAlertInfo(alert.message);
    } else if (alert.Message != "") {
      console.log("Unknown alert", alert);
    }
  }
}

/*******************************************************
 *                    Parsing functions
 *******************************************************/
function parseFormToJSON(id) {
  var o = {};
  var a = $(id).serializeArray();
  $.each(a, function () {
    if (o[this.name]) {
      if (!o[this.name].push) {
        o[this.name] = [o[this.name]];
      }
      o[this.name].push(this.value || "");
    } else {
      o[this.name] = this.value || "";
    }
  });
  return o;
}
