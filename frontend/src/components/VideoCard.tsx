'use client'

import React from 'react'
import Link from 'next/link'

interface VideoCardProps {
  id: string
  uploadedAt: string
}

export default function VideoCard({ id, uploadedAt }: VideoCardProps) {
  const handleDelete = async (e: React.MouseEvent) => {
    e.preventDefault()
    if (!confirm('Are you sure you want to delete this video?')) {
      return
    }

    try {
      const response = await fetch(`/api/delete/${id}`, {
        method: 'DELETE',
      })

      if (response.ok) {
        // Refresh the page to update the video list
        window.location.reload()
      } else {
        const error = await response.json()
        alert(`Delete failed: ${error.error}`)
      }
    } catch {
      alert('Delete failed. Please try again.')
    }
  }

  return (
    <div className="relative group">
      <Link href={`/videos/${id}`} className="block">
        <div className="aspect-video bg-gray-100 dark:bg-yt-hover rounded-xl mb-3 overflow-hidden">
          <img 
            src={`/api/content/${id}/thumbnail.jpg`} 
            className="w-full h-full object-cover" 
            onError={(e) => {
              const target = e.target as HTMLImageElement
              target.onerror = null
              target.parentElement!.innerHTML = `
                <div class="w-full h-full flex items-center justify-center">
                  <svg class="w-12 h-12 text-gray-400 dark:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                </div>
              `
            }}
            alt={id}
          />
        </div>
        <div>
          <h3 className="text-base font-medium text-gray-900 dark:text-white group-hover:text-blue-500 line-clamp-2">{id}</h3>
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{uploadedAt}</p>
        </div>
      </Link>
      <button 
        onClick={handleDelete}
        className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity p-2 bg-red-600 text-white rounded-full hover:bg-red-700 transition-colors"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
        </svg>
      </button>
    </div>
  )
} 