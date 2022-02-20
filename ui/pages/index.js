import Head from 'next/head'
import Router from 'next/router'
import { useEffect, useRef } from 'react'
import { BiSearchAlt } from "react-icons/bi"

export default function Home() {
  const inputRef = useRef();

  useEffect(() => {
    inputRef.current.focus()
  }, [])

  const handleOnChange = (e) => {
    Router.push("/search?q=" + e.target.value)
  }

  return (
    <>
      <Head>
        <meta name='viewport' content='width=device-width,initial-scale=1' />
        <title>Heline - Code search for modern developer 🚀.</title>
        <link rel='icon' type='image/png' href='/favicon.png' />
        <script defer data-domain="heline.dev" src="https://plausible.io/js/plausible.js"></script>
      </Head>

      <div className="max-w-7xl min-h-[calc(100vh-4.75rem)] mx-auto py-14 sm:py-20 md:py-24 px-4 sm:px-6 lg:px-8 sm:flex sm:flex-col sm:items-center">
        <div className="text-center space-y-4 w-full max-w-xl p-8">
          <h1 className="text-5xl font-extrabold text-gray-900 dark:text-white sm:text-center inline-flex items-center select-none">
            <span className="bg-clip-text text-transparent bg-gradient-to-r from-emerald-500 to-teal-600">
              heline
            </span>
            <span>.</span>
            <span className="text-gray-700">dev</span>
          </h1>
          <p className="text-lg font-normal text-gray-400">
            Search Engine for Modern Developers.
          </p>
          <div className="self-center py-2 mt-8 bg-white dark:bg-black shadow-md w-full max-w-xl rounded overflow-hidden content border">
            <div className="px-4 py-2 font-mono flex items-center">
              <span className="text-gray-500 select-none text-xl pr-3"><BiSearchAlt /></span>
              <input
                ref={inputRef}
                onChange={handleOnChange}
                type="text"
                name="search-string"
                id="search-string"
                autoFocus
                autoComplete="off"
                autoCorrect="off"
                autoCapitalize="off"
                spellCheck="false"
                placeholder="Search for function, variable, snippets etc."
                className="flex-grow text-gray-900 dark:bg-black dark:text-white border-none outline-none focus:outline-none focus:ring-0 autofill:shadow-fill-white dark:autofill:shadow-fill-black"
                aria-label="search"
              />
            </div>
          </div>
          <div className="text-gray-500">
            Created by <a
              className="text-emerald-500 font-medium hover:underline"
              target="_blank"
              href="https://ahmadrosid.com">@_ahmadrosid</a
            >
          </div>
        </div>
      </div>

    </>
  )
}
