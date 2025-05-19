// Lab 7, 8, 9: Use these templates to render the web pages

package web

const indexHTML = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>TritonTube</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet">
  </head>
  <body class="bg-light">
    <nav class="navbar navbar-expand-lg navbar-dark bg-primary mb-4">
      <div class="container">
        <a class="navbar-brand" href="/">TritonTube</a>
      </div>
    </nav>
    <div class="container">
      <div class="row mb-4">
        <div class="col-md-8">
          <h2>Watchlist</h2>
          <div class="row row-cols-1 row-cols-md-2 g-4">
            {{range .}}
            <div class="col">
              <div class="card shadow-sm h-100">
                <div class="card-body">
                  <h5 class="card-title">{{.Id}}</h5>
                  <p class="card-text"><small class="text-muted">Uploaded: {{.UploadTime}}</small></p>
                  <a href="/videos/{{.EscapedId}}" class="btn btn-primary">Watch</a>
                </div>
              </div>
            </div>
            {{else}}
            <div class="col">
              <div class="alert alert-info">No videos uploaded yet.</div>
            </div>
            {{end}}
          </div>
        </div>
        <div class="col-md-4">
          <h2>Upload an MP4 Video</h2>
          <form action="/upload" method="post" enctype="multipart/form-data" class="p-4 border rounded bg-white shadow-sm">
            <div class="mb-3">
              <input type="file" class="form-control" name="file" accept="video/mp4" required />
            </div>
            <button type="submit" class="btn btn-success w-100">Upload</button>
          </form>
        </div>
      </div>
    </div>
  </body>
</html>
`

const videoHTML = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>{{.Id}} - TritonTube</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://cdn.dashjs.org/latest/dash.all.min.js"></script>
  </head>
  <body class="bg-light">
    <nav class="navbar navbar-expand-lg navbar-dark bg-primary mb-4">
      <div class="container">
        <a class="navbar-brand" href="/">TritonTube</a>
      </div>
    </nav>
    <div class="container">
      <div class="card shadow-sm mx-auto" style="max-width: 700px;">
        <div class="card-body">
          <h2 class="card-title">{{.Id}}</h2>
          <p class="text-muted">Uploaded at: {{.UploadedAt}}</p>
          <div class="ratio ratio-16x9 mb-3">
            <video id="dashPlayer" controls class="w-100 h-100 bg-dark rounded"></video>
          </div>
          <script>
            var url = "/content/{{.Id}}/manifest.mpd";
            var player = dashjs.MediaPlayer().create();
            player.initialize(document.querySelector("#dashPlayer"), url, false);
          </script>
          <a href="/" class="btn btn-outline-primary mt-3">Back to Home</a>
        </div>
      </div>
    </div>
  </body>
</html>
`
