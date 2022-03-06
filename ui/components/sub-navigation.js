
import { BiCode, BiNote, BiNews } from "react-icons/bi"
import { BsStackOverflow } from "react-icons/bs"

export default function SubNavigation({ updateMatchingSearch, tbm = "" }) {
  const active = (val) => {
    if (val === tbm) {
      return "flex items-center py-2 gap-2 border-b-[3px] border-emerald-500 text-emerald-600"
    } else {
      return "flex items-center py-2 gap-2 text-gray-700 border-white border-b-[3px] hover:border-emerald-500 hover:text-emerald-600 hover:transition-all duration-500"
    }
  }

  return (
    <div className='w-full flex'>
      <div className="w-full max-w-[25%]"></div>
      <div className='w-full px-8 pt-3 flex gap-5'>
        <button onClick={() => updateMatchingSearch("code")} className={active(tbm === "" ? "" : "code")}>
          <BiCode className='text-xl' />
          <span className='text-xs'>Code</span>
        </button>
        <button onClick={() => updateMatchingSearch("docs")} className={active("docs")}>
          <BiNote className='text-xl' />
          <span className='text-xs'>Documentation</span>
        </button>
        <button onClick={() => updateMatchingSearch("stf")} className={active("stf")}>
          <BsStackOverflow className='text-xl' />
          <span className='text-xs'>Stack Overflow</span>
        </button>
        <button onClick={() => updateMatchingSearch("blog")} className={active("blog")}>
          <BiNews className='text-xl' />
          <span className='text-xs'>Blog</span>
        </button>
      </div>
    </div>
  )
}