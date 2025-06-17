'use client'

import React, { useState } from 'react'

export default function UploadForm() {
  const [isUploading, setIsUploading] = useState(false)
  const [fileName, setFileName] = useState('Drop your MP4 video here or click to browse')

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      const maxSize = 10 * 1024 * 1024 // 10MB in bytes
      if (file.size > maxSize) {
        alert('File size exceeds 10MB limit. Please choose a smaller file for demo purposes.')
        e.target.value = '' // Clear the file input
        setFileName('Drop your MP4 video here or click to browse')
        return
      }
      setFileName(file.name)
    } else {
      setFileName('Drop your MP4 video here or click to browse')
    }
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setIsUploading(true)

    const formData = new FormData(e.currentTarget)
    
    try {
      const response = await fetch('/api/upload', {
        method: 'POST',
        body: formData,
      })

      if (response.ok) {
        // Refresh the page to show the new video
        window.location.reload()
      } else {
        const error = await response.json()
        alert(`Upload failed: ${error.error}`)
      }
    } catch {
      alert('Upload failed. Please try again.')
    } finally {
      setIsUploading(false)
    }
  }

  return (
    <div className="mb-8 p-6 bg-gray-100 dark:bg-gray-800 rounded-lg shadow-md">
      <h2 className="text-xl font-semibold text-gray-800 dark:text-white mb-4">Upload a Video</h2>
      <div className="mb-4">
        <p className="text-sm text-gray-800 dark:text-white">
          <span className="font-medium">Note:</span> For demo purposes, please upload videos smaller than 10MB. Larger files may take significant time to process.
        </p>
      </div>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="flex items-center justify-center w-full">
          <label className="w-full flex flex-col items-center px-4 py-6 bg-gray-100 dark:bg-gray-800 rounded-lg shadow-lg tracking-wide border-2 border-dashed border-gray-300 dark:border-gray-600 cursor-pointer hover:border-blue-500 dark:hover:border-blue-500">
            <svg className="w-8 h-8 text-gray-500 dark:text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
            </svg>
            <span className="mt-2 text-sm text-gray-600 dark:text-gray-300">{fileName}</span>
            <input 
              type="file" 
              name="file" 
              accept="video/mp4" 
              required 
              className="hidden" 
              onChange={handleFileChange}
            />
          </label>
        </div>
        <div className="flex justify-end">
          <button 
            type="submit" 
            disabled={isUploading}
            className="bg-blue-600 text-white px-6 py-2 rounded-full hover:bg-blue-700 transition duration-300 flex items-center space-x-2 disabled:opacity-75 disabled:cursor-not-allowed"
          >
            <span>{isUploading ? 'Uploading...' : 'Upload Video'}</span>
            {isUploading && (
              <svg className="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            )}
          </button>
        </div>
      </form>
    </div>
  )
} 