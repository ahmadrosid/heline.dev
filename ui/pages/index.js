import Head from "next/head";
import Router from "next/router";
import { useEffect, useRef } from "react";
import { BiSearchAlt } from "react-icons/bi";

export default function Home() {
  const inputRef = useRef();

  useEffect(() => {
    inputRef.current.focus();
  }, []);

  const handleOnChange = (e) => {
    Router.push("/search?q=" + e.target.value);
  };

  return (
    <>
      <Head>
        <meta name="viewport" content="width=device-width,initial-scale=1" />
        <title>Heline.dev - Code search for modern developer 🚀.</title>
        <link rel="icon" type="image/png" href="/favicon.png" />
        {/* <script
          defer
          data-domain="heline.dev"
          src="https://plausible.io/js/plausible.js"
        ></script> */}
      </Head>

      <div className="bg-gray-50 min-h-screen">
        <div class="absolute inset-0 bg-[url(/grid.svg)] bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))]"></div>
        <div className="max-w-7xl min-h-[calc(100vh-4.75rem)] mx-auto py-14 sm:py-20 md:py-24 px-4 sm:px-6 lg:px-8 sm:flex sm:flex-col sm:items-center relative">
          <div className="text-center space-y-4 w-full max-w-xl p-8">
            <h1 className="text-6xl font-extrabold text-gray-900 sm:text-center inline-flex items-center select-none leading-tight">
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-sky-500 to-blue-600">
                heline
              </span>
              <span>.</span>
              <span className="text-gray-700">dev</span>
            </h1>
            <p className="text-lg font-normal text-gray-400 tracking-wider">
              Search Engine for Modern Developers.
            </p>
            <div className="self-center py-2 mt-8 bg-white shadow-md w-full max-w-xl rounded-xl overflow-hidden content border">
              <div className="px-4 flex items-center">
                <span className="text-gray-500 select-none text-xl pr-3">
                  <BiSearchAlt />
                </span>
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
                  className="py-2 flex-grow text-gray-900 border-none outline-none focus:outline-none focus:ring-0 autofill:shadow-fill-white"
                  aria-label="search"
                />
              </div>
            </div>
            <div className="text-gray-500 font-light font-sm">
              Created by{" "}
              <a
                className="text-sky-500 font-medium hover:underline"
                target="_blank"
                href="https://ahmadrosid.com"
              >
                @_ahmadrosid
              </a>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
