function logout() {
  clearCookie("session_token");
  document.location.href = "/";
}

function clearCookie(name) {
  document.cookie = name + "=; Max-age=0; path=/; domain=" + location.host;
}

function getCookieValue(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length == 2) return parts.pop().split(";").shift();
}

function onClickRedirect(url) {
  window.location.href = url
}

/*******************************************************
 *                    Like/Dislike Section
 *******************************************************/
function likeOnClick(obj, name) {
  var ele = $(obj);
  if (ele.hasClass("fas")) {
    //Unlike
    postLink("/app/" + name + "/dislike", (d) => {
      ele.parent().html("<i class=\"far fa-thumbs-up mr-1 hover-pointer hover-150\" onclick=\"likeOnClick(this, '" + name +  "')\"></i>" + d)
    })
  } else {
    //like
    postLink("/app/" + name + "/like", (d) => {
      ele.parent().html("<i class=\"fas fa-thumbs-up mr-1 hover-pointer hover-150\" onclick=\"likeOnClick(this, '" + name +  "')\"></i>" + d)
    })
  }
}


/*******************************************************
 *                    Ajax Section
 *******************************************************/
function bindForm(id, url, callback) {
  id = "#" + id;
  $(id).on("submit", (e) => {
    e.preventDefault();
    var data = parseFormToJSON(id);
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
      statusCode: {
        401: () => {
          document.location.href = "/login";
        },
        200: (data) => {
          if (data != null && data != undefined 
              && data.Alert != null && data.Alert != undefined &&
              data.Alert.Message != "") {
            showAlertCondition(data.Alert);
          } 
          if (callback != null && callback != undefined) {
            callback(data.Data);
          }
        },
      },
      error: (xhr, status, error) => {
        console.log(xhr.status + ":" + xhr.statusText);
        showAlertDanger("Failed to receive success response, please try again later.");
      },
    });
  });
}

function postJSONData(url, data, callback) {
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
    statusCode: {
      401: () => {
        document.location.href = "/login";
      },
      200: (data) => {
        if (data != null && data != undefined 
          && data.Alert != null && data.Alert != undefined &&
          data.Alert.Message != "") {
          showAlertCondition(data.Alert);
        } 
        if (callback != null && callback != undefined) {
          callback(data.Data);
        }
      },
    },
    error: (xhr, status, error) => {
      console.log(xhr.status + ":" + xhr.statusText);
      showAlertDanger("Failed to receive success response, please try again later.");
    },
  });
}

function postLink(url, callback) {
  $.ajax({
    type: "POST",
    url: url,
    dataType: "json",
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
    statusCode: {
      401: () => {
        showAlertDanger("Please login first");
      },
      200: (data) => {
        if (data != null && data != undefined 
          && data.Alert != null && data.Alert != undefined &&
          data.Alert.Message != "") {
          showAlertCondition(data.Alert);
        } 
        if (callback != null && callback != undefined) {
          callback(data.Data);
        }
      },
    },
    error: (xhr, status, error) => {
      console.log(xhr.status + ":" + xhr.statusText);
      showAlertDanger("Failed to receive success response, please try again later.");
    },
  });
}

/*******************************************************
 *                    Dynamic Modal Section
 *******************************************************/
function getModalHTML(id, title, content, primaryButtonOnClick, primaryButtonText) {
  modal = '<div class="modal fade" id="' + id + '"><div class="modal-dialog" role="document"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">' + title + '</h5><button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button></div><div class="modal-body">' + content + '</div><div class="modal-footer">';
  if (!(primaryButtonOnClick == null || primaryButtonOnClick == undefined)) {
    modal += '<button type="button" class="btn btn-primary" onclick="' + primaryButtonOnClick + '">' + primaryButtonText + '</button>';
  }
  modal += '<button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button></div></div></div></div>';
  return modal;
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

function hideAllAlerts() {  
  $(".alert").hide();
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
      showAlertInfo(alert.Message);
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

/*******************************************************
 *                    Chart colors
 *******************************************************/
window.chartColors = {
  black: "rgb(0,0,0)",
  navy: "rgb(0,0,128)",
  blue: "rgb(0,0,255)",
  green: "rgb(0,128,0)",
  teal: "rgb(0,128,128)",
  lime: "rgb(0,255,0)",
  aqua: "rgb(0,255,255)",
  maroon: "rgb(128,0,0)",
  purple: "rgb(128,0,128)",
  olive: "rgb(128,128,0)",
  gray: "rgb(128,128,128)",
  silver: "rgb(192,192,192)",
  red: "rgb(255,0,0)",
  fuchsia: "rgb(255,0,255)",
  yellow: "rgb(255,255,0)",
  white: "rgb(255,255,255)"
};
