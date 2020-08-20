<!DOCTYPE html>
<html>
<head runat="server">
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, minimum-scale=1, user-scalable=no" />
    <meta name="apple-mobile-web-app-capable" content="yes" />
    <meta name="mobile-web-app-capable" content="yes" />
    <title>ONLYOFFICE</title>
    <link rel="icon"
          href="images/{{.icon}}.ico"
          type="image/x-icon" />
    <link rel="stylesheet" type="text/css" href="stylesheets/editor.css" />
</head>
<body style="height: 100%; margin: 0;">
<div id="yumioffice" style="height: 100%"></div>

<script type="text/javascript" src = "{{.apiUrl}}"></script>
<script type="text/javascript" language="javascript">
    let innerAlert = function (message) {
        if (console && console.log)
            console.log(message);
    };

    const onAppReady = function () {
        innerAlert("Document editor ready");
    };
    const onDocumentStateChange = function (event) {
        var title = document.title.replace(/\*$/g, "");
        document.title = title + (event.data ? "*" : "");
    };
    const onRequestEditRights = function () {
        location.href = location.href.replace(RegExp("mode=view\&?", "i"), "");
    };
    const onError = function (event) {
        if (event)
            innerAlert(event.data);
    };
    const onRequestHistory = function () {
        docEditor.refreshHistory(
            {
                currentVersion: "{{.Version}}",
                history: JSON.stringify({{.History}})
            }
        )
    }
    const onRequestHistoryClose = function () {
        document.location.reload();
    };
    const onRequestHistoryData = function (event) {
        version = event.data

        docEditor.setHistoryData({{.HistoryData}}[version-1])
    }
    const onOutdatedVersion = function () {
        location.reload(true);
    };

    let docEditor;
    const connectEditor = function () {
        docEditor = new DocsAPI.DocEditor(
            "yumioffice",
            {
                "document": {{.config}},
                "events": {
                    "onAppReady": onAppReady,
                    "onDocumentStateChange": onDocumentStateChange,
                    'onRequestEditRights': onRequestEditRights,
                    "onError": onError,
                    "onRequestHistory":  onRequestHistory,
                    "onRequestHistoryData": onRequestHistoryData,
                    "onRequestHistoryClose": onRequestHistoryClose,
                    "onOutdatedVersion": onOutdatedVersion,
                }
            }
        )

        fixSize();
    }

    let fixSize = function () {
        var wrapEl = document.getElementsByClassName("form");
        if (wrapEl.length) {
            wrapEl[0].style.height = screen.availHeight + "px";
            window.scrollTo(0, -1);
            wrapEl[0].style.height = window.innerHeight + "px";
        }
    };

    if (window.addEventListener) {
        window.addEventListener("load", connectEditor);
        window.addEventListener("resize", fixSize);
    } else if (window.attachEvent) {
        window.attachEvent("onload", connectEditor);
        window.attachEvent("onresize", fixSize);
    }
</script>
</body>
</html>