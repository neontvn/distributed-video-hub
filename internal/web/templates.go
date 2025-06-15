// Lab 7, 8, 9: Use these templates to render the web pages

package web

const indexHTML = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TritonTube</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
      tailwind.config = {
        darkMode: 'class',
        theme: {
          extend: {
            colors: {
              'yt-dark': '#0F0F0F',
              'yt-light': '#FFFFFF',
              'yt-sidebar': '#212121',
              'yt-hover': '#272727'
            }
          }
        }
      }

      // Check system dark mode preference
      if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        document.documentElement.classList.add('dark')
      }
    </script>
  </head>
  <body class="bg-gray-50 dark:bg-yt-dark min-h-screen">
    <!-- Navigation -->
    <nav class="fixed top-0 z-50 w-full bg-white dark:bg-yt-dark border-b dark:border-gray-700">
      <div class="px-3 py-3 lg:px-5 lg:pl-3">
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <a href="/" class="flex items-center">
              <h1 class="text-xl font-bold text-gray-800 dark:text-white">TritonTube</h1>
            </a>
          </div>
          <button id="toggleDarkMode" class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-yt-hover">
            <svg class="w-6 h-6 text-gray-600 dark:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path>
            </svg>
          </button>
        </div>
      </div>
    </nav>

    <!-- Main content -->
    <div class="p-4 pt-20">
      <!-- Upload Section -->
      <div class="mb-8 p-6 bg-white dark:bg-yt-sidebar rounded-lg shadow-md">
        <h2 class="text-xl font-semibold text-gray-800 dark:text-white mb-4">Upload a Video</h2>
        <div class="mb-4">
          <p class="text-sm text-gray-800 dark:text-white">
            <span class="font-medium">Note:</span> For demo purposes, please upload videos smaller than 10MB. Larger files may take significant time to process.
          </p>
        </div>
        <form action="/upload" method="post" enctype="multipart/form-data" class="space-y-4" id="uploadForm">
          <div class="flex items-center justify-center w-full">
            <label class="w-full flex flex-col items-center px-4 py-6 bg-white dark:bg-yt-hover rounded-lg shadow-lg tracking-wide border-2 border-dashed border-gray-300 dark:border-gray-600 cursor-pointer hover:border-blue-500 dark:hover:border-blue-500">
              <svg class="w-8 h-8 text-gray-500 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
              </svg>
              <span id="file-name" class="mt-2 text-sm text-gray-600 dark:text-gray-300">Drop your MP4 video here or click to browse</span>
              <input type="file" name="file" accept="video/mp4" required class="hidden" onchange="updateFileName(this)" />
            </label>
          </div>
          <div class="flex justify-end">
            <button type="submit" id="uploadButton" class="bg-blue-600 text-white px-6 py-2 rounded-full hover:bg-blue-700 transition duration-300 flex items-center space-x-2">
              <span id="uploadText">Upload Video</span>
              <svg id="uploadSpinner" class="hidden animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </button>
          </div>
        </form>
      </div>

      <!-- Videos Grid -->
      <div class="bg-white dark:bg-yt-sidebar rounded-lg shadow-md p-6">
        <h2 class="text-xl font-semibold text-gray-800 dark:text-white mb-6">Videos</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {{range .}}
          <div class="relative group">
            <a href="/videos/{{.EscapedId}}" class="block">
              <div class="aspect-video bg-gray-100 dark:bg-yt-hover rounded-xl mb-3 overflow-hidden">
                <img src="/content/{{.EscapedId}}/thumbnail.jpg" class="w-full h-full object-cover" onerror="this.onerror=null; this.parentElement.innerHTML='<div class=\'w-full h-full flex items-center justify-center\'><svg class=\'w-12 h-12 text-gray-400 dark:text-gray-600\' fill=\'none\' stroke=\'currentColor\' viewBox=\'0 0 24 24\'><path stroke-linecap=\'round\' stroke-linejoin=\'round\' stroke-width=\'2\' d=\'M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z\'></path><path stroke-linecap=\'round\' stroke-linejoin=\'round\' stroke-width=\'2\' d=\'M21 12a9 9 0 11-18 0 9 9 0 0118 0z\'></path></svg></div>'">
              </div>
              <div>
                <h3 class="text-base font-medium text-gray-900 dark:text-white group-hover:text-blue-500 line-clamp-2">{{.Id}}</h3>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{.UploadTime}}</p>
              </div>
            </a>
            <form action="/delete/{{.EscapedId}}" method="post" class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <button type="submit" onclick="return confirm('Are you sure you want to delete this video?')" class="p-2 bg-red-600 text-white rounded-full hover:bg-red-700 transition-colors">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                </svg>
              </button>
            </form>
          </div>
          {{else}}
          <div class="col-span-full flex flex-col items-center justify-center py-12 text-gray-500 dark:text-gray-400">
            <svg class="w-16 h-16 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
            </svg>
            <p class="text-lg">No videos uploaded yet</p>
            <p class="text-sm mt-2">Your uploaded videos will appear here</p>
          </div>
          {{end}}
        </div>
      </div>
    </div>

    <script>
      // Toggle dark mode
      document.getElementById('toggleDarkMode').addEventListener('click', () => {
        document.documentElement.classList.toggle('dark');
      });

      // Update filename display
      function updateFileName(input) {
        const fileNameDisplay = document.getElementById('file-name');
        if (input.files.length > 0) {
          fileNameDisplay.textContent = input.files[0].name;
        } else {
          fileNameDisplay.textContent = 'Drop your MP4 video here or click to browse';
        }
      }

      // Form submission handler
      document.getElementById('uploadForm').addEventListener('submit', function(e) {
        const button = document.getElementById('uploadButton');
        const text = document.getElementById('uploadText');
        const spinner = document.getElementById('uploadSpinner');
        
        // Disable the button and show spinner
        button.disabled = true;
        button.classList.add('opacity-75', 'cursor-not-allowed');
        text.textContent = 'Uploading...';
        spinner.classList.remove('hidden');
      });

      // Add file size check
      document.querySelector('input[type="file"]').addEventListener('change', function(e) {
        const file = e.target.files[0];
        if (file) {
          const maxSize = 10 * 1024 * 1024; // 10MB in bytes
          if (file.size > maxSize) {
            alert('File size exceeds 10MB limit. Please choose a smaller file for demo purposes.');
            e.target.value = ''; // Clear the file input
            document.getElementById('file-name').textContent = 'Drop your MP4 video here or click to browse';
          }
        }
      });
    </script>
  </body>
</html>
`

const videoHTML = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Id}} - TritonTube</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://cdn.dashjs.org/latest/dash.all.min.js"></script>
    <script>
      tailwind.config = {
        darkMode: 'class',
        theme: {
          extend: {
            colors: {
              'yt-dark': '#0F0F0F',
              'yt-light': '#FFFFFF',
              'yt-sidebar': '#212121',
              'yt-hover': '#272727'
            }
          }
        }
      }

      if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        document.documentElement.classList.add('dark')
      }
    </script>
  </head>
  <body class="bg-gray-50 dark:bg-yt-dark min-h-screen">
    <!-- Navigation -->
    <nav class="fixed top-0 z-50 w-full bg-white dark:bg-yt-dark border-b dark:border-gray-700">
      <div class="px-3 py-3 lg:px-5 lg:pl-3">
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <a href="/" class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-yt-hover">
              <svg class="w-6 h-6 text-gray-600 dark:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
              </svg>
            </a>
            <h1 class="text-xl font-bold text-gray-800 dark:text-white ml-2">TritonTube</h1>
          </div>
          <button id="toggleDarkMode" class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-yt-hover">
            <svg class="w-6 h-6 text-gray-600 dark:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path>
            </svg>
          </button>
        </div>
      </div>
    </nav>

    <!-- Main content -->
    <div class="pt-16 px-4 md:px-8 max-w-[1800px] mx-auto">
      <div class="grid grid-cols-1">
        <!-- Video Player -->
        <div>
          <div class="aspect-video bg-black rounded-xl overflow-hidden">
            <video id="dashPlayer" controls class="w-full h-full"></video>
          </div>
          
          <div class="mt-4 bg-white dark:bg-yt-sidebar rounded-xl p-4">
            <div class="flex justify-between items-start">
              <div>
                <h1 class="text-xl font-bold text-gray-900 dark:text-white mb-2">{{.Id}}</h1>
                <div class="flex items-center space-x-2">
                  <div class="h-9 w-9 rounded-full bg-gray-200 dark:bg-yt-hover flex items-center justify-center">
                    <svg class="w-5 h-5 text-gray-500 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                    </svg>
                  </div>
                  <div>
                    <p class="text-sm text-gray-500 dark:text-gray-400">Uploaded at: {{.UploadedAt}}</p>
                  </div>
                </div>
              </div>
              <form action="/delete/{{.Id}}" method="post">
                <button type="submit" onclick="return confirm('Are you sure you want to delete this video?')" class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors flex items-center space-x-2">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                  </svg>
                  <span>Delete Video</span>
                </button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>

    <script>
      // Initialize video player
      var url = "/content/{{.Id}}/manifest.mpd";
      var player = dashjs.MediaPlayer().create();
      player.initialize(document.querySelector("#dashPlayer"), url, false);

      // Toggle dark mode
      document.getElementById('toggleDarkMode').addEventListener('click', () => {
        document.documentElement.classList.toggle('dark');
      });
    </script>
  </body>
</html>
`
