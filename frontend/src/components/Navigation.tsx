'use client'

import React from 'react'
import Link from 'next/link'

interface NavigationProps {
  showBackButton?: boolean
}

export default function Navigation({ showBackButton = false }: NavigationProps) {
  return (
    <nav className="fixed top-0 z-50 w-full bg-gray-100 dark:bg-gray-800 border-b dark:border-gray-700">
      <div className="px-3 py-3 lg:px-5 lg:pl-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center">
            {showBackButton ? (
              <Link href="/" className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-yt-hover">
                <svg className="w-6 h-6 text-gray-600 dark:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
                </svg>
              </Link>
            ) : null}
            <Link href="/" className="flex items-center">
              <h1 className="text-xl font-bold text-white-800  ml-2">TritonTube</h1>
            </Link>
          </div>
        </div>
      </div>
    </nav>
  )
} 