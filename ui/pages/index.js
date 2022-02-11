import Head from 'next/head'
import Router from 'next/router'
import { useEffect, useRef } from 'react'

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

      <div className="grid place-items-center w-full h-screen">
        <div className="text-center space-y-4 w-full max-w-xl p-8">
          <h1 className="text-6xl font-bold uppercase text-green-500">Heline</h1>
          <p className="text-lg font-normal text-gray-500">
            Code search for function names, code snippets etc 🚀.
          </p>
          <div className="py-6">
            <input
              ref={inputRef}
              onChange={handleOnChange}
              type="text"
              placeholder="Search"
              className="px-5 py-3 rounded-lg text-md w-full border-[1.5px] border-green-600 focus:outline-none"
            />
          </div>
          <div className="text-gray-500">
            Created by <a
              className="text-green-500 font-medium hover:underline"
              target="_blank"
              href="https://ahmadrosid.com">@_ahmadrosid</a
            >
          </div>
        </div>
      </div>

    </>
  )
}
