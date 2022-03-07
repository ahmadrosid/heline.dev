import renderArray from '../lib/render-array'

function DocHighlight({ data }) {
  return (
    <div className="w-full max-w-7xl mx-auto flex">
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
      <div className='px-4'>
        {renderArray(data.hits.map(item => (
          <div className='py-2'>
            <h2 className='text-gray-800 font-medium pb-1 text-2xl'>{item.title?.raw}</h2>
            <div className='bg-white p-2 border border-zinc-300 rounded text-sm leading-normal w-full text-gray-700'>
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