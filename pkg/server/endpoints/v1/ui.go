package v1

import (
	"html/template"
	"net/http"

	"github.com/jrasell/sherpa/pkg/build"
)

type UIServer struct {
	build buildInfo
}

// buildInfo is the Sherpa server build information used to display on the UI.
type buildInfo struct {
	Version string
}

func NewUIServer() *UIServer {
	return &UIServer{
		build: buildInfo{
			Version: build.Version,
		},
	}
}

func (s *UIServer) Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui", http.StatusSeeOther)
}

func (s *UIServer) Get(w http.ResponseWriter, r *http.Request) {
	_ = tmplScaleEvent.ExecuteTemplate(w, "scaling-events", s.build)
}

var tmplScaleEvent = template.Must(template.New("scaling-events").Parse(`
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Sherpa</title>
    <script type="text/javascript" src="https://code.jquery.com/jquery-3.3.1.min.js"></script>
    <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.100.2/css/materialize.min.css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.100.2/js/materialize.min.js"></script>
    <meta name="viewport" content="width=device-width, initial-scale=1.0"/>

    <style type="text/css">
        td.tags { display: none; }
        .footer { padding-top: 10px; }
        .logo { height: 32px; margin: 0 auto; display: block; }

        @media (min-width: 78em) {
            td.tags{ display: table-cell; }
        }
    </style>
</head>
<body>

<ul id="overrides" class="dropdown-content"></ul>

<nav class="top-nav teal accent-4">
    <div class="container">
        <div class="nav-wrapper">
            <a href="/ui" class="brand-logo">Sherpa</a>
            <ul id="nav-mobile" class="right hide-on-med-and-down">
				<li><a href="https://github.com/jrasell/sherpa/releases">{{.Version}}</a></li>
				<li><a href="https://cloud.docker.com/u/jrasell/repository/docker/jrasell/sherpa">DockerHub</a></li>
            </ul>
        </div>
    </div>
</nav>

<div class="container">
    <div class="section">
        <table class="events highlight"></table>
    </div>
</div>

<script>
    $(function(){
        var params={};window.location.search.replace(/[?&]+([^=&]+)=([^&]*)/gi,function(str,key,value){params[key] = value;});
        function renderEvents(events) {
            var $table = $('table.events');
            var thead = '<thead><tr>';
            thead += '<th>ID</th>';
            thead += '<th>Job:Group</th>';
            thead += '<th>Direction</th>';
            thead += '<th>Count</th>';
            thead += '<th>Status</th>';
            thead += '<th>Time</th>';
            thead += '</tr></thead>';
            var $tbody = $('<tbody />');
            for (var [id, job] of Object.entries(events)) {
                for (var [jbname, event] of Object.entries(job)) {
                    var $tr = $('<tr />');
                    $tr.append($('<td />').text(id));
                    $tr.append($('<td />').text(jbname));
                    $tr.append($('<td />').text(event.Details.Direction))
                    $tr.append($('<td />').text(event.Details.Count))
                    $tr.append($('<td />').text(event.Status));
                    $tr.append($('<td />').text(timeConverter(event.Time)));
                    $tr.appendTo($tbody);
                }
            }

            $table.empty().append($(thead)).append($tbody);
        }

        function timeConverter(ts){
            var a = new Date(ts/1000000);
            var year = a.getUTCFullYear();
            var month = a.getUTCMonth();
            var date = a.getUTCDate();
            var hour = a.getUTCHours();
            var min = a.getUTCMinutes();
            var sec = a.getUTCSeconds();
            var ms = a.getUTCMilliseconds();
            return year + "-" + month + "-" + date + " " + hour + ':' + min + ':' + sec + "." + ms + " +0000 UTC"
        }		

        $.get("/v1/scale/status", function(data) {
            renderEvents(data);
        });
    })
</script>

</body>
</html>
`))
