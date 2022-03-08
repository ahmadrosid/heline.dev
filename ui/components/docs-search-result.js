import { useState, useEffect } from "react";
import { useRouter } from 'next/router'
import renderArray from '../lib/render-array'
import useGetDocumentByID from '../lib/useGetDocumentByID'
import Modal from './modal'
import { CgCloseO } from 'react-icons/cg'

function DocHighlight({ data }) {
  const router = useRouter()
  const { q = "" } = router.query

  const [detail, setHits] = useState(null)
  const [notFound, setNotFound] = useState(null)
  const [isLoading, setIsLoading] = useState(null)
  const fetchDocumentByID = useGetDocumentByID({ setHits, setNotFound, setIsLoading })

  const getDetailDocs = ({ id, link }) => {
    if (link.raw !== "") {
      let url = link.raw
      if (!url.includes("https://")) {
        url = "https://" + link.raw
      }

      if (window) {
        window.open(url, '_blank').focus();
        return
      }
    }

    fetchDocumentByID(id.raw)
  }

  const useKeyPress = (targetKey) => {
    const [keyPressed, setKeyPressed] = useState(false);
    function downHandler({ key }) {
      if (key === targetKey) {
        setKeyPressed(true);
        setHits(null)
      }
    }
    const upHandler = ({ key }) => {
      if (key === targetKey) {
        setKeyPressed(false);
      }
    };
    useEffect(() => {
      window.addEventListener("keydown", downHandler);
      window.addEventListener("keyup", upHandler);
      return () => {
        window.removeEventListener("keydown", downHandler);
        window.removeEventListener("keyup", upHandler);
      };
    }, []);
    return keyPressed;
  }

  const _ = useKeyPress("Escape")
  
  return (
    <>
      {detail && (
        <Modal>
          <div className='relative'>
            <div className='flex justify-between items-center bg-gray-100 -m-6 p-2 rounded-t-md'>
              <div className='px-2'>
                <h1 className='font-medium text-xl'>{detail.title}</h1>
              </div>
              <div className='cursor-pointer p-2' onClick={() => setHits(null)}>
                <CgCloseO className='text-2xl text-gray-700' />
              </div>
            </div>
            <div
              className='modal-docset'
              dangerouslySetInnerHTML={{ __html: detail.content.join("").replace(q, `<mark>${q}</mark>`) }}></div>
          </div>
        </Modal>
      )}
      <div className="w-full max-w-7xl mx-auto flex py-2">
        <div className="w-full min-w-[250px] max-w-[25%] py-4 space-y-4 pl-4">
          <div className="space-y-1">
            <h3 className="text-gray-800 font-medium uppercase">Document</h3>
            <div className="py-2 space-y-1">
              {renderArray(
                data.facets?.document?.buckets.map((item, index) => (
                  <div className="flex justify-between items-center text-gray-600 pr-1">
                    <div className="inline-flex gap-2 items-center truncate">
                      <input id={item.val} className="p-2" type="checkbox" />
                      <label className='truncate' htmlFor={item.val}>{item.val}</label>
                    </div>
                    <div className="text-sm">{item.count}</div>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
        <div className='px-8'>
          {renderArray(data.hits.map(item => (
            <div className='py-2 space-y-1 pb-2'>
              {item.link.raw && (
                <a href={`https://${item.link.raw}`} target="_blank" className="text-gray-500 text-sm truncate max-w-[250px] inline-block -mb-2 cursor-pointer">{item.id.raw}</a>
              )}

              {!item.link.raw && (
                <a onClick={() => getDetailDocs(item)} className="text-gray-500 text-sm truncate max-w-[250px] inline-block -mb-2 cursor-pointer">{item.id.raw}</a>
              )}
              <h2 className='text-gray-800 font-mono font-semibold text-xl'>{item.title?.raw}</h2>
              <div
                onClick={() => getDetailDocs(item)}
                className='highlight-docset cursor-pointer'>
                {renderArray(item.content.snippet.map((source, i) => (
                  <>
                    <div dangerouslySetInnerHTML={{ __html: source }}></div>
                    {(i < (item.content.snippet.length - 1)) && <div className="bg-zinc-100 h-4 my-[8px] -mx-2"></div>}
                  </>
                )))}

              </div>
            </div>
          )))}
        </div>
      </div>
    </>
  )
}

export default function DocSearchResult({ hits, isLoading = false }) {
  return (
    <>
      {isLoading && (
        <div className="flex flex-col">
          <div className="relative w-full bg-gray-200">
            <div style={{ width: "100%" }} className="absolute top-0 h-1 shim-red"></div>
          </div>
        </div>
      )}

      {!isLoading && (<div className='h-1' />)}

      {hits && <DocHighlight data={hits} />}
    </>
  )
}