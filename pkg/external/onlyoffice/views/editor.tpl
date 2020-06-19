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
          href="images/<%- editor.documentType %>.ico"
          type="image/x-icon" />
    <link rel="stylesheet" type="text/css" href="stylesheets/editor.css" />
</head>
<body style="height: 100%; margin: 0;">
<div id="yumioffice" style="height: 100%"></div>

<script type="text/javascript"></script>
<script type="text/javascript">
    const onAppReady = function () {
        console.log("ONLYOFFICE Document Editor is ready");
    };
    const onCollaborativeChanges = function () {
        console.log("The document changed by collaborative user");
    };
    const onDocumentReady = function() {
        console.log("Document is loaded");
    };
    const onDocumentStateChange = function (event) {
        if (event.data) {
            console.log("The document changed");
        } else {
            console.log("Changes are collected on document editing service");
        }
    };
    const onDownloadAs = function (event) {
        console.log("ONLYOFFICE Document Editor create file: " + event.data);
        window.top.postMessage(event.data);
        createAndDownloadFile("test.docx", event.data)
    };
    const onRequestInsertImage = function (event) {
        console.log("ONLYOFFICE Document Editor insertImage" + event.data);
        docEditor.insertImage({
            "fileType": "png",
            "url": "http://192.168.99.1/attachment/20190728测试上传文件名修改/2020January/1580363537940306800_small.png"
        });
    };
    const onError = function (event) {
        console.log("ONLYOFFICE Document Editor reports an error: code " + event.data.errorCode + ", description " + event.data.errorDescription);
    };
    const onOutdatedVersion = function () {
        location.reload(true);
    };
    const onRequestEditRights = function () {
        console.log("ONLYOFFICE Document Editor requests editing rights");
        // document.location.reload();
        var he = location.href.replace("view", "edit");
        location.href = he;
    };
    const onRequestHistory = function () {
    }
    const onRequestHistoryClose = function () {
        document.location.reload();
    };
    const onRequestHistoryData = function (event) {
    }
</script>


</body>
</html>