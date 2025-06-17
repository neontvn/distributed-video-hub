'use client'

import React, { useEffect, useState, useCallback } from 'react'
import Navigation from '@/components/Navigation'
import VideoPlayer from '@/components/VideoPlayer'

interface VideoPageProps {
  params: Promise<{
    id: string
  }>
}

export default function VideoPage({ params }: VideoPageProps) {
  const [video, setVideo] = useState<{ id: string; uploadedAt: string } | null>(null)
  const [loading, setLoading] = useState(true)
  
  // Unwrap the params promise using React.use()
  const { id } = React.use(params)

  const fetchVideo = useCallback(async () => {
    try {
      const response = await fetch(`/api/videos/${id}`)
      if (response.ok) {
        const data = await response.json()
        if (data.success) {
          setVideo(data.data)
        }
      }
    } catch (error) {
      console.error('Error fetching video:', error)
    } finally {
      setLoading(false)
    }
  }, [id])

  useEffect(() => {
    fetchVideo()
  }, [fetchVideo])

  if (loading) {
    return (
      <>
        <Navigation showBackButton />
        <div className="flex justify-center items-center min-h-screen">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
      </>
    )
  }

  if (!video) {
    return (
      <>
        <Navigation showBackButton />
        <div className="flex justify-center items-center min-h-screen">
          <div className="text-center">
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">Video not found</h1>
            <p className="text-gray-600 dark:text-gray-400">The video you&apos;re looking for doesn&apos;t exist.</p>
          </div>
        </div>
      </>
    )
  }

  return (
    <>
      <Navigation showBackButton />
      <VideoPlayer videoId={video.id} uploadedAt={video.uploadedAt} />
    </>
  )
} 