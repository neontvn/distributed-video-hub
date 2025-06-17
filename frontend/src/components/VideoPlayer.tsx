'use client'

import React, { useEffect, useRef } from 'react'

interface VideoPlayerProps {
  videoId: string
  uploadedAt: string
}

interface DashJSPlayer {
  destroy(): void
  initialize(element: HTMLVideoElement, url: string, autoplay: boolean): void
  on(event: string, callback: () => void): void
}

interface DashJS {
  MediaPlayer(): {
    create(): DashJSPlayer
  }
}

declare global {
  interface Window {
    dashjs: DashJS
  }
}

export default function VideoPlayer({ videoId, uploadedAt }: VideoPlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null)
  const playerRef = useRef<DashJSPlayer | null>(null)

  useEffect(() => {
    // Load DASH.js script dynamically
    const script = document.createElement('script')
    script.src = 'https://cdn.dashjs.org/latest/dash.all.min.js'
    script.onload = () => {
      if (videoRef.current && window.dashjs) {
        // Destroy existing player if it exists
        if (playerRef.current) {
          playerRef.current.destroy()
        }
        
        // Create new player
        playerRef.current = window.dashjs.MediaPlayer().create()
        playerRef.current.initialize(videoRef.current, `/api/content/${videoId}/manifest.mpd`, false)
        // Ensure video starts from the beginning
        playerRef.current.on('streamInitialized', function () {
          if (videoRef.current) {
            videoRef.current.currentTime = 0
          }
        })
      }
    }
    document.head.appendChild(script)

    return () => {
      // Cleanup: destroy player and remove script
      if (playerRef.current) {
        try {
          playerRef.current.destroy()
        } catch {
          console.log('Player cleanup error')
        }
        playerRef.current = null
      }
      
      if (script.parentNode) {
        document.head.removeChild(script)
      }
    }
  }, [videoId])

  const handleDelete = async () => {
    if (!confirm('Are you sure you want to delete this video?')) {
      return
    }

    try {
      const response = await fetch(`/api/delete/${videoId}`, {
        method: 'DELETE',
      })

      if (response.ok) {
        // Redirect to home page after successful deletion
        window.location.href = '/'
      } else {
        const error = await response.json()
        alert(`Delete failed: ${error.error}`)
      }
    } catch {
      alert('Delete failed. Please try again.')
    }
  }

  return (
    <div className="pt-16 px-4 md:px-8 w-full">
      <div className="grid grid-cols-1">
        {/* Video Player */}
        <div>
          <div className="aspect-[16/9] bg-black rounded-xl overflow-hidden max-h-[60vh] w-full">
            <video ref={videoRef} controls className="w-full h-full"></video>
          </div>
          
          <div className="mt-4 bg-gray-100 dark:bg-gray-800 rounded-xl p-4">
            <div className="flex justify-between items-start">
              <div>
                <h1 className="text-xl font-bold text-gray-900 dark:text-white mb-2">{videoId}</h1>
                <div className="flex items-center space-x-2">
                  <div className="h-9 w-9 rounded-full bg-gray-200 dark:bg-yt-hover flex items-center justify-center">
                    <svg className="w-5 h-5 text-gray-500 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                    </svg>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Uploaded at: {uploadedAt}</p>
                  </div>
                </div>
              </div>
              <button 
                onClick={handleDelete}
                className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors flex items-center space-x-2"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                </svg>
                <span>Delete Video</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
} 