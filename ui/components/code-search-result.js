import renderArray from '../lib/render-array'

export default function CodeSearchResult({ hits, filter, updateFilter }) {

  const updateFilterRepo = (index) => {
    if (hits.facets.repo.buckets.length == 0) {
      return
    }
    updateFilter('repo', index)
  }

  const updateFilterPath = (index) => {
    if (hits.facets.path.buckets.length == 0) {
      return
    }
    updateFilter('path', index)
  }

  const updateFilterLang = (index) => {
    if (hits.facets.lang.buckets.length == 0) {
      return
    }
    updateFilter('lang', index)
  }

  const getPath = (val) => {
    let path = String(val)
    if (path.includes("/")) {
      let paths = path.split("/")
      return paths.slice(Math.max(paths.length - 2, 1)).join("/")
    }

    return val
  }

  return (
    <>
      {hits && (
        <div className="w-full max-w-7xl mx-auto flex">
          <div className="w-full min-w-[250px] max-w-[25%] py-8 space-y-4 pl-4">
            <div className="space-y-2">
              <h3 className="text-gray-800 uppercase">Repository</h3>
              <div className="py-2 space-y-1">
                {renderArray(
                  hits.facets?.repo?.buckets.map((item, index) => (
                    <div className="flex justify-between items-center text-gray-600 pr-1">
                      <div className="inline-flex gap-2 items-center truncate">
                        <input onChange={() => updateFilterRepo(index)} id={item.val} className="p-2" type="checkbox" checked={filter.repo?.includes(item.val)} />
                        <label className='truncate' htmlFor={item.val}>{item.val}</label>
                      </div>
                      <div className="text-sm">{item.count}</div>
                    </div>
                  ))
                )}
              </div>
            </div>
            <div className="space-y-2 pt-2">
              <h3 className="text-gray-800 uppercase">Path</h3>
              <div className="py-2 space-y-1">
                {renderArray(
                  hits.facets?.path?.buckets.map((item, index) => (
                    <div className="flex justify-between items-center text-gray-600 pr-1">
                      <div className="flex gap-2 items-center truncate">
                        <input onChange={() => updateFilterPath(index)} className="p-2" type="checkbox" checked={filter.path?.includes(item.val)} />
                        <label className='truncate'>{getPath(item.val)}</label>
                      </div>
                      <div className="text-sm">{item.count}</div>
                    </div>
                  ))
                )}
              </div>
            </div>
            <div className="space-y-2 pt-2">
              <h3 className="text-gray-800 uppercase">Language</h3>
              <div className="py-2 space-y-1">
                {renderArray(
                  hits.facets?.lang?.buckets.map((item, index) => (
                    <div className="flex justify-between items-center text-gray-600 pr-1">
                      <div className="flex gap-2 items-center">
                        <input onChange={() => updateFilterLang(index)} className="p-2" type="checkbox" checked={filter.lang?.includes(item.val)} />
                        <label>{item.val}</label>
                      </div>
                      <div className="text-sm">{item.count}</div>
                    </div>
                  ))
                )}
              </div>
            </div>
          </div>
          <div className="w-full max-w-[75%] p-8 pr-6">
            <div className="pb-2">
              <p className="text-gray-700">Total: {hits.total}</p>
            </div>
            {renderArray(
              hits.hits?.map(item => {
                if (item.content.snippet === null) {
                  return;
                }

                return (
                  <div className="py-2">
                    <div>
                      <a target="_blank" href={`https://github.com/${item.repo.raw}`} className="flex gap-2 items-center">
                        <img className="repo-avatar rounded-full border" src={`https://avatars.githubusercontent.com/u/${item.owner_id.raw}?s=60&amp;v=4`} alt="" width="22" height="22" />
                        <span className="text-gray-700 font-light">{item.repo.raw}</span>
                      </a>

                      <a target="_blank" href={`https://github.com/${item.repo.raw}/blob/${item.branch.raw}/${item.file_id.raw.split("/").slice(4, 100).join("/")}`} className="flex gap-1 items-center">
                        <span className="text-emerald-500 pl-8 truncate">{item.file_id.raw.split("/").slice(2, 100).join("/")}</span>
                      </a>
                    </div>
                    <div className="border border-zinc-200 rounded-md bg-white p-2 my-2">
                      {renderArray(
                        item.content.snippet?.map((content, parentIndex) => {
                          if (content.length === 0) return;

                          let contents = []
                          let index = 0;
                          let chunk = "";
                          content.split("\n").forEach((item) => {
                            index++;
                            chunk = chunk + item;
                            if (index === 4) {
                              if (chunk.includes("<mark>")) {
                                contents.push(chunk);
                              }
                              chunk = "";
                              index = 0;
                            }
                          })

                          if (chunk !== "" && chunk.includes("<mark>")) {
                            contents.push(chunk)
                          }

                          if (contents.length == 0) return;

                          // Take only 3 element to render!
                          contents = contents.slice(0, 3);

                          return (
                            <>
                              {renderArray(
                                contents.map((source, i) => {
                                  return (
                                    <>
                                      <table className="highlight-table">
                                        <tbody dangerouslySetInnerHTML={{ __html: source }}></tbody>
                                      </table>
                                      {(i < contents.length && parentIndex < (item.content.snippet.length - 1)) && <div className="bg-zinc-100 h-8 my-[8px] -mx-2"></div>}
                                    </>
                                  )
                                })
                              )}
                            </>
                          )
                        })
                      )}

                    </div>
                  </div>
                )
              })
            )}
          </div>
        </div>
      )}
    </>
  )
}