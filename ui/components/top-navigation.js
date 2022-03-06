import { BiSearchAlt } from "react-icons/bi"

export default function TopNavigation({ setVal, q }) {
  return (
    <div className='flex items-center gap-4 w-full'>
      <div className="w-full max-w-[25%]">
        <div className="text-emerald-500 flex items-center gap-x-2 px-4">
          <a href="/">
            <h1 className="text-2xl font-extrabold text-gray-900 dark:text-white sm:text-center inline-flex items-center select-none">
              <span className="bg-clip-text text-transparent bg-gradient-to-r from-emerald-500 to-teal-600">
                heline
              </span>
              <span>.</span>
              <span className="text-gray-700">dev</span>
            </h1>
          </a>
        </div>
      </div>
      <div className="flex px-4 rounded shadow-sm border bg-white items-center justify-between w-full mr-8 mx-4">
        <span className="text-gray-500 select-none text-xl pr-4"><BiSearchAlt /></span>

        <input
          onChange={(el) => {
            setVal(encodeURIComponent(el.target.value))
          }}
          autoFocus={true}
          spellCheck={false}
          defaultValue={q}
          type="text"
          placeholder="Search"
          className="py-3 flex-grow text-gray-900 dark:bg-black dark:text-white border-none outline-none focus:outline-none focus:ring-0 autofill:shadow-fill-white dark:autofill:shadow-fill-black"
        />
      </div>
    </div>
  )
}