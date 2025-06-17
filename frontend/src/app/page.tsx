'use client'

import React, { useEffect, useState } from 'react'
import Navigation from '@/components/Navigation'
import UploadForm from '@/components/UploadForm'
import VideoCard from '@/components/VideoCard'

interface Video {
  id: string
  uploadedAt: string
}

export default function HomePage() {
  const [videos, setVideos] = useState<Video[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchVideos()
  }, [])

  const fetchVideos = async () => {
    try {
      const response = await fetch('/api/videos')
      if (response.ok) {
        const data = await response.json()
        if (data.success) {
          setVideos(Array.isArray(data.data) ? data.data : [])
        }
      }
    } catch (error) {
      console.error('Error fetching videos:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <Navigation />
      
      {/* Main content */}
      <div className="p-4 pt-20">
        {/* Upload Section */}
        <UploadForm />

        {/* Videos Grid */}
        <div className="bg-gray-100 dark:bg-gray-800 rounded-lg shadow-md p-6">
          <h2 className="text-xl font-semibold text-gray-800 dark:text-white mb-6">Videos</h2>
          {loading ? (
            <div className="flex justify-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
          ) : videos.length > 0 ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
              {videos.map((video) => (
                <VideoCard 
                  key={video.id} 
                  id={video.id} 
                  uploadedAt={video.uploadedAt} 
                />
              ))}
            </div>
          ) : (
            <div className="col-span-full flex flex-col items-center justify-center py-12 text-gray-500 dark:text-gray-400">
              <svg className="w-16 h-16 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
              </svg>
              <p className="text-lg">No videos uploaded yet</p>
              <p className="text-sm mt-2">Your uploaded videos will appear here</p>
            </div>
          )}
        </div>
      </div>
    </>
  )
}
