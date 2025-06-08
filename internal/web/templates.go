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
            <button id="toggleSidebar" class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-yt-hover">
              <svg class="w-6 h-6 text-gray-600 dark:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
              </svg>
            </button>
            <a href="/" class="flex ml-2 md:mr-24 items-center">
              <h1 class="text-xl font-bold text-gray-800 dark:text-white ml-2">TritonTube</h1>
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

    <!-- Sidebar -->
    <aside id="sidebar" class="fixed top-0 left-0 z-40 w-64 h-screen pt-16 transition-transform -translate-x-full bg-white dark:bg-yt-sidebar border-r dark:border-gray-700 md:translate-x-0">
      <div class="h-full px-3 py-4 overflow-y-auto">
        <ul class="space-y-2">
          <li>
            <a href="/" class="flex items-center p-2 text-gray-900 dark:text-white rounded-lg hover:bg-gray-100 dark:hover:bg-yt-hover group">
              <svg class="w-5 h-5 text-gray-600 dark:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path>
              </svg>
              <span class="ml-3">Home</span>
            </a>
          </li>
          <li>
            <button type="button" class="flex items-center w-full p-2 text-gray-900 dark:text-white rounded-lg hover:bg-gray-100 dark:hover:bg-yt-hover group">
              <svg class="w-5 h-5 text-gray-600 dark:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
              </svg>
              <span class="ml-3">Upload Video</span>
            </button>
          </li>
        </ul>
      </div>
    </aside>

    <!-- Main content -->
    <div class="p-4 md:ml-64 pt-20">
      <!-- Upload Section -->
      <div class="mb-8 p-6 bg-white dark:bg-yt-sidebar rounded-lg shadow-md">
        <h2 class="text-xl font-semibold text-gray-800 dark:text-white mb-4">Upload a Video</h2>
        <form action="/upload" method="post" enctype="multipart/form-data" class="space-y-4">
          <div class="flex items-center justify-center w-full">
            <label class="w-full flex flex-col items-center px-4 py-6 bg-white dark:bg-yt-hover rounded-lg shadow-lg tracking-wide border-2 border-dashed border-gray-300 dark:border-gray-600 cursor-pointer hover:border-blue-500 dark:hover:border-blue-500">
              <svg class="w-8 h-8 text-gray-500 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
              </svg>
              <span class="mt-2 text-sm text-gray-600 dark:text-gray-300">Drop your MP4 video here or click to browse</span>
              <input type="file" name="file" accept="video/mp4" required class="hidden" />
            </label>
          </div>
          <div class="flex justify-end">
            <button type="submit" class="bg-blue-600 text-white px-6 py-2 rounded-full hover:bg-blue-700 transition duration-300">
              Upload Video
            </button>
          </div>
        </form>
      </div>

      <!-- Videos Grid -->
      <div class="bg-white dark:bg-yt-sidebar rounded-lg shadow-md p-6">
        <h2 class="text-xl font-semibold text-gray-800 dark:text-white mb-6">Videos</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {{range .}}
          <a href="/videos/{{.EscapedId}}" class="block group">
            <div class="aspect-video bg-gray-100 dark:bg-yt-hover rounded-xl mb-3 overflow-hidden">
              <div class="w-full h-full flex items-center justify-center">
                <svg class="w-12 h-12 text-gray-400 dark:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
              </div>
            </div>
            <div>
              <h3 class="text-base font-medium text-gray-900 dark:text-white group-hover:text-blue-500 line-clamp-2">{{.Id}}</h3>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{.UploadTime}}</p>
            </div>
          </a>
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
      // Toggle sidebar
      document.getElementById('toggleSidebar').addEventListener('click', () => {
        const sidebar = document.getElementById('sidebar');
        sidebar.classList.toggle('-translate-x-full');
      });

      // Toggle dark mode
      document.getElementById('toggleDarkMode').addEventListener('click', () => {
        document.documentElement.classList.toggle('dark');
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
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Video Player -->
        <div class="lg:col-span-2">
          <div class="aspect-video bg-black rounded-xl overflow-hidden">
            <video id="dashPlayer" controls class="w-full h-full"></video>
          </div>
          
          <div class="mt-4 bg-white dark:bg-yt-sidebar rounded-xl p-4">
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
        </div>

        <!-- Sidebar - Can be used for related videos in the future -->
        <div class="hidden lg:block">
          <div class="bg-white dark:bg-yt-sidebar rounded-xl p-4">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">More Videos</h2>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              Coming soon...
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
